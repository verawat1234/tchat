//
//  StateSyncService.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import Combine
import Network

/// Cross-platform state synchronization service for iOS
/// Implements T049: Cross-platform state synchronization service
public class StateSyncService: ObservableObject {

    // MARK: - Types

    public enum SyncStatus {
        case offline
        case connecting
        case connected
        case syncing
        case synced
        case error(String)
    }

    public enum SyncStrategy {
        case optimistic    // Apply changes immediately, rollback on conflict
        case conservative  // Wait for server confirmation
        case automatic     // Choose based on network conditions
    }

    // MARK: - Published Properties

    @Published public var syncStatus: SyncStatus = .offline
    @Published public var lastSyncTimestamp: Date?
    @Published public var pendingSyncCount: Int = 0
    @Published public var isOnline: Bool = false

    // MARK: - Private Properties

    private let navigationSyncClient: NavigationSyncAPIClient
    private let uiComponentSyncClient: UIComponentSyncAPIClient
    private let networkMonitor = NWPathMonitor()
    private let syncQueue = DispatchQueue(label: "com.tchat.statesync", qos: .utility)
    private var cancellables = Set<AnyCancellable>()

    // Sync configuration
    private let maxRetries = 3
    private let syncInterval: TimeInterval = 30.0
    private let batchSize = 10

    // Pending sync operations
    private var pendingEvents: [SyncEvent] = []
    private var conflictResolutionQueue: [SyncEvent] = []
    private var syncTimer: Timer?

    // MARK: - Initialization

    public init(
        navigationSyncClient: NavigationSyncAPIClient,
        uiComponentSyncClient: UIComponentSyncAPIClient
    ) {
        self.navigationSyncClient = navigationSyncClient
        self.uiComponentSyncClient = uiComponentSyncClient

        setupNetworkMonitoring()
        setupPeriodicSync()
    }

    deinit {
        stopNetworkMonitoring()
        stopPeriodicSync()
    }

    // MARK: - Public Interface

    /// Starts the state synchronization service
    public func start() {
        syncStatus = .connecting

        Task {
            await performInitialSync()
        }
    }

    /// Stops the state synchronization service
    public func stop() {
        syncStatus = .offline
        stopPeriodicSync()
    }

    /// Queues a sync event for processing
    public func queueSync(_ event: SyncEvent) {
        syncQueue.async { [weak self] in
            guard let self = self else { return }

            self.pendingEvents.append(event)

            DispatchQueue.main.async {
                self.pendingSyncCount = self.pendingEvents.count
            }

            // Trigger immediate sync for high-priority events
            if event.requiresAck {
                Task {
                    await self.processPendingEvents()
                }
            }
        }
    }

    /// Forces an immediate sync of all pending events
    public func forcSync() async {
        guard isOnline else {
            updateSyncStatus(.error("Cannot sync while offline"))
            return
        }

        updateSyncStatus(.syncing)
        await processPendingEvents()
    }

    /// Retrieves the current sync statistics
    public func getSyncStatistics() -> SyncStatistics {
        return SyncStatistics(
            pendingEventCount: pendingSyncCount,
            lastSyncTimestamp: lastSyncTimestamp,
            isOnline: isOnline,
            conflictsCount: conflictResolutionQueue.count
        )
    }

    // MARK: - Network Monitoring

    private func setupNetworkMonitoring() {
        networkMonitor.pathUpdateHandler = { [weak self] path in
            DispatchQueue.main.async {
                self?.isOnline = path.status == .satisfied

                if path.status == .satisfied {
                    self?.onNetworkAvailable()
                } else {
                    self?.onNetworkUnavailable()
                }
            }
        }

        let queue = DispatchQueue(label: "NetworkMonitor")
        networkMonitor.start(queue: queue)
    }

    private func stopNetworkMonitoring() {
        networkMonitor.cancel()
    }

    private func onNetworkAvailable() {
        guard syncStatus == .offline else { return }

        syncStatus = .connecting

        Task {
            await performInitialSync()
        }
    }

    private func onNetworkUnavailable() {
        syncStatus = .offline
        stopPeriodicSync()
    }

    // MARK: - Periodic Sync

    private func setupPeriodicSync() {
        syncTimer = Timer.scheduledTimer(withTimeInterval: syncInterval, repeats: true) { [weak self] _ in
            Task {
                await self?.processPendingEvents()
            }
        }
    }

    private func stopPeriodicSync() {
        syncTimer?.invalidate()
        syncTimer = nil
    }

    // MARK: - Sync Processing

    private func performInitialSync() async {
        do {
            // Sync navigation state
            _ = try await navigationSyncClient.getCurrentNavigationState(
                platform: "ios",
                userId: getCurrentUserId(),
                sessionId: getCurrentSessionId()
            )

            // Sync UI component state
            _ = try await uiComponentSyncClient.getComponentStates(
                platform: "ios",
                userId: getCurrentUserId()
            )

            updateSyncStatus(.synced)
            updateLastSyncTimestamp()

            // Process any pending events after initial sync
            await processPendingEvents()

        } catch {
            updateSyncStatus(.error(error.localizedDescription))
        }
    }

    private func processPendingEvents() async {
        guard isOnline, !pendingEvents.isEmpty else { return }

        updateSyncStatus(.syncing)

        await syncQueue.sync {
            // Process events in batches
            let batchesToProcess = pendingEvents.chunked(into: batchSize)

            for batch in batchesToProcess {
                Task {
                    await processBatch(batch)
                }
            }
        }
    }

