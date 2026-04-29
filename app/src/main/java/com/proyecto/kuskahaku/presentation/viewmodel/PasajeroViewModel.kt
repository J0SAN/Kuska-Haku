package com.proyecto.kuskahaku.presentation.viewmodel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.proyecto.kuskahaku.data.repository.TransporteRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

data class PasajeroUiState(
    val isLoading: Boolean = false,
    val searchQuery: String = "",
    val nombreRutaActiva: String = "Línea P13 - Kuska Haku", // Ya actualizado a P13
    val tiempoEstimadoLlegada: String = "Conectando...",
    val estadoRuta: String = "Buscando señal..."
)

class PasajeroViewModel : ViewModel() {
    private val repository = TransporteRepository()

    private val _uiState = MutableStateFlow(PasajeroUiState())
    val uiState: StateFlow<PasajeroUiState> = _uiState.asStateFlow()

    init {
        // Apenas nace el ViewModel, empezamos a observar los buses
        observarBusesEnTiempoReal()
    }

    fun onSearchQueryChanged(newQuery: String) {
        _uiState.value = _uiState.value.copy(searchQuery = newQuery)
    }

    private fun observarBusesEnTiempoReal() {
        viewModelScope.launch {
            // Recolectamos el flujo de datos que viene desde Firebase
            repository.obtenerBusesP13().collect { listaDeBuses ->
                if (listaDeBuses.isNotEmpty()) {
                    // Para este MVP, tomamos el primer bus disponible (bus_01)
                    val busPrincipal = listaDeBuses[0]

                    _uiState.value = _uiState.value.copy(
                        estadoRuta = busPrincipal.estado,
                        tiempoEstimadoLlegada = "Actualizado en vivo" // Luego calcularemos el ETA real con la lat/lon
                    )
                } else {
                    // Si no hay buses en el JSON (o está vacío)
                    _uiState.value = _uiState.value.copy(
                        estadoRuta = "Sin buses activos",
                        tiempoEstimadoLlegada = "--"
                    )
                }
            }
        }
    }
}