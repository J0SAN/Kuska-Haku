package com.proyecto.kuskahaku.domain.model

data class Paradero (
    val id: String,
    val idRuta: String,
    val nombre: String,
    val longitud: Double,
    val latitud: Double,
    val ordenSecuencia: Int
)