    private func processBatch(_ events: [SyncEvent]) async {
        for event in events {
            do {
                try await processEvent(event)

                // Remove successfully processed event
                syncQueue.async { [weak self] in
                    self?.pendingEvents.removeAll { $0.id == event.id }

                    DispatchQueue.main.async {
                        self?.pendingSyncCount = self?.pendingEvents.count ?? 0
                    }
                }

            } catch {
                await handleSyncError(event, error: error)
            }
        }

        updateSyncStatus(.synced)
        updateLastSyncTimestamp()
    }

    private func processEvent(_ event: SyncEvent) async throws {
        switch event.type {
        case .navigation:
            try await processNavigationEvent(event)
        case .stateUpdate:
            try await processStateUpdateEvent(event)
        case .dataChange:
            try await processDataChangeEvent(event)
        }
    }

    private func processNavigationEvent(_ event: SyncEvent) async throws {
        guard let fromRoute = event.getPayloadValue("fromRoute", as: String.self),
              let toRoute = event.getPayloadValue("toRoute", as: String.self) else {
            throw StateSyncError.invalidEventPayload("Missing navigation routes")
        }

        _ = try await navigationSyncClient.recordNavigationEvent(
            fromRoute: fromRoute,
            toRoute: toRoute,
            trigger: event.getPayloadValue("trigger", as: String.self) ?? "user",
            userId: event.userId,
            sessionId: event.sessionId,
            platform: "ios"
        )
    }

    private func processStateUpdateEvent(_ event: SyncEvent) async throws {
        _ = try await uiComponentSyncClient.updateComponentState(
            componentId: event.getPayloadValue("componentId", as: String.self) ?? "",
            state: event.payload,
            userId: event.userId,
            sessionId: event.sessionId,
            platform: "ios"
        )
    }

    private func processDataChangeEvent(_ event: SyncEvent) async throws {
        // Handle generic data change events
        // This could route to different services based on entity type
        let entity = event.getPayloadValue("entity", as: String.self) ?? ""
        let action = event.getPayloadValue("action", as: String.self) ?? ""

        // Log the data change for now - in a real implementation,
        // this would route to the appropriate service
        print("Data change: \(entity) - \(action)")
    }

    // MARK: - Error Handling

    private func handleSyncError(_ event: SyncEvent, error: Error) async {
        var updatedEvent = event

        do {
            updatedEvent = try updatedEvent.updateStatus(to: .failed, errorMessage: error.localizedDescription)

            if updatedEvent.shouldRetry {
                // Add back to queue with incremented retry count
                let retriedEvent = try updatedEvent.incrementRetry()

                syncQueue.async { [weak self] in
                    self?.pendingEvents.append(retriedEvent)
                }
            } else {
                // Max retries exceeded - add to conflict resolution queue
                syncQueue.async { [weak self] in
                    self?.conflictResolutionQueue.append(updatedEvent)
                }
            }

        } catch {
            // If we can't even update the event status, log the error
            print("Failed to handle sync error for event \(event.id): \(error)")
        }
    }

    // MARK: - Conflict Resolution

    /// Resolves sync conflicts using last-write-wins strategy
    private func resolveConflicts() async {
        guard !conflictResolutionQueue.isEmpty else { return }

        for event in conflictResolutionQueue {
            // For now, implement last-write-wins
            // In a production app, this might show user prompts
            await retryEvent(event)
        }

        conflictResolutionQueue.removeAll()
    }

    private func retryEvent(_ event: SyncEvent) async {
        var retriedEvent = event
        retriedEvent.retryCount = 0 // Reset retry count for conflict resolution

        do {
            retriedEvent = try retriedEvent.updateStatus(to: .pending)
            queueSync(retriedEvent)
        } catch {
            print("Failed to retry conflicted event \(event.id): \(error)")
        }
    }

    // MARK: - Helper Methods

    private func updateSyncStatus(_ status: SyncStatus) {
        DispatchQueue.main.async { [weak self] in
            self?.syncStatus = status
        }
    }

    private func updateLastSyncTimestamp() {
        DispatchQueue.main.async { [weak self] in
            self?.lastSyncTimestamp = Date()
        }
    }

    private func getCurrentUserId() -> String {
        // In a real app, this would get the current user ID from auth service
        return "current_user_id"
    }

    private func getCurrentSessionId() -> String {
        // In a real app, this would get the current session ID
        return "current_session_id"
    }
}

// MARK: - Supporting Types

public struct SyncStatistics {
    public let pendingEventCount: Int
    public let lastSyncTimestamp: Date?
    public let isOnline: Bool
    public let conflictsCount: Int
}

public enum StateSyncError: LocalizedError {
    case invalidEventPayload(String)
    case networkUnavailable
    case syncInProgress
    case authenticationRequired

    public var errorDescription: String? {
        switch self {
        case .invalidEventPayload(let message):
            return "Invalid event payload: \(message)"
        case .networkUnavailable:
            return "Network is not available"
        case .syncInProgress:
            return "Sync operation already in progress"
        case .authenticationRequired:
            return "Authentication required for sync"
        }
    }
}

// MARK: - Array Extension

private extension Array {
    func chunked(into size: Int) -> [[Element]] {
        return stride(from: 0, to: count, by: size).map {
            Array(self[$0..<Swift.min($0 + size, count)])
        }
    }
}