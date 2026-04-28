package com.proyecto.kuskahaku.presentation.ui.screens
import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.DirectionsBus
import androidx.compose.material.icons.filled.Warning
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

@Composable
fun AdminScreen(onLogout: () -> Unit) {
    Column(
        modifier = Modifier.fillMaxSize().padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text("Panel de Administración", style = MaterialTheme.typography.headlineSmall)
        Spacer(modifier = Modifier.height(24.dp))

        // Tarjeta de métricas
        Card(modifier = Modifier.fillMaxWidth().padding(8.dp)) {
            Row(modifier = Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
                Icon(Icons.Default.DirectionsBus, contentDescription = null, tint = MaterialTheme.colorScheme.primary)
                Spacer(modifier = Modifier.width(16.dp))
                Column {
                    Text("Unidades Activas", style = MaterialTheme.typography.titleMedium)
                    Text("14 buses en ruta actualmente")
                }
            }
        }

        Card(modifier = Modifier.fillMaxWidth().padding(8.dp)) {
            Row(modifier = Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
                Icon(Icons.Default.Warning, contentDescription = null, tint = MaterialTheme.colorScheme.error)
                Spacer(modifier = Modifier.width(16.dp))
                Column {
                    Text("Alertas del Sistema", style = MaterialTheme.typography.titleMedium)
                    Text("Sin incidencias reportadas")
                }
            }
        }

        Spacer(modifier = Modifier.weight(1f))

        Button(onClick = onLogout, colors = ButtonDefaults.buttonColors(containerColor = MaterialTheme.colorScheme.error)) {
            Text("Cerrar Sesión Admin")
        }
    }
}