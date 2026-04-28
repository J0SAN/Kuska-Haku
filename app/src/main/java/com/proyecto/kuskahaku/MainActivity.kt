package com.proyecto.kuskahaku

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.ui.Modifier
import com.proyecto.kuskahaku.presentation.theme.KuskaHakuTheme
import com.proyecto.kuskahaku.presentation.ui.KuskaApp // <-- Tu nuevo import

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        setContent {
            KuskaHakuTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    // ¡Llamamos al contenedor principal!
                    KuskaApp()
                }
            }
        }
    }
}