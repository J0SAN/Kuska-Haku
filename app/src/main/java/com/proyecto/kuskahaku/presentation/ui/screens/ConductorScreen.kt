package com.proyecto.kuskahaku.presentation.ui.screens
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.DirectionsBus
import androidx.compose.material.icons.filled.PowerSettingsNew
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

@Composable
fun ConductorScreen(onLogout: () -> Unit) {
    var transmitiendo by remember { mutableStateOf(false) }

    // Fondo crema de tu marca
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(Color(0xFFF4F1DE))
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        // Cabecera superior
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column {
                Text(
                    text = "Panel de Control",
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Bold,
                    color = Color(0xFF301934)
                )
                Text(
                    text = "Unidad: Ruta P13",
                    fontSize = 16.sp,
                    color = Color.Gray
                )
            }
            Icon(
                imageVector = Icons.Default.DirectionsBus,
                contentDescription = "Bus",
                tint = Color(0xFF5D327C),
                modifier = Modifier.size(40.dp)
            )
        }

        Spacer(modifier = Modifier.height(48.dp))

        // Tarjeta indicadora de estado
        Card(
            colors = CardDefaults.cardColors(
                containerColor = Color.White
            ),
            elevation = CardDefaults.cardElevation(defaultElevation = 4.dp),
            modifier = Modifier.fillMaxWidth()
        ) {
            Column(
                modifier = Modifier.padding(24.dp),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Text(
                    text = "ESTADO DEL GPS",
                    fontSize = 14.sp,
                    fontWeight = FontWeight.Bold,
                    color = Color.Gray
                )
                Spacer(modifier = Modifier.height(8.dp))
                Text(
                    text = if (transmitiendo) "Transmitiendo en vivo" else "Fuera de servicio",
                    fontSize = 20.sp,
                    fontWeight = FontWeight.Bold,
                    color = if (transmitiendo) Color(0xFF2E7D32) else Color(0xFFD32F2F)
                )
            }
        }

        Spacer(modifier = Modifier.height(64.dp))

        // Botón masivo de acción (Tipo botón de pánico/encendido)
        Button(
            onClick = { transmitiendo = !transmitiendo },
            shape = CircleShape,
            colors = ButtonDefaults.buttonColors(
                containerColor = if (transmitiendo) Color(0xFFD32F2F) else Color(0xFF2E7D32)
            ),
            modifier = Modifier.size(200.dp)
        ) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Icon(
                    imageVector = Icons.Default.PowerSettingsNew,
                    contentDescription = null,
                    modifier = Modifier.size(64.dp)
                )
                Spacer(modifier = Modifier.height(8.dp))
                Text(
                    text = if (transmitiendo) "DETENER" else "INICIAR",
                    fontSize = 20.sp,
                    fontWeight = FontWeight.Bold
                )
            }
        }

        Spacer(modifier = Modifier.weight(1f))

        // Botón de cerrar sesión más sutil
        OutlinedButton(
            onClick = onLogout,
            modifier = Modifier.fillMaxWidth(),
            shape = RoundedCornerShape(12.dp)
        ) {
            Text("Cerrar Sesión", color = Color(0xFFD32F2F))
        }
    }
}