package com.proyecto.kuskahaku.presentation.ui.components

import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.ui.Modifier
import com.google.android.gms.maps.CameraUpdateFactory
import com.google.android.gms.maps.model.CameraPosition
import com.google.android.gms.maps.model.LatLng
import com.google.maps.android.compose.GoogleMap
import com.google.maps.android.compose.MapProperties
import com.google.maps.android.compose.MapUiSettings
import com.google.maps.android.compose.Marker
import com.google.maps.android.compose.MarkerState
import com.google.maps.android.compose.rememberCameraPositionState

@Composable
fun GoogleMapView(
    userLat: Double,
    userLng: Double,
    modifier: Modifier = Modifier
) {
    // 1. Configuramos la posición inicial de la cámara
    // Si las coordenadas son 0.0 (aún no carga el GPS), centramos por defecto en Cajamarca
    val ubicacionActual = if (userLat != 0.0 && userLng != 0.0) {
        LatLng(userLat, userLng)
    } else {
        LatLng(-7.1617, -78.5128) // Plaza de Armas de Cajamarca
    }

    val cameraPositionState = rememberCameraPositionState {
        position = CameraPosition.fromLatLngZoom(ubicacionActual, 15f)
    }

    // 2. Si el GPS se actualiza, movemos la cámara suavemente a la nueva posición
    LaunchedEffect(userLat, userLng) {
        if (userLat != 0.0 && userLng != 0.0) {
            cameraPositionState.animate(
                update = CameraUpdateFactory.newLatLng(LatLng(userLat, userLng)),
                durationMs = 1000
            )
        }
    }

    // 3. Dibujamos el mapa de Google
    GoogleMap(
        modifier = modifier.fillMaxSize(),
        cameraPositionState = cameraPositionState,
        properties = MapProperties(
            isMyLocationEnabled = true // Esto habilita el punto azul nativo de Google (si hay permisos)
        ),
        uiSettings = MapUiSettings(
            zoomControlsEnabled = false, // Ocultamos los botones de zoom para que se vea limpio
            myLocationButtonEnabled = false // Ocultamos el botón de "mi ubicación" para controlarlo nosotros
        )
    ) {
        // 4. (Opcional) Dibujamos un marcador personalizado en tu ubicación
        if (userLat != 0.0 && userLng != 0.0) {
            Marker(
                state = MarkerState(position = LatLng(userLat, userLng)),
                title = "Tú",
                snippet = "Tu ubicación actual"
            )
        }
    }
}