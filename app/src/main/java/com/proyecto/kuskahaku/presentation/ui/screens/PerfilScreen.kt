package com.proyecto.kuskahaku.presentation.ui.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountCircle
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

@Composable
fun PerfilScreen(onLogout: () -> Unit) { // <-- 1. Añadimos el parámetro aquí
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(Color(0xFFF4F1DE)) // ¡Aquí aplicamos tu color de fondo explícito!
            .padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Icon(
            imageVector = Icons.Default.AccountCircle,
            contentDescription = "Foto de perfil",
            modifier = Modifier.size(120.dp),
            tint = Color(0xFF5D327C) // Chicha Purple
        )
        Spacer(modifier = Modifier.height(24.dp))

        Text(
            text = "Usuario",
            fontSize = 28.sp,
            fontWeight = FontWeight.Bold,
            color = Color(0xFF301934)
        )
        Text(
            text = "Estudiante / Pasajero",
            fontSize = 16.sp,
            color = Color.Gray
        )

        Spacer(modifier = Modifier.height(40.dp))

        Button(
            onClick = { onLogout() }, // <-- 2. Cambiamos el comentario por la acción real
            colors = ButtonDefaults.buttonColors(containerColor = Color(0xFFD32F2F)), // Rojo para cierre
            modifier = Modifier.fillMaxWidth(0.8f).height(50.dp)
        ) {
            Text(text = "Cerrar Sesión", fontSize = 16.sp)
        }
    }
}