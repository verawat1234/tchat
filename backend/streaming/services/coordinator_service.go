package services

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// CoordinatorService defines the interface for horizontal scaling coordination
type CoordinatorService interface {
	// Viewer coordination
	PublishViewerJoin(streamID, viewerID uuid.UUID, serverID string) error
	PublishViewerLeave(streamID, viewerID uuid.UUID, serverID string) error
	SyncViewerCount(streamID uuid.UUID) (int, error)

	// Load balancing
	SelectLeastLoadedServer(streamID uuid.UUID) (string, error)
	RegisterServer(serverID string, capacity int) error
	HeartbeatServer(serverID string) error
	GetServerLoad(serverID string) (int, error)

	// Lifecycle
	StartListening(ctx context.Context) error
	Shutdown() error
}

// coordinatorService implements the CoordinatorService interface using Redis Pub/Sub
type coordinatorService struct {
	redis         *redis.Client
	serverID      string
	mu            sync.RWMutex
	viewerCounts  map[uuid.UUID]int           // streamID -> viewer count
	serverLoads   map[string]int              // serverID -> current load
	serverHealth  map[string]time.Time        // serverID -> last heartbeat
	subscribers   map[string]*redis.PubSub    // channel -> subscription
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// Event types for Redis Pub/Sub
const (
	EventViewerJoin  = "viewer_join"
	EventViewerLeave = "viewer_leave"
	EventServerSync  = "server_sync"
)

// ViewerEvent represents a viewer join/leave event
type ViewerEvent struct {
	Type     string    `json:"type"`
	StreamID uuid.UUID `json:"stream_id"`
	ViewerID uuid.UUID `json:"viewer_id"`
	ServerID string    `json:"server_id"`
	Time     time.Time `json:"time"`
}

// ServerHealthEvent represents a server health status event
type ServerHealthEvent struct {
	ServerID  string    `json:"server_id"`
	Load      int       `json:"load"`
	Capacity  int       `json:"capacity"`
	Time      time.Time `json:"time"`
}

// Configuration constants
const (
	MaxViewersPerServer   = 50000 // Maximum viewers per stream distributed across servers
	HeartbeatInterval     = 10 * time.Second
	HeartbeatTimeout      = 30 * time.Second
	SyncInterval          = 5 * time.Second
	RedisKeyTTL           = 1 * time.Hour
)

// NewCoordinatorService creates a new coordinator service instance
func NewCoordinatorService(redisClient *redis.Client, serverID string) CoordinatorService {
	ctx, cancel := context.WithCancel(context.Background())
	return &coordinatorService{
		redis:        redisClient,
		serverID:     serverID,
		viewerCounts: make(map[uuid.UUID]int),
		serverLoads:  make(map[string]int),
		serverHealth: make(map[string]time.Time),
		subscribers:  make(map[string]*redis.PubSub),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Channel naming functions
func streamControlChannel(streamID uuid.UUID) string {
	return fmt.Sprintf("stream:%s:control", streamID)
}

func streamEventsChannel(streamID uuid.UUID) string {
	return fmt.Sprintf("stream:%s:events", streamID)
}

func serverHealthChannel(serverID string) string {
	return fmt.Sprintf("server:%s:health", serverID)
}

func globalControlChannel() string {
	return "streaming:global:control"
}

// Redis key functions
func viewerCountKey(streamID uuid.UUID) string {
	return fmt.Sprintf("stream:%s:viewer_count", streamID)
}

func serverLoadKey(serverID string) string {
	return fmt.Sprintf("server:%s:load", serverID)
}

func serverCapacityKey(serverID string) string {
	return fmt.Sprintf("server:%s:capacity", serverID)
}

func serverRegistryKey() string {
	return "streaming:servers:registry"
}

// PublishViewerJoin publishes a viewer join event
func (c *coordinatorService) PublishViewerJoin(streamID, viewerID uuid.UUID, serverID string) error {
	event := ViewerEvent{
		Type:     EventViewerJoin,
		StreamID: streamID,
		ViewerID: viewerID,
		ServerID: serverID,
		Time:     time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal viewer join event: %w", err)
	}

	// Publish to stream events channel
	channel := streamEventsChannel(streamID)
	if err := c.redis.Publish(c.ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish viewer join event: %w", err)
	}

	// Increment local viewer count
	c.mu.Lock()
	c.viewerCounts[streamID]++
	localCount := c.viewerCounts[streamID]
	c.mu.Unlock()

	// Increment Redis viewer count
	if err := c.redis.Incr(c.ctx, viewerCountKey(streamID)).Err(); err != nil {
		return fmt.Errorf("failed to increment viewer count: %w", err)
	}

	// Set TTL on viewer count key
	c.redis.Expire(c.ctx, viewerCountKey(streamID), RedisKeyTTL)

	// Increment server load
	if err := c.redis.Incr(c.ctx, serverLoadKey(serverID)).Err(); err != nil {
		return fmt.Errorf("failed to increment server load: %w", err)
	}

	// Update local server load
	c.mu.Lock()
	c.serverLoads[serverID] = localCount
	c.mu.Unlock()

	return nil
}

// PublishViewerLeave publishes a viewer leave event
func (c *coordinatorService) PublishViewerLeave(streamID, viewerID uuid.UUID, serverID string) error {
	event := ViewerEvent{
		Type:     EventViewerLeave,
		StreamID: streamID,
		ViewerID: viewerID,
		ServerID: serverID,
		Time:     time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal viewer leave event: %w", err)
	}

	// Publish to stream events channel
	channel := streamEventsChannel(streamID)
	if err := c.redis.Publish(c.ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish viewer leave event: %w", err)
	}

	// Decrement local viewer count
	c.mu.Lock()
	if c.viewerCounts[streamID] > 0 {
		c.viewerCounts[streamID]--
	}
	localCount := c.viewerCounts[streamID]
	c.mu.Unlock()

	// Decrement Redis viewer count
	if err := c.redis.Decr(c.ctx, viewerCountKey(streamID)).Err(); err != nil {
		return fmt.Errorf("failed to decrement viewer count: %w", err)
	}

	// Decrement server load
	loadKey := serverLoadKey(serverID)
	currentLoad, err := c.redis.Get(c.ctx, loadKey).Int()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to get server load: %w", err)
	}

	if currentLoad > 0 {
		if err := c.redis.Decr(c.ctx, loadKey).Err(); err != nil {
			return fmt.Errorf("failed to decrement server load: %w", err)
		}
	}

	// Update local server load
	c.mu.Lock()
	c.serverLoads[serverID] = localCount
	c.mu.Unlock()

	return nil
}

// SyncViewerCount synchronizes viewer count from Redis
func (c *coordinatorService) SyncViewerCount(streamID uuid.UUID) (int, error) {
	count, err := c.redis.Get(c.ctx, viewerCountKey(streamID)).Int()
	if err != nil {
		if err == redis.Nil {
			// Key doesn't exist, return 0
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get viewer count: %w", err)
	}

	// Update local cache
	c.mu.Lock()
	c.viewerCounts[streamID] = count
	c.mu.Unlock()

	return count, nil
}

// SelectLeastLoadedServer selects the server with the lowest load using consistent hashing
func (c *coordinatorService) SelectLeastLoadedServer(streamID uuid.UUID) (string, error) {
	// Get all registered servers
	servers, err := c.redis.SMembers(c.ctx, serverRegistryKey()).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get server registry: %w", err)
	}

	if len(servers) == 0 {
		return "", fmt.Errorf("no servers available")
	}

	// Filter out unhealthy servers and get loads
	type serverInfo struct {
		id       string
		load     int
		capacity int
		healthy  bool
	}

	serverInfos := make([]serverInfo, 0, len(servers))
	now := time.Now()

	for _, serverID := range servers {
		// Check health
		c.mu.RLock()
		lastHeartbeat, exists := c.serverHealth[serverID]
		c.mu.RUnlock()

		healthy := exists && now.Sub(lastHeartbeat) < HeartbeatTimeout

		// Get load
		load, err := c.redis.Get(c.ctx, serverLoadKey(serverID)).Int()
		if err != nil && err != redis.Nil {
			continue // Skip this server on error
		}

		// Get capacity
		capacity, err := c.redis.Get(c.ctx, serverCapacityKey(serverID)).Int()
		if err != nil && err != redis.Nil {
			capacity = MaxViewersPerServer // Default capacity
		}

		serverInfos = append(serverInfos, serverInfo{
			id:       serverID,
			load:     load,
			capacity: capacity,
			healthy:  healthy,
		})
	}

	// Filter only healthy servers
	healthyServers := make([]serverInfo, 0, len(serverInfos))
	for _, info := range serverInfos {
		if info.healthy && info.load < info.capacity {
			healthyServers = append(healthyServers, info)
		}
	}

	if len(healthyServers) == 0 {
		return "", fmt.Errorf("no healthy servers available")
	}

	// Use consistent hashing for geographic-aware routing
	// This ensures viewers from the same region tend to land on the same server
	streamHash := hashStreamID(streamID)

	// Sort servers by ID for consistent ordering
	sort.Slice(healthyServers, func(i, j int) bool {
		return healthyServers[i].id < healthyServers[j].id
	})

	// Find the server closest to the hash value with capacity
	// First, try to find a server with low load near the hash
	selectedIdx := int(streamHash) % len(healthyServers)

	// Check if the selected server is overloaded
	if healthyServers[selectedIdx].load >= healthyServers[selectedIdx].capacity*80/100 {
		// Server is near capacity, find the least loaded one
		minLoad := healthyServers[0].load
		minIdx := 0
		for i, info := range healthyServers {
			if info.load < minLoad {
				minLoad = info.load
				minIdx = i
			}
		}
		selectedIdx = minIdx
	}

	return healthyServers[selectedIdx].id, nil
}

// hashStreamID creates a hash value from stream ID for consistent hashing
func hashStreamID(streamID uuid.UUID) uint32 {
	h := fnv.New32a()
	h.Write(streamID[:])
	return h.Sum32()
}

// RegisterServer registers a server in the cluster
func (c *coordinatorService) RegisterServer(serverID string, capacity int) error {
	// Add to server registry
	if err := c.redis.SAdd(c.ctx, serverRegistryKey(), serverID).Err(); err != nil {
		return fmt.Errorf("failed to register server: %w", err)
	}

	// Set server capacity
	if err := c.redis.Set(c.ctx, serverCapacityKey(serverID), capacity, RedisKeyTTL).Err(); err != nil {
		return fmt.Errorf("failed to set server capacity: %w", err)
	}

	// Initialize server load to 0
	if err := c.redis.Set(c.ctx, serverLoadKey(serverID), 0, RedisKeyTTL).Err(); err != nil {
		return fmt.Errorf("failed to initialize server load: %w", err)
	}

	// Update local state
	c.mu.Lock()
	c.serverLoads[serverID] = 0
	c.serverHealth[serverID] = time.Now()
	c.mu.Unlock()

	return nil
}

// HeartbeatServer sends a heartbeat for a server
func (c *coordinatorService) HeartbeatServer(serverID string) error {
	now := time.Now()

	// Get current load
	load, err := c.redis.Get(c.ctx, serverLoadKey(serverID)).Int()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to get server load: %w", err)
	}

	// Get capacity
	capacity, err := c.redis.Get(c.ctx, serverCapacityKey(serverID)).Int()
	if err != nil && err != redis.Nil {
		capacity = MaxViewersPerServer
	}

	// Publish health event
	event := ServerHealthEvent{
		ServerID: serverID,
		Load:     load,
		Capacity: capacity,
		Time:     now,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal health event: %w", err)
	}

	channel := serverHealthChannel(serverID)
	if err := c.redis.Publish(c.ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish health event: %w", err)
	}

	// Update local health status
	c.mu.Lock()
	c.serverHealth[serverID] = now
	c.serverLoads[serverID] = load
	c.mu.Unlock()

	// Refresh TTLs
	c.redis.Expire(c.ctx, serverLoadKey(serverID), RedisKeyTTL)
	c.redis.Expire(c.ctx, serverCapacityKey(serverID), RedisKeyTTL)

	return nil
}

// GetServerLoad returns the current load of a server
func (c *coordinatorService) GetServerLoad(serverID string) (int, error) {
	// Try local cache first
	c.mu.RLock()
	if load, exists := c.serverLoads[serverID]; exists {
		c.mu.RUnlock()
		return load, nil
	}
	c.mu.RUnlock()

	// Fetch from Redis
	load, err := c.redis.Get(c.ctx, serverLoadKey(serverID)).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get server load: %w", err)
	}

	// Update local cache
	c.mu.Lock()
	c.serverLoads[serverID] = load
	c.mu.Unlock()

	return load, nil
}

