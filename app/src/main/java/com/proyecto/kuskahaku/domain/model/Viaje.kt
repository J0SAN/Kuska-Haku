package com.proyecto.kuskahaku.domain.model

enum class EstadoViaje {
    PROGRAMADO, EN_CURSO, FINALIZADO, CANCELADO
}

data class Viaje (
    val id: String,
    val idConductor: String,
    val idVehiculo: String,
    val idRuta: String,
    val estado: EstadoViaje,
    val horaInicio: Long,
    val horaFin: Long
)