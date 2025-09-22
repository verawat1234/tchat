package main

import (
	"sync"
	"time"
)

// RegisterService adds a service to the registry
func (sr *ServiceRegistry) RegisterService(service *ServiceInstance) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	service.LastSeen = time.Now()
	sr.services[service.ID] = service
}

// DeregisterService removes a service from the registry
func (sr *ServiceRegistry) DeregisterService(serviceID string) bool {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	if _, exists := sr.services[serviceID]; exists {
		delete(sr.services, serviceID)
		return true
	}
	return false
}

// GetService retrieves a service by ID
func (sr *ServiceRegistry) GetService(serviceID string) *ServiceInstance {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	if service, exists := sr.services[serviceID]; exists {
		return service
	}
	return nil
}

// GetServiceByName retrieves the first healthy service by name
func (sr *ServiceRegistry) GetServiceByName(serviceName string) *ServiceInstance {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	for _, service := range sr.services {
		if service.Name == serviceName && service.Health == string(Healthy) {
			return service
		}
	}
	return nil
}

// GetHealthyService returns a healthy service instance using load balancing
func (sr *ServiceRegistry) GetHealthyService(serviceName string) *ServiceInstance {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	var healthyServices []*ServiceInstance
	for _, service := range sr.services {
		if service.Name == serviceName && service.Health == string(Healthy) {
			healthyServices = append(healthyServices, service)
		}
	}

	if len(healthyServices) == 0 {
		return nil
	}

	// Simple round-robin load balancing
	// In production, this could be more sophisticated
	return healthyServices[time.Now().Unix()%int64(len(healthyServices))]
}

// GetAllServices returns all registered services
func (sr *ServiceRegistry) GetAllServices() []*ServiceInstance {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	services := make([]*ServiceInstance, 0, len(sr.services))
	for _, service := range sr.services {
		services = append(services, service)
	}
	return services
}

// GetServicesByName returns all services with the given name
func (sr *ServiceRegistry) GetServicesByName(serviceName string) []*ServiceInstance {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	var services []*ServiceInstance
	for _, service := range sr.services {
		if service.Name == serviceName {
			services = append(services, service)
		}
	}
	return services
}

// GetHealthyServiceCount returns the number of healthy services
func (sr *ServiceRegistry) GetHealthyServiceCount() int {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	count := 0
	for _, service := range sr.services {
		if service.Health == string(Healthy) {
			count++
		}
	}
	return count
}

// GetServiceCount returns the total number of registered services
func (sr *ServiceRegistry) GetServiceCount() int {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	return len(sr.services)
}

// UpdateServiceHealth updates the health status of a service
func (sr *ServiceRegistry) UpdateServiceHealth(serviceID string, health HealthStatus) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	if service, exists := sr.services[serviceID]; exists {
		service.Health = string(health)
		service.LastSeen = time.Now()
	}
}

// CleanupStaleServices removes services that haven't been seen for a while
func (sr *ServiceRegistry) CleanupStaleServices(timeout time.Duration) int {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	cutoff := time.Now().Add(-timeout)
	removed := 0

	for id, service := range sr.services {
		if service.LastSeen.Before(cutoff) {
			delete(sr.services, id)
			removed++
		}
	}

	return removed
}

// LoadBalancer implements load balancing strategies
type LoadBalancer struct {
	serviceName string
	strategy    LoadBalancingStrategy
	counter     uint64
	mu          sync.Mutex
}

// LoadBalancingStrategy defines load balancing strategy
type LoadBalancingStrategy string

const (
	RoundRobin        LoadBalancingStrategy = "round_robin"
	LeastConnections  LoadBalancingStrategy = "least_connections"
	WeightedRoundRobin LoadBalancingStrategy = "weighted_round_robin"
	Random            LoadBalancingStrategy = "random"
)

// NewLoadBalancer creates a new load balancer for a service
func NewLoadBalancer(serviceName string, strategy LoadBalancingStrategy) *LoadBalancer {
	return &LoadBalancer{
		serviceName: serviceName,
		strategy:    strategy,
		counter:     0,
	}
}

// SelectService selects a service instance based on the load balancing strategy
func (lb *LoadBalancer) SelectService(registry *ServiceRegistry) *ServiceInstance {
	services := registry.GetServicesByName(lb.serviceName)
	if len(services) == 0 {
		return nil
	}

	// Filter healthy services
	var healthyServices []*ServiceInstance
	for _, service := range services {
		if service.Health == string(Healthy) {
			healthyServices = append(healthyServices, service)
		}
	}

	if len(healthyServices) == 0 {
		return nil
	}

	switch lb.strategy {
	case RoundRobin:
		return lb.roundRobinSelect(healthyServices)
	case Random:
		return lb.randomSelect(healthyServices)
	default:
		return lb.roundRobinSelect(healthyServices)
	}
}

func (lb *LoadBalancer) roundRobinSelect(services []*ServiceInstance) *ServiceInstance {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	selected := services[lb.counter%uint64(len(services))]
	lb.counter++
	return selected
}

func (lb *LoadBalancer) randomSelect(services []*ServiceInstance) *ServiceInstance {
	index := time.Now().UnixNano() % int64(len(services))
	return services[index]
}