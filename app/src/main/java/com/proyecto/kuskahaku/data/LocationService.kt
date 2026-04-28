package com.proyecto.kuskahaku.data

import android.annotation.SuppressLint
import android.content.Context
import android.os.Looper
import android.util.Log
import com.google.android.gms.location.LocationCallback
import com.google.android.gms.location.LocationRequest
import com.google.android.gms.location.LocationResult
import com.google.android.gms.location.LocationServices
import com.google.android.gms.location.Priority
import kotlinx.coroutines.channels.awaitClose
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.callbackFlow

class LocationService(private val context: Context) {

    // Se suprime la advertencia porque pediremos el permiso en la vista antes de llamar a esto
    @SuppressLint("MissingPermission")
    fun obtenerUbicacionEnTiempoReal(): Flow<Pair<Double, Double>> = callbackFlow {
        val client = LocationServices.getFusedLocationProviderClient(context)

        // Configuramos alta precisión y que nos actualice cada 5 segundos
        val locationRequest = LocationRequest.Builder(Priority.PRIORITY_HIGH_ACCURACY, 5000L)
            .setMinUpdateIntervalMillis(2000L)
            .build()

        val locationCallback = object : LocationCallback() {
            override fun onLocationResult(result: LocationResult) {
                result.lastLocation?.let { location ->
                    val lat = location.latitude
                    val lng = location.longitude

                    // Imprimimos en consola para demostrar que funciona
                    Log.d("GPS_KUSKA", "Ubicación real capturada -> Lat: $lat, Lng: $lng")

                    // Emitimos los datos al ViewModel
                    trySend(Pair(lat, lng))
                }
            }
        }

        // Encendemos el GPS
        client.requestLocationUpdates(locationRequest, locationCallback, Looper.getMainLooper())

        // Apagamos el GPS si el flujo se cancela (ahorro de batería)
        awaitClose {
            client.removeLocationUpdates(locationCallback)
        }
    }
}