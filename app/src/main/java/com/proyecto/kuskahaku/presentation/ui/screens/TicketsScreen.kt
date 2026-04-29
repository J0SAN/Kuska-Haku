package com.proyecto.kuskahaku.presentation.ui.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountBalanceWallet
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.DirectionsBus
import androidx.compose.material.icons.filled.QrCode2
import androidx.compose.material.icons.filled.ReceiptLong
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

@Composable
fun TicketsScreen() {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(Color(0xFFF4F1DE)) // Fondo Crema Inca
            .padding(horizontal = 24.dp, vertical = 16.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        // Cabecera
        Text(
            text = "Mi Billetera",
            fontSize = 28.sp,
            fontWeight = FontWeight.Bold,
            color = Color(0xFF301934),
            modifier = Modifier.fillMaxWidth(),
            textAlign = TextAlign.Start
        )
        Text(
            text = "Gestiona tu saldo y pasajes",
            fontSize = 14.sp,
            color = Color.Gray,
            modifier = Modifier.fillMaxWidth(),
            textAlign = TextAlign.Start
        )

        Spacer(modifier = Modifier.height(24.dp))

        // 1. Tarjeta de Saldo Principal
        Card(
            modifier = Modifier.fillMaxWidth(),
            colors = CardDefaults.cardColors(containerColor = Color(0xFF5D327C)), // Fondo Morado Chicha
            elevation = CardDefaults.cardElevation(defaultElevation = 8.dp),
            shape = RoundedCornerShape(24.dp)
        ) {
            Column(
                modifier = Modifier.padding(24.dp),
                horizontalAlignment = Alignment.Start
            ) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(
                        imageVector = Icons.Default.AccountBalanceWallet,
                        contentDescription = "Billetera",
                        tint = Color.White,
                        modifier = Modifier.size(24.dp)
                    )
                    Spacer(modifier = Modifier.width(8.dp))
                    Text("Saldo Disponible", color = Color.White.copy(alpha = 0.8f), fontSize = 14.sp)
                }

                Spacer(modifier = Modifier.height(16.dp))

                Text(
                    text = "S/ 12.50",
                    fontSize = 40.sp,
                    fontWeight = FontWeight.Bold,
                    color = Color.White
                )

                Spacer(modifier = Modifier.height(24.dp))

                // Botón para recargar saldo
                Button(
                    onClick = { /* Acción de recarga */ },
                    colors = ButtonDefaults.buttonColors(containerColor = Color.White),
                    shape = RoundedCornerShape(12.dp),
                    modifier = Modifier.fillMaxWidth().height(48.dp)
                ) {
                    Icon(imageVector = Icons.Default.Add, contentDescription = "Recargar", tint = Color(0xFF5D327C))
                    Spacer(modifier = Modifier.width(8.dp))
                    Text("Recargar Saldo", color = Color(0xFF5D327C), fontWeight = FontWeight.Bold)
                }
            }
        }

        Spacer(modifier = Modifier.height(32.dp))

        // 2. Tarjeta del Código QR para Pagar
        Card(
            modifier = Modifier.fillMaxWidth(),
            colors = CardDefaults.cardColors(containerColor = Color.White),
            elevation = CardDefaults.cardElevation(defaultElevation = 4.dp),
            shape = RoundedCornerShape(16.dp)
        ) {
            Column(
                modifier = Modifier.padding(24.dp),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Text(
                    text = "Muestra este código al subir al bus",
                    fontSize = 14.sp,
                    color = Color.Gray,
                    textAlign = TextAlign.Center
                )
                Spacer(modifier = Modifier.height(16.dp))

                // Simulación visual de un QR (Solo el ícono gigante por ahora)
                Icon(
                    imageVector = Icons.Default.QrCode2,
                    contentDescription = "Código QR de Pago",
                    tint = Color(0xFF301934),
                    modifier = Modifier.size(150.dp)
                )
            }
        }

        Spacer(modifier = Modifier.height(32.dp))

        // 3. Mini historial (Simulado)
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text("Últimos Movimientos", fontSize = 18.sp, fontWeight = FontWeight.Bold, color = Color(0xFF301934))
            TextButton(onClick = { /* Ver todos */ }) {
                Text("Ver todo", color = Color(0xFF5D327C))
            }
        }

        // Un ítem de historial
        MovimientoItem(
            ruta = "Ruta P13",
            fecha = "Hoy, 08:30 AM",
            monto = "-S/ 1.50",
            icon = Icons.Default.DirectionsBus
        )
        Spacer(modifier = Modifier.height(12.dp))
        MovimientoItem(
            ruta = "Recarga Yape",
            fecha = "Ayer, 18:45 PM",
            monto = "+S/ 10.00",
            icon = Icons.Default.ReceiptLong,
            esIngreso = true
        )
    }
}

// Componente para los ítems del historial
@Composable
fun MovimientoItem(ruta: String, fecha: String, monto: String, icon: ImageVector, esIngreso: Boolean = false) {
    Row(
        modifier = Modifier.fillMaxWidth().padding(vertical = 8.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            shape = RoundedCornerShape(12.dp),
            color = if (esIngreso) Color(0xFF388E3C).copy(alpha = 0.1f) else Color.Gray.copy(alpha = 0.1f),
            modifier = Modifier.size(48.dp)
        ) {
            Icon(
                imageVector = icon,
                contentDescription = null,
                tint = if (esIngreso) Color(0xFF388E3C) else Color.Gray,
                modifier = Modifier.padding(12.dp)
            )
        }
        Spacer(modifier = Modifier.width(16.dp))
        Column(modifier = Modifier.weight(1f)) {
            Text(text = ruta, fontWeight = FontWeight.SemiBold, color = Color(0xFF301934))
            Text(text = fecha, fontSize = 12.sp, color = Color.Gray)
        }
        Text(
            text = monto,
            fontWeight = FontWeight.Bold,
            color = if (esIngreso) Color(0xFF388E3C) else Color(0xFF301934),
            fontSize = 16.sp
        )
    }
}

