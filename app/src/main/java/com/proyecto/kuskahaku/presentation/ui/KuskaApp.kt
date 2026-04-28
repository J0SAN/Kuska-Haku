package com.proyecto.kuskahaku.presentation.ui

import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.proyecto.kuskahaku.presentation.ui.components.KuskaBottomNavigation
import com.proyecto.kuskahaku.presentation.ui.components.NavItem
import com.proyecto.kuskahaku.presentation.ui.screens.PasajeroDashboardScreen
import com.proyecto.kuskahaku.presentation.ui.screens.PerfilScreen

@Composable
fun KuskaApp() {
    // El controlador que gestiona los viajes entre pantallas
    val navController = rememberNavController()

    Scaffold(
        bottomBar = { KuskaBottomNavigation(navController = navController) }
    ) { paddingValues ->
        // El espacio central que cambiará su contenido
        NavHost(
            navController = navController,
            startDestination = NavItem.Explorar.route,
            modifier = Modifier.padding(paddingValues) // Respeta el espacio de la barra inferior
        ) {
            // Declaramos qué función pintar para cada ruta
            composable(NavItem.Explorar.route) {
                PasajeroDashboardScreen()
            }
            composable(NavItem.Rutas.route) {
                // Pantalla temporal vacía por ahora
            }
            composable(NavItem.Tickets.route) {
                // Pantalla temporal vacía por ahora
            }
            composable(NavItem.Perfil.route) {
                PerfilScreen()
            }
        }
    }
}