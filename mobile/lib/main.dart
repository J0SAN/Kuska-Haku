import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

void main() {
  runApp(const KuskaHakuApp());
}

class KuskaHakuApp extends StatelessWidget {
  const KuskaHakuApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Kuska Haku Pasajero',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.green),
        useMaterial3: true,
      ),
      home: const MapScreen(),
    );
  }
}

class MapScreen extends StatefulWidget {
  const MapScreen({super.key});

  @override
  State<MapScreen> createState() => _MapScreenState();
}

class _MapScreenState extends State<MapScreen> {
  // Ajuste de IP según el README: 10.0.2.2:8080 para Emulador Android
  final String _backendHost = '10.254.83.86:8080';
  late WebSocketChannel _channel;

  // Centro inicial en Cajamarca
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
    final marker = Marker(
      markerId: markerId,
      position: LatLng(lat, lng),
      // Marcador verde como pide la documentación
      icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueGreen),
    );

    setState(() {
      _markers[markerId] = marker;
    });
  }

  @override
  void dispose() {
    _channel.sink.close();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Kuska Haku - Fase 1'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ),
      body: GoogleMap(
        initialCameraPosition: _initialPosition,
        markers: Set<Marker>.of(_markers.values),
        onMapCreated: (GoogleMapController controller) {},
      ),
    );
  }
}
