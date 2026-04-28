package com.proyecto.kuskahaku.presentation.viewmodel
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

// 1. Cambiamos X e Y (píxeles) por Latitud y Longitud (coordenadas reales)
data class MapState(
    val userLat: Double = 0.0,
    val userLng: Double = 0.0
)

class MapViewModel : ViewModel() {

    // Variables reactivas que la UI va a observar
    private val _mapState = MutableStateFlow(MapState())
    val mapState: StateFlow<MapState> = _mapState.asStateFlow()

    // 2. Ya no hay 'init' ni 'simularMovimientoBus'.
    // Ahora usamos esta función que será llamada cada vez que el GPS detecte un cambio.
    fun actualizarUbicacion(lat: Double, lng: Double) {
        _mapState.value = MapState(userLat = lat, userLng = lng)
    }
}