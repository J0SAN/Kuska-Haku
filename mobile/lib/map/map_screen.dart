import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

import '../auth/api_client.dart';
import '../config.dart';
import '../transport_routes/transport_routes_api.dart';

class LiveMapScreen extends StatefulWidget {
  final AuthResult session;
  final VoidCallback? onLogout;

  const LiveMapScreen({
    super.key,
    required this.session,
    this.onLogout,
  });

  @override
  State<LiveMapScreen> createState() => _LiveMapScreenState();
}

class _LiveMapScreenState extends State<LiveMapScreen> {
  static const _cajamarcaCenter = LatLng(-7.1611, -78.5127);
  static const _greenColor = Color(0xFF2E7D32);

  late final RoutesApi _api;
  GoogleMapController? _mapController;
  WebSocketChannel? _channel;
  Timer? _etaTimer;

  RouteInfo? _route;
  final Map<String, Marker> _markers = {};
  List<BusEta> _etas = [];
  String _status = 'Cargando ruta...';

  @override
  void initState() {
    super.initState();
    _api = RoutesApi(token: widget.session.token);
    _initialize();
  }

  Future<void> _initialize() async {
    try {
      final route = await _api.getActiveRoute();
      if (!mounted) {
        return;
      }
      if (route == null) {
        setState(() => _status = 'No hay rutas activas por el momento');
        return;
      }
      setState(() {
        _route = route;
        _status = 'Esperando buses en ${route.name}';
      });
      _fitCamera();
      _connectWebSocket(route);
      await _updateEta();
      _etaTimer = Timer.periodic(
        const Duration(seconds: 5),
        (_) => _updateEta(),
      );
    } on RoutesApiException catch (e) {
      if (mounted) {
        setState(() => _status = e.message);
      }
    } catch (_) {
      if (mounted) {
        setState(() => _status = 'No se pudo conectar con el backend');
      }
    }
  }

  void _connectWebSocket(RouteInfo route) {
    final encodedRoute = Uri.encodeComponent(route.name);
    final uri = Uri.parse(
      'ws://${AppConfig.backendHost}/ws/routes/$encodedRoute',
    ).replace(queryParameters: {'token': widget.session.token});

    _channel = WebSocketChannel.connect(uri);
    _channel!.stream.listen(
      (raw) => _onPosition(jsonDecode(raw as String) as Map<String, dynamic>),
      onError: (_) {
        if (mounted) {
          setState(() => _status = 'Conexión en vivo interrumpida');
        }
      },
      onDone: () {
        if (mounted) {
          setState(() => _status = 'Desconectado del servidor');
        }
      },
    );
  }

  void _onPosition(Map<String, dynamic> data) {
    final busId = (data['bus_id'] ?? data['busId']) as String;
    final position = LatLng(
      ((data['latitude'] ?? data['lat']) as num).toDouble(),
      ((data['longitude'] ?? data['lng']) as num).toDouble(),
    );
    final heading = (data['heading'] as num?)?.toDouble() ?? 0;

    setState(() {
      _status = 'Bus $busId en movimiento';
      _markers[busId] = Marker(
        markerId: MarkerId(busId),
        position: position,
        rotation: heading,
        anchor: const Offset(0.5, 0.5),
        flat: true,
        icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueGreen),
        infoWindow: InfoWindow(title: busId),
      );
    });
  }

  Future<void> _updateEta() async {
    final route = _route;
    if (route == null) {
      return;
    }
    try {
      final etas = await _api.getEta(route.id);
      if (mounted) {
        setState(() => _etas = etas);
      }
    } catch (_) {
      // Fallo puntual de red: se conserva el último ETA mostrado.
    }
  }

  void _fitCamera() {
    final route = _route;
    final controller = _mapController;
    if (route == null || controller == null || route.path.length < 2) {
      return;
    }

    double minLat = route.path.first.latitude, maxLat = minLat;
    double minLng = route.path.first.longitude, maxLng = minLng;
    for (final p in route.path) {
      if (p.latitude < minLat) {
        minLat = p.latitude;
      }
      if (p.latitude > maxLat) {
        maxLat = p.latitude;
      }
      if (p.longitude < minLng) {
        minLng = p.longitude;
      }
      if (p.longitude > maxLng) {
        maxLng = p.longitude;
      }
    }
    try {
      controller.animateCamera(
        CameraUpdate.newLatLngBounds(
          LatLngBounds(
            southwest: LatLng(minLat, minLng),
            northeast: LatLng(maxLat, maxLng),
          ),
          64,
        ),
      );
    } catch (_) {}
  }

  String _formatEta(int seconds) {
    if (seconds < 60) {
      return '$seconds s';
    }
    final min = seconds ~/ 60;
    final sec = seconds % 60;
    return sec == 0 ? '$min min' : '$min min $sec s';
  }

  Future<void> _openSuggestionDialog() async {
    String draftText = '';
    final message = await showDialog<String>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Enviar sugerencia'),
        content: TextField(
          autofocus: true,
          maxLines: 4,
          maxLength: 1000,
          onChanged: (value) => draftText = value,
          decoration: const InputDecoration(
            hintText: 'Cuéntanos qué podemos mejorar',
            border: OutlineInputBorder(),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cancelar'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, draftText),
            child: const Text('Enviar'),
          ),
        ],
      ),
    );

    if (message == null) {
      return;
    }

    String feedback = 'Gracias por tu sugerencia';
    try {
      await _api.sendSuggestion(message);
    } on RoutesApiException catch (e) {
      feedback = e.message;
    } catch (_) {
      feedback = 'No se pudo enviar la sugerencia';
    }
    if (!mounted) {
      return;
    }
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(feedback)),
    );
  }

  @override
  void dispose() {
    _etaTimer?.cancel();
    _channel?.sink.close();
    _mapController?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final route = _route;
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.session.user.fullName),
        backgroundColor: _greenColor,
        foregroundColor: Colors.white,
        actions: [
          IconButton(
            tooltip: 'Enviar sugerencia',
            icon: const Icon(Icons.rate_review_outlined),
            onPressed: _openSuggestionDialog,
          ),
          if (widget.onLogout != null)
            IconButton(
              tooltip: 'Cerrar sesión',
              icon: const Icon(Icons.logout),
              onPressed: widget.onLogout,
            ),
        ],
      ),
      body: Stack(
        children: [
          GoogleMap(
            initialCameraPosition: const CameraPosition(
              target: _cajamarcaCenter,
              zoom: 14,
            ),
            markers: _markers.values.toSet(),
            polylines: route == null
                ? {}
                : {
                    Polyline(
                      polylineId: PolylineId(route.id),
                      points: route.path,
                      color: _greenColor,
                      width: 5,
                    ),
                  },
            onMapCreated: (controller) {
              _mapController = controller;
              _fitCamera();
            },
          ),
          Positioned(
            top: 12,
            left: 12,
            right: 12,
            child: Card(
              elevation: 3,
              child: Padding(
                padding: const EdgeInsets.all(10),
                child: Text(_status),
              ),
            ),
          ),
          if (_etas.isNotEmpty)
            Positioned(
              left: 12,
              right: 12,
              bottom: 12,
              child: Card(
                elevation: 4,
                child: Padding(
                  padding: const EdgeInsets.all(12),
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Próximas llegadas',
                        style: TextStyle(
                          fontWeight: FontWeight.bold,
                          color: _greenColor,
                        ),
                      ),
                      const SizedBox(height: 6),
                      for (final bus in _etas) _buildEtaRow(bus),
                    ],
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildEtaRow(BusEta bus) {
    final next = bus.nextStop;
    final text = next == null
        ? '${bus.busId}: recorrido completado'
        : '${bus.busId} llega a ${next.name} en '
            '${_formatEta(next.etaSeconds ?? 0)} '
            '(${next.distanceMeters.round()} m)';
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        children: [
          const Icon(Icons.directions_bus, size: 18, color: _greenColor),
          const SizedBox(width: 8),
          Expanded(child: Text(text)),
        ],
      ),
    );
  }
}
