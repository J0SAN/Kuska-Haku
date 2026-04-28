package com.proyecto.kuskahaku.presentation.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Schedule
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

@Composable
fun RouteInfoCard(
    routeName: String,
    eta: String,
    status: String,
    onViewRouteClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        shape = RoundedCornerShape(topStart = 24.dp, topEnd = 24.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
    ) {
        Column(
            modifier = Modifier.fillMaxWidth().padding(24.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            // El "handle" o tirador gris superior
            Box(
                modifier = Modifier.width(40.dp).height(4.dp)
                    .background(Color.LightGray, RoundedCornerShape(2.dp))
            )
            Spacer(modifier = Modifier.height(20.dp))

            // Título de la ruta
            Text(
                text = routeName, fontSize = 22.sp, fontWeight = FontWeight.Bold,
                color = Color(0xFF301934), modifier = Modifier.fillMaxWidth()
            )
            Spacer(modifier = Modifier.height(16.dp))

            // Fila de etiquetas
            Row(modifier = Modifier.fillMaxWidth()) {
                // Etiqueta ETA
                Surface(shape = RoundedCornerShape(16.dp), color = Color(0xFFF4F1DE)) {
                    Row(modifier = Modifier.padding(horizontal = 12.dp, vertical = 6.dp), verticalAlignment = Alignment.CenterVertically) {
                        Icon(Icons.Default.Schedule, contentDescription = null, modifier = Modifier.size(16.dp), tint = Color.DarkGray)
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(text = eta, fontSize = 14.sp, color = Color.DarkGray)
                    }
                }
                Spacer(modifier = Modifier.width(12.dp))
                // Etiqueta Estado
                Surface(shape = RoundedCornerShape(16.dp), color = Color(0xFFC8E6C9)) {
                    Text(text = status, fontSize = 14.sp, color = Color(0xFF2E7D32), modifier = Modifier.padding(horizontal = 12.dp, vertical = 6.dp))
                }
            }
            Spacer(modifier = Modifier.height(24.dp))

            Button(
                onClick = onViewRouteClick,
                modifier = Modifier.fillMaxWidth().height(50.dp),
                colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF5D327C))
            ) {
                Text(text = "Ver recorrido completo", fontSize = 16.sp)
            }
        }
    }
}