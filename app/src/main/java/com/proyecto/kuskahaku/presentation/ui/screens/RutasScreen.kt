package com.proyecto.kuskahaku.presentation.ui.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.DirectionsBus
import androidx.compose.material.icons.filled.Map
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

// 1. Creamos un modelo de datos temporal para representar una Ruta
data class RutaBus(
    val codigo: String,
    val nombreEmpresa: String,
    val trayecto: String,
    val operativa: Boolean = true
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun RutasScreen() {
    var searchQuery by remember { mutableStateOf("") }

    // 2. Simulamos una base de datos de rutas (Luego esto vendrá de Firebase)
    val todasLasRutas = listOf(
        RutaBus("P13", "Kuska Haku", "Baños del Inca - UPN - Nuevo Cajamarca"),
        RutaBus("P24", "Turismo Cajamarca", "Plaza de Armas - Fonavi - Mollepampa"),
        RutaBus("P08", "Los Andes", "San Ramón - Mercado - Lucmacucho"),
        RutaBus("P11", "Inca Express", "Aeropuerto - Centro Histórico - Otuzco"),
        RutaBus("P33", "El Correcaminos", "Pacasmayo - Evitamiento - Huacariz", operativa = false)
    )

    // Filtramos la lista basándonos en lo que el usuario escriba en el buscador
    val rutasFiltradas = todasLasRutas.filter {
        it.codigo.contains(searchQuery, ignoreCase = true) ||
                it.trayecto.contains(searchQuery, ignoreCase = true)
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(Color(0xFFF4F1DE)) // Fondo Crema Inca
            .padding(horizontal = 24.dp, vertical = 16.dp)
    ) {
        // Cabecera
        Text(
            text = "Rutas Disponibles",
            fontSize = 28.sp,
            fontWeight = FontWeight.Bold,
            color = Color(0xFF301934)
        )
        Text(
            text = "Encuentra tu línea de transporte ideal",
            fontSize = 14.sp,
            color = Color.Gray
        )

        Spacer(modifier = Modifier.height(24.dp))

        // Buscador
        OutlinedTextField(
            value = searchQuery,
            onValueChange = { searchQuery = it },
            placeholder = { Text("Buscar por ruta o zona...") },
            leadingIcon = { Icon(Icons.Default.Search, contentDescription = "Buscar") },
            modifier = Modifier.fillMaxWidth(),
            shape = RoundedCornerShape(16.dp),
            colors = OutlinedTextFieldDefaults.colors(
                focusedBorderColor = Color(0xFF5D327C), // Morado Chicha
                focusedContainerColor = Color.White,
                unfocusedContainerColor = Color.White
            )
        )

        Spacer(modifier = Modifier.height(24.dp))

        // 3. LazyColumn para la lista deslizable y eficiente
        LazyColumn(
            verticalArrangement = Arrangement.spacedBy(16.dp), // Espacio entre tarjetas
            contentPadding = PaddingValues(bottom = 80.dp) // Espacio al final para que la BottomBar no tape la última tarjeta
        ) {
            items(rutasFiltradas) { ruta ->
                RutaItemCard(ruta)
            }
        }
    }
}

// 4. El diseño de cada tarjeta individual
@Composable
fun RutaItemCard(ruta: RutaBus) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp),
        shape = RoundedCornerShape(16.dp)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Círculo con el ícono del bus
            Surface(
                shape = RoundedCornerShape(12.dp),
                color = if (ruta.operativa) Color(0xFF5D327C).copy(alpha = 0.1f) else Color.LightGray.copy(alpha = 0.3f),
                modifier = Modifier.size(60.dp)
            ) {
                Column(
                    verticalArrangement = Arrangement.Center, // <-- ESTA ES LA CORRECCIÓN
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Icon(
                        imageVector = Icons.Default.DirectionsBus,
                        contentDescription = "Bus",
                        tint = if (ruta.operativa) Color(0xFF5D327C) else Color.Gray,
                        modifier = Modifier.size(24.dp)
                    )
                    Text(
                        text = ruta.codigo,
                        fontWeight = FontWeight.Bold,
                        color = if (ruta.operativa) Color(0xFF5D327C) else Color.Gray,
                        fontSize = 14.sp
                    )
                }
            }

            Spacer(modifier = Modifier.width(16.dp))

            // Información de la ruta
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = ruta.nombreEmpresa,
                    fontSize = 16.sp,
                    fontWeight = FontWeight.Bold,
                    color = if (ruta.operativa) Color(0xFF301934) else Color.Gray
                )
                Text(
                    text = ruta.trayecto,
                    fontSize = 12.sp,
                    color = Color.Gray,
                    lineHeight = 16.sp
                )
                if (!ruta.operativa) {
                    Text(
                        text = "Fuera de servicio hoy",
                        fontSize = 12.sp,
                        color = Color(0xFFD32F2F),
                        fontWeight = FontWeight.SemiBold,
                        modifier = Modifier.padding(top = 4.dp)
                    )
                }
            }

            // Botón para ver en el mapa (Solo visual por ahora)
            if (ruta.operativa) {
                IconButton(onClick = { /* Acción futura para ver la ruta en el mapa */ }) {
                    Icon(
                        imageVector = Icons.Default.Map,
                        contentDescription = "Ver en mapa",
                        tint = Color(0xFF388E3C) // Verde
                    )
                }
            }
        }
    }
}