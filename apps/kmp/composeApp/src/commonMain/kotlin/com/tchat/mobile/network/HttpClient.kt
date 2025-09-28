package com.tchat.mobile.network

import io.ktor.client.*
import io.ktor.client.plugins.*
import io.ktor.client.plugins.auth.*
import io.ktor.client.plugins.auth.providers.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.plugins.logging.*
import io.ktor.client.plugins.websocket.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.json.Json

/**
 * HTTP Client factory for cross-platform API communication
 */
object HttpClientFactory {

    fun create(baseUrl: String): HttpClient {
        return HttpClient {
            // JSON content negotiation
            install(ContentNegotiation) {
                json(Json {
                    ignoreUnknownKeys = true
                    encodeDefaults = true
                    prettyPrint = true
                })
            }

            // Logging
            install(Logging) {
                logger = Logger.SIMPLE
                level = LogLevel.INFO
            }

            // WebSocket support
            install(WebSockets)

            // Default request configuration
            install(DefaultRequest) {
                url(baseUrl)
            }

            // HTTP timeout configuration
            install(HttpTimeout) {
                requestTimeoutMillis = 30000
                connectTimeoutMillis = 10000
                socketTimeoutMillis = 30000
            }

            // Auth configuration (will be set up later)
            install(Auth) {
                bearer {
                    loadTokens {
                        // Will be implemented with token storage
                        null
                    }
                }
            }
        }
    }
}