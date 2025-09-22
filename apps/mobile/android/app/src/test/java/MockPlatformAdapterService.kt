import kotlinx.coroutines.delay
import java.util.Date

/**
 * Mock service for Platform Adapter testing
 */
class MockPlatformAdapterService {

    suspend fun getPlatformCapabilities(platform: String): PlatformCapabilitiesResponse {
        // Simulate network delay
        delay(100)

        return PlatformCapabilitiesResponse(
            platform = platform,
            capabilities = mapOf(
                "ui" to mapOf(
                    "components" to listOf("button", "input", "card"),
                    "theming" to true,
                    "accessibility" to true
                ),
                "navigation" to mapOf(
                    "tabNavigation" to true,
                    "deepLinking" to true,
                    "backNavigation" to true
                ),
                "storage" to mapOf(
                    "localStorage" to true,
                    "secureStorage" to true,
                    "offline" to true
                )
            ),
            version = "1.0.0",
            lastUpdated = Date()
        )
    }

    suspend fun validatePlatformCompatibility(request: CompatibilityRequest): CompatibilityResponse {
        // Simulate network delay
        delay(100)

        return CompatibilityResponse(
            compatible = true,
            message = "Platform compatibility validated successfully",
            supportedFeatures = request.requiredFeatures,
            unsupportedFeatures = emptyList(),
            timestamp = Date()
        )
    }
}

/**
 * Data models for testing
 */

data class PlatformCapabilitiesResponse(
    val platform: String,
    val capabilities: Map<String, Any>,
    val version: String,
    val lastUpdated: Date
)

data class CompatibilityRequest(
    val platform: String,
    val requiredFeatures: List<String>,
    val targetVersion: String
)

data class CompatibilityResponse(
    val compatible: Boolean,
    val message: String,
    val supportedFeatures: List<String>,
    val unsupportedFeatures: List<String>,
    val timestamp: Date
)