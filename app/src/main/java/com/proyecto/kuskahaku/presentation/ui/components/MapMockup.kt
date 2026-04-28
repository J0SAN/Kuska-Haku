package com.proyecto.kuskahaku.presentation.ui.components

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
@Composable
fun MapMockup(
    busX: Float, // 1. Recibimos la posición X en tiempo real
    busY: Float, // 2. Recibimos la posición Y en tiempo real
    modifier: Modifier = Modifier
) {
    Canvas(modifier = modifier.fillMaxSize()) {
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

        // Dibujamos un "Paradero" fijo (Círculo morado)
        drawCircle(
            color = Color(0xFF5D327C),
            radius = 25f,
            center = Offset(x = canvasWidth * 0.4f, y = canvasHeight * 0.5f)
        )

        // Dibujamos un "Bus" móvil (Círculo verde)
        drawCircle(
            color = Color(0xFF2E7D32),
            radius = 30f,
            // 3. AQUÍ está la clave: usamos las variables en lugar de valores fijos
            center = Offset(x = busX, y = busY)
        )
    }
}