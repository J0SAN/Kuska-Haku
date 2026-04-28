package com.proyecto.kuskahaku.presentation.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.lifecycle.viewmodel.compose.viewModel
import com.proyecto.kuskahaku.presentation.ui.components.MapMockup
import com.proyecto.kuskahaku.presentation.ui.components.RouteInfoCard
import com.proyecto.kuskahaku.presentation.ui.components.SearchBarInca
import com.proyecto.kuskahaku.presentation.viewmodel.PasajeroViewModel

@Composable
fun PasajeroDashboardScreen(
    viewModel: PasajeroViewModel = viewModel()
) {
    val state by viewModel.uiState.collectAsState()

    // Usamos Box para que los elementos se apilen (Z-Index)
    Box(modifier = Modifier.fillMaxSize()) {

        // 1. Capa de fondo: El Mapa Mockup
        MapMockup()

        // 2. Capa de Interfaz: Barra y Tarjeta
        Column(
            modifier = Modifier.fillMaxSize(),
            verticalArrangement = Arrangement.SpaceBetween
        ) {
            // Barra de búsqueda (flota arriba)
            SearchBarInca(
                query = state.searchQuery,
                onQueryChange = { viewModel.onSearchQueryChanged(it) }
            )

            // Tarjeta de información (flota abajo)
            RouteInfoCard(
                routeName = state.nombreRutaActiva,
                eta = state.tiempoEstimadoLlegada,
                status = state.estadoRuta,
                onViewRouteClick = { /* Acción */ }
            )
        }
    }
}