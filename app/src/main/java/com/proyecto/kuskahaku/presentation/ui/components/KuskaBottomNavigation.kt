package com.proyecto.kuskahaku.presentation.ui.components

import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountCircle
import androidx.compose.material.icons.filled.ConfirmationNumber
import androidx.compose.material.icons.filled.DirectionsBus
import androidx.compose.material.icons.filled.Explore
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.unit.dp
import androidx.navigation.NavHostController
import androidx.navigation.compose.currentBackStackEntryAsState

// 1. Añadimos el parámetro 'route' a la clase
sealed class NavItem(val title: String, val icon: ImageVector, val route: String) {
    object Explorar : NavItem("Explorar", Icons.Default.Explore, "explorar")
    object Rutas : NavItem("Rutas", Icons.Default.DirectionsBus, "rutas")
    object Tickets : NavItem("Tickets", Icons.Default.ConfirmationNumber, "tickets")
    object Perfil : NavItem("Perfil", Icons.Default.AccountCircle, "perfil")
}

@Composable
fun KuskaBottomNavigation(navController: NavHostController) {
    val items = listOf(NavItem.Explorar, NavItem.Rutas, NavItem.Tickets, NavItem.Perfil)

    // 2. Escuchamos en qué pantalla estamos actualmente
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    val currentRoute = navBackStackEntry?.destination?.route

    NavigationBar(
        containerColor = Color.White,
        tonalElevation = 8.dp
    ) {
        items.forEach { item ->
            NavigationBarItem(
                selected = currentRoute == item.route, // Se marca si coincide la ruta
                onClick = {
                    navController.navigate(item.route) {
                        // Evita que se acumulen pantallas infinitas al presionar el mismo botón
                        popUpTo(navController.graph.startDestinationId) { saveState = true }
                        launchSingleTop = true
                        restoreState = true
                    }
                },
                label = { Text(item.title) },
                icon = { Icon(item.icon, contentDescription = item.title) },
                colors = NavigationBarItemDefaults.colors(
                    selectedIconColor = Color(0xFF5D327C),
                    selectedTextColor = Color(0xFF5D327C),
                    indicatorColor = Color(0xFFF4F1DE),
                    unselectedIconColor = Color.Gray,
                    unselectedTextColor = Color.Gray
                )
            )
        }
    }
}