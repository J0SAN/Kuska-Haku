package com.proyecto.kuskahaku.presentation.ui.screens

import android.Manifest
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.lifecycle.viewmodel.compose.viewModel
import com.proyecto.kuskahaku.presentation.ui.components.GoogleMapView
import com.proyecto.kuskahaku.presentation.ui.components.MapMockup
import com.proyecto.kuskahaku.presentation.ui.components.RouteInfoCard
import com.proyecto.kuskahaku.presentation.ui.components.SearchBarInca
import com.proyecto.kuskahaku.presentation.viewmodel.MapViewModel
import com.proyecto.kuskahaku.presentation.viewmodel.PasajeroViewModel
import kotlinx.coroutines.launch

@Composable
fun PasajeroDashboardScreen(
    viewModel: PasajeroViewModel = viewModel(),
    mapViewModel: MapViewModel = viewModel()
) {
    val context = LocalContext.current

    // 1. CREA EL SCOPE DE COMPOSE AQUÍ (justo debajo del context)
    val coroutineScope = rememberCoroutineScope()

    val permissionLauncher = rememberLauncherForActivityResult(
        contract = ActivityResultContracts.RequestMultiplePermissions()
    ) { permissions ->
        val gpsOtorgado = permissions[Manifest.permission.ACCESS_FINE_LOCATION] == true ||
                permissions[Manifest.permission.ACCESS_COARSE_LOCATION] == true

        if (gpsOtorgado) {
            val locationService = com.proyecto.kuskahaku.data.LocationService(context)

            // 2. REEMPLAZA EL VIEWMODELSCOPE POR TU NUEVO COROUTINESCOPE
            coroutineScope.launch {
                locationService.obtenerUbicacionEnTiempoReal().collect { coordenadas ->
                    // AHORA SÍ: Enviamos las coordenadas reales al ViewModel
                    mapViewModel.actualizarUbicacion(lat = coordenadas.first, lng = coordenadas.second)
                }
            }
        }
    }

    // Pedimos permiso al abrir la app
    LaunchedEffect(Unit) {
        permissionLauncher.launch(
            arrayOf(
                Manifest.permission.ACCESS_FINE_LOCATION,
                Manifest.permission.ACCESS_COARSE_LOCATION
            )
        )
    }

    val uiState by viewModel.uiState.collectAsState()
    val mapState by mapViewModel.mapState.collectAsState()

    Box(modifier = Modifier.fillMaxSize()) {

        // --- EL CAMBIO PRINCIPAL ---
        // Borramos MapMockup y ponemos el Google Map
        GoogleMapView(
            userLat = mapState.userLat,
            userLng = mapState.userLng
        )

        Column(
            modifier = Modifier.fillMaxSize(),
            verticalArrangement = Arrangement.SpaceBetween
        ) {
            SearchBarInca(
                query = uiState.searchQuery,
                onQueryChange = { viewModel.onSearchQueryChanged(it) }
            )

            RouteInfoCard(
                routeName = uiState.nombreRutaActiva,
                eta = uiState.tiempoEstimadoLlegada,
                status = uiState.estadoRuta,
                onViewRouteClick = { /* Acción */ }
            )
        }
    }
}