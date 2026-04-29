package com.proyecto.kuskahaku.presentation.ui.screens
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.DirectionsBus
import androidx.compose.material.icons.filled.Warning
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

@Composable
fun AdminScreen(onLogout: () -> Unit) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(Color(0xFFF4F1DE)) // Fondo Crema Inca
            .padding(24.dp)
    ) {
        // Cabecera
        Text(
            text = "Panel de Administración",
            fontSize = 28.sp,
            fontWeight = FontWeight.Bold,
            color = Color(0xFF301934)
        )
        Text(
            text = "Resumen del sistema en tiempo real",
            fontSize = 14.sp,
            color = Color.Gray
        )

        Spacer(modifier = Modifier.height(32.dp))

        // Tarjeta de métricas 1: Unidades
        AdminStatCard(
            title = "Unidades Activas",
            value = "14 / 20",
            subtitle = "Buses circulando actualmente",
            icon = Icons.Default.DirectionsBus,
            iconColor = Color(0xFF5D327C) // Morado Chicha
        )

        Spacer(modifier = Modifier.height(16.dp))

        // Tarjeta de métricas 2: Alertas
        AdminStatCard(
            title = "Alertas del Sistema",
            value = "0",
            subtitle = "Sin incidencias reportadas",
            icon = Icons.Default.Warning,
            iconColor = Color(0xFF388E3C) // Verde indicando que todo está bien (o Rojo si hubiera alertas)
        )

        Spacer(modifier = Modifier.weight(1f))

        // Botón de salir (Discreto pero visible)
        OutlinedButton(
            onClick = onLogout,
            modifier = Modifier.fillMaxWidth().height(50.dp),
            shape = RoundedCornerShape(12.dp),
            colors = ButtonDefaults.outlinedButtonColors(contentColor = Color(0xFFD32F2F))
        ) {
            Text("Cerrar Sesión Administrativa", fontWeight = FontWeight.Bold)
        }
    }
}

// Componente reutilizable para crear las tarjetas del Dashboard fácilmente
@Composable
fun AdminStatCard(
    title: String,
    value: String,
    subtitle: String,
    icon: ImageVector,
    iconColor: Color
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        elevation = CardDefaults.cardElevation(defaultElevation = 4.dp),
        shape = RoundedCornerShape(16.dp)
    ) {
        Row(
            modifier = Modifier.padding(20.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Contenedor del ícono
            Surface(
                shape = RoundedCornerShape(12.dp),
                color = iconColor.copy(alpha = 0.1f), // Fondo clarito del mismo color del ícono
                modifier = Modifier.size(56.dp)
            ) {
                Icon(
                    imageVector = icon,
                    contentDescription = null,
                    tint = iconColor,
                    modifier = Modifier.padding(12.dp)
                )
            }

            Spacer(modifier = Modifier.width(20.dp))

            // Textos de la tarjeta
            Column {
                Text(text = title, fontSize = 14.sp, color = Color.Gray, fontWeight = FontWeight.SemiBold)
                Text(text = value, fontSize = 24.sp, fontWeight = FontWeight.Bold, color = Color(0xFF301934))
                Text(text = subtitle, fontSize = 12.sp, color = Color.Gray)
            }
        }
    }
}