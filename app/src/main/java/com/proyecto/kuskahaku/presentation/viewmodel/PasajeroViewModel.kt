package com.proyecto.kuskahaku.presentation.viewmodel

import androidx.lifecycle.ViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

// Añadimos las variables que representarán la realidad
data class PasajeroUiState(
    val isLoading: Boolean = false,
    val searchQuery: String = "",
    val nombreRutaActiva: String = "P13 - Kuska Haku",
    val tiempoEstimadoLlegada: String = "Calculando...", // Esto cambiará con Firebase
    val estadoRuta: String = "Buscando bus..."
)

class PasajeroViewModel : ViewModel() {
    private val _uiState = MutableStateFlow(PasajeroUiState())
    val uiState: StateFlow<PasajeroUiState> = _uiState.asStateFlow()

    fun onSearchQueryChanged(newQuery: String) {
        _uiState.value = _uiState.value.copy(searchQuery = newQuery)
    }
}