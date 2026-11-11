import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:kuska_haku_pasajero/auth/api_client.dart';
import 'package:kuska_haku_pasajero/auth/register_screen.dart';

class _Harness extends StatelessWidget {
  final AuthApiClient apiClient;
  final void Function(AuthResult) onRegistered;

  const _Harness({required this.apiClient, required this.onRegistered});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: Scaffold(
        body: Builder(
          builder: (context) => Center(
            child: ElevatedButton(
              key: const Key('open-register'),
              onPressed: () => Navigator.of(context).push(
                MaterialPageRoute(
                  builder: (_) => RegisterScreen(
                    apiClient: apiClient,
                    onRegistered: onRegistered,
                  ),
                ),
              ),
              child: const Text('Abrir registro'),
            ),
          ),
        ),
      ),
    );
  }
}

AuthApiClient _client(MockClient mock) =>
    AuthApiClient(backendHost: 'localhost:8080', client: mock);

Future<void> _open(
  WidgetTester tester,
  AuthApiClient client, {
  void Function(AuthResult)? onRegistered,
}) async {
  tester.view.physicalSize = const Size(1200, 3000);
  tester.view.devicePixelRatio = 1.0;
  addTearDown(tester.view.resetPhysicalSize);
  addTearDown(tester.view.resetDevicePixelRatio);

  await tester.pumpWidget(
    _Harness(apiClient: client, onRegistered: onRegistered ?? (_) {}),
  );
  await tester.tap(find.byKey(const Key('open-register')));
  await tester.pumpAndSettle();
}

Future<void> _fillValidForm(WidgetTester tester, {bool acceptTerms = true}) async {
  await tester.enterText(find.byKey(const Key('dni-field')), '12345678');
  await tester.enterText(find.byKey(const Key('check-digit-field')), '1');
  await tester.enterText(
    find.byKey(const Key('first-names-field')),
    'Roxana Karina',
  );
  await tester.enterText(
    find.byKey(const Key('last-names-field')),
    'Delgado Huamán',
  );
  await tester.enterText(
    find.byKey(const Key('email-field')),
    'roxana@correo.com',
  );
  await tester.enterText(
    find.byKey(const Key('password-field')),
    'secreto123',
  );
  await tester.enterText(
    find.byKey(const Key('confirm-password-field')),
    'secreto123',
  );
  if (acceptTerms) {
    await tester.tap(find.byKey(const Key('accept-terms')));
  }
  await tester.pump();
}

ElevatedButton _submitButton(WidgetTester tester) =>
    tester.widget<ElevatedButton>(find.byKey(const Key('submit-button')));

http.Response _createdResponse() => http.Response(
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
  testWidgets('el botón está deshabilitado hasta tener todo válido', (
    tester,
  ) async {
    await _open(tester, _client(MockClient((_) async => _createdResponse())));

    expect(_submitButton(tester).onPressed, isNull);

    await _fillValidForm(tester, acceptTerms: false);
    expect(_submitButton(tester).onPressed, isNull);

    await tester.tap(find.byKey(const Key('accept-terms')));
    await tester.pump();
    expect(_submitButton(tester).onPressed, isNotNull);
  });

  testWidgets('un dígito incorrecto muestra el error en vivo', (tester) async {
    await _open(tester, _client(MockClient((_) async => _createdResponse())));

    await tester.enterText(find.byKey(const Key('dni-field')), '12345678');
    await tester.enterText(find.byKey(const Key('check-digit-field')), '2');
    await tester.pump();

    expect(
      find.text('El dígito verificador no coincide. Revisa el número de tu DNI.'),
      findsOneWidget,
    );
  });

  testWidgets('un 409 de DNI muestra el enlace Iniciar sesión', (tester) async {
    final mock = MockClient(
      (_) async => http.Response(
        jsonEncode({
          'error': 'Ya existe una cuenta con este DNI.',
          'message':
              'Ya existe una cuenta con este DNI. Si no la creaste tú, contacta a soporte.',
        }),
        409,
        headers: {'content-type': 'application/json'},
      ),
    );

    await _open(tester, _client(mock));
    await _fillValidForm(tester);
    await tester.tap(find.byKey(const Key('submit-button')));
    await tester.pumpAndSettle();

    expect(
      find.text(
        'Ya existe una cuenta con este DNI. Si no la creaste tú, contacta a soporte.',
      ),
      findsOneWidget,
    );
    expect(find.text('Iniciar sesión'), findsOneWidget);
  });

  testWidgets('al registrarse se llama a onRegistered y se cierra', (
    tester,
  ) async {
    AuthResult? registered;

    await _open(
      tester,
      _client(MockClient((_) async => _createdResponse())),
      onRegistered: (result) => registered = result,
    );
    await _fillValidForm(tester);
    await tester.tap(find.byKey(const Key('submit-button')));
    await tester.pumpAndSettle();

    expect(registered, isNotNull);
    expect(registered!.user.role, 'passenger');
    expect(find.byType(RegisterScreen), findsNothing);
  });
}
