import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

import '../auth/api_client.dart';
import '../config.dart';
import '../theme/app_colors.dart';
import '../transport_routes/transport_routes_api.dart';
import 'navigation/passenger_drawer.dart';
import 'navigation/passenger_help_dialog.dart';
import 'navigation/passenger_menu.dart';
import 'navigation/passenger_profile_sheet.dart';
import 'navigation/passenger_route_picker.dart';
import 'navigation/passenger_route_stops_sheet.dart';
import 'navigation/passenger_suggestion_dialog.dart';
import 'widgets/eta_card.dart';

/// Mapa en vivo del pasajero. Orquesta el estado (ruta, WebSocket, ETA y
/// selección del menú) y delega la UI secundaria a `navigation/` y `widgets/`.
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

  late final RoutesApi _api;
  GoogleMapController? _mapController;
  WebSocketChannel? _channel;
  Timer? _etaTimer;

  RouteInfo? _route;
  final Map<String, Marker> _markers = {};
  List<BusEta> _etas = [];
  String _status = 'Cargando ruta...';
  PassengerMenuItem _selectedItem = PassengerMenuItem.liveMap;

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
      await _switchRoute(route);
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

  /// Cambia la ruta activa cerrando la conexión anterior antes de abrir la nueva.
  Future<void> _switchRoute(RouteInfo route) async {
    _etaTimer?.cancel();
    _etaTimer = null;

    final previousChannel = _channel;
    _channel = null;
    await previousChannel?.sink.close();

    if (!mounted) {
      return;
    }

    setState(() {
      _markers.clear();
      _etas = [];
      _route = route;
      _status = 'Esperando buses en ${route.name}';
    });

    _fitCamera();
    _connectWebSocket(route);
    await _updateEta();

    if (!mounted || _route?.id != route.id) {
      return;
    }
    _etaTimer = Timer.periodic(
      const Duration(seconds: 5),
      (_) => _updateEta(),
    );
  }

  void _connectWebSocket(RouteInfo route) {
    // El canal en Redis es `route:<routeID>:positions` y la ingesta usa el UUID
    // de la ruta; suscribirse por nombre abre un canal donde nadie publica.
    final encodedRoute = Uri.encodeComponent(route.id);
    final uri = Uri.parse(
      '${AppConfig.wsScheme}://${AppConfig.backendHost}/ws/routes/$encodedRoute',
    ).replace(queryParameters: {'token': widget.session.token});

    final channel = WebSocketChannel.connect(uri);
    _channel = channel;
    channel.stream.listen(
      (raw) {
        if (!identical(_channel, channel)) {
          return;
        }
        _onPosition(jsonDecode(raw as String) as Map<String, dynamic>);
      },
      onError: (_) {
        if (mounted && identical(_channel, channel)) {
          setState(() => _status = 'Conexión en vivo interrumpida');
        }
      },
      onDone: () {
        if (mounted && identical(_channel, channel)) {
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

    if (!mounted) {
      return;
    }
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
      if (mounted && _route?.id == route.id) {
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

  Future<void> _openSuggestionDialog() async {
    final message = await showPassengerSuggestionDialog(context);
    if (message == null || !mounted) {
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

  Future<void> _showRoutePicker() async {
    final selectedRoute = await showPassengerRoutePicker(
      context,
      routes: _api.listRoutes(),
      currentRouteId: _route?.id,
    );
    if (!mounted || selectedRoute == null) {
      return;
    }
    if (selectedRoute.id == _route?.id) {
      return;
    }
    await _switchRoute(selectedRoute);
  }

  void _onMenuSelected(PassengerMenuItem item) {
    switch (item) {
      case PassengerMenuItem.liveMap:
        setState(() => _selectedItem = PassengerMenuItem.liveMap);
        _fitCamera();
      case PassengerMenuItem.arrivals:
        showRouteStopsSheet(context, route: _route, etas: _etas);
      case PassengerMenuItem.routes:
        setState(() => _selectedItem = PassengerMenuItem.routes);
        _showRoutePicker();
      case PassengerMenuItem.suggestions:
        _openSuggestionDialog();
      case PassengerMenuItem.help:
        showPassengerHelpDialog(context);
    }
  }

  @override
  void dispose() {
    _etaTimer?.cancel();
    final channel = _channel;
    _channel = null;
    channel?.sink.close();
    _mapController?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final route = _route;
    return Scaffold(
      drawer: PassengerDrawer(
        user: widget.session.user,
        selected: _selectedItem,
        onSelect: _onMenuSelected,
        onProfileTap: () =>
            showPassengerProfileSheet(context, widget.session.user),
        onLogout: widget.onLogout,
      ),
      appBar: AppBar(
        title: Text(_route?.name ?? 'Kuska Haku'),
        backgroundColor: AppColors.primary,
        foregroundColor: AppColors.onPrimary,
        actions: [
          IconButton(
            tooltip: 'Centrar ruta',
            icon: const Icon(Icons.my_location),
            onPressed: _fitCamera,
          ),
          IconButton(
            tooltip: 'Sugerencia general',
            icon: const Icon(Icons.rate_review_outlined),
            onPressed: _openSuggestionDialog,
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
                      color: AppColors.primary,
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
              child: EtaCard(etas: _etas),
            ),
        ],
      ),
    );
  }
}
