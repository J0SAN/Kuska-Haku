package com.proyecto.kuskahaku.presentation.theme

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

private val DarkColorScheme = darkColorScheme(
    primary = ChichaPurpleLight,
    secondary = CocaGreen,
    tertiary = IntiGold,
    background = StoneDark,
    surface = StoneDark
)

private val LightColorScheme = lightColorScheme(
    primary = ChichaPurple,
    secondary = CocaGreen,
    tertiary = IntiGold,
    background = StoneWhite,
    surface = Color.White,
    onBackground = Color (0xFF2B2D42),
    onSurface = Color(0xFF2B2D42)

    /* Other default colors to override
    background = Color(0xFFFFFBFE),
    surface = Color(0xFFFFFBFE),
    onPrimary = Color.White,
    onSecondary = Color.White,
    onTertiary = Color.White,
    onBackground = Color(0xFF1C1B1F),
    onSurface = Color(0xFF1C1B1F),
    */
)

@Composable
fun KuskaHakuTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    // Dynamic color is available on Android 12+
    dynamicColor: Boolean = true,
    content: @Composable () -> Unit
) {
    val colorScheme = if (darkTheme) DarkColorScheme else LightColorScheme

    MaterialTheme(
        colorScheme = colorScheme,
        typography = Typography,
        content = content
    )
}