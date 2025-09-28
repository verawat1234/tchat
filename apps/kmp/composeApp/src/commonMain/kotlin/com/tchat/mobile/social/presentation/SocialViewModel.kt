package com.tchat.mobile.social.presentation

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.tchat.mobile.social.data.repository.SocialRepository
import com.tchat.mobile.social.data.repository.SyncState
import com.tchat.mobile.social.domain.models.*
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.launch
import kotlinx.datetime.Clock

/**
 * KMP Social ViewModel
 *
 * Cross-platform ViewModel for social features with:
 * - Reactive state management
 * - Offline-first operations
 * - Southeast Asian regional features
 * - Mobile performance optimization
 */
class SocialViewModel(
    private val repository: SocialRepository,
    private val currentUserId: String = "" // Should come from user session
) : ViewModel() {

    // UI State
    private val _uiState = MutableStateFlow(SocialUiState())
    val uiState: StateFlow<SocialUiState> = _uiState.asStateFlow()

    // Feed State
    private val _feedState = MutableStateFlow(FeedState())
    val feedState: StateFlow<FeedState> = _feedState.asStateFlow()

    // Profile State
    private val _profileState = MutableStateFlow(ProfileState())
    val profileState: StateFlow<ProfileState> = _profileState.asStateFlow()

    // Discovery State
    private val _discoveryState = MutableStateFlow(DiscoveryState())
    val discoveryState: StateFlow<DiscoveryState> = _discoveryState.asStateFlow()

    // Sync State
    val syncState: StateFlow<SyncState> = repository.syncState

    init {
        // Load initial data
        loadProfile(currentUserId)
        loadHomeFeed()
        loadDiscoveryFeed()
    }

    // Profile Management
    fun loadProfile(userId: String) {
        viewModelScope.launch {
            _profileState.update { it.copy(isLoading = true, error = null) }

            repository.getProfileFlow(userId)
                .catch { error ->
                    _profileState.update { it.copy(isLoading = false, error = error.message) }
                }
                .collect { profile ->
                    _profileState.update {
                        it.copy(
                            isLoading = false,
                            profile = profile,
                            error = null
                        )
                    }
                }
        }
    }

    fun updateProfile(request: UpdateProfileRequest) {
        viewModelScope.launch {
            _profileState.update { it.copy(isUpdating = true, error = null) }

            repository.updateProfile(currentUserId, request)
                .fold(
                    onSuccess = { updatedProfile ->
                        _profileState.update {
                            it.copy(
                                isUpdating = false,
                                profile = updatedProfile,
                                error = null
                            )
                        }
                        _uiState.update { it.copy(message = "Profile updated successfully") }
                    },
                    onFailure = { error ->
                        _profileState.update {
                            it.copy(
                                isUpdating = false,
                                error = error.message
                            )
                        }
                    }
                )
        }
    }

    // Feed Management
    fun loadHomeFeed(refresh: Boolean = false) {
        viewModelScope.launch {
            if (refresh) {
                _feedState.update { it.copy(isRefreshing = true) }
            } else {
                _feedState.update { it.copy(isLoading = true, error = null) }
            }

            repository.getFeedFlow(currentUserId, "home", getCurrentRegion())
                .catch { error ->
                    _feedState.update {
                        it.copy(
                            isLoading = false,
                            isRefreshing = false,
                            error = error.message
                        )
                    }
                }
                .collect { feed ->
                    _feedState.update {
                        it.copy(
                            isLoading = false,
                            isRefreshing = false,
                            homeFeed = feed,
                            error = null
                        )
                    }
                }
        }
    }

    fun loadDiscoveryFeed() {
        viewModelScope.launch {
            _discoveryState.update { it.copy(isLoading = true, error = null) }

            repository.getDiscoveryFeed(currentUserId, getCurrentRegion())
                .fold(
                    onSuccess = { profiles ->
                        _discoveryState.update {
                            it.copy(
                                isLoading = false,
                                discoveryProfiles = profiles,
                                error = null
                            )
                        }
                    },
                    onFailure = { error ->
                        _discoveryState.update {
                            it.copy(
                                isLoading = false,
                                error = error.message
                            )
                        }
                    }
                )
        }
    }

    fun refreshFeed() {
        viewModelScope.launch {
            _feedState.update { it.copy(isRefreshing = true) }

            repository.refreshFeed(currentUserId, "home", getCurrentRegion())
                .fold(
                    onSuccess = { feed ->
                        _feedState.update {
                            it.copy(
                                isRefreshing = false,
                                homeFeed = feed,
                                error = null
                            )
                        }
                    },
                    onFailure = { error ->
                        _feedState.update {
                            it.copy(
                                isRefreshing = false,
                                error = error.message
                            )
                        }
                    }
                )
        }
    }

    // Post Management
    fun createPost(content: String, contentType: String = "text", mediaUrls: List<String> = emptyList()) {
        viewModelScope.launch {
            _uiState.update { it.copy(isPosting = true, error = null) }

            val request = CreatePostRequest(
                content = content,
                contentType = contentType,
                mediaUrls = mediaUrls,
                language = getCurrentLanguage(),
                region = getCurrentRegion()
            )

            repository.createPost(request)
                .fold(
                    onSuccess = { post ->
                        _uiState.update {
                            it.copy(
                                isPosting = false,
                                message = "Post created successfully"
                            )
                        }
                        // Refresh feed to show new post
                        loadHomeFeed(refresh = true)
                    },
                    onFailure = { error ->
                        _uiState.update {
                            it.copy(
                                isPosting = false,
                                error = error.message
                            )
                        }
                    }
                )
        }
    }

    fun likePost(postId: String) {
        viewModelScope.launch {
            repository.likePost(postId, currentUserId)
                .fold(
                    onSuccess = {
                        // Update post in feed
                        updatePostInFeed(postId) { post ->
                            post.copy(
                                isLikedByUser = !post.isLikedByUser,
                                likesCount = if (post.isLikedByUser) post.likesCount - 1 else post.likesCount + 1
                            )
                        }
                    },
                    onFailure = { error ->
                        _uiState.update { it.copy(error = error.message) }
                    }
                )
        }
    }

    fun bookmarkPost(postId: String) {
        viewModelScope.launch {
            repository.bookmarkPost(postId, currentUserId)
                .fold(
                    onSuccess = {
                        // Update post in feed
                        updatePostInFeed(postId) { post ->
                            post.copy(isBookmarkedByUser = !post.isBookmarkedByUser)
                        }
                        _uiState.update {
                            it.copy(
                                message = if (_feedState.value.homeFeed?.posts?.find { it.id == postId }?.isBookmarkedByUser == true)
                                    "Post removed from bookmarks"
                                else
                                    "Post bookmarked"
                            )
                        }
                    },
                    onFailure = { error ->
                        _uiState.update { it.copy(error = error.message) }
                    }
                )
        }
    }

    // Follow Management
    fun followUser(userId: String) {
        viewModelScope.launch {
            repository.followUser(currentUserId, userId)
                .fold(
                    onSuccess = {
                        _uiState.update { it.copy(message = "User followed") }
                        // Update discovery profiles
                        updateDiscoveryProfile(userId) { profile ->
                            profile.copy(
                                profile = profile.profile.copy(
                                    followersCount = profile.profile.followersCount + 1
                                )
                            )
                        }
                    },
                    onFailure = { error ->
                        _uiState.update { it.copy(error = error.message) }
                    }
                )
        }
    }

    fun unfollowUser(userId: String) {
        viewModelScope.launch {
            repository.unfollowUser(currentUserId, userId)
                .fold(
                    onSuccess = {
                        _uiState.update { it.copy(message = "User unfollowed") }
                        // Update discovery profiles
                        updateDiscoveryProfile(userId) { profile ->
                            profile.copy(
                                profile = profile.profile.copy(
                                    followersCount = maxOf(0, profile.profile.followersCount - 1)
                                )
                            )
                        }
                    },
                    onFailure = { error ->
                        _uiState.update { it.copy(error = error.message) }
                    }
                )
        }
    }

    fun loadFollowing() {
        viewModelScope.launch {
            repository.getFollowingFlow(currentUserId)
                .catch { error ->
                    _profileState.update { it.copy(error = error.message) }
                }
                .collect { following ->
                    _profileState.update { it.copy(following = following) }
                }
        }
    }

    fun loadFollowers() {
        viewModelScope.launch {
            repository.getFollowersFlow(currentUserId)
                .catch { error ->
                    _profileState.update { it.copy(error = error.message) }
                }
                .collect { followers ->
                    _profileState.update { it.copy(followers = followers) }
                }
        }
    }

    // Sync Management
    fun performSync() {
        viewModelScope.launch {
            _uiState.update { it.copy(isSyncing = true) }

            val lastSyncAt = _profileState.value.profile?.lastSyncAt
            repository.performIncrementalSync(currentUserId, lastSyncAt)
                .fold(
                    onSuccess = { syncResponse ->
                        _uiState.update {
                            it.copy(
                                isSyncing = false,
                                message = "Sync completed successfully"
                            )
                        }
                        // Refresh data after sync
                        loadHomeFeed(refresh = true)
                        loadProfile(currentUserId)
                    },
                    onFailure = { error ->
                        _uiState.update {
                            it.copy(
                                isSyncing = false,
                                error = error.message
                            )
                        }
                    }
                )
        }
    }

    // UI State Management
    fun clearError() {
        _uiState.update { it.copy(error = null) }
        _feedState.update { it.copy(error = null) }
        _profileState.update { it.copy(error = null) }
        _discoveryState.update { it.copy(error = null) }
    }

    fun clearMessage() {
        _uiState.update { it.copy(message = null) }
    }

    // Regional and Localization
    fun changeRegion(region: String) {
        viewModelScope.launch {
            _uiState.update { it.copy(currentRegion = region) }
            // Reload feeds for new region
            loadHomeFeed(refresh = true)
            loadDiscoveryFeed()
        }
    }

    fun changeLanguage(language: String) {
        viewModelScope.launch {
            _uiState.update { it.copy(currentLanguage = language) }
            // Reload localized content
            loadHomeFeed(refresh = true)
        }
    }

    // Private helper methods
    private fun getCurrentRegion(): String {
        return _uiState.value.currentRegion
    }

    private fun getCurrentLanguage(): String {
        return _uiState.value.currentLanguage
    }

    private fun updatePostInFeed(postId: String, transform: (SocialPost) -> SocialPost) {
        _feedState.update { state ->
            state.homeFeed?.let { feed ->
                val updatedPosts = feed.posts.map { post ->
                    if (post.id == postId) transform(post) else post
                }
                state.copy(
                    homeFeed = feed.copy(posts = updatedPosts)
                )
            } ?: state
        }
    }

    private fun updateDiscoveryProfile(userId: String, transform: (DiscoveryProfile) -> DiscoveryProfile) {
        _discoveryState.update { state ->
            val updatedProfiles = state.discoveryProfiles.map { profile ->
                if (profile.profile.id == userId) transform(profile) else profile
            }
            state.copy(discoveryProfiles = updatedProfiles)
        }
    }
}

// UI State Data Classes
data class SocialUiState(
    val isPosting: Boolean = false,
    val isSyncing: Boolean = false,
    val error: String? = null,
    val message: String? = null,
    val currentRegion: String = "TH",
    val currentLanguage: String = "en"
)

data class FeedState(
    val isLoading: Boolean = false,
    val isRefreshing: Boolean = false,
    val homeFeed: SocialFeed? = null,
    val error: String? = null
)

data class ProfileState(
    val isLoading: Boolean = false,
    val isUpdating: Boolean = false,
    val profile: SocialProfile? = null,
    val following: List<SocialProfile> = emptyList(),
    val followers: List<SocialProfile> = emptyList(),
    val error: String? = null
)

data class DiscoveryState(
    val isLoading: Boolean = false,
    val discoveryProfiles: List<DiscoveryProfile> = emptyList(),
    val error: String? = null
)