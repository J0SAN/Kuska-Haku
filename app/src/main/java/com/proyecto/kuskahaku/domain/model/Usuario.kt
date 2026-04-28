package com.proyecto.kuskahaku.domain.model

enum class RolUsuario {
    PASAJERO, CONDUCTOR, ADMIN
}

data class Usuario (
    val id: String,
    val nombreCompleto: String,
    val correo: String,
    val telefono: String,
    val rol: RolUsuario
)