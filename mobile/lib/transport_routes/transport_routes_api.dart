import 'dart:convert';

import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:http/http.dart' as http;

import '../config.dart';

class RoutesApiException implements Exception {
  final String message;
  RoutesApiException(this.message);

  @override
  String toString() => message;
}

class RouteInfo {
  final String id;
  final String name;
  final List<LatLng> path;

  RouteInfo({
    required this.id,
    required this.name,
    required this.path,
  });

  factory RouteInfo.fromJson(Map<String, dynamic> json) => RouteInfo(
        id: json['id'] as String,
        name: json['name'] as String,
        path: ((json['path'] as List<dynamic>?) ?? [])
            .map(
              (p) => LatLng(
                ((p['latitude'] ?? p['lat']) as num).toDouble(),
                ((p['longitude'] ?? p['lng']) as num).toDouble(),
              ),
            )
            .toList(),
      );
}

class StopEta {
  final String name;
  final int stopOrder;
  final double distanceMeters;
  final int? etaSeconds;
  final bool passed;

  StopEta({
    required this.name,
    required this.stopOrder,
    required this.distanceMeters,
    required this.etaSeconds,
    required this.passed,
  });

  factory StopEta.fromJson(Map<String, dynamic> json) => StopEta(
        name: json['name'] as String,
        stopOrder: json['stopOrder'] as int,
        distanceMeters: (json['distanceMeters'] as num).toDouble(),
        etaSeconds: json['etaSeconds'] as int?,
        passed: json['passed'] as bool,
      );
}

class BusEta {
  final String busId;
  final bool offRoute;
  final List<StopEta> stops;

  BusEta({
    required this.busId,
    required this.offRoute,
    required this.stops,
  });

  factory BusEta.fromJson(Map<String, dynamic> json) => BusEta(
        busId: json['busId'] as String,
        offRoute: (json['offRoute'] as bool?) ??
            (json['fueraDeRuta'] as bool?) ??
            false,
        stops: ((json['stops'] as List<dynamic>?) ?? [])
            .map((s) => StopEta.fromJson(s as Map<String, dynamic>))
            .toList(),
      );

  /// Primer paradero que el bus todavía no ha pasado.
  StopEta? get nextStop {
    for (final s in stops) {
      if (!s.passed) return s;
    }
    return null;
  }
}

class RoutesApi {
  final String token;
  RoutesApi({required this.token});

  Map<String, String> get _headers => {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      };

  Uri _uri(String path) => Uri.parse('http://${AppConfig.backendHost}$path');

  void _verify(http.Response res, {int expected = 200}) {
    if (res.statusCode == expected) return;
    String message = 'Error ${res.statusCode}';
    try {
      final body = jsonDecode(res.body) as Map<String, dynamic>;
      message = (body['error'] ?? body['message'] ?? message) as String;
    } catch (_) {}
    throw RoutesApiException(message);
  }

  /// Obtiene la primera ruta activa disponible.
  Future<RouteInfo?> getActiveRoute() async {
    final res = await http.get(_uri('/api/v1/routes'), headers: _headers);
    _verify(res);
    final list = jsonDecode(res.body) as List<dynamic>;
    if (list.isEmpty) return null;
    return RouteInfo.fromJson(list.first as Map<String, dynamic>);
  }

  Future<List<BusEta>> getEta(String routeId) async {
    final res = await http.get(
      _uri('/api/v1/routes/$routeId/eta'),
      headers: _headers,
    );
    _verify(res);
    final body = jsonDecode(res.body) as Map<String, dynamic>;
    return ((body['buses'] as List<dynamic>?) ?? [])
        .map((b) => BusEta.fromJson(b as Map<String, dynamic>))
        .toList();
  }

  Future<void> sendSuggestion(String message) async {
    final res = await http.post(
      _uri('/api/v1/suggestions'),
      headers: _headers,
      body: jsonEncode({'message': message}),
    );
    _verify(res, expected: 201);
  }
}
