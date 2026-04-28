package com.proyecto.kuskahaku.presentation.ui.screens
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp

@Composable
fun ConductorScreen(onLogout: () -> Unit) {
    var transmitiendo by remember { mutableStateOf(false) }

    Column(
        modifier = Modifier.fillMaxSize().padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text("Panel de Conductor - Ruta P13", style = MaterialTheme.typography.titleLarge)
        Spacer(modifier = Modifier.height(48.dp))

        // Botón principal de transmisión
        Button(
            onClick = { transmitiendo = !transmitiendo },
            colors = ButtonDefaults.buttonColors(
                containerColor = if (transmitiendo) Color(0xFFD32F2F) else Color(0xFF388E3C)
            ),
            modifier = Modifier.size(200.dp, 60.dp)
        ) {
            Text(if (transmitiendo) "Detener Transmisión" else "Iniciar Transmisión GPS")
        }

        Spacer(modifier = Modifier.height(24.dp))

        if (transmitiendo) {
            Text("Transmitiendo ubicación en vivo...", color = Color(0xFF388E3C))
        }

        Spacer(modifier = Modifier.weight(1f))

        Button(onClick = onLogout, colors = ButtonDefaults.buttonColors(containerColor = Color.Gray)) {
            Text("Cerrar Sesión")
        }
    }
}