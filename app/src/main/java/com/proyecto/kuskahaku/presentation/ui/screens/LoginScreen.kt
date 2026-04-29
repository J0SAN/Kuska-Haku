package com.proyecto.kuskahaku.presentation.ui.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.DirectionsBus
import androidx.compose.material.icons.filled.Email
import androidx.compose.material.icons.filled.Lock
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LoginScreen(onLoginSuccess: (String) -> Unit) {
    var email by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }

    // 1. El fondo de toda la pantalla usando tu color Crema Inca
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(Color(0xFFF4F1DE)),
        contentAlignment = Alignment.Center
    ) {

        Card(
            modifier = Modifier
                .fillMaxWidth(0.9f) // Ocupa el 90% del ancho
                .padding(16.dp),
            colors = CardDefaults.cardColors(containerColor = Color.White),
            elevation = CardDefaults.cardElevation(defaultElevation = 8.dp),
            shape = RoundedCornerShape(24.dp) // Bordes muy redondeados
        ) {
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(32.dp),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                // 3. Logo o Ícono Superior
                Icon(
                    imageVector = Icons.Default.DirectionsBus,
                    contentDescription = "Logo Kuska Haku",
                    tint = Color(0xFF5D327C), // Morado Chicha
                    modifier = Modifier.size(64.dp)
                )

                Spacer(modifier = Modifier.height(16.dp))


                Text(
                    text = "Bienvenido a Kuska Haku",
                    fontSize = 22.sp,
                    fontWeight = FontWeight.Bold,
                    color = Color(0xFF301934)
                )
                Text(
                    text = "Inicia sesión para continuar",
                    fontSize = 14.sp,
                    color = Color.Gray
                )

                Spacer(modifier = Modifier.height(32.dp))

                // 5. Campos de texto modernos (Outlined)
                OutlinedTextField(
                    value = email,
                    onValueChange = { email = it },
                    label = { Text("Correo Electrónico") },
                    leadingIcon = { Icon(Icons.Default.Email, contentDescription = null) },
                    modifier = Modifier.fillMaxWidth(),
                    shape = RoundedCornerShape(12.dp),
                    colors = TextFieldDefaults.colors(
                        focusedIndicatorColor = Color(0xFF5D327C),
                        focusedLabelColor = Color(0xFF5D327C)
                    )
                )

                Spacer(modifier = Modifier.height(16.dp))

                OutlinedTextField(
                    value = password,
                    onValueChange = { password = it },
                    label = { Text("Contraseña") },
                    leadingIcon = { Icon(Icons.Default.Lock, contentDescription = null) },
                    visualTransformation = PasswordVisualTransformation(), // Oculta los puntos de la contraseña
                    modifier = Modifier.fillMaxWidth(),
                    shape = RoundedCornerShape(12.dp),
                    colors = TextFieldDefaults.colors(
                        focusedIndicatorColor = Color(0xFF5D327C),
                        focusedLabelColor = Color(0xFF5D327C)
                    )
                )

                Spacer(modifier = Modifier.height(32.dp))

                // 6. Botones de Simulación de Roles con Jerarquía Visual

                // Pasajero (El botón más importante y grande)
                Button(
                    onClick = { onLoginSuccess("pasajero") },
                    modifier = Modifier.fillMaxWidth().height(50.dp),
                    shape = RoundedCornerShape(12.dp),
                    colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF5D327C))
                ) {
                    Text("Entrar como Pasajero", fontSize = 16.sp, fontWeight = FontWeight.Bold)
                }

                Spacer(modifier = Modifier.height(12.dp))

                // Conductor (Color secundario)
                Button(
                    onClick = { onLoginSuccess("conductor") },
                    modifier = Modifier.fillMaxWidth().height(45.dp),
                    shape = RoundedCornerShape(12.dp),
                    colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF388E3C))
                ) {
                    Text("Entrar como Conductor")
                }

                Spacer(modifier = Modifier.height(8.dp))

                // Admin (Botón con borde, menos intrusivo)
                OutlinedButton(
                    onClick = { onLoginSuccess("admin") },
                    modifier = Modifier.fillMaxWidth().height(45.dp),
                    shape = RoundedCornerShape(12.dp)
                ) {
                    Text("Entrar como Administrador", color = Color.Gray)
                }
            }
        }
    }
}