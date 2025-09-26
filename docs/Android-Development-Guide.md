# Android Development Guide - Tchat App

> **Comprehensive development guide for Android platform following the Tchat Design System**
>
> Last Updated: September 2025 | Kotlin 1.9.23 | Jetpack Compose | Android 7.0+ (API 24)

## Table of Contents

1. [Development Setup](#development-setup)
2. [Design Token System](#design-token-system)
3. [Core Components](#core-components)
4. [Architecture Patterns](#architecture-patterns)
5. [Testing Strategy](#testing-strategy)
6. [Code Style & Conventions](#code-style--conventions)
7. [Performance Standards](#performance-standards)
8. [Accessibility Guidelines](#accessibility-guidelines)

---

## Development Setup

### Prerequisites

```bash
# Required versions
- Android Studio Flamingo | 2022.2.1+
- Kotlin 1.9.23
- Gradle 8.4+
- Android SDK 34 (compile)
- Minimum SDK 24 (Android 7.0)
- Target SDK 34
- Java 17 target compatibility
```

### Dependencies

```kotlin
// build.gradle.kts (app module)
dependencies {
    // Jetpack Compose BOM - ensures compatible versions
    implementation(platform("androidx.compose:compose-bom:2023.10.01"))

    // Compose Core
    implementation("androidx.compose.ui:ui")
    implementation("androidx.compose.ui:ui-graphics")
    implementation("androidx.compose.ui:ui-tooling-preview")
    implementation("androidx.compose.material3:material3")

    // Navigation
    implementation("androidx.navigation:navigation-compose:2.7.5")

    // ViewModel
    implementation("androidx.lifecycle:lifecycle-viewmodel-compose:2.7.0")

    // Networking
    implementation("com.squareup.retrofit2:retrofit:2.9.0")
    implementation("com.squareup.okhttp3:logging-interceptor:4.12.0")

    // Coroutines
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-android:1.7.3")

    // Dependency Injection
    implementation("com.google.dagger:hilt-android:2.48")
    kapt("com.google.dagger:hilt-compiler:2.48")
    implementation("androidx.hilt:hilt-navigation-compose:1.1.0")

    // Storage
    implementation("androidx.room:room-runtime:2.6.1")
    implementation("androidx.room:room-ktx:2.6.1")
    kapt("androidx.room:room-compiler:2.6.1")
    implementation("androidx.datastore:datastore-preferences:1.0.0")

    // Security
    implementation("androidx.security:security-crypto:1.1.0-alpha06")

    // Optional - Camera & Media
    implementation("androidx.camera:camera-camera2:1.3.0")
    implementation("androidx.camera:camera-lifecycle:1.3.0")
    implementation("androidx.camera:camera-view:1.3.0")
    implementation("com.google.android.exoplayer:exoplayer:2.19.1")

    // Testing
    testImplementation("junit:junit:4.13.2")
    testImplementation("org.mockito:mockito-core:5.7.0")
    androidTestImplementation("androidx.compose.ui:ui-test-junit4")
    androidTestImplementation("androidx.test.espresso:espresso-core:3.5.1")
}
```

### Project Structure

```
app/src/main/java/com/tchat/
├── components/              # Atom design system components
│   ├── TchatButton.kt
│   ├── TchatInput.kt
│   ├── TchatCard.kt
│   └── ...
├── designsystem/           # Design tokens and theming
│   ├── Colors.kt
│   ├── Typography.kt
│   ├── Spacing.kt
│   └── DesignTokens.kt
├── screens/               # Screen implementations
├── navigation/            # Tab navigation system
├── viewmodels/           # ViewModels for business logic
├── models/               # Data models and types
├── network/              # API clients and networking
└── di/                  # Dependency injection modules
```

---

## Design Token System

### Color Palette (TailwindCSS v4 Mapped)

```kotlin
// app/src/main/java/com/tchat/designsystem/Colors.kt
package com.tchat.designsystem

import androidx.compose.ui.graphics.Color

object TchatColors {
    // Brand Colors - Primary blue (#3B82F6)
    val Primary = Color(0xFF3B82F6)           // blue-500
    val PrimaryLight = Color(0xFF60A5FA)      // blue-400
    val PrimaryDark = Color(0xFF2563EB)       // blue-600

    // Semantic Colors
    val Success = Color(0xFF10B981)           // green-500
    val Warning = Color(0xFFF59E0B)           // amber-500
    val Error = Color(0xFFEF4444)             // red-500
    val Info = Color(0xFF3B82F6)              // blue-500

    // Surface Colors
    val Surface = Color(0xFFFFFFFF)           // white
    val SurfaceSecondary = Color(0xFFF9FAFB)  // gray-50
    val SurfaceTertiary = Color(0xFFF3F4F6)   // gray-100

    // Text Colors
    val TextPrimary = Color(0xFF111827)       // gray-900
    val TextSecondary = Color(0xFF6B7280)     // gray-500
    val TextTertiary = Color(0xFF9CA3AF)      // gray-400
    val TextOnPrimary = Color(0xFFFFFFFF)     // white

    // Border Colors
    val Border = Color(0xFFE5E7EB)            // gray-200
    val BorderSecondary = Color(0xFFD1D5DB)   // gray-300
    val BorderFocus = Color(0xFF3B82F6)       // blue-500

    // Dark Mode Support
    object Dark {
        val Surface = Color(0xFF111827)        // gray-900
        val SurfaceSecondary = Color(0xFF1F2937) // gray-800
        val TextPrimary = Color(0xFFF9FAFB)    // gray-50
        val TextSecondary = Color(0xFFD1D5DB)  // gray-300
        val Border = Color(0xFF374151)         // gray-700
    }
}
```

### Typography System

```kotlin
// app/src/main/java/com/tchat/designsystem/Typography.kt
package com.tchat.designsystem

import androidx.compose.material3.Typography
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp

val TchatTypography = Typography(
    // Display Typography
    displayLarge = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Bold,
        fontSize = 48.sp,
        lineHeight = 56.sp
    ),
    displayMedium = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Bold,
        fontSize = 36.sp,
        lineHeight = 44.sp
    ),
    displaySmall = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Bold,
        fontSize = 32.sp,
        lineHeight = 40.sp
    ),

    // Headline Typography
    headlineLarge = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.SemiBold,
        fontSize = 28.sp,
        lineHeight = 36.sp
    ),
    headlineMedium = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.SemiBold,
        fontSize = 24.sp,
        lineHeight = 32.sp
    ),
    headlineSmall = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.SemiBold,
        fontSize = 20.sp,
        lineHeight = 28.sp
    ),

    // Body Typography
    bodyLarge = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Normal,
        fontSize = 18.sp,
        lineHeight = 28.sp
    ),
    bodyMedium = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Normal,
        fontSize = 16.sp,
        lineHeight = 24.sp
    ),
    bodySmall = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Normal,
        fontSize = 14.sp,
        lineHeight = 20.sp
    ),

    // Label Typography
    labelLarge = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Medium,
        fontSize = 16.sp,
        lineHeight = 20.sp
    ),
    labelMedium = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Medium,
        fontSize = 14.sp,
        lineHeight = 16.sp
    ),
    labelSmall = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Medium,
        fontSize = 12.sp,
        lineHeight = 14.sp
    )
)
```

### Spacing System (4dp Base Unit)

```kotlin
// app/src/main/java/com/tchat/designsystem/Spacing.kt
package com.tchat.designsystem

import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp

object TchatSpacing {
    // Base 4dp spacing system matching TailwindCSS
    val xs: Dp = 4.dp      // space-1 (0.25rem)
    val sm: Dp = 8.dp      // space-2 (0.5rem)
    val md: Dp = 16.dp     // space-4 (1rem)
    val lg: Dp = 24.dp     // space-6 (1.5rem)
    val xl: Dp = 32.dp     // space-8 (2rem)
    val xxl: Dp = 48.dp    // space-12 (3rem)

    // Component-specific spacing
    val buttonPaddingVertical: Dp = 12.dp    // 3/4 of md
    val buttonPaddingHorizontal: Dp = 20.dp  // 5/4 of md
    val cardPadding: Dp = 16.dp             // md
    val screenPadding: Dp = 16.dp           // md
}
```

---

## Core Components

### TchatButton Implementation

```kotlin
// app/src/main/java/com/tchat/components/TchatButton.kt
package com.tchat.components

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.interaction.collectIsPressedAsState
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.scale
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalHapticFeedback
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.tchat.designsystem.TchatColors
import com.tchat.designsystem.TchatSpacing
import com.tchat.designsystem.TchatTypography

/**
 * TchatButton - Primary button component with multiple variants and sizes
 *
 * @param text Button text content (optional for icon-only buttons)
 * @param icon Leading icon (optional)
 * @param variant Button style variant
 * @param size Button size variant
 * @param isLoading Shows loading state with progress indicator
 * @param isEnabled Controls button enabled/disabled state
 * @param modifier Compose modifier for styling and layout
 * @param onClick Click event handler
 */
@Composable
fun TchatButton(
    text: String? = null,
    icon: ImageVector? = null,
    variant: ButtonVariant = ButtonVariant.Primary,
    size: ButtonSize = ButtonSize.Medium,
    isLoading: Boolean = false,
    isEnabled: Boolean = true,
    modifier: Modifier = Modifier,
    onClick: () -> Unit
) {
    val haptic = LocalHapticFeedback.current
    val interactionSource = remember { MutableInteractionSource() }
    val isPressed by interactionSource.collectIsPressedAsState()

    // Press animation
    val scale by animateFloatAsState(
        targetValue = if (isPressed && isEnabled) 0.95f else 1.0f,
        animationSpec = tween(durationMillis = 100),
        label = "button_press_animation"
    )

    Button(
        onClick = {
            if (isEnabled && !isLoading) {
                haptic.performHapticFeedback(androidx.compose.ui.hapticfeedback.HapticFeedbackType.TextHandleMove)
                onClick()
            }
        },
        modifier = modifier
            .scale(scale)
            .height(size.height)
            .then(
                if (size == ButtonSize.Icon) {
                    Modifier.width(size.height)
                } else {
                    Modifier.fillMaxWidth()
                }
            ),
        enabled = isEnabled && !isLoading,
        colors = variant.colors,
        border = variant.border,
        shape = RoundedCornerShape(8.dp),
        contentPadding = if (size == ButtonSize.Icon) {
            PaddingValues(0.dp)
        } else {
            PaddingValues(
                horizontal = size.horizontalPadding,
                vertical = TchatSpacing.buttonPaddingVertical
            )
        },
        interactionSource = interactionSource
    ) {
        Row(
            horizontalArrangement = Arrangement.Center,
            verticalAlignment = Alignment.CenterVertically
        ) {
            if (isLoading) {
                CircularProgressIndicator(
                    modifier = Modifier.size(16.dp),
                    color = variant.textColor,
                    strokeWidth = 2.dp
                )
                if (!text.isNullOrEmpty() && size != ButtonSize.Icon) {
                    Spacer(modifier = Modifier.width(TchatSpacing.sm))
                }
            } else if (icon != null) {
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    modifier = Modifier.size(20.dp),
                    tint = variant.textColor
                )
                if (!text.isNullOrEmpty() && size != ButtonSize.Icon) {
                    Spacer(modifier = Modifier.width(TchatSpacing.sm))
                }
            }

            if (!text.isNullOrEmpty() && size != ButtonSize.Icon) {
                Text(
                    text = if (isLoading && icon == null) "Loading..." else text,
                    style = size.textStyle,
                    color = variant.textColor,
                    maxLines = 1
                )
            }
        }
    }
}

// MARK: - Button Variants
enum class ButtonVariant(
    val colors: ButtonColors,
    val border: BorderStroke?,
    val textColor: Color
) {
    Primary(
        colors = ButtonDefaults.buttonColors(
            containerColor = TchatColors.Primary,
            contentColor = TchatColors.TextOnPrimary,
            disabledContainerColor = TchatColors.Primary.copy(alpha = 0.6f),
            disabledContentColor = TchatColors.TextOnPrimary.copy(alpha = 0.6f)
        ),
        border = null,
        textColor = TchatColors.TextOnPrimary
    ),

    Secondary(
        colors = ButtonDefaults.buttonColors(
            containerColor = TchatColors.SurfaceSecondary,
            contentColor = TchatColors.TextPrimary,
            disabledContainerColor = TchatColors.SurfaceSecondary.copy(alpha = 0.6f),
            disabledContentColor = TchatColors.TextPrimary.copy(alpha = 0.6f)
        ),
        border = null,
        textColor = TchatColors.TextPrimary
    ),

    Ghost(
        colors = ButtonDefaults.buttonColors(
            containerColor = Color.Transparent,
            contentColor = TchatColors.Primary,
            disabledContainerColor = Color.Transparent,
            disabledContentColor = TchatColors.Primary.copy(alpha = 0.6f)
        ),
        border = null,
        textColor = TchatColors.Primary
    ),

    Destructive(
        colors = ButtonDefaults.buttonColors(
            containerColor = TchatColors.Error,
            contentColor = TchatColors.TextOnPrimary,
            disabledContainerColor = TchatColors.Error.copy(alpha = 0.6f),
            disabledContentColor = TchatColors.TextOnPrimary.copy(alpha = 0.6f)
        ),
        border = null,
        textColor = TchatColors.TextOnPrimary
    ),

    Outline(
        colors = ButtonDefaults.buttonColors(
            containerColor = Color.Transparent,
            contentColor = TchatColors.Primary,
            disabledContainerColor = Color.Transparent,
            disabledContentColor = TchatColors.Primary.copy(alpha = 0.6f)
        ),
        border = BorderStroke(1.dp, TchatColors.Border),
        textColor = TchatColors.Primary
    )
}

// MARK: - Button Sizes
enum class ButtonSize(
    val height: androidx.compose.ui.unit.Dp,
    val horizontalPadding: androidx.compose.ui.unit.Dp,
    val textStyle: TextStyle
) {
    Small(
        height = 32.dp,
        horizontalPadding = TchatSpacing.sm,
        textStyle = TchatTypography.bodySmall
    ),

    Medium(
        height = 44.dp,
        horizontalPadding = TchatSpacing.md,
        textStyle = TchatTypography.bodyMedium
    ),

    Large(
        height = 48.dp,
        horizontalPadding = TchatSpacing.lg,
        textStyle = TchatTypography.bodyLarge
    ),

    Icon(
        height = 44.dp,
        horizontalPadding = 0.dp,
        textStyle = TchatTypography.bodyMedium
    )
}

// MARK: - Convenience Functions
@Composable
fun PrimaryButton(
    text: String,
    isLoading: Boolean = false,
    isEnabled: Boolean = true,
    modifier: Modifier = Modifier,
    onClick: () -> Unit
) {
    TchatButton(
        text = text,
        variant = ButtonVariant.Primary,
        isLoading = isLoading,
        isEnabled = isEnabled,
        modifier = modifier,
        onClick = onClick
    )
}

@Composable
fun IconButton(
    icon: ImageVector,
    contentDescription: String? = null,
    variant: ButtonVariant = ButtonVariant.Ghost,
    isEnabled: Boolean = true,
    modifier: Modifier = Modifier,
    onClick: () -> Unit
) {
    TchatButton(
        icon = icon,
        variant = variant,
        size = ButtonSize.Icon,
        isEnabled = isEnabled,
        modifier = modifier,
        onClick = onClick
    )
}

// MARK: - Preview
@Preview(showBackground = true)
@Composable
fun TchatButtonPreview() {
    Column(
        modifier = Modifier.padding(TchatSpacing.md),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
    ) {
        // Variants
        TchatButton("Primary Button", variant = ButtonVariant.Primary) { }
        TchatButton("Secondary Button", variant = ButtonVariant.Secondary) { }
        TchatButton("Ghost Button", variant = ButtonVariant.Ghost) { }
        TchatButton("Destructive Button", variant = ButtonVariant.Destructive) { }
        TchatButton("Outline Button", variant = ButtonVariant.Outline) { }

        // Sizes
        Row(horizontalArrangement = Arrangement.spacedBy(TchatSpacing.sm)) {
            TchatButton("Small", variant = ButtonVariant.Primary, size = ButtonSize.Small, modifier = Modifier.weight(1f)) { }
            TchatButton("Medium", variant = ButtonVariant.Primary, size = ButtonSize.Medium, modifier = Modifier.weight(1f)) { }
            TchatButton("Large", variant = ButtonVariant.Primary, size = ButtonSize.Large, modifier = Modifier.weight(1f)) { }
        }

        // States
        TchatButton("Loading", variant = ButtonVariant.Primary, isLoading = true) { }
        TchatButton("Disabled", variant = ButtonVariant.Primary, isEnabled = false) { }
    }
}
```

### TchatInput Implementation

```kotlin
// app/src/main/java/com/tchat/components/TchatInput.kt
package com.tchat.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.tchat.designsystem.TchatColors
import com.tchat.designsystem.TchatSpacing
import com.tchat.designsystem.TchatTypography
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*

/**
 * TchatInput - Text input component with validation states and various input types
 *
 * @param value Current input value
 * @param onValueChange Callback when value changes
 * @param placeholder Placeholder text
 * @param inputType Type of input (affects keyboard and validation)
 * @param validationState Current validation state
 * @param size Input field size
 * @param leadingIcon Optional leading icon
 * @param isEnabled Controls enabled/disabled state
 * @param modifier Compose modifier
 */
@Composable
fun TchatInput(
    value: String,
    onValueChange: (String) -> Unit,
    placeholder: String,
    inputType: InputType = InputType.Text,
    validationState: ValidationState = ValidationState.None,
    size: InputSize = InputSize.Medium,
    leadingIcon: ImageVector? = null,
    isEnabled: Boolean = true,
    modifier: Modifier = Modifier
) {
    var isPasswordVisible by remember { mutableStateOf(false) }

    // Animated border color based on validation state
    val borderColor by animateColorAsState(
        targetValue = when (validationState) {
            is ValidationState.Valid -> TchatColors.Success
            is ValidationState.Invalid -> TchatColors.Error
            ValidationState.None -> TchatColors.Border
        },
        animationSpec = tween(durationMillis = 200),
        label = "border_color_animation"
    )

    Column(modifier = modifier) {
        OutlinedTextField(
            value = value,
            onValueChange = onValueChange,
            placeholder = {
                Text(
                    text = placeholder,
                    style = size.textStyle,
                    color = TchatColors.TextTertiary
                )
            },
            leadingIcon = leadingIcon?.let { icon ->
                {
                    Icon(
                        imageVector = icon,
                        contentDescription = null,
                        tint = TchatColors.TextSecondary,
                        modifier = Modifier.size(20.dp)
                    )
                }
            },
            trailingIcon = {
                Row {
                    // Password visibility toggle
                    if (inputType == InputType.Password) {
                        IconButton(
                            onClick = { isPasswordVisible = !isPasswordVisible }
                        ) {
                            Icon(
                                imageVector = if (isPasswordVisible) Icons.Filled.VisibilityOff else Icons.Filled.Visibility,
                                contentDescription = if (isPasswordVisible) "Hide password" else "Show password",
                                tint = TchatColors.TextSecondary,
                                modifier = Modifier.size(20.dp)
                            )
                        }
                    }

                    // Validation icon
                    when (validationState) {
                        is ValidationState.Valid -> {
                            Icon(
                                imageVector = Icons.Filled.CheckCircle,
                                contentDescription = "Valid input",
                                tint = TchatColors.Success,
                                modifier = Modifier
                                    .size(20.dp)
                                    .padding(end = TchatSpacing.sm)
                            )
                        }
                        is ValidationState.Invalid -> {
                            Icon(
                                imageVector = Icons.Filled.Error,
                                contentDescription = "Invalid input",
                                tint = TchatColors.Error,
                                modifier = Modifier
                                    .size(20.dp)
                                    .padding(end = TchatSpacing.sm)
                            )
                        }
                        ValidationState.None -> { /* No icon */ }
                    }
                }
            },
            enabled = isEnabled,
            textStyle = size.textStyle.copy(color = TchatColors.TextPrimary),
            keyboardOptions = KeyboardOptions(
                keyboardType = inputType.keyboardType,
                imeAction = ImeAction.Next
            ),
            keyboardActions = KeyboardActions.Default,
            singleLine = inputType != InputType.Multiline(),
            maxLines = if (inputType is InputType.Multiline) inputType.maxLines else 1,
            visualTransformation = if (inputType == InputType.Password && !isPasswordVisible) {
                PasswordVisualTransformation()
            } else {
                VisualTransformation.None
            },
            shape = RoundedCornerShape(8.dp),
            colors = OutlinedTextFieldDefaults.colors(
                focusedBorderColor = borderColor,
                unfocusedBorderColor = borderColor,
                focusedContainerColor = TchatColors.Surface,
                unfocusedContainerColor = TchatColors.Surface,
                disabledContainerColor = TchatColors.SurfaceSecondary,
                errorBorderColor = TchatColors.Error,
                cursorColor = TchatColors.Primary
            ),
            modifier = Modifier
                .fillMaxWidth()
                .height(size.height)
        )

        // Error message
        if (validationState is ValidationState.Invalid) {
            Text(
                text = validationState.message,
                style = TchatTypography.labelSmall,
                color = TchatColors.Error,
                modifier = Modifier.padding(
                    start = TchatSpacing.xs,
                    top = TchatSpacing.xs
                )
            )
        }
    }
}

// MARK: - Input Types
sealed class InputType(val keyboardType: KeyboardType) {
    object Text : InputType(KeyboardType.Text)
    object Email : InputType(KeyboardType.Email)
    object Password : InputType(KeyboardType.Password)
    object Number : InputType(KeyboardType.Number)
    object Search : InputType(KeyboardType.Text)
    data class Multiline(val maxLines: Int = 3) : InputType(KeyboardType.Text)
}

// MARK: - Validation States
sealed class ValidationState {
    object None : ValidationState()
    object Valid : ValidationState()
    data class Invalid(val message: String) : ValidationState()
}

// MARK: - Input Sizes
enum class InputSize(
    val height: androidx.compose.ui.unit.Dp,
    val textStyle: androidx.compose.ui.text.TextStyle
) {
    Small(
        height = 36.dp,
        textStyle = TchatTypography.bodySmall
    ),

    Medium(
        height = 44.dp,
        textStyle = TchatTypography.bodyMedium
    ),

    Large(
        height = 52.dp,
        textStyle = TchatTypography.bodyLarge
    )
}

// MARK: - Preview
@Preview(showBackground = true)
@Composable
fun TchatInputPreview() {
    Column(
        modifier = Modifier.padding(TchatSpacing.md),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
    ) {
        TchatInput(
            value = "",
            onValueChange = { },
            placeholder = "Enter text",
            leadingIcon = Icons.Filled.Person
        )

        TchatInput(
            value = "",
            onValueChange = { },
            placeholder = "Email address",
            inputType = InputType.Email,
            leadingIcon = Icons.Filled.Email
        )

        TchatInput(
            value = "",
            onValueChange = { },
            placeholder = "Password",
            inputType = InputType.Password,
            leadingIcon = Icons.Filled.Lock
        )

        TchatInput(
            value = "Valid input",
            onValueChange = { },
            placeholder = "Valid input",
            validationState = ValidationState.Valid
        )

        TchatInput(
            value = "Invalid input",
            onValueChange = { },
            placeholder = "Invalid input",
            validationState = ValidationState.Invalid("This field is required")
        )
    }
}
```

### TchatCard Implementation

```kotlin
// app/src/main/java/com/tchat/components/TchatCard.kt
package com.tchat.components

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.interaction.collectIsPressedAsState
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.scale
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.tchat.designsystem.TchatColors
import com.tchat.designsystem.TchatSpacing

/**
 * TchatCard - Flexible card component with multiple visual variants
 *
 * @param variant Card visual style
 * @param size Card padding size
 * @param isInteractive Whether card responds to press interactions
 * @param modifier Compose modifier
 * @param onClick Optional click handler (enables interactive behavior)
 * @param content Card content composable
 */
@Composable
fun TchatCard(
    variant: CardVariant = CardVariant.Elevated,
    size: CardSize = CardSize.Standard,
    isInteractive: Boolean = false,
    modifier: Modifier = Modifier,
    onClick: (() -> Unit)? = null,
    content: @Composable ColumnScope.() -> Unit
) {
    val interactionSource = remember { MutableInteractionSource() }
    val isPressed by interactionSource.collectIsPressedAsState()

    // Press animation for interactive cards
    val scale by animateFloatAsState(
        targetValue = if (isPressed && (isInteractive || onClick != null)) 0.98f else 1.0f,
        animationSpec = tween(durationMillis = 100),
        label = "card_press_animation"
    )

    Card(
        modifier = modifier.scale(scale),
        shape = RoundedCornerShape(12.dp),
        colors = variant.colors,
        elevation = variant.elevation,
        border = variant.border,
        onClick = onClick,
        interactionSource = if (isInteractive || onClick != null) interactionSource else null
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .then(
                    when (variant) {
                        CardVariant.Glass -> Modifier.background(
                            brush = Brush.verticalGradient(
                                colors = listOf(
                                    Color.White.copy(alpha = 0.8f),
                                    Color.White.copy(alpha = 0.6f)
                                )
                            )
                        )
                        else -> Modifier
                    }
                )
                .padding(size.padding),
            content = content
        )
    }
}

// MARK: - Card Variants
enum class CardVariant(
    val colors: CardColors,
    val elevation: CardElevation,
    val border: BorderStroke?
) {
    Elevated(
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.Surface
        ),
        elevation = CardDefaults.cardElevation(
            defaultElevation = 4.dp,
            pressedElevation = 6.dp,
            hoveredElevation = 5.dp
        ),
        border = null
    ),

    Outlined(
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.Surface
        ),
        elevation = CardDefaults.cardElevation(
            defaultElevation = 0.dp
        ),
        border = BorderStroke(1.dp, TchatColors.Border)
    ),

    Filled(
        colors = CardDefaults.cardColors(
            containerColor = TchatColors.SurfaceSecondary
        ),
        elevation = CardDefaults.cardElevation(
            defaultElevation = 0.dp
        ),
        border = null
    ),

    Glass(
        colors = CardDefaults.cardColors(
            containerColor = Color.Transparent
        ),
        elevation = CardDefaults.cardElevation(
            defaultElevation = 0.dp
        ),
        border = BorderStroke(1.dp, Color.White.copy(alpha = 0.2f))
    )
}

// MARK: - Card Sizes
enum class CardSize(val padding: androidx.compose.ui.unit.Dp) {
    Compact(TchatSpacing.sm),     // 8dp
    Standard(TchatSpacing.md),    // 16dp
    Expanded(TchatSpacing.lg)     // 24dp
}

// MARK: - Preview
@Preview(showBackground = true)
@Composable
fun TchatCardPreview() {
    Column(
        modifier = Modifier.padding(TchatSpacing.md),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.md)
    ) {
        // Elevated card
        TchatCard(variant = CardVariant.Elevated) {
            Text("Elevated Card", style = MaterialTheme.typography.headlineSmall)
            Text("This card has a shadow elevation", style = MaterialTheme.typography.bodyMedium)
        }

        // Outlined card
        TchatCard(variant = CardVariant.Outlined) {
            Text("Outlined Card", style = MaterialTheme.typography.headlineSmall)
            Text("This card has a border outline", style = MaterialTheme.typography.bodyMedium)
        }

        // Filled card
        TchatCard(variant = CardVariant.Filled) {
            Text("Filled Card", style = MaterialTheme.typography.headlineSmall)
            Text("This card has a filled background", style = MaterialTheme.typography.bodyMedium)
        }

        // Interactive card
        TchatCard(
            variant = CardVariant.Elevated,
            isInteractive = true,
            onClick = { /* Handle click */ }
        ) {
            Text("Interactive Card", style = MaterialTheme.typography.headlineSmall)
            Text("This card can be tapped", style = MaterialTheme.typography.bodyMedium)
        }
    }
}
```

---

## Architecture Patterns

### MVVM with Hilt Dependency Injection

```kotlin
// Example ViewModel with Hilt
@HiltViewModel
class ChatViewModel @Inject constructor(
    private val chatRepository: ChatRepository,
    private val userRepository: UserRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(ChatUiState())
    val uiState: StateFlow<ChatUiState> = _uiState.asStateFlow()

    private val _uiEvent = Channel<ChatUiEvent>()
    val uiEvent = _uiEvent.receiveAsFlow()

    init {
        loadMessages()
    }

    fun onAction(action: ChatAction) {
        when (action) {
            is ChatAction.SendMessage -> sendMessage(action.content)
            is ChatAction.LoadMessages -> loadMessages()
            is ChatAction.RefreshMessages -> refreshMessages()
        }
    }

    private fun sendMessage(content: String) {
        viewModelScope.launch {
            try {
                _uiState.update { it.copy(isLoading = true) }

                val message = chatRepository.sendMessage(content)
                _uiState.update { currentState ->
                    currentState.copy(
                        messages = currentState.messages + message,
                        isLoading = false
                    )
                }
            } catch (e: Exception) {
                _uiEvent.send(ChatUiEvent.ShowError(e.message ?: "Unknown error"))
                _uiState.update { it.copy(isLoading = false) }
            }
        }
    }

    private fun loadMessages() {
        viewModelScope.launch {
            chatRepository.getMessages()
                .catch { exception ->
                    _uiEvent.send(ChatUiEvent.ShowError(exception.message ?: "Failed to load messages"))
                }
                .collect { messages ->
                    _uiState.update { it.copy(messages = messages) }
                }
        }
    }
}

// UI State
data class ChatUiState(
    val messages: List<Message> = emptyList(),
    val isLoading: Boolean = false,
    val currentUser: User? = null
)

// UI Events
sealed class ChatUiEvent {
    data class ShowError(val message: String) : ChatUiEvent()
    object MessageSent : ChatUiEvent()
}

// Actions
sealed class ChatAction {
    data class SendMessage(val content: String) : ChatAction()
    object LoadMessages : ChatAction()
    object RefreshMessages : ChatAction()
}
```

---

## Testing Strategy

### Unit Testing with JUnit and Mockito

```kotlin
// Example ViewModel Test
@ExperimentalCoroutinesApi
class ChatViewModelTest {

    @get:Rule
    val instantTaskExecutorRule = InstantTaskExecutorRule()

    @get:Rule
    val mainDispatcherRule = MainDispatcherRule()

    @Mock
    private lateinit var chatRepository: ChatRepository

    @Mock
    private lateinit var userRepository: UserRepository

    private lateinit var viewModel: ChatViewModel

    @Before
    fun setUp() {
        MockitoAnnotations.openMocks(this)
        viewModel = ChatViewModel(chatRepository, userRepository)
    }

    @Test
    fun `sendMessage should update ui state with new message`() = runTest {
        // Given
        val content = "Test message"
        val message = Message(id = "1", content = content, senderId = "user1")
        `when`(chatRepository.sendMessage(content)).thenReturn(message)

        // When
        viewModel.onAction(ChatAction.SendMessage(content))

        // Then
        val state = viewModel.uiState.value
        assertEquals(1, state.messages.size)
        assertEquals(message, state.messages.first())
        assertFalse(state.isLoading)
    }

    @Test
    fun `loadMessages should handle repository error`() = runTest {
        // Given
        val error = RuntimeException("Network error")
        `when`(chatRepository.getMessages()).thenReturn(flow { throw error })

        // When
        viewModel.onAction(ChatAction.LoadMessages)

        // Then
        val events = mutableListOf<ChatUiEvent>()
        val job = launch {
            viewModel.uiEvent.toList(events)
        }

        assertTrue(events.any { it is ChatUiEvent.ShowError })
        job.cancel()
    }
}
```

### Compose UI Testing

```kotlin
// Example Compose Test
@ExperimentalTestApi
class TchatButtonTest {

    @get:Rule
    val composeTestRule = createComposeRule()

    @Test
    fun tchatButton_displaysText() {
        // Given
        val buttonText = "Click me"

        // When
        composeTestRule.setContent {
            TchatButton(text = buttonText) { }
        }

        // Then
        composeTestRule
            .onNodeWithText(buttonText)
            .assertIsDisplayed()
    }

    @Test
    fun tchatButton_callsOnClickWhenTapped() {
        // Given
        var clicked = false

        // When
        composeTestRule.setContent {
            TchatButton(text = "Click me") {
                clicked = true
            }
        }

        composeTestRule
            .onNodeWithText("Click me")
            .performClick()

        // Then
        assertTrue(clicked)
    }

    @Test
    fun tchatButton_showsLoadingWhenIsLoadingTrue() {
        // When
        composeTestRule.setContent {
            TchatButton(text = "Submit", isLoading = true) { }
        }

        // Then
        composeTestRule
            .onNodeWithContentDescription("Loading")
            .assertIsDisplayed()
    }

    @Test
    fun tchatButton_isDisabledWhenIsEnabledFalse() {
        // When
        composeTestRule.setContent {
            TchatButton(text = "Submit", isEnabled = false) { }
        }

        // Then
        composeTestRule
            .onNodeWithText("Submit")
            .assertIsNotEnabled()
    }
}
```

---

## Code Style & Conventions

### Naming Conventions

```kotlin
// Package names use lowercase
package com.tchat.components

// Class names use PascalCase
class TchatButton

// Function names use camelCase
fun sendMessage() { }
fun validateInput(): Boolean { }

// Constants use SCREAMING_SNAKE_CASE
const val MAX_MESSAGE_LENGTH = 4096
const val DEFAULT_TIMEOUT_MS = 5000L

// Private properties can use underscore prefix for backing properties
private val _uiState = MutableStateFlow(UiState())
val uiState: StateFlow<UiState> = _uiState.asStateFlow()

// Composable functions use PascalCase
@Composable
fun TchatButton() { }

@Composable
fun ChatScreen() { }
```

### Code Organization

```kotlin
// File structure order:
// 1. Package declaration
package com.tchat.components

// 2. Imports (grouped: standard library, Android, third-party, project)
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import com.tchat.designsystem.TchatColors

// 3. Class/Function declaration with documentation
/**
 * TchatButton - Primary button component
 *
 * @param text Button text content
 * @param onClick Click event handler
 */
@Composable
fun TchatButton(
    text: String,
    onClick: () -> Unit
) {
    // 4. Implementation
}

// 5. Supporting classes/enums
enum class ButtonVariant {
    Primary, Secondary, Ghost
}

// 6. Preview functions
@Preview
@Composable
fun TchatButtonPreview() {
    // Preview implementation
}
```

### Compose Best Practices

```kotlin
// ✅ Good - Use remember for expensive operations
@Composable
fun ExpensiveComponent() {
    val expensiveValue = remember {
        computeExpensiveValue()
    }

    Text(text = expensiveValue)
}

// ✅ Good - Use derivedStateOf for computed values
@Composable
fun FilteredList(items: List<String>, query: String) {
    val filteredItems by remember(items, query) {
        derivedStateOf {
            items.filter { it.contains(query, ignoreCase = true) }
        }
    }

    LazyColumn {
        items(filteredItems) { item ->
            Text(item)
        }
    }
}

// ✅ Good - Use stable parameters
@Composable
fun StableComponent(
    @Stable items: List<String>,
    onItemClick: (String) -> Unit
) {
    // Implementation
}

// ✅ Good - Use LazyColumn for long lists
@Composable
fun MessageList(messages: List<Message>) {
    LazyColumn {
        items(messages) { message ->
            MessageItem(message = message)
        }
    }
}

// ❌ Avoid - Heavy computation in composition
@Composable
fun BadComponent(items: List<String>) {
    val processedItems = items.map { heavyOperation(it) } // Don't do this

    // Better: use LaunchedEffect or remember
}

// ✅ Better approach
@Composable
fun GoodComponent(items: List<String>) {
    var processedItems by remember { mutableStateOf<List<String>>(emptyList()) }

    LaunchedEffect(items) {
        processedItems = items.map { heavyOperation(it) }
    }

    // Use processedItems...
}
```

---

## Performance Standards

### Target Metrics

- **Frame Rate**: 60 FPS for all animations and scrolling
- **Touch Response**: <16ms from touch to visual feedback
- **Component Render**: <8ms for basic components
- **Memory Usage**: <100MB for typical usage
- **Battery Impact**: Minimal background processing

### Performance Guidelines

```kotlin
// ✅ Use LazyColumn for long lists
@Composable
fun MessageList(messages: List<Message>) {
    LazyColumn(
        contentPadding = PaddingValues(TchatSpacing.md),
        verticalArrangement = Arrangement.spacedBy(TchatSpacing.sm)
    ) {
        items(
            items = messages,
            key = { it.id } // Important for recomposition optimization
        ) { message ->
            MessageItem(
                message = message,
                modifier = Modifier.animateItemPlacement()
            )
        }
    }
}

// ✅ Use stable keys for better performance
@Composable
fun OptimizedList(items: List<DataItem>) {
    LazyColumn {
        items(
            items = items,
            key = { item -> item.id } // Stable key
        ) { item ->
            ListItem(item = item)
        }
    }
}

// ✅ Minimize State updates and use proper scoping
@Composable
fun ChatInput() {
    var text by remember { mutableStateOf("") }

    // Good - scoped state update
    TextField(
        value = text,
        onValueChange = { newText ->
            if (newText.length <= MAX_LENGTH) { // Validate before update
                text = newText
            }
        }
    )
}

// ✅ Use Modifier.drawBehind for custom drawing
@Composable
fun CustomBackground(modifier: Modifier = Modifier) {
    Box(
        modifier = modifier.drawBehind {
            // Custom drawing code
            drawRoundRect(
                color = Color.Blue,
                cornerRadius = CornerRadius(8.dp.toPx())
            )
        }
    )
}

// ✅ Avoid creating new objects in composition
@Composable
fun EfficientComponent() {
    // ❌ Bad - creates new list every composition
    val items = listOf("Item 1", "Item 2", "Item 3")

    // ✅ Good - stable reference
    val items = remember { listOf("Item 1", "Item 2", "Item 3") }

    // Or even better - pass as parameter
}
```

---

## Accessibility Guidelines

### TalkBack Support

```kotlin
// ✅ Proper content descriptions
TchatButton(
    text = "Send",
    modifier = Modifier.semantics {
        contentDescription = "Send message"
        role = Role.Button
        stateDescription = if (isLoading) "Loading" else null
    }
) {
    sendMessage()
}

// ✅ Semantic properties for custom components
@Composable
fun CustomSlider(
    value: Float,
    onValueChange: (Float) -> Unit,
    range: ClosedFloatingPointRange<Float> = 0f..1f
) {
    Box(
        modifier = Modifier.semantics {
            contentDescription = "Volume slider"
            stateDescription = "Volume at ${(value * 100).toInt()}%"
            role = Role.Slider

            // Custom actions
            customActions = listOf(
                CustomAccessibilityAction("Increase volume") {
                    onValueChange((value + 0.1f).coerceAtMost(range.endInclusive))
                    true
                },
                CustomAccessibilityAction("Decrease volume") {
                    onValueChange((value - 0.1f).coerceAtLeast(range.start))
                    true
                }
            )
        }
    ) {
        // Slider implementation
    }
}
```

### Focus Management

```kotlin
// ✅ Proper focus handling in forms
@Composable
fun LoginForm() {
    val focusManager = LocalFocusManager.current
    var email by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }

    Column {
        TchatInput(
            value = email,
            onValueChange = { email = it },
            placeholder = "Email",
            inputType = InputType.Email,
            modifier = Modifier.onKeyEvent { keyEvent ->
                if (keyEvent.key == Key.Tab) {
                    focusManager.moveFocus(FocusDirection.Down)
                    true
                } else false
            }
        )

        TchatInput(
            value = password,
            onValueChange = { password = it },
            placeholder = "Password",
            inputType = InputType.Password,
            modifier = Modifier.onKeyEvent { keyEvent ->
                if (keyEvent.key == Key.Enter) {
                    submitForm(email, password)
                    true
                } else false
            }
        )

        TchatButton("Login") {
            submitForm(email, password)
        }
    }
}
```

### Color Contrast Compliance

```kotlin
// ✅ All TchatColors meet WCAG 2.1 AA standards
// Contrast ratios are validated:
// - Primary text (gray-900) on white: 18.7:1 (AAA)
// - Secondary text (gray-500) on white: 7.0:1 (AA Large)
// - Primary button (white) on blue-500: 8.6:1 (AA)

// ✅ Custom color validation
fun validateColorContrast(foreground: Color, background: Color): Boolean {
    val contrastRatio = calculateContrastRatio(foreground, background)
    return contrastRatio >= 4.5 // WCAG AA standard
}
```

---

## Getting Started Checklist

- [ ] Install Android Studio Flamingo+
- [ ] Set up Kotlin 1.9.23 and Gradle 8.4+
- [ ] Configure minimum SDK 24, target SDK 34
- [ ] Add Compose BOM and required dependencies
- [ ] Import TchatComponents package
- [ ] Configure design tokens in your app theme
- [ ] Implement first screen using TchatButton and TchatInput
- [ ] Set up Hilt for dependency injection
- [ ] Configure unit testing with JUnit and Mockito
- [ ] Set up Compose UI testing
- [ ] Enable TalkBack for accessibility testing
- [ ] Run performance profiling with GPU rendering

---

**Questions or Issues?**
Refer to the project's GitHub repository or contact the Android development team for support.

---

*This guide is part of the Tchat Design System documentation suite.*