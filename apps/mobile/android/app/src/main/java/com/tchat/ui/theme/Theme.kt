package com.tchat.ui.theme

import android.app.Activity
import android.os.Build
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.dynamicDarkColorScheme
import androidx.compose.material3.dynamicLightColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.SideEffect
import androidx.compose.ui.graphics.toArgb
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.platform.LocalView
import androidx.core.view.WindowCompat
import com.tchat.designsystem.Colors
import com.tchat.designsystem.DarkColors

private val DarkColorScheme = darkColorScheme(
    primary = DarkColors.primary,
    secondary = DarkColors.secondary,
    tertiary = DarkColors.accent,
    background = DarkColors.background,
    surface = DarkColors.surface,
    onPrimary = DarkColors.textPrimary,
    onSecondary = DarkColors.textSecondary,
    onTertiary = DarkColors.textTertiary,
    onBackground = DarkColors.textPrimary,
    onSurface = DarkColors.textPrimary,
    error = DarkColors.error,
    onError = Colors.textOnPrimary
)

private val LightColorScheme = lightColorScheme(
    primary = Colors.primary,
    secondary = Colors.secondary,
    tertiary = Colors.accent,
    background = Colors.background,
    surface = Colors.surface,
    onPrimary = Colors.textOnPrimary,
    onSecondary = Colors.textSecondary,
    onTertiary = Colors.textTertiary,
    onBackground = Colors.textPrimary,
    onSurface = Colors.textPrimary,
    error = Colors.error,
    onError = Colors.textOnPrimary
)

@Composable
fun TchatTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    // Dynamic color is available on Android 12+
    dynamicColor: Boolean = true,
    content: @Composable () -> Unit
) {
    val colorScheme = when {
        dynamicColor && Build.VERSION.SDK_INT >= Build.VERSION_CODES.S -> {
            val context = LocalContext.current
            if (darkTheme) dynamicDarkColorScheme(context) else dynamicLightColorScheme(context)
        }

        darkTheme -> DarkColorScheme
        else -> LightColorScheme
    }
    val view = LocalView.current
    if (!view.isInEditMode) {
        SideEffect {
            val window = (view.context as Activity).window
            window.statusBarColor = colorScheme.primary.toArgb()
            WindowCompat.getInsetsController(window, view).isAppearanceLightStatusBars = darkTheme
        }
    }

    MaterialTheme(
        colorScheme = colorScheme,
        typography = Typography,
        content = content
    )
}