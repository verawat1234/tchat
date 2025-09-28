package com.tchat.mobile.screens

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.tchat.mobile.components.TchatButton
import com.tchat.mobile.components.TchatInput
import com.tchat.mobile.components.TchatInputType
import com.tchat.mobile.components.TchatCard
import com.tchat.mobile.designsystem.TchatColors
import com.tchat.mobile.designsystem.TchatSpacing
import com.tchat.mobile.models.*
import com.tchat.mobile.services.AuthService
import kotlinx.coroutines.launch
import org.koin.compose.koinInject

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AuthScreen(
    onAuthSuccess: (User) -> Unit,
    modifier: Modifier = Modifier
) {
    val authService: AuthService = koinInject()
    val scope = rememberCoroutineScope()

    // State management
    val authState by authService.authState.collectAsState()
    val authStep by authService.authStep.collectAsState()
    val authMethod by authService.authMethod.collectAsState()
    val isLoading by authService.isLoading.collectAsState()

    // Form state
    var email by remember { mutableStateOf("demo@tchat.app") }
    var password by remember { mutableStateOf("demo123") }
    var phoneNumber by remember { mutableStateOf("") }
    var otpCode by remember { mutableStateOf("") }
    var showPassword by remember { mutableStateOf(false) }

    // Handle authentication success
    LaunchedEffect(authState) {
        when (val state = authState) {
            is AuthState.Authenticated -> {
                onAuthSuccess(state.user)
            }
            else -> {} // Do nothing for other states
        }
    }

    // Show verification step
    if (authStep == AuthStep.VERIFY) {
        VerificationScreen(
            authMethod = authMethod,
            phoneNumber = phoneNumber,
            email = email,
            otpCode = otpCode,
            onOtpCodeChange = { otpCode = it },
            isLoading = isLoading,
            onVerifyClick = {
                scope.launch {
                    when (authMethod) {
                        AuthMethod.PHONE -> {
                            authService.verifyOtp(phoneNumber, otpCode)
                        }
                        AuthMethod.EMAIL -> {
                            // Email verification not implemented yet
                        }
                    }
                }
            },
            onBackClick = {
                authService.goBackToInput()
            },
            authState = authState,
            modifier = modifier
        )
        return
    }

    // Main authentication screen
    Column(
        modifier = modifier
            .fillMaxSize()
            .padding(TchatSpacing.lg),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        // Hero Section
        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
            modifier = Modifier.padding(bottom = TchatSpacing.xl)
        ) {
            // App Icon
            Box(
                modifier = Modifier
                    .size(80.dp)
                    .clip(CircleShape),
                contentAlignment = Alignment.Center
            ) {
                Card(
                    modifier = Modifier.fillMaxSize(),
                    colors = CardDefaults.cardColors(
                        containerColor = TchatColors.primary
                    )
                ) {
                    Box(
                        modifier = Modifier.fillMaxSize(),
                        contentAlignment = Alignment.Center
                    ) {
                        Icon(
                            imageVector = Icons.Default.Message,
                            contentDescription = "Tchat Logo",
                            tint = Color.White,
                            modifier = Modifier.size(40.dp)
                        )
                    }
                }
            }

            Spacer(modifier = Modifier.height(TchatSpacing.md))

            // App Title
            Text(
                text = "Tchat KMP",
                fontSize = 28.sp,
                fontWeight = FontWeight.Bold,
                color = TchatColors.onSurface
            )

            Text(
                text = "Cloud messaging, payments, and social commerce\nbuilt for Southeast Asia",
                fontSize = 16.sp,
                color = TchatColors.onSurfaceVariant,
                textAlign = TextAlign.Center,
                modifier = Modifier.padding(top = TchatSpacing.sm)
            )
        }

        // Feature Highlights
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = TchatSpacing.md),
            horizontalArrangement = Arrangement.SpaceEvenly
        ) {
            FeatureChip(
                icon = Icons.Default.Shield,
                text = "Encrypted",
                color = TchatColors.success
            )
            FeatureChip(
                icon = Icons.Default.Bolt,
                text = "Fast",
                color = TchatColors.warning
            )
            FeatureChip(
                icon = Icons.Default.Payment,
                text = "Payments",
                color = TchatColors.primary
            )
            FeatureChip(
                icon = Icons.Default.Language,
                text = "SEA",
                color = TchatColors.primary
            )
        }

        Spacer(modifier = Modifier.height(TchatSpacing.lg))

        // Auth Form Card
        TchatCard(
            variant = com.tchat.mobile.components.TchatCardVariant.Elevated,
            modifier = Modifier.fillMaxWidth()
        ) {
            Column(
                modifier = Modifier.padding(TchatSpacing.lg)
            ) {
                // Title
                Text(
                    text = "Sign In",
                    fontSize = 20.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = TchatColors.onSurface,
                    modifier = Modifier.padding(bottom = TchatSpacing.sm)
                )

                Text(
                    text = "Choose your preferred sign-in method",
                    fontSize = 14.sp,
                    color = TchatColors.onSurfaceVariant,
                    modifier = Modifier.padding(bottom = TchatSpacing.lg)
                )

                // Method Selection Tabs
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
                ) {
                    TchatButton(
                        text = "📧 Email",
                        onClick = { authService.setAuthMethod(AuthMethod.EMAIL) },
                        variant = if (authMethod == AuthMethod.EMAIL)
                            com.tchat.mobile.components.TchatButtonVariant.Primary
                        else
                            com.tchat.mobile.components.TchatButtonVariant.Outline,
                        size = com.tchat.mobile.components.TchatButtonSize.Small,
                        modifier = Modifier.weight(1f)
                    )
                    TchatButton(
                        text = "📱 Phone",
                        onClick = { authService.setAuthMethod(AuthMethod.PHONE) },
                        variant = if (authMethod == AuthMethod.PHONE)
                            com.tchat.mobile.components.TchatButtonVariant.Primary
                        else
                            com.tchat.mobile.components.TchatButtonVariant.Outline,
                        size = com.tchat.mobile.components.TchatButtonSize.Small,
                        modifier = Modifier.weight(1f)
                    )
                }

                Spacer(modifier = Modifier.height(TchatSpacing.lg))

                // Auth Form Content
                when (authMethod) {
                    AuthMethod.EMAIL -> {
                        EmailAuthForm(
                            email = email,
                            password = password,
                            showPassword = showPassword,
                            onEmailChange = { email = it },
                            onPasswordChange = { password = it },
                            onShowPasswordChange = { showPassword = it },
                            isLoading = isLoading,
                            onSignInClick = {
                                scope.launch {
                                    authService.loginWithEmail(email, password)
                                }
                            }
                        )
                    }
                    AuthMethod.PHONE -> {
                        PhoneAuthForm(
                            phoneNumber = phoneNumber,
                            onPhoneNumberChange = { phoneNumber = it },
                            isLoading = isLoading,
                            onSendOtpClick = {
                                scope.launch {
                                    authService.sendOtp(phoneNumber)
                                }
                            }
                        )
                    }
                }

                // Error Display
                when (val state = authState) {
                    is AuthState.Error -> {
                        Spacer(modifier = Modifier.height(TchatSpacing.md))
                        Card(
                            colors = CardDefaults.cardColors(
                                containerColor = TchatColors.error.copy(alpha = 0.1f)
                            ),
                            modifier = Modifier.fillMaxWidth()
                        ) {
                            Text(
                                text = state.message,
                                color = TchatColors.error,
                                fontSize = 14.sp,
                                modifier = Modifier.padding(TchatSpacing.md)
                            )
                        }
                    }
                    else -> {} // Do nothing for other states
                }

                // Privacy Notice
                Spacer(modifier = Modifier.height(TchatSpacing.lg))
                Text(
                    text = "By continuing, you agree to our Terms of Service and Privacy Policy. Built for PDPA (TH/MY), PDP (ID) compliance.",
                    fontSize = 12.sp,
                    color = TchatColors.onSurfaceVariant,
                    textAlign = TextAlign.Center,
                    modifier = Modifier.fillMaxWidth()
                )
            }
        }
    }
}

