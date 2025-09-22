import kotlinx.coroutines.delay
import java.util.Date

/**
 * Mock service for UI Component Sync testing
 */
class MockUIComponentSyncService {

    suspend fun getComponentRegistry(platform: String): ComponentRegistry {
        // Simulate network delay
        delay(100)

        return ComponentRegistry(
            platform = platform,
            components = listOf(
                UIComponent(
                    id = "tchat-button",
                    name = "TchatButton",
                    version = "1.0.0",
                    properties = mapOf("type" to "button", "variants" to listOf("primary", "secondary"))
                ),
                UIComponent(
                    id = "tchat-input",
                    name = "TchatInput",
                    version = "1.0.0",
                    properties = mapOf("type" to "input", "validation" to true)
                )
            ),
            lastUpdated = Date()
        )
    }

    suspend fun syncComponentData(request: ComponentSyncRequest): ComponentSyncResponse {
        // Simulate network delay
        delay(100)

        return ComponentSyncResponse(
            success = true,
            message = "Components synced successfully",
            syncedComponents = request.components.size,
            timestamp = Date()
        )
    }
}

/**
 * Data models for testing
 */

data class ComponentRegistry(
    val platform: String,
    val components: List<UIComponent>,
    val lastUpdated: Date
)

data class UIComponent(
    val id: String,
    val name: String,
    val version: String,
    val properties: Map<String, Any>
)

data class ComponentSyncRequest(
    val platform: String,
    val components: List<UIComponent>,
    val timestamp: Date
)

data class ComponentSyncResponse(
    val success: Boolean,
    val message: String,
    val syncedComponents: Int,
    val timestamp: Date
)