// StartListening starts listening to Redis Pub/Sub channels
func (c *coordinatorService) StartListening(ctx context.Context) error {
	// Subscribe to global control channel
	globalSub := c.redis.Subscribe(ctx, globalControlChannel())
	c.subscribers[globalControlChannel()] = globalSub

	// Start heartbeat goroutine
	c.wg.Add(1)
	go c.heartbeatLoop()

	// Start listener goroutine
	c.wg.Add(1)
	go c.listenToMessages(globalSub)

	return nil
}

// heartbeatLoop sends periodic heartbeats for this server
func (c *coordinatorService) heartbeatLoop() {
	defer c.wg.Done()

	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.HeartbeatServer(c.serverID); err != nil {
				// Log error but continue
				fmt.Printf("Heartbeat failed for server %s: %v\n", c.serverID, err)
			}
		}
	}
}

// listenToMessages listens to messages from Redis Pub/Sub
func (c *coordinatorService) listenToMessages(sub *redis.PubSub) {
	defer c.wg.Done()

	ch := sub.Channel()
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}

			// Handle different message types
			c.handleMessage(msg.Channel, msg.Payload)
		}
	}
}

// handleMessage processes incoming Pub/Sub messages
func (c *coordinatorService) handleMessage(channel, payload string) {
	// Try to parse as viewer event
	var viewerEvent ViewerEvent
	if err := json.Unmarshal([]byte(payload), &viewerEvent); err == nil {
		c.handleViewerEvent(viewerEvent)
		return
	}

	// Try to parse as server health event
	var healthEvent ServerHealthEvent
	if err := json.Unmarshal([]byte(payload), &healthEvent); err == nil {
		c.handleHealthEvent(healthEvent)
		return
	}
}

// handleViewerEvent handles viewer join/leave events
func (c *coordinatorService) handleViewerEvent(event ViewerEvent) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch event.Type {
	case EventViewerJoin:
		c.viewerCounts[event.StreamID]++
		c.serverLoads[event.ServerID]++
	case EventViewerLeave:
		if c.viewerCounts[event.StreamID] > 0 {
			c.viewerCounts[event.StreamID]--
		}
		if c.serverLoads[event.ServerID] > 0 {
			c.serverLoads[event.ServerID]--
		}
	}
}

// handleHealthEvent handles server health status events
func (c *coordinatorService) handleHealthEvent(event ServerHealthEvent) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.serverHealth[event.ServerID] = event.Time
	c.serverLoads[event.ServerID] = event.Load
}

// Shutdown gracefully shuts down the coordinator service
func (c *coordinatorService) Shutdown() error {
	// Cancel context to stop all goroutines
	c.cancel()

	// Close all subscribers
	for _, sub := range c.subscribers {
		if err := sub.Close(); err != nil {
			return fmt.Errorf("failed to close subscriber: %w", err)
		}
	}

	// Remove server from registry
	if err := c.redis.SRem(c.ctx, serverRegistryKey(), c.serverID).Err(); err != nil {
		return fmt.Errorf("failed to unregister server: %w", err)
	}

	// Wait for all goroutines to finish
	c.wg.Wait()

	return nil
}