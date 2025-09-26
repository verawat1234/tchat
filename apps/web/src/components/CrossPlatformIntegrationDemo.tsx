/**
 * Cross-Platform Integration Services Demo Component
 * Demonstrates T051-T056 integration with real-time validation and monitoring
 */

import React, { useState, useEffect } from 'react';
import { TchatCard, TchatCardHeader, TchatCardContent } from './TchatCard';

// Import all the integration services
import {
  useGetAllComponentsQuery,
  useValidateComponentConsistencyMutation,
  useSyncComponentMutation,
  useGetSyncConflictsQuery,
  ComponentSyncUtils
} from '../services/componentSync';
import {
  consistencyValidator,
  ConsistencyValidatorUtils
} from '../utils/consistencyValidator';
import {
  useGetAllImplementationStatusesQuery,
  useGetOverallHealthMetricsQuery,
  useGetRealtimeTrackingStatusQuery,
  ImplementationTrackerUtils
} from '../services/implementationTracker';
import {
  useValidateComponentComprehensiveMutation,
  useValidateComponentPerformanceMutation,
  useValidateComponentAccessibilityMutation,
  useValidateVisualConsistencyMutation,
  ValidationUtils
} from '../services/performanceValidator';

// Demo component interface
interface CrossPlatformIntegrationDemoProps {
  selectedComponentId?: string;
  enableRealTimeMonitoring?: boolean;
}

