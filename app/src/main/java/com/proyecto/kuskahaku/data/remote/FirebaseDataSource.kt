package com.proyecto.kuskahaku.data.remote

import android.util.Log
import com.google.firebase.database.DataSnapshot
import com.google.firebase.database.DatabaseError
import com.google.firebase.database.FirebaseDatabase
import com.google.firebase.database.ValueEventListener
import kotlinx.coroutines.channels.awaitClose
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.callbackFlow

// 1. Creamos un DTO (Data Transfer Object) específico para Firebase.
// Nota: Firebase exige que las variables tengan valores por defecto (ej. = 0.0)
// para poder construir el objeto automáticamente desde el JSON.
data class BusFirebaseDto(
    val latitud: Double = 0.0,
    val longitud: Double = 0.0,
    val estado: String = "",
    val ultima_actualizacion: Long = 0L
)

class FirebaseDataSource {
    // 2. Apuntamos directamente a la "carpeta" de nuestra Línea P13
    private val databaseRef = FirebaseDatabase.getInstance().getReference("transporte/linea_P13/buses_activos")

    // 3. Usamos un Flow (Flujo de datos) porque es información en Tiempo Real
    fun getBusesActivos(): Flow<List<BusFirebaseDto>> = callbackFlow {

        // Creamos el "escuchador" que estará atento a cualquier cambio en la base de datos
        val listener = object : ValueEventListener {

            // Esta función se dispara SOLA cada vez que un dato cambia en Firebase
            override fun onDataChange(snapshot: DataSnapshot) {
                val listaBuses = mutableListOf<BusFirebaseDto>()

                // Recorremos todos los "hijos" (bus_01, bus_02, etc.)
                for (busSnapshot in snapshot.children) {
                    // Traducimos el JSON mágico de Firebase a nuestra clase Kotlin
                    val bus = busSnapshot.getValue(BusFirebaseDto::class.java)
                    if (bus != null) {
                        listaBuses.add(bus)
                    }
                }

                // Enviamos la lista fresca a quien sea que esté escuchando (el ViewModel)
                trySend(listaBuses)
            }

            override fun onCancelled(error: DatabaseError) {
                // Si falla la conexión o no hay permisos, enviamos el error
                Log.e("FirebaseData", "Error al leer datos: ${error.message}")
                close(error.toException())
            }
        }

        // Atachamos el escuchador a nuestra referencia de Firebase
        databaseRef.addValueEventListener(listener)

        // IMPORTANTE: Cuando el usuario cierre la pantalla, destruimos el escuchador
        // para no gastar sus datos móviles ni la batería de su teléfono.
        awaitClose {
            databaseRef.removeEventListener(listener)
        }
    }
}
