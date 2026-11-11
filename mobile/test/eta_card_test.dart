import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:kuska_haku_pasajero/map/widgets/eta_card.dart';
import 'package:kuska_haku_pasajero/transport_routes/transport_routes_api.dart';

BusEta _bus({
  String busId = 'bus-01',
  bool offRoute = false,
  StopEta? nextStop,
}) {
  return BusEta(
    busId: busId,
    offRoute: offRoute,
    stops: nextStop == null ? const [] : [nextStop],
  );
}

Future<void> _pumpCard(WidgetTester tester, List<BusEta> etas) {
  return tester.pumpWidget(
    MaterialApp(
      home: Scaffold(body: EtaCard(etas: etas)),
    ),
  );
}

void main() {
  testWidgets('muestra la próxima llegada de cada unidad', (tester) async {
    await _pumpCard(tester, [
      _bus(
        nextStop: StopEta(
          name: 'Paradero Central',
          stopOrder: 1,
          distanceMeters: 120,
          etaSeconds: 90,
          passed: false,
        ),
      ),
    ]);

    expect(find.text('Próximas llegadas'), findsOneWidget);
    expect(
      find.textContaining('bus-01 llega a Paradero Central en 1 min 30 s'),
      findsOneWidget,
    );
  });

  testWidgets('avisa cuando la unidad está fuera de ruta', (tester) async {
    await _pumpCard(tester, [_bus(offRoute: true)]);

    expect(
      find.textContaining('bus-01: fuera de ruta'),
      findsOneWidget,
    );
  });

  testWidgets('no dibuja nada sin unidades', (tester) async {
    await _pumpCard(tester, []);

    expect(find.text('Próximas llegadas'), findsNothing);
  });
}