@Composable
private fun FeatureChip(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    text: String,
    color: Color
) {
    Card(
        colors = CardDefaults.cardColors(
            containerColor = color.copy(alpha = 0.1f)
        ),
        modifier = Modifier.padding(TchatSpacing.xs)
    ) {
        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
            modifier = Modifier.padding(TchatSpacing.sm)
        ) {
            Icon(
                imageVector = icon,
                contentDescription = text,
                tint = color,
                modifier = Modifier.size(20.dp)
            )
            Text(
                text = text,
                fontSize = 10.sp,
                color = color,
                fontWeight = FontWeight.Medium,
                modifier = Modifier.padding(top = 2.dp)
            )
        }
    }
}

@Composable
private fun EmailAuthForm(
    email: String,
    password: String,
    showPassword: Boolean,
    onEmailChange: (String) -> Unit,
    onPasswordChange: (String) -> Unit,
    onShowPasswordChange: (Boolean) -> Unit,
    isLoading: Boolean,
    onSignInClick: () -> Unit
) {
    Column {
        TchatInput(
            value = email,
            onValueChange = onEmailChange,
            placeholder = "Email address",
            leadingIcon = Icons.Default.Email,
            type = TchatInputType.Email,
            enabled = !isLoading
        )

        Spacer(modifier = Modifier.height(TchatSpacing.md))

        TchatInput(
            value = password,
            onValueChange = onPasswordChange,
            placeholder = "Password",
            leadingIcon = Icons.Default.Lock,
            trailingIcon = if (showPassword) Icons.Default.VisibilityOff else Icons.Default.Visibility,
            onTrailingIconClick = { onShowPasswordChange(!showPassword) },
            type = TchatInputType.Password,
            enabled = !isLoading
        )

        Text(
            text = "Demo credentials: demo@tchat.app / demo123",
            fontSize = 12.sp,
            color = TchatColors.onSurfaceVariant,
            modifier = Modifier.padding(top = TchatSpacing.sm)
        )

        Spacer(modifier = Modifier.height(TchatSpacing.lg))

        TchatButton(
            text = if (isLoading) "Signing in..." else "Sign In",
            onClick = onSignInClick,
            variant = com.tchat.mobile.components.TchatButtonVariant.Primary,
            enabled = !isLoading && email.isNotBlank() && password.isNotBlank(),
            modifier = Modifier.fillMaxWidth(),
            loading = isLoading
        )
    }
}

