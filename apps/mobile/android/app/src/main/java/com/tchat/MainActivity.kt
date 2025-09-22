package com.tchat

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.viewModels
import androidx.compose.animation.core.tween
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ChatBubble
import androidx.compose.material.icons.filled.Email
import androidx.compose.material.icons.filled.Lock
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.core.splashscreen.SplashScreen.Companion.installSplashScreen
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import com.tchat.components.*
import com.tchat.designsystem.Colors
import com.tchat.designsystem.Spacing
import com.tchat.navigation.TabNavigationComposable
import com.tchat.state.AppState
import com.tchat.state.UserModel
import com.tchat.state.UserPreferences
import com.tchat.ui.theme.TchatTheme

/**
 * Main activity for the Tchat Android app
 */
class MainActivity : ComponentActivity() {

    private val appState: AppState by viewModels()

    override fun onCreate(savedInstanceState: Bundle?) {
        // Handle splash screen
        installSplashScreen()

        super.onCreate(savedInstanceState)

        // Configure app
        configureApp()

        setContent {
            TchatTheme {
                TchatApp(appState = appState)
            }
        }
    }

    private fun configureApp() {
        // Configure design system
        configureDesignSystem()

        // Load initial data
        loadInitialData()

        // Setup analytics if needed
        setupAnalytics()
    }

    private fun configureDesignSystem() {
        // Apply global theme configurations
        println("Configuring design system...")
    }

    private fun loadInitialData() {
        // Load any necessary initial data
        println("Loading initial app data...")
    }

    private fun setupAnalytics() {
        // Setup analytics tracking if needed
        println("Setting up analytics...")
    }
}

/**
 * Main application composable
 */
@Composable
fun TchatApp(appState: AppState) {
    val isAuthenticated by appState.isAuthenticated.collectAsStateWithLifecycle()

    // Animate between authentication and main app
    androidx.compose.animation.AnimatedVisibility(
        visible = isAuthenticated,
        enter = fadeIn(animationSpec = tween(300)),
        exit = fadeOut(animationSpec = tween(300))
    ) {
        // Main app interface
        TabNavigationComposable()
    }

    androidx.compose.animation.AnimatedVisibility(
        visible = !isAuthenticated,
        enter = fadeIn(animationSpec = tween(300)),
        exit = fadeOut(animationSpec = tween(300))
    ) {
        // Authentication flow
        AuthenticationScreen(appState = appState)
    }
}

/**
 * Authentication screen composable
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AuthenticationScreen(appState: AppState) {
    var email by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }

    Scaffold { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(Spacing.lg),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(Spacing.lg)
        ) {
            // App logo/branding
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(Spacing.md)
            ) {
                Icon(
                    imageVector = Icons.Default.ChatBubble,
                    contentDescription = "Tchat Logo",
                    tint = Colors.primary,
                    modifier = Modifier.size(80.dp)
                )

                Text(
                    text = "Tchat",
                    fontSize = 32.sp,
                    fontWeight = FontWeight.Bold,
                    color = Colors.textPrimary
                )

                Text(
                    text = "Connect with the world",
                    fontSize = 16.sp,
                    color = Colors.textSecondary
                )
            }

            Spacer(modifier = Modifier.weight(1f))

            // Login form
            Column(
                verticalArrangement = Arrangement.spacedBy(Spacing.md)
            ) {
                TchatInput(
                    value = email,
                    onValueChange = { email = it },
                    placeholder = "Email",
                    type = TchatInputType.Email,
                    leadingIcon = Icons.Default.Email
                )

                TchatInput(
                    value = password,
                    onValueChange = { password = it },
                    placeholder = "Password",
                    type = TchatInputType.Password,
                    leadingIcon = Icons.Default.Lock
                )

                Spacer(modifier = Modifier.height(Spacing.sm))

                TchatButton(
                    text = "Sign In",
                    onClick = { authenticateUser(appState, email) },
                    variant = TchatButtonVariant.Primary,
                    size = TchatButtonSize.Large,
                    modifier = Modifier.fillMaxWidth()
                )
            }

            Spacer(modifier = Modifier.weight(1f))

            // Demo login for testing
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(Spacing.sm)
            ) {
                Text(
                    text = "Demo Mode",
                    fontSize = 12.sp,
                    color = Colors.textTertiary
                )

                TchatButton(
                    text = "Continue as Demo User",
                    onClick = { authenticateAsDemoUser(appState) },
                    variant = TchatButtonVariant.Outline,
                    size = TchatButtonSize.Medium
                )
            }
        }
    }
}

// MARK: - Authentication Logic

private fun authenticateUser(appState: AppState, email: String) {
    // Simulate authentication
    val username = email.substringBefore("@").ifEmpty { "user" }
    val user = UserModel(
        id = "user_123",
        username = username,
        email = email,
        displayName = "User Name",
        preferences = UserPreferences()
    )

    appState.updateUser(user)
}

private fun authenticateAsDemoUser(appState: AppState) {
    val demoUser = UserModel(
        id = "demo_user",
        username = "demo",
        email = "demo@tchat.app",
        displayName = "Demo User",
        preferences = UserPreferences()
    )

    appState.updateUser(demoUser)
}

// MARK: - Previews

@Preview(showBackground = true)
@Composable
fun AuthenticationScreenPreview() {
    TchatTheme {
        AuthenticationScreen(appState = AppState())
    }
}

@Preview(showBackground = true)
@Composable
fun TchatAppPreview() {
    TchatTheme {
        TchatApp(appState = AppState())
    }
}