import 'dart:convert';

import 'package:flutter/material.dart';

import 'package:google_maps_flutter/google_maps_flutter.dart';

import 'package:web_socket_channel/web_socket_channel.dart';

import 'package:shared_preferences/shared_preferences.dart';

import 'auth/login_screen.dart';

void main() async {
  // Obligatorio para ejecutar código asíncrono nativo antes de arrancar la app

  WidgetsFlutterBinding.ensureInitialized();

  // Verificamos la sesión directamente en el punto de entrada

  final prefs = await SharedPreferences.getInstance();

  final token = prefs.getString('jwt_token');

  final initialRoute = (token != null) ? '/map' : '/login';

  runApp(KuskaHakuApp(initialRoute: initialRoute));
}

class KuskaHakuApp extends StatelessWidget {
  final String initialRoute;

  const KuskaHakuApp({super.key, required this.initialRoute});

  @override
  Widget build(BuildContext context) {
    // Un único MaterialApp construido desde el inicio

    return MaterialApp(
      title: 'Kuska Haku Pasajero',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.green),
        useMaterial3: true,
      ),
      initialRoute: initialRoute,
      routes: {
        '/login': (context) => const LoginScreen(),
        '/map': (context) => const MapScreen(),
      },
    );
  }
}

// ============================================================================

// TU CÓDIGO INTACTO DE LA FASE 1 A PARTIR DE AQUÍ

// ============================================================================

class MapScreen extends StatefulWidget {
  const MapScreen({super.key});

  @override
  State<MapScreen> createState() => _MapScreenState();
}

class _MapScreenState extends State<MapScreen> {
  final String _backendHost = '10.118.10.86:8080';

  late WebSocketChannel _channel;

  GoogleMapController? _mapController; // Controlador para mover la cámara

  // Centro inicial temporal (se actualizará al recibir el primer dato)

  static const CameraPosition _initialPosition = CameraPosition(
    target: LatLng(-7.162, -78.500),
    zoom: 15,
  );

  final Map<MarkerId, Marker> _markers = {};

  @override
  void initState() {
    super.initState();

    _connectWebSocket();
  }

  void _connectWebSocket() {
    _channel = WebSocketChannel.connect(Uri.parse('ws://$_backendHost/ws'));

    _channel.stream.listen((message) {
      final data = jsonDecode(message);

      _updateMarker(data['bus_id'], data['latitude'], data['longitude']);
    });
  }

  void _updateMarker(String busId, double lat, double lng) {
    final markerId = MarkerId(busId);

    final position = LatLng(lat, lng);

    final marker = Marker(
      markerId: markerId,
      position: position,
      icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueGreen),
    );

    setState(() {
      _markers[markerId] = marker;
    });

    // Mueve la cámara automáticamente hacia donde está el bus

    _mapController?.animateCamera(CameraUpdate.newLatLng(position));
  }

  @override
  void dispose() {
    _channel.sink.close();

    _mapController?.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Kuska Haku'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () async {
              final prefs = await SharedPreferences.getInstance();

              await prefs.remove('jwt_token');

              if (!context.mounted) return;

              Navigator.pushReplacementNamed(context, '/login');
            },
          )
        ],
      ),
      body: GoogleMap(
        initialCameraPosition: _initialPosition,

        markers: Set<Marker>.of(_markers.values),

        // Guardamos la referencia del controlador cuando el mapa termina de cargar

        onMapCreated: (GoogleMapController controller) {
          _mapController = controller;
        },
      ),
    );
  }
}
