/**
 * Service Health Monitoring for Microservices and Gateway
 *
 * Monitors the health and availability of backend microservices
 * and API gateway, with fallback strategies for development.
 */

import { useState, useEffect } from 'react';
import { SERVICE_CONFIG, HEALTH_CHECK_CONFIG, type ServiceHealth } from './serviceConfig';

interface GatewayHealth {
  service: 'gateway';
  available: boolean;
  lastChecked: number;
  responseTime?: number;
  error?: string;
}

class ServiceHealthMonitor {
  private healthStatus: Map<string, ServiceHealth> = new Map();
  private gatewayHealth: GatewayHealth | null = null;
  private healthCheckInterval: NodeJS.Timeout | null = null;

  constructor() {
    this.initializeHealthStatus();
  }

  /**
   * Initialize health status for all services
   */
  private initializeHealthStatus() {
    Object.keys(SERVICE_CONFIG.services).forEach(service => {
      this.healthStatus.set(service, {
        service,
        available: true, // Assume available initially
        lastChecked: Date.now()
      });
    });
  }

  /**
   * Check gateway health
   */
  async checkGatewayHealth(): Promise<GatewayHealth> {
    const startTime = Date.now();
    const gatewayUrl = `${SERVICE_CONFIG.gateway.baseUrl}/health`;

    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), HEALTH_CHECK_CONFIG.timeout);

      const response = await fetch(gatewayUrl, {
        method: 'GET',
        signal: controller.signal,
        headers: {
          'Content-Type': 'application/json'
        }
      });

      clearTimeout(timeoutId);
      const responseTime = Date.now() - startTime;

      const health: GatewayHealth = {
        service: 'gateway',
        available: response.ok,
        lastChecked: Date.now(),
        responseTime,
        error: response.ok ? undefined : `HTTP ${response.status}: ${response.statusText}`
      };

      this.gatewayHealth = health;
      return health;

    } catch (error) {
      const responseTime = Date.now() - startTime;
      const health: GatewayHealth = {
        service: 'gateway',
        available: false,
        lastChecked: Date.now(),
        responseTime,
        error: error instanceof Error ? error.message : 'Unknown error'
      };

      this.gatewayHealth = health;
      return health;
    }
  }

  /**
   * Test API endpoint availability through gateway
   */
  async testApiEndpoint(endpoint: string): Promise<{
    available: boolean;
    status?: number;
    error?: string;
    data?: any;
  }> {
    const fullUrl = `${SERVICE_CONFIG.gateway.baseUrl}${SERVICE_CONFIG.gateway.apiPrefix}${endpoint}`;

    try {
      const response = await fetch(fullUrl, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        signal: AbortSignal.timeout(5000),
      });

      if (response.ok) {
        try {
          const data = await response.json();
          return {
            available: true,
            status: response.status,
            data,
          };
        } catch {
          return {
            available: true,
            status: response.status,
          };
        }
      } else {
        return {
          available: false,
          status: response.status,
          error: `HTTP ${response.status}: ${response.statusText}`,
        };
      }
    } catch (error) {
      return {
        available: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  /**
   * Check health of a specific service
   */
  async checkServiceHealth(serviceName: string): Promise<ServiceHealth> {
    const serviceConfig = SERVICE_CONFIG.services[serviceName as keyof typeof SERVICE_CONFIG.services];
    if (!serviceConfig) {
      throw new Error(`Unknown service: ${serviceName}`);
    }

    const startTime = Date.now();

    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), HEALTH_CHECK_CONFIG.timeout);

      const response = await fetch(`${serviceConfig.baseUrl}${serviceConfig.healthEndpoint}`, {
        method: 'GET',
        signal: controller.signal,
        headers: {
          'Content-Type': 'application/json'
        }
      });

      clearTimeout(timeoutId);
      const responseTime = Date.now() - startTime;

      const health: ServiceHealth = {
        service: serviceName,
        available: response.ok,
        lastChecked: Date.now(),
        responseTime,
        error: response.ok ? undefined : `HTTP ${response.status}: ${response.statusText}`
      };

      this.healthStatus.set(serviceName, health);
      return health;

    } catch (error) {
      const responseTime = Date.now() - startTime;
      const health: ServiceHealth = {
        service: serviceName,
        available: false,
        lastChecked: Date.now(),
        responseTime,
        error: error instanceof Error ? error.message : 'Unknown error'
      };

      this.healthStatus.set(serviceName, health);
      return health;
    }
  }

  /**
   * Check health of all services including gateway
   */
  async checkAllServicesHealth(): Promise<ServiceHealth[]> {
    // Check gateway health first
    await this.checkGatewayHealth();

    const serviceNames = Object.keys(SERVICE_CONFIG.services);
    const healthChecks = serviceNames.map(service => this.checkServiceHealth(service));

    return Promise.all(healthChecks);
  }

  /**
   * Get gateway health status
   */
  getGatewayHealth(): GatewayHealth | null {
    return this.gatewayHealth;
  }

  /**
   * Check if gateway is available
   */
  isGatewayAvailable(): boolean {
    if (!this.gatewayHealth) return false;

    const staleThreshold = HEALTH_CHECK_CONFIG.interval * 2;
    const isStale = Date.now() - this.gatewayHealth.lastChecked > staleThreshold;

    return this.gatewayHealth.available && !isStale;
  }

  /**
   * Get current health status of a service
   */
  getServiceHealth(serviceName: string): ServiceHealth | undefined {
    return this.healthStatus.get(serviceName);
  }

  /**
   * Get health status of all services
   */
  getAllServicesHealth(): ServiceHealth[] {
    return Array.from(this.healthStatus.values());
  }

  /**
   * Check if a service is available
   */
  isServiceAvailable(serviceName: string): boolean {
    const health = this.healthStatus.get(serviceName);
    if (!health) return false;

    // Consider service unavailable if last check was too long ago
    const staleThreshold = HEALTH_CHECK_CONFIG.interval * 2;
    const isStale = Date.now() - health.lastChecked > staleThreshold;

    return health.available && !isStale;
  }

  /**
   * Start periodic health monitoring
   */
  startHealthMonitoring() {
    if (this.healthCheckInterval) {
      clearInterval(this.healthCheckInterval);
    }

    // Do initial health check
    this.checkAllServicesHealth().catch(error => {
      console.warn('Initial health check failed:', error);
    });

    // Set up periodic monitoring
    this.healthCheckInterval = setInterval(async () => {
      try {
        await this.checkAllServicesHealth();
      } catch (error) {
        console.warn('Periodic health check failed:', error);
      }
    }, HEALTH_CHECK_CONFIG.interval);
  }

  /**
   * Stop health monitoring
   */
  stopHealthMonitoring() {
    if (this.healthCheckInterval) {
      clearInterval(this.healthCheckInterval);
      this.healthCheckInterval = null;
    }
  }

  /**
   * Get service availability summary
   */
  getServiceSummary() {
    const services = this.getAllServicesHealth();
    const total = services.length;
    const available = services.filter(s => s.available).length;
    const avgResponseTime = services
      .filter(s => s.responseTime !== undefined)
      .reduce((sum, s) => sum + (s.responseTime || 0), 0) / services.length;

    return {
      total,
      available,
      unavailable: total - available,
      availabilityPercentage: (available / total) * 100,
      averageResponseTime: Math.round(avgResponseTime) || 0
    };
  }
}

// Create singleton instance
export const serviceHealthMonitor = new ServiceHealthMonitor();

/**
 * Hook for React components to monitor service health
 */
export function useServiceHealth() {
  const [healthData, setHealthData] = useState(serviceHealthMonitor.getAllServicesHealth());

  useEffect(() => {
    // Start monitoring when component mounts
    serviceHealthMonitor.startHealthMonitoring();

    // Set up periodic updates for React state
    const updateInterval = setInterval(() => {
      setHealthData(serviceHealthMonitor.getAllServicesHealth());
    }, 5000); // Update every 5 seconds

    return () => {
      clearInterval(updateInterval);
      // Note: Don't stop health monitoring here as other components might need it
    };
  }, []);

  return {
    services: healthData,
    summary: serviceHealthMonitor.getServiceSummary(),
    isServiceAvailable: (serviceName: string) => serviceHealthMonitor.isServiceAvailable(serviceName),
    refreshHealth: () => serviceHealthMonitor.checkAllServicesHealth()
  };
}

/**
 * Service availability checker for API calls
 */
export async function ensureServiceAvailable(serviceName: string): Promise<boolean> {
  const health = await serviceHealthMonitor.checkServiceHealth(serviceName);
  return health.available;
}

// Auto-start health monitoring in development
if (import.meta.env.DEV) {
  serviceHealthMonitor.startHealthMonitoring();
}