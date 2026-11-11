import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:kuska_haku_pasajero/auth/api_client.dart';
import 'package:kuska_haku_pasajero/map/navigation/passenger_profile_sheet.dart';

AuthUser _user({String role = 'passenger'}) => AuthUser(
      id: 'user-1',
      email: 'jose@correo.com',
      fullName: 'José Pérez',
      role: role,
    );

void main() {
  testWidgets('muestra nombre, email, inicial y rol legible', (tester) async {
    await tester.pumpWidget(
      MaterialApp(
        home: Scaffold(body: PassengerProfileSheet(user: _user())),
      ),
    );

    expect(find.text('José Pérez'), findsOneWidget);
    expect(find.text('jose@correo.com'), findsOneWidget);
    expect(find.text('J'), findsOneWidget);
    expect(find.text('Rol activo: Pasajero'), findsOneWidget);
  });

  testWidgets('traduce el rol de conductor', (tester) async {
    await tester.pumpWidget(
      MaterialApp(
        home:
            Scaffold(body: PassengerProfileSheet(user: _user(role: 'driver'))),
      ),
    );

    expect(find.text('Rol activo: Conductor'), findsOneWidget);
  });
}