@Composable
private fun PhoneAuthForm(
    phoneNumber: String,
    onPhoneNumberChange: (String) -> Unit,
    isLoading: Boolean,
    onSendOtpClick: () -> Unit
) {
    Column {
        // Country Code Chips
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
        ) {
            CountryCodeChip(
                flag = "🇹🇭",
                code = "+66",
                onClick = { onPhoneNumberChange("+66 ") }
            )
            CountryCodeChip(
                flag = "🇮🇩",
                code = "+62",
                onClick = { onPhoneNumberChange("+62 ") }
            )
            CountryCodeChip(
                flag = "🇵🇭",
                code = "+63",
                onClick = { onPhoneNumberChange("+63 ") }
            )
            CountryCodeChip(
                flag = "🇻🇳",
                code = "+84",
                onClick = { onPhoneNumberChange("+84 ") }
            )
        }

        Spacer(modifier = Modifier.height(TchatSpacing.md))

        TchatInput(
            value = phoneNumber,
            onValueChange = onPhoneNumberChange,
            placeholder = "+66 XX XXX XXXX",
            leadingIcon = Icons.Default.Phone,
            type = TchatInputType.Number,
            enabled = !isLoading
        )

        Text(
            text = "We'll send you a 6-digit OTP via SMS",
            fontSize = 12.sp,
            color = TchatColors.onSurfaceVariant,
            modifier = Modifier.padding(top = TchatSpacing.sm)
        )

        Spacer(modifier = Modifier.height(TchatSpacing.lg))

        TchatButton(
            text = if (isLoading) "Sending..." else "Send OTP",
            onClick = onSendOtpClick,
            variant = com.tchat.mobile.components.TchatButtonVariant.Primary,
            enabled = !isLoading && phoneNumber.isNotBlank(),
            modifier = Modifier.fillMaxWidth(),
            loading = isLoading
        )
    }
}

@Composable
private fun CountryCodeChip(
    flag: String,
    code: String,
    onClick: () -> Unit
) {
    OutlinedCard(
        onClick = onClick,
        modifier = Modifier.padding(2.dp)
    ) {
        Row(
            verticalAlignment = Alignment.CenterVertically,
            modifier = Modifier.padding(horizontal = TchatSpacing.sm, vertical = TchatSpacing.xs)
        ) {
            Text(text = flag, fontSize = 16.sp)
            Spacer(modifier = Modifier.width(4.dp))
            Text(
                text = code,
                fontSize = 12.sp,
                fontWeight = FontWeight.Medium,
                color = TchatColors.onSurface
            )
        }
    }
}

