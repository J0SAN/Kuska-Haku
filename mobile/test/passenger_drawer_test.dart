import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:kuska_haku_pasajero/auth/api_client.dart';
import 'package:kuska_haku_pasajero/map/navigation/passenger_drawer.dart';
import 'package:kuska_haku_pasajero/map/navigation/passenger_menu.dart';

AuthUser _testUser() => AuthUser(
      id: 'user-1',
      email: 'jose@correo.com',
      fullName: 'José Pérez',
      role: 'passenger',
    );

Future<void> _pumpDrawer(
  WidgetTester tester, {
  PassengerMenuItem selected = PassengerMenuItem.liveMap,
  ValueChanged<PassengerMenuItem>? onSelect,
  VoidCallback? onProfileTap,
  VoidCallback? onLogout,
}) async {
  final scaffoldKey = GlobalKey<ScaffoldState>();
  await tester.pumpWidget(
    MaterialApp(
      home: Scaffold(
        key: scaffoldKey,
        drawer: PassengerDrawer(
          user: _testUser(),
          selected: selected,
          onSelect: onSelect ?? (_) {},
          onProfileTap: onProfileTap ?? () {},
          onLogout: onLogout,
        ),
        body: const SizedBox.expand(),
      ),
    ),
  );
  scaffoldKey.currentState!.openDrawer();
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('muestra nombre, email e inicial del usuario', (tester) async {
    await _pumpDrawer(tester);

    expect(find.text('José Pérez'), findsOneWidget);
    expect(find.text('jose@correo.com'), findsOneWidget);
    expect(find.widgetWithText(CircleAvatar, 'J'), findsOneWidget);
  });

  testWidgets('no muestra estrellas ni calificación inventada', (tester) async {
    await _pumpDrawer(tester);

    expect(find.byIcon(Icons.star), findsNothing);
    expect(find.textContaining('5.0'), findsNothing);
  });

  testWidgets('tocar Rutas llama a onSelect y cierra el drawer',
      (tester) async {
    PassengerMenuItem? selectedItem;
    await _pumpDrawer(
      tester,
      onSelect: (item) => selectedItem = item,
    );

    await tester.tap(find.text('Rutas'));
    await tester.pumpAndSettle();

    expect(selectedItem, PassengerMenuItem.routes);
    expect(find.byType(PassengerDrawer), findsNothing);
  });

  testWidgets('tocar la cabecera llama a onProfileTap', (tester) async {
    var profileTaps = 0;
    await _pumpDrawer(
      tester,
      onProfileTap: () => profileTaps++,
    );

    await tester.tap(find.text('José Pérez'));
    await tester.pumpAndSettle();

    expect(profileTaps, 1);
    expect(find.byType(PassengerDrawer), findsNothing);
  });

  testWidgets('Cerrar sesión llama a onLogout', (tester) async {
    var logoutTaps = 0;
    await _pumpDrawer(
      tester,
      onLogout: () => logoutTaps++,
    );

    await tester.tap(find.text('Cerrar sesión'));
    await tester.pumpAndSettle();

    expect(logoutTaps, 1);
    expect(find.byType(PassengerDrawer), findsNothing);
  });

  testWidgets('el botón Cerrar sesión no aparece sin onLogout', (tester) async {
    await _pumpDrawer(tester);

    expect(find.text('Cerrar sesión'), findsNothing);
  });

  testWidgets('marca como seleccionado solo el ítem activo', (tester) async {
    await _pumpDrawer(tester, selected: PassengerMenuItem.arrivals);

    final selectedTile = tester.widget<ListTile>(
      find.widgetWithText(ListTile, 'Próximas llegadas'),
    );
    final otherTile = tester.widget<ListTile>(
      find.widgetWithText(ListTile, 'Mapa en vivo'),
    );

    expect(selectedTile.selected, isTrue);
    expect(otherTile.selected, isFalse);
  });

  testWidgets('muestra la versión de la app', (tester) async {
    await _pumpDrawer(tester);

    expect(find.textContaining('v1.0.0'), findsOneWidget);
  });
}