export const CrossPlatformIntegrationDemo: React.FC<CrossPlatformIntegrationDemoProps> = ({
  selectedComponentId = 'tchat-button',
  enableRealTimeMonitoring = true
}) => {
  const [activeTab, setActiveTab] = useState<'sync' | 'validation' | 'tracking' | 'performance'>('sync');
  const [validationResults, setValidationResults] = useState<any>(null);
  const [realTimeUpdates, setRealTimeUpdates] = useState<any[]>([]);

  // RTK Query hooks for real-time data
  const { data: components } = useGetAllComponentsQuery();
  const { data: implementationStatuses } = useGetAllImplementationStatusesQuery({});
  const { data: overallHealth } = useGetOverallHealthMetricsQuery({});
  const { data: realTimeStatus } = useGetRealtimeTrackingStatusQuery(undefined, {
    skip: !enableRealTimeMonitoring,
    pollingInterval: 10000
  });
  const { data: syncConflicts } = useGetSyncConflictsQuery({});

  // Mutation hooks
  const [validateConsistency] = useValidateComponentConsistencyMutation();
  const [syncComponent] = useSyncComponentMutation();
  const [validateComprehensive] = useValidateComponentComprehensiveMutation();
  const [validatePerformance] = useValidateComponentPerformanceMutation();
  const [validateAccessibility] = useValidateComponentAccessibilityMutation();
  const [validateVisualConsistency] = useValidateVisualConsistencyMutation();

  // Get current component data
  const currentComponent = components?.find(c => c.id === selectedComponentId);
  const currentImplementations = implementationStatuses?.filter(
    impl => impl.componentId === selectedComponentId
  );

  // Real-time updates effect
  useEffect(() => {
    if (realTimeStatus?.isActive) {
      const newUpdate = {
        timestamp: new Date().toISOString(),
        component: selectedComponentId,
        status: 'monitoring_active',
        details: `Monitoring ${realTimeStatus.watchedComponents?.length || 0} components`
      };
      setRealTimeUpdates(prev => [newUpdate, ...prev.slice(0, 9)]); // Keep last 10 updates
    }
  }, [realTimeStatus, selectedComponentId]);

  // Handle comprehensive validation
  const handleComprehensiveValidation = async () => {
    try {
      const result = await validateComprehensive({
        componentId: selectedComponentId
      }).unwrap();

      setValidationResults(result);

      // Add real-time update
      setRealTimeUpdates(prev => [{
        timestamp: new Date().toISOString(),
        component: selectedComponentId,
        status: result.constitutionalCompliance ? 'validation_passed' : 'validation_failed',
        details: `Overall score: ${(result.overallScore * 100).toFixed(1)}%`
      }, ...prev.slice(0, 9)]);
    } catch (error) {
      console.error('Validation failed:', error);
    }
  };

  // Handle component sync
  const handleComponentSync = async () => {
    try {
      const result = await syncComponent({
        componentId: selectedComponentId,
        targetPlatforms: ['web', 'ios', 'android'],
        syncType: 'full',
        forceSync: false
      }).unwrap();

      setRealTimeUpdates(prev => [{
        timestamp: new Date().toISOString(),
        component: selectedComponentId,
        status: result.success ? 'sync_completed' : 'sync_failed',
        details: `Synced to ${result.platformsUpdated.length} platforms in ${result.syncDuration}ms`
      }, ...prev.slice(0, 9)]);
    } catch (error) {
      console.error('Sync failed:', error);
    }
  };

  // Handle consistency validation
  const handleConsistencyValidation = async () => {
    try {
      const result = await validateConsistency({
        componentId: selectedComponentId,
        includePerformance: true
      }).unwrap();

      setRealTimeUpdates(prev => [{
        timestamp: new Date().toISOString(),
        component: selectedComponentId,
        status: result.meetsConstitutionalRequirement ? 'consistency_compliant' : 'consistency_violation',
        details: `Consistency: ${ConsistencyValidatorUtils.formatConsistencyScore(result.overallScore)}`
      }, ...prev.slice(0, 9)]);
    } catch (error) {
      console.error('Consistency validation failed:', error);
    }
  };

  const renderSyncTab = () => (
    <div className="space-y-4">
      {/* Component Sync Status */}
      <TchatCard variant="outlined">
        <TchatCardHeader
          title="Component Sync Status"
          subtitle={`${selectedComponentId} - Cross-Platform Synchronization`}
        />
        <TchatCardContent>
          <div className="grid grid-cols-3 gap-4">
            {currentComponent?.platforms.map(platform => (
              <div key={platform.platform} className="text-center p-3 bg-gray-50 rounded">
                <div className="font-medium">{platform.platform.toUpperCase()}</div>
                <div className="text-sm text-gray-600">{platform.version}</div>
                <div className="text-xs mt-1">
                  {ConsistencyValidatorUtils.formatConsistencyScore(platform.visualConsistencyScore)}
                </div>
                <div className={`text-xs px-2 py-1 rounded mt-2 ${
                  platform.implementationStatus === 'implemented' ? 'bg-green-100 text-green-800' :
                  platform.implementationStatus === 'in_progress' ? 'bg-yellow-100 text-yellow-800' :
                  'bg-red-100 text-red-800'
                }`}>
                  {platform.implementationStatus.replace('_', ' ')}
                </div>
              </div>
            )) || (
              <div className="col-span-3 text-center text-gray-500">
                No component data available
              </div>
            )}
          </div>

          <div className="mt-4 flex gap-2">
            <button
              onClick={handleComponentSync}
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
            >
              Sync Component
            </button>
            <button
              onClick={handleConsistencyValidation}
              className="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600 transition-colors"
            >
              Validate Consistency
            </button>
          </div>
        </TchatCardContent>
      </TchatCard>

      {/* Sync Conflicts */}
      {syncConflicts && syncConflicts.length > 0 && (
        <TchatCard variant="filled">
          <TchatCardHeader
            title="Sync Conflicts"
            subtitle={`${syncConflicts.length} conflicts requiring attention`}
          />
          <TchatCardContent>
            <div className="space-y-2">
              {syncConflicts.slice(0, 3).map(conflict => (
                <div key={conflict.componentId} className="p-3 border rounded">
                  <div className="flex justify-between items-start">
                    <div>
                      <div className="font-medium">{conflict.conflictType.replace('_', ' ')}</div>
                      <div className="text-sm text-gray-600">{conflict.description}</div>
                    </div>
                    <span className={`px-2 py-1 text-xs rounded ${
                      conflict.severity === 'critical' ? 'bg-red-100 text-red-800' :
                      conflict.severity === 'high' ? 'bg-orange-100 text-orange-800' :
                      'bg-yellow-100 text-yellow-800'
                    }`}>
                      {conflict.severity}
                    </span>
                  </div>
                  <div className="text-xs text-gray-500 mt-1">
                    Platforms: {conflict.platforms.join(', ')}
                  </div>
                </div>
              ))}
            </div>
          </TchatCardContent>
        </TchatCard>
      )}
    </div>
  );

  const renderValidationTab = () => (
    <div className="space-y-4">
      {/* Constitutional Compliance Dashboard */}
      <TchatCard variant="elevated">
        <TchatCardHeader
          title="Constitutional Compliance"
          subtitle="97% Cross-Platform Consistency Requirement"
        />
        <TchatCardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="text-center p-3 bg-blue-50 rounded">
              <div className="text-2xl font-bold text-blue-600">
                {currentComponent ?
                  ConsistencyValidatorUtils.formatConsistencyScore(currentComponent.crossPlatformConsistencyScore) :
                  '---'
                }
              </div>
              <div className="text-sm text-gray-600">Cross-Platform</div>
            </div>
            <div className="text-center p-3 bg-green-50 rounded">
              <div className="text-2xl font-bold text-green-600">
                {overallHealth ?
                  ValidationUtils.formatPerformanceScore(overallHealth.overallProgress) :
                  '---'
                }
              </div>
              <div className="text-sm text-gray-600">Performance</div>
            </div>
            <div className="text-center p-3 bg-purple-50 rounded">
              <div className="text-2xl font-bold text-purple-600">AA</div>
              <div className="text-sm text-gray-600">Accessibility</div>
            </div>
            <div className="text-center p-3 bg-orange-50 rounded">
              <div className="text-2xl font-bold text-orange-600">
                {currentComponent?.constitutionalCompliance ? '✅' : '❌'}
              </div>
              <div className="text-sm text-gray-600">Compliant</div>
            </div>
          </div>

          <button
            onClick={handleComprehensiveValidation}
            className="mt-4 w-full px-4 py-2 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded hover:from-blue-600 hover:to-purple-700 transition-all"
          >
            Run Comprehensive Validation
          </button>

          {validationResults && (
            <div className="mt-4 p-4 bg-gray-50 rounded">
              <h4 className="font-medium mb-2">Validation Results</h4>
              <div className="text-sm space-y-1">
                <div>Overall Score: {ValidationUtils.formatPerformanceScore(validationResults.overallScore)}</div>
                <div>Constitutional Compliance: {validationResults.constitutionalCompliance ? '✅ Yes' : '❌ No'}</div>
                <div>Total Issues: {validationResults.summary.totalIssues}</div>
                <div>Critical Issues: {validationResults.summary.criticalIssues}</div>
                {validationResults.summary.constitutionalViolations > 0 && (
                  <div className="text-red-600 font-medium">
                    Constitutional Violations: {validationResults.summary.constitutionalViolations}
                  </div>
                )}
              </div>
            </div>
          )}
        </TchatCardContent>
      </TchatCard>
    </div>
  );

  const renderTrackingTab = () => (
    <div className="space-y-4">
      {/* Implementation Status Overview */}
      <TchatCard variant="outlined">
        <TchatCardHeader
          title="Implementation Status Tracking"
          subtitle="Real-time cross-platform development monitoring"
        />
        <TchatCardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {currentImplementations?.map(impl => (
              <div key={`${impl.componentId}-${impl.platform}`} className="p-3 border rounded">
                <div className="flex justify-between items-center mb-2">
                  <span className="font-medium">{impl.platform.toUpperCase()}</span>
                  <span className={ImplementationTrackerUtils.getStatusColor(impl.status)}>
                    {ImplementationTrackerUtils.formatImplementationStatus(impl.status)}
                  </span>
                </div>
                <div className="text-sm space-y-1">
                  <div>Progress: {ImplementationTrackerUtils.calculateProgressPercentage(impl.status)}%</div>
                  <div>Test Coverage: {(impl.testCoverage * 100).toFixed(1)}%</div>
                  <div>Performance: {ValidationUtils.formatPerformanceScore(impl.performanceScore)}</div>
                  <div>Build: {impl.buildStatus.status}</div>
                </div>
              </div>
            )) || (
              <div className="col-span-3 text-center text-gray-500">
                No implementation data available
              </div>
            )}
          </div>

          {/* Real-time monitoring status */}
          <div className="mt-4 p-3 bg-blue-50 rounded">
            <div className="flex items-center justify-between">
              <span className="font-medium">Real-time Monitoring</span>
              <span className={`px-2 py-1 text-xs rounded ${
                realTimeStatus?.isActive ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
              }`}>
                {realTimeStatus?.isActive ? 'Active' : 'Inactive'}
              </span>
            </div>
            {realTimeStatus?.isActive && (
              <div className="text-sm text-gray-600 mt-1">
                Watching {realTimeStatus.watchedComponents?.length || 0} components,
                {realTimeStatus.pendingSyncs?.length || 0} pending syncs
              </div>
            )}
          </div>
        </TchatCardContent>
      </TchatCard>

      {/* Overall Health Metrics */}
      {overallHealth && (
        <TchatCard variant="filled">
          <TchatCardHeader
            title="System Health Overview"
            subtitle="Cross-platform ecosystem status"
          />
          <TchatCardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="text-center">
                <div className="text-lg font-bold">
                  {ValidationUtils.formatPerformanceScore(overallHealth.overallProgress)}
                </div>
                <div className="text-sm text-gray-600">Overall Progress</div>
              </div>
              <div className="text-center">
                <div className="text-lg font-bold">
                  {ConsistencyValidatorUtils.formatConsistencyScore(overallHealth.crossPlatformConsistency)}
                </div>
                <div className="text-sm text-gray-600">Consistency</div>
              </div>
              <div className="text-center">
                <div className="text-lg font-bold">{overallHealth.totalIssues}</div>
                <div className="text-sm text-gray-600">Total Issues</div>
              </div>
              <div className="text-center">
                <div className="text-lg font-bold text-red-600">{overallHealth.criticalIssues}</div>
                <div className="text-sm text-gray-600">Critical Issues</div>
              </div>
            </div>

            {overallHealth.blockedComponents.length > 0 && (
              <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded">
                <div className="font-medium text-red-800">Blocked Components</div>
                <div className="text-sm text-red-600">
                  {overallHealth.blockedComponents.join(', ')}
                </div>
              </div>
            )}
          </TchatCardContent>
        </TchatCard>
      )}
    </div>
  );

  const renderPerformanceTab = () => (
    <div className="space-y-4">
      <TchatCard variant="elevated">
        <TchatCardHeader
          title="Performance Validation"
          subtitle="Constitutional requirement: <200ms load time, 60fps animations"
        />
        <TchatCardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <h4 className="font-medium">Performance Thresholds</h4>
              <div className="text-sm space-y-1">
                <div>Load Time: <span className="font-mono">≤200ms</span></div>
                <div>Animation: <span className="font-mono">60fps</span></div>
                <div>Bundle Size: <span className="font-mono">≤500KB</span></div>
                <div>Memory Usage: <span className="font-mono">≤100MB mobile</span></div>
              </div>
            </div>
            <div className="space-y-2">
              <h4 className="font-medium">Accessibility Standards</h4>
              <div className="text-sm space-y-1">
                <div>WCAG Level: <span className="font-mono">AA (Constitutional)</span></div>
                <div>Contrast Ratio: <span className="font-mono">≥4.5:1</span></div>
                <div>Keyboard Navigation: <span className="font-mono">Full Support</span></div>
                <div>Screen Reader: <span className="font-mono">Compatible</span></div>
              </div>
            </div>
          </div>

          <div className="mt-4 flex gap-2 flex-wrap">
            <button
              onClick={async () => {
                try {
                  const result = await validatePerformance({
                    componentId: selectedComponentId
                  }).unwrap();
                  console.log('Performance validation:', result);
                } catch (error) {
                  console.error('Performance validation failed:', error);
                }
              }}
              className="px-3 py-2 bg-blue-500 text-white rounded text-sm hover:bg-blue-600 transition-colors"
            >
              Test Performance
            </button>
            <button
              onClick={async () => {
                try {
                  const result = await validateAccessibility({
                    componentId: selectedComponentId
                  }).unwrap();
                  console.log('Accessibility validation:', result);
                } catch (error) {
                  console.error('Accessibility validation failed:', error);
                }
              }}
              className="px-3 py-2 bg-green-500 text-white rounded text-sm hover:bg-green-600 transition-colors"
            >
              Test Accessibility
            </button>
            <button
              onClick={async () => {
                try {
                  const result = await validateVisualConsistency({
                    componentId: selectedComponentId
                  }).unwrap();
                  console.log('Visual consistency validation:', result);
                } catch (error) {
                  console.error('Visual consistency validation failed:', error);
                }
              }}
              className="px-3 py-2 bg-purple-500 text-white rounded text-sm hover:bg-purple-600 transition-colors"
            >
              Test Visual Consistency
            </button>
          </div>
        </TchatCardContent>
      </TchatCard>
    </div>
  );

  return (
    <div className="max-w-6xl mx-auto p-6 space-y-6">
      {/* Header */}
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Cross-Platform Integration Services
        </h1>
        <p className="text-gray-600">
          T051-T056 Implementation: Component Sync, Consistency Validation & Performance Monitoring
        </p>
        <div className="mt-2 text-sm text-blue-600">
          Constitutional Requirement: 97% Cross-Platform Consistency
        </div>
      </div>

      {/* Real-time Updates Ticker */}
      {realTimeUpdates.length > 0 && (
        <TchatCard variant="glass">
          <TchatCardHeader title="Real-time Updates" />
          <TchatCardContent>
            <div className="space-y-2 max-h-32 overflow-y-auto">
              {realTimeUpdates.slice(0, 5).map((update, index) => (
                <div key={index} className="flex items-center justify-between text-sm p-2 bg-white/50 rounded">
                  <div className="flex items-center gap-2">
                    <span className="text-xs text-gray-500">
                      {new Date(update.timestamp).toLocaleTimeString()}
                    </span>
                    <span className="font-medium">{update.component}</span>
                    <span className={`px-2 py-1 text-xs rounded ${
                      update.status.includes('passed') || update.status.includes('completed') || update.status.includes('compliant') ?
                        'bg-green-100 text-green-800' :
                      update.status.includes('failed') || update.status.includes('violation') ?
                        'bg-red-100 text-red-800' :
                        'bg-blue-100 text-blue-800'
                    }`}>
                      {update.status.replace('_', ' ')}
                    </span>
                  </div>
                  <span className="text-xs text-gray-600">{update.details}</span>
                </div>
              ))}
            </div>
          </TchatCardContent>
        </TchatCard>
      )}

      {/* Tab Navigation */}
      <div className="flex space-x-1 bg-gray-100 p-1 rounded-lg">
        {[
          { id: 'sync', label: 'Component Sync' },
          { id: 'validation', label: 'Validation' },
          { id: 'tracking', label: 'Implementation Tracking' },
          { id: 'performance', label: 'Performance & A11y' }
        ].map(tab => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id as any)}
            className={`flex-1 px-4 py-2 rounded-md text-sm font-medium transition-all ${
              activeTab === tab.id
                ? 'bg-white text-gray-900 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      <div className="min-h-96">
        {activeTab === 'sync' && renderSyncTab()}
        {activeTab === 'validation' && renderValidationTab()}
        {activeTab === 'tracking' && renderTrackingTab()}
        {activeTab === 'performance' && renderPerformanceTab()}
      </div>

      {/* Footer */}
      <div className="text-center text-sm text-gray-500 border-t pt-4">
        <div>
          Cross-Platform Integration Services Demo -
          Component: <span className="font-mono">{selectedComponentId}</span>
        </div>
        <div className="mt-1">
          Services: Component Sync (T051) • Design Token Validation (T052) •
          Implementation Tracking (T053) • Performance & Accessibility Validation (T054-T056)
        </div>
      </div>
    </div>
  );
};

export default CrossPlatformIntegrationDemo;