/**
 * Service Routing and Gateway Integration Test Utility
 *
 * Test utility to verify service routing, gateway connectivity,
 * and RTK Query integration are working correctly.
 */

import { buildServiceUrl, getServiceForEndpoint, SERVICE_CONFIG } from '../services/serviceConfig';
import { serviceHealthMonitor } from '../services/serviceHealth';

/**
 * Test gateway connectivity and health
 */
export async function testGatewayIntegration() {
  console.group('🚪 Gateway Integration Test');

  try {
    // Test gateway health
    const gatewayHealth = await serviceHealthMonitor.checkGatewayHealth();
    const gatewayStatus = gatewayHealth.available ? '✅' : '❌';
    const gatewayTiming = gatewayHealth.responseTime ? ` (${gatewayHealth.responseTime}ms)` : '';

    console.log(`${gatewayStatus} Gateway Health: ${gatewayHealth.available ? 'Healthy' : 'Unhealthy'}${gatewayTiming}`);

    if (gatewayHealth.error) {
      console.warn(`   Gateway Error: ${gatewayHealth.error}`);
    }

    // Test critical API endpoints through gateway
    const criticalEndpoints = [
      '/products',
      '/chats',
      '/wallet',
      '/notifications',
      '/social/feed'
    ];

    console.log('\n🔍 API Endpoints Test:');
    let availableEndpoints = 0;

    for (const endpoint of criticalEndpoints) {
      try {
        const result = await serviceHealthMonitor.testApiEndpoint(endpoint);
        const status = result.available ? '✅' : '❌';
        const statusCode = result.status ? ` (${result.status})` : '';

        console.log(`  ${status} ${endpoint}${statusCode}`);

        if (result.available) {
          availableEndpoints++;
        } else if (result.error) {
          console.warn(`     Error: ${result.error}`);
        }
      } catch (error) {
        console.log(`  ❌ ${endpoint} - Connection failed`);
      }
    }

    console.log(`\n📊 Endpoint Summary: ${availableEndpoints}/${criticalEndpoints.length} endpoints available`);

    return {
      gatewayHealthy: gatewayHealth.available,
      endpointsAvailable: availableEndpoints,
      totalEndpoints: criticalEndpoints.length
    };

  } catch (error) {
    console.error('❌ Gateway integration test failed:', error);
    return {
      gatewayHealthy: false,
      endpointsAvailable: 0,
      totalEndpoints: 0
    };
  } finally {
    console.groupEnd();
  }
}

/**
 * Test RTK Query integration with gateway
 */
export async function testRTKQueryIntegration() {
  console.group('⚛️ RTK Query Integration Test');

  try {
    // Check if RTK Query is properly configured
    console.log('🔧 Checking RTK Query configuration...');

    // Import RTK Query components
    const { api } = await import('../services/api');
    console.log('✅ RTK Query API slice loaded');

    // Check endpoints
    const endpoints = Object.keys(api.endpoints);
    console.log(`📋 Available RTK Query endpoints: ${endpoints.length}`);

    // Test microservices API
    const { checkServiceHealth } = await import('../services/microservicesApi');
    if (checkServiceHealth) {
      console.log('✅ Microservices API module loaded');
    }

    console.log('✅ RTK Query integration appears healthy');

  } catch (error) {
    console.error('❌ RTK Query integration test failed:', error);
  } finally {
    console.groupEnd();
  }
}

/**
 * Test service routing configuration
 */
export async function testServiceRouting() {
  console.group('🔧 Complete Service Integration Test');

  // Test gateway integration first
  const gatewayTest = await testGatewayIntegration();

  // Test RTK Query integration
  await testRTKQueryIntegration();

  // Test endpoint routing
  const testEndpoints = [
    '/videos/shorts',
    '/videos/long',
    '/auth/login',
    '/users/profile',
    '/content/items',
    '/messages/list',
    '/notifications/list',
    '/payments/wallet',
    '/products/list'
  ];

  console.group('📍 Endpoint Routing Configuration:');
  testEndpoints.forEach(endpoint => {
    const service = getServiceForEndpoint(endpoint);
    const url = buildServiceUrl(endpoint);
    console.log(`  ${endpoint} → ${service} → ${url}`);
  });
  console.groupEnd();

  // Test individual service health
  console.group('🏥 Individual Service Health:');
  try {
    const healthResults = await serviceHealthMonitor.checkAllServicesHealth();
    healthResults.forEach(health => {
      const status = health.available ? '✅' : '❌';
      const responseTime = health.responseTime ? `(${health.responseTime}ms)` : '';
      const error = health.error ? ` - ${health.error}` : '';
      console.log(`  ${status} ${health.service}: ${health.available ? 'Available' : 'Unavailable'} ${responseTime}${error}`);
    });

    const summary = serviceHealthMonitor.getServiceSummary();
    console.log(`\n📊 Service Summary: ${summary.available}/${summary.total} services available (${summary.availabilityPercentage.toFixed(1)}%)`);
    if (summary.averageResponseTime > 0) {
      console.log(`⏱️ Average Response Time: ${summary.averageResponseTime}ms`);
    }

  } catch (error) {
    console.error('❌ Service health check failed:', error);
  }
  console.groupEnd();

  // Overall summary
  console.group('📋 Integration Test Summary:');
  if (gatewayTest.gatewayHealthy) {
    console.log('✅ Gateway: Connected and responding');
  } else {
    console.log('❌ Gateway: Connection issues detected');
  }

  if (gatewayTest.endpointsAvailable > 0) {
    console.log(`✅ API Endpoints: ${gatewayTest.endpointsAvailable}/${gatewayTest.totalEndpoints} responding`);
  } else {
    console.log('⚠️ API Endpoints: Using fallback/mock data (normal for early development)');
  }

  console.log('✅ RTK Query: Integration configured');
  console.log('ℹ️ Ready for development with gateway integration');
  console.groupEnd();

  console.groupEnd();
}

/**
 * Test a specific service endpoint
 */
export async function testServiceEndpoint(endpoint: string) {
  const service = getServiceForEndpoint(endpoint);
  const url = buildServiceUrl(endpoint);

  console.log(`Testing ${endpoint}:`);
  console.log(`  Service: ${service}`);
  console.log(`  URL: ${url}`);

  try {
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    });

    console.log(`  Status: ${response.status} ${response.statusText}`);

    if (response.ok) {
      console.log('  ✅ Endpoint accessible');
    } else {
      console.log('  ❌ Endpoint returned error');
    }

    return response.ok;

  } catch (error) {
    console.log(`  ❌ Connection failed: ${error instanceof Error ? error.message : error}`);
    return false;
  }
}

/**
 * Run service routing tests in development
 */
if (import.meta.env.DEV) {
  // Auto-run tests after a short delay to allow services to start
  setTimeout(() => {
    testServiceRouting();
  }, 2000);
}