@Composable
private fun VerificationScreen(
    authMethod: AuthMethod,
    phoneNumber: String,
    email: String,
    otpCode: String,
    onOtpCodeChange: (String) -> Unit,
    isLoading: Boolean,
    onVerifyClick: () -> Unit,
    onBackClick: () -> Unit,
    authState: AuthState,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxSize()
            .padding(TchatSpacing.lg),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        TchatCard(
            variant = com.tchat.mobile.components.TchatCardVariant.Elevated,
            modifier = Modifier.fillMaxWidth()
        ) {
            Column(
                modifier = Modifier.padding(TchatSpacing.lg),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                // Icon
                Box(
                    modifier = Modifier
                        .size(64.dp)
                        .clip(CircleShape),
                    contentAlignment = Alignment.Center
                ) {
                    Card(
                        modifier = Modifier.fillMaxSize(),
                        colors = CardDefaults.cardColors(
                            containerColor = TchatColors.primary
                        )
                    ) {
                        Box(
                            modifier = Modifier.fillMaxSize(),
                            contentAlignment = Alignment.Center
                        ) {
                            Icon(
                                imageVector = Icons.Default.Message,
                                contentDescription = "Verification",
                                tint = Color.White,
                                modifier = Modifier.size(32.dp)
                            )
                        }
                    }
                }

                Spacer(modifier = Modifier.height(TchatSpacing.lg))

                // Title
                Text(
                    text = "Verify Your ${if (authMethod == AuthMethod.PHONE) "Phone" else "Email"}",
                    fontSize = 20.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = TchatColors.onSurface,
                    textAlign = TextAlign.Center
                )

                // Description
                Text(
                    text = if (authMethod == AuthMethod.PHONE) {
                        "We sent a code to $phoneNumber"
                    } else {
                        "Check your email at $email"
                    },
                    fontSize = 14.sp,
                    color = TchatColors.onSurfaceVariant,
                    textAlign = TextAlign.Center,
                    modifier = Modifier.padding(top = TchatSpacing.sm)
                )

                Spacer(modifier = Modifier.height(TchatSpacing.xl))

                // OTP Input
                TchatInput(
                    value = otpCode,
                    onValueChange = { if (it.length <= 6) onOtpCodeChange(it) },
                    placeholder = "Enter verification code",
                    type = TchatInputType.Number,
                    enabled = !isLoading,
                    modifier = Modifier.fillMaxWidth()
                )

                Spacer(modifier = Modifier.height(TchatSpacing.lg))

                // Verify Button
                TchatButton(
                    text = if (isLoading) "Verifying..." else "Verify & Continue",
                    onClick = onVerifyClick,
                    variant = com.tchat.mobile.components.TchatButtonVariant.Primary,
                    enabled = !isLoading && otpCode.length >= 4,
                    modifier = Modifier.fillMaxWidth(),
                    loading = isLoading
                )

                Spacer(modifier = Modifier.height(TchatSpacing.md))

                // Back Button
                TchatButton(
                    text = "Back to ${if (authMethod == AuthMethod.PHONE) "Phone" else "Email"}",
                    onClick = onBackClick,
                    variant = com.tchat.mobile.components.TchatButtonVariant.Ghost,
                    enabled = !isLoading,
                    modifier = Modifier.fillMaxWidth()
                )

                // Error Display
                when (val state = authState) {
                    is AuthState.Error -> {
                        Spacer(modifier = Modifier.height(TchatSpacing.md))
                        Card(
                            colors = CardDefaults.cardColors(
                                containerColor = TchatColors.error.copy(alpha = 0.1f)
                            ),
                            modifier = Modifier.fillMaxWidth()
                        ) {
                            Text(
                                text = state.message,
                                color = TchatColors.error,
                                fontSize = 14.sp,
                                modifier = Modifier.padding(TchatSpacing.md),
                                textAlign = TextAlign.Center
                            )
                        }
                    }
                    else -> {} // Do nothing for other states
                }
            }
        }
    }
}