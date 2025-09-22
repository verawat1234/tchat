//
//  UIStateManager.swift
//  TchatApp
//
//  Created by Claude on 22/09/2024.
//

import Foundation
import SwiftUI
import Combine

/// UI state manager for managing component states and synchronization
@MainActor
public class UIStateManager: ObservableObject {

    // MARK: - Published Properties

    @Published public var componentStates: [String: ComponentState] = [:]
    @Published public var isSyncing: Bool = false
    @Published public var lastSyncTimestamp: Date?
    @Published public var syncErrors: [UIStateError] = []

    // MARK: - Private Properties

    private let syncService: ComponentSyncService
    private let persistenceService: UIStatePersistence
    private let conflictResolver: StateConflictResolver
    private var cancellables = Set<AnyCancellable>()
    private var syncTimer: Timer?

    // MARK: - Initialization

    public init(
        syncService: ComponentSyncService = ComponentSyncService(),
        persistenceService: UIStatePersistence = UIStatePersistence(),
        conflictResolver: StateConflictResolver = StateConflictResolver()
    ) {
        self.syncService = syncService
        self.persistenceService = persistenceService
        self.conflictResolver = conflictResolver

        setupAutoSync()
        loadComponentStates()
    }

    // MARK: - Public Methods

    /// Get component state by instance ID
    public func getComponentState(_ instanceId: String) -> ComponentState? {
        return componentStates[instanceId]
    }

    /// Update component state
    public func updateComponentState(
        instanceId: String,
        updates: [String: Any],
        shouldSync: Bool = true
    ) async throws {
        guard var state = componentStates[instanceId] else {
            throw UIStateError.componentStateNotFound(instanceId)
        }

        // Update state
        state.updateProperties(updates)
        componentStates[instanceId] = state

        // Persist changes
        try await persistenceService.saveComponentState(state)

        // Sync if requested
        if shouldSync {
            try await syncComponentState(instanceId)
        }
    }

    /// Create new component state
    public func createComponentState(
        _ state: ComponentState,
        shouldSync: Bool = true
    ) async throws {
        // Check for existing state
        if componentStates[state.instanceId] != nil {
            throw UIStateError.componentStateAlreadyExists(state.instanceId)
        }

        // Store state
        componentStates[state.instanceId] = state

        // Persist
        try await persistenceService.saveComponentState(state)

        // Sync if requested
        if shouldSync {
            try await syncComponentState(state.instanceId)
        }
    }

    /// Remove component state
    public func removeComponentState(
        _ instanceId: String,
        shouldSync: Bool = true
    ) async throws {
        guard let state = componentStates[instanceId] else {
            throw UIStateError.componentStateNotFound(instanceId)
        }

        // Remove from memory
        componentStates.removeValue(forKey: instanceId)

        // Remove from persistence
        try await persistenceService.removeComponentState(instanceId)

        // Sync removal if requested
        if shouldSync {
            // TODO: Implement state removal sync
        }
    }

    /// Merge component state with external state
    public func mergeComponentState(
        instanceId: String,
        externalState: [String: Any],
        strategy: MergeStrategy = .overwrite
    ) async throws {
        guard var state = componentStates[instanceId] else {
            throw UIStateError.componentStateNotFound(instanceId)
        }

        // Merge states
        state.mergeState(externalState, strategy: strategy)
        componentStates[instanceId] = state

        // Persist changes
        try await persistenceService.saveComponentState(state)
    }

    /// Sync specific component state
    public func syncComponentState(_ instanceId: String) async throws {
        guard let state = componentStates[instanceId] else {
            throw UIStateError.componentStateNotFound(instanceId)
        }

        isSyncing = true
        defer { isSyncing = false }

        do {
            let syncRequest = state.createSyncRequest()
            let response = try await syncService.syncComponentState(request: syncRequest)

            if response.success {
                var updatedState = state
                updatedState = updatedState.applySyncResponse(response)
                componentStates[instanceId] = updatedState

                // Update persistence
                try await persistenceService.saveComponentState(updatedState)

                lastSyncTimestamp = Date()
            }

        } catch {
            let uiError = UIStateError.syncFailed(instanceId, error)
            syncErrors.append(uiError)
            throw uiError
        }
    }

    /// Sync all component states
    public func syncAllComponentStates() async throws {
        guard !componentStates.isEmpty else { return }

        isSyncing = true
        defer { isSyncing = false }

        let states = Array(componentStates.values)
        let syncRequest = ComponentStateSyncRequest(
            userId: getCurrentUserId(),
            sessionId: getCurrentSessionId(),
            platform: "ios",
            componentStates: states,
            timestamp: Date(),
            syncVersion: 1
        )

        do {
            let response = try await syncService.syncComponentState(request: syncRequest)

            if response.success {
                // Handle conflicts if any
                if !response.conflictsResolved.isEmpty {
                    try await handleSyncConflicts(response.conflictsResolved)
                }

                lastSyncTimestamp = Date()
            }

        } catch {
            let uiError = UIStateError.syncAllFailed(error)
            syncErrors.append(uiError)
            throw uiError
        }
    }

    /// Get states by component type
    public func getComponentStates(for componentId: String) -> [ComponentState] {
        return componentStates.values.filter { $0.componentId == componentId }
    }

    /// Get all synchronized states
    public func getSynchronizedStates() -> [ComponentState] {
        return componentStates.values.filter { $0.isSynchronized }
    }

    /// Get unsynchronized states
    public func getUnsynchronizedStates() -> [ComponentState] {
        return componentStates.values.filter { !$0.isSynchronized }
    }

    /// Clear all errors
    public func clearErrors() {
        syncErrors.removeAll()
    }

    /// Reset all component states
    public func resetAllStates() async throws {
        componentStates.removeAll()
        try await persistenceService.clearAllStates()
        lastSyncTimestamp = nil
        syncErrors.removeAll()
    }

    // MARK: - Private Methods

    private func setupAutoSync() {
        // Auto-sync every 30 seconds
        syncTimer = Timer.scheduledTimer(withTimeInterval: 30.0, repeats: true) { [weak self] _ in
            Task { @MainActor in
                try? await self?.syncUnsynchronizedStates()
            }
        }

        // Sync when app becomes active
        NotificationCenter.default.publisher(for: UIApplication.didBecomeActiveNotification)
            .sink { [weak self] _ in
                Task { @MainActor in
                    try? await self?.syncUnsynchronizedStates()
                }
            }
            .store(in: &cancellables)
    }

    private func syncUnsynchronizedStates() async throws {
        let unsyncedStates = getUnsynchronizedStates()

        for state in unsyncedStates {
            try? await syncComponentState(state.instanceId)
        }
    }

    private func loadComponentStates() {
        Task {
            do {
                let loadedStates = try await persistenceService.loadAllComponentStates()

                for state in loadedStates {
                    componentStates[state.instanceId] = state
                }

            } catch {
                // Log error but don't fail initialization
                print("Failed to load component states: \(error)")
            }
        }
    }

    private func handleSyncConflicts(_ conflictIds: [String]) async throws {
        for conflictId in conflictIds {
            guard let localState = componentStates[conflictId] else { continue }

            // TODO: Fetch remote state and resolve conflict
            // For now, we'll keep the local state
        }
    }

    private func getCurrentUserId() -> String {
        // TODO: Get from authentication service
        return "current_user_id"
    }

    private func getCurrentSessionId() -> String {
        // TODO: Get from session service
        return UUID().uuidString
    }

    deinit {
        syncTimer?.invalidate()
    }
}

// MARK: - Supporting Services

/// Component sync service
public class ComponentSyncService {

    public init() {}

    public func syncComponentState(request: ComponentStateSyncRequest) async throws -> ComponentStateSyncResponse {
        // TODO: Implement actual sync with backend
        return ComponentStateSyncResponse(
            success: true,
            syncVersion: request.syncVersion + 1,
            conflictsResolved: [],
            timestamp: Date()
        )
    }
}

/// UI state persistence service
public class UIStatePersistence {

    public init() {}

    public func saveComponentState(_ state: ComponentState) async throws {
        // TODO: Implement persistence (UserDefaults, Core Data, etc.)
    }

    public func loadComponentState(_ instanceId: String) async throws -> ComponentState? {
        // TODO: Implement loading
        return nil
    }

    public func loadAllComponentStates() async throws -> [ComponentState] {
        // TODO: Implement loading all states
        return []
    }

    public func removeComponentState(_ instanceId: String) async throws {
        // TODO: Implement removal
    }

    public func clearAllStates() async throws {
        // TODO: Implement clearing all states
    }
}

/// State conflict resolver
public class StateConflictResolver {

    public init() {}

    public func resolveConflict(
        local: ComponentState,
        remote: ComponentState,
        strategy: ConflictResolutionStrategy = .newestWins
    ) -> ComponentState {
        var resolvedState = local

        switch strategy {
        case .newestWins:
            if remote.timestamp > local.timestamp {
                resolvedState = remote
            }
        case .highestVersionWins:
            if remote.version > local.version {
                resolvedState = remote
            }
        case .merge:
            resolvedState.mergeState(remote.state, strategy: .merge)
        case .manual:
            // Manual resolution required - return local for now
            break
        }

        return resolvedState
    }
}

// MARK: - Error Types

public enum UIStateError: Error, LocalizedError {
    case componentStateNotFound(String)
    case componentStateAlreadyExists(String)
    case syncFailed(String, Error)
    case syncAllFailed(Error)
    case persistenceFailed(Error)
    case conflictResolutionFailed(String)

    public var errorDescription: String? {
        switch self {
        case .componentStateNotFound(let id):
            return "Component state not found: \(id)"
        case .componentStateAlreadyExists(let id):
            return "Component state already exists: \(id)"
        case .syncFailed(let id, let error):
            return "Sync failed for component \(id): \(error.localizedDescription)"
        case .syncAllFailed(let error):
            return "Sync all failed: \(error.localizedDescription)"
        case .persistenceFailed(let error):
            return "Persistence failed: \(error.localizedDescription)"
        case .conflictResolutionFailed(let id):
            return "Conflict resolution failed for component: \(id)"
        }
    }
}