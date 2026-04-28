package com.proyecto.kuskahaku.presentation.ui.components

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color

@Composable
fun MapMockup() {
    Canvas(modifier = Modifier.fillMaxSize()) {
        val canvasWidth = size.width
        val canvasHeight = size.height

        // Dibujamos una cuadrícula gris suave que simule "calles"
        val step = 150f
        for (i in 0..(canvasWidth / step).toInt()) {
            drawLine(
                color = Color.LightGray.copy(alpha = 0.5f),
                start = Offset(x = i * step, y = 0f),
                end = Offset(x = i * step, y = canvasHeight),
                strokeWidth = 20f
            )
        }
        for (i in 0..(canvasHeight / step).toInt()) {
            drawLine(
                color = Color.LightGray.copy(alpha = 0.5f),
                start = Offset(x = 0f, y = i * step),
                end = Offset(x = canvasWidth, y = i * step),
                strokeWidth = 20f
            )
        }

        // Dibujamos un "Paradero" (Círculo morado)
        drawCircle(
            color = Color(0xFF5D327C),
            radius = 25f,
            center = Offset(x = canvasWidth * 0.4f, y = canvasHeight * 0.5f)
        )

        // Dibujamos un "Bus" (Círculo verde)
        drawCircle(
            color = Color(0xFF2E7D32),
            radius = 30f,
            center = Offset(x = canvasWidth * 0.7f, y = canvasHeight * 0.3f)
        )
    }
}