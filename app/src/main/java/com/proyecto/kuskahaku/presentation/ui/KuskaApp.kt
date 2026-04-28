package com.proyecto.kuskahaku.presentation.ui

import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.proyecto.kuskahaku.presentation.ui.components.KuskaBottomNavigation
import com.proyecto.kuskahaku.presentation.ui.components.NavItem
import com.proyecto.kuskahaku.presentation.ui.screens.AdminScreen
import com.proyecto.kuskahaku.presentation.ui.screens.ConductorScreen
import com.proyecto.kuskahaku.presentation.ui.screens.LoginScreen
import com.proyecto.kuskahaku.presentation.ui.screens.PasajeroDashboardScreen
import com.proyecto.kuskahaku.presentation.ui.screens.PerfilScreen
import androidx.compose.runtime.getValue
import androidx.compose.runtime.setValue

@Composable
fun KuskaApp() {
    // El controlador que gestiona los viajes entre pantallas
    val navController = rememberNavController()

    // 1. Añadimos el estado para saber si el usuario está logueado y qué rol tiene
    var isLoggedIn by remember { mutableStateOf(false) }
    var userRole by remember { mutableStateOf("") } // Puede ser "pasajero", "conductor", o "admin"

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

            // 2. Modificamos la ruta de Perfil para que exija Login
            composable(NavItem.Perfil.route) {
                if (!isLoggedIn) {
                    // Si no está logueado, mostramos la pantalla de Login
                    LoginScreen(
                        onLoginSuccess = { role ->
                            isLoggedIn = true
                            userRole = role
                        }
                    )
                } else {
                    // Si ya está logueado, mostramos el panel según su rol
                    when (userRole) {
                        "conductor" -> ConductorScreen(onLogout = { isLoggedIn = false })
                        "admin" -> AdminScreen(onLogout = { isLoggedIn = false })
                        else -> PerfilScreen(onLogout = { isLoggedIn = false }) // El perfil normal del pasajero
                    }
                }
            }
        }
    }
}