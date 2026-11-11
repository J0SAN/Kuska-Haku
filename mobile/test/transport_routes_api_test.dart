import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:kuska_haku_pasajero/transport_routes/transport_routes_api.dart';

RoutesApi _api(MockClient client) =>
    RoutesApi(token: 'token-prueba', client: client);

void main() {
  group('RoutesApi.listRoutes', () {
    test('devuelve todas las rutas de la respuesta', () async {
      final client = MockClient((request) async {
        expect(request.url.path, '/api/v1/routes');
        expect(request.headers['Authorization'], 'Bearer token-prueba');
        return http.Response(
          jsonEncode([
            {
              'id': 'ruta-01',
              'name': 'Ruta 01',
              'path': [
                {'latitude': -7.1611, 'longitude': -78.5127},
              ],
            },
            {
              'id': 'ruta-02',
              'name': 'Ruta 02',
              'path': <Map<String, dynamic>>[],
            },
          ]),
          200,
          headers: {'content-type': 'application/json'},
        );
      });

      final routes = await _api(client).listRoutes();

      expect(routes, hasLength(2));
      expect(routes.first.id, 'ruta-01');
      expect(routes.first.name, 'Ruta 01');
      expect(routes.first.path, hasLength(1));
      expect(routes.last.id, 'ruta-02');
    });

    test('devuelve una lista vacía cuando no hay rutas', () async {
      final client = MockClient(
        (_) async => http.Response('[]', 200),
      );

      expect(await _api(client).listRoutes(), isEmpty);
    });

    test('lanza RoutesApiException con el mensaje del backend', () async {
      final client = MockClient(
        (_) async => http.Response(
          jsonEncode({'error': 'No autorizado'}),
          401,
        ),
      );

      expect(
        () => _api(client).listRoutes(),
        throwsA(
          isA<RoutesApiException>().having(
            (error) => error.message,
            'message',
            'No autorizado',
          ),
        ),
      );
    });
  });

  group('RoutesApi.getActiveRoute', () {
    test('devuelve la primera ruta de la lista', () async {
      final client = MockClient(
        (_) async => http.Response(
          jsonEncode([
            {
              'id': 'ruta-01',
              'name': 'Ruta 01',
              'path': <Map<String, dynamic>>[]
            },
            {
              'id': 'ruta-02',
              'name': 'Ruta 02',
              'path': <Map<String, dynamic>>[]
            },
          ]),
          200,
        ),
      );

      final route = await _api(client).getActiveRoute();

      expect(route?.id, 'ruta-01');
    });

    test('devuelve null cuando la lista está vacía', () async {
      final client = MockClient(
        (_) async => http.Response('[]', 200),
      );

      expect(await _api(client).getActiveRoute(), isNull);
    });
  });
}
