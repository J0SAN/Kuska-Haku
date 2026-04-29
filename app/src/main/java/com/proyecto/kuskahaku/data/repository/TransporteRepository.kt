package com.proyecto.kuskahaku.data.repository

import com.proyecto.kuskahaku.data.remote.BusFirebaseDto
import com.proyecto.kuskahaku.data.remote.FirebaseDataSource
import kotlinx.coroutines.flow.Flow

class TransporteRepository {
    // Instanciamos nuestra fuente de datos
    private val firebaseDataSource = FirebaseDataSource()

    // En un proyecto más avanzado, aquí "traduciríamos" el BusFirebaseDto
    // a un modelo puro de tu carpeta domain.model. Por ahora, pasaremos
    // el flujo directamente para ver la conexión funcionar.
    fun obtenerBusesP13(): Flow<List<BusFirebaseDto>> {
        return firebaseDataSource.getBusesActivos()
    }
}