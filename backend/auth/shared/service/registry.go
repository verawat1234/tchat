package service

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/gorm"
	"tchat.dev/shared/config"
)

// DefaultServiceRegistry implements ServiceRegistry
type DefaultServiceRegistry struct {
	mu         sync.RWMutex
	components map[string]ServiceComponent
}

// NewDefaultServiceRegistry creates a new service registry
func NewDefaultServiceRegistry() ServiceRegistry {
	return &DefaultServiceRegistry{
		components: make(map[string]ServiceComponent),
	}
}

// Register registers a service component
func (r *DefaultServiceRegistry) Register(name string, component ServiceComponent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.components[name]; exists {
		return fmt.Errorf("component %s already registered", name)
	}

	r.components[name] = component
	return nil
}

// Unregister removes a service component
func (r *DefaultServiceRegistry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.components[name]; !exists {
		return fmt.Errorf("component %s not found", name)
	}

	delete(r.components, name)
	return nil
}

// Get retrieves a service component
func (r *DefaultServiceRegistry) Get(name string) ServiceComponent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.components[name]
}

// List returns all registered components
func (r *DefaultServiceRegistry) List() map[string]ServiceComponent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]ServiceComponent)
	for name, component := range r.components {
		result[name] = component
	}
	return result
}

// StartAll starts all registered components
func (r *DefaultServiceRegistry) StartAll(ctx context.Context) error {
	r.mu.RLock()
	components := r.List()
	r.mu.RUnlock()

	for name, component := range components {
		if err := component.Start(ctx); err != nil {
			// Try to stop already started components
			r.stopComponents(ctx, components, name)
			return fmt.Errorf("failed to start component %s: %w", name, err)
		}
	}

	return nil
}

// StopAll stops all registered components
func (r *DefaultServiceRegistry) StopAll(ctx context.Context) error {
	r.mu.RLock()
	components := r.List()
	r.mu.RUnlock()

	return r.stopComponents(ctx, components, "")
}

// stopComponents stops components, excluding the excludeName if provided
func (r *DefaultServiceRegistry) stopComponents(ctx context.Context, components map[string]ServiceComponent, excludeName string) error {
	var lastErr error

	for name, component := range components {
		if name == excludeName {
			continue
		}

		if err := component.Stop(ctx); err != nil {
			lastErr = fmt.Errorf("failed to stop component %s: %w", name, err)
		}
	}

	return lastErr
}

// BaseServiceComponent provides a base implementation for ServiceComponent
type BaseServiceComponent struct {
	name      string
	healthy   bool
	mu        sync.RWMutex
	startTime string
}

// NewBaseServiceComponent creates a new base service component
func NewBaseServiceComponent(name string) *BaseServiceComponent {
	return &BaseServiceComponent{
		name:    name,
		healthy: false,
	}
}

// Name returns the component name
func (c *BaseServiceComponent) Name() string {
	return c.name
}

// IsHealthy returns the health status
func (c *BaseServiceComponent) IsHealthy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.healthy
}

// SetHealthy sets the health status
func (c *BaseServiceComponent) SetHealthy(healthy bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.healthy = healthy
}

// Initialize provides a default initialization (can be overridden)
func (c *BaseServiceComponent) Initialize(ctx context.Context, cfg *config.Config, db *gorm.DB) error {
	// Default implementation - services can override this
	return nil
}

// Start provides a default start implementation (can be overridden)
func (c *BaseServiceComponent) Start(ctx context.Context) error {
	c.SetHealthy(true)
	return nil
}

// Stop provides a default stop implementation (can be overridden)
func (c *BaseServiceComponent) Stop(ctx context.Context) error {
	c.SetHealthy(false)
	return nil
}

// BackgroundService represents a service that runs background tasks
type BackgroundService struct {
	*BaseServiceComponent
	taskFunc    func(ctx context.Context) error
	stopChannel chan struct{}
	done        chan struct{}
}

// NewBackgroundService creates a new background service
func NewBackgroundService(name string, taskFunc func(ctx context.Context) error) *BackgroundService {
	return &BackgroundService{
		BaseServiceComponent: NewBaseServiceComponent(name),
		taskFunc:            taskFunc,
		stopChannel:         make(chan struct{}),
		done:               make(chan struct{}),
	}
}

// Start starts the background service
func (s *BackgroundService) Start(ctx context.Context) error {
	if err := s.BaseServiceComponent.Start(ctx); err != nil {
		return err
	}

	go func() {
		defer close(s.done)

		for {
			select {
			case <-s.stopChannel:
				return
			case <-ctx.Done():
				return
			default:
				if err := s.taskFunc(ctx); err != nil {
					// Log error but continue running
					fmt.Printf("Background service %s error: %v\n", s.Name(), err)
				}
			}
		}
	}()

	return nil
}

// Stop stops the background service
func (s *BackgroundService) Stop(ctx context.Context) error {
	close(s.stopChannel)

	// Wait for the service to stop or context to timeout
	select {
	case <-s.done:
		// Service stopped gracefully
	case <-ctx.Done():
		// Context timeout
		return fmt.Errorf("background service %s did not stop within timeout", s.Name())
	}

	return s.BaseServiceComponent.Stop(ctx)
}