import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:kuska_haku_pasajero/auth/api_client.dart';
import 'package:kuska_haku_pasajero/auth/registration.dart';

AuthApiClient _client(MockClient mock) =>
    AuthApiClient(backendHost: 'localhost:8080', client: mock);

RegisterInput _input() => const RegisterInput(
      dni: '12345678',
      checkDigit: '1',
      firstNames: 'Roxana Karina',
      lastNames: 'Delgado Huamán',
      email: 'roxana@correo.com',
      password: 'secreto123',
      acceptTerms: true,
    );

http.Response _authResponse() => http.Response(
      jsonEncode({
        'message': 'Usuario creado exitosamente',
        'token': 'token-1',
        'user': {
          'id': 'user-1',
          'email': 'roxana@correo.com',
          'fullName': 'Roxana Karina Delgado Huamán',
          'role': 'passenger',
        },
      }),
      201,
      headers: {'content-type': 'application/json'},
    );

void main() {
  group('AuthApiClient.register', () {
    test('201 devuelve AuthResult y envía el body correcto', () async {
      late http.Request captured;
      final mock = MockClient((request) async {
        captured = request;
        return _authResponse();
      });

      final result = await _client(mock).register(_input());

      expect(result.token, 'token-1');
      expect(result.user.role, 'passenger');
      expect(captured.url.path, '/api/v1/auth/register');

      final body = jsonDecode(captured.body) as Map<String, dynamic>;
      expect(body['dni'], '12345678');
      expect(body['checkDigit'], '1');
      expect(body['acceptTerms'], isTrue);
      expect(body.containsKey('role'), isFalse);
      expect(body.containsKey('fullName'), isFalse);
    });

    for (final status in [400, 409, 429]) {
      test('$status lanza AuthException con el mensaje del backend', () async {
        final mock = MockClient(
          (_) async => http.Response(
            jsonEncode({
              'error': 'mensaje-$status',
              'message': 'mensaje-$status',
            }),
            status,
            headers: {'content-type': 'application/json'},
          ),
        );

        try {
          await _client(mock).register(_input());
          fail('debía lanzar AuthException');
        } on AuthException catch (e) {
          expect(e.statusCode, status);
          expect(e.message, 'mensaje-$status');
        }
      });
    }

    test('error de red se traduce a un mensaje de conexión', () async {
      final mock = MockClient((_) async => throw Exception('sin red'));

      await expectLater(
        _client(mock).register(_input()),
        throwsA(
          isA<AuthException>().having(
            (e) => e.message,
            'message',
            contains('No pudimos conectar'),
          ),
        ),
      );
    });
  });

  group('AuthApiClient.login', () {
    test('200 devuelve AuthResult', () async {
      final mock = MockClient(
        (_) async => http.Response(
          jsonEncode({
            'token': 'token-login',
            'user': {
              'id': 'user-2',
              'email': 'pasajero@correo.com',
              'fullName': 'Pasajero Uno',
              'role': 'passenger',
            },
          }),
          200,
          headers: {'content-type': 'application/json'},
        ),
      );

      final result = await _client(mock).login('pasajero@correo.com', 'x');

      expect(result.token, 'token-login');
      expect(result.user.email, 'pasajero@correo.com');
    });

    test('401 lanza AuthException con statusCode', () async {
      final mock = MockClient(
        (_) async => http.Response(
          jsonEncode({'error': 'Credenciales inválidas', 'message': 'Credenciales inválidas'}),
          401,
          headers: {'content-type': 'application/json'},
        ),
      );

      await expectLater(
        _client(mock).login('nadie@correo.com', 'x'),
        throwsA(
          isA<AuthException>()
              .having((e) => e.statusCode, 'statusCode', 401)
              .having((e) => e.message, 'message', 'Credenciales inválidas'),
        ),
      );
    });
  });
}
