import 'package:flutter_test/flutter_test.dart';
import 'package:kuska_haku_pasajero/main.dart';

void main() {
  testWidgets('Renderiza la pantalla inicial correctamente',
      (WidgetTester tester) async {
    // Pasamos el parámetro requerido initialRoute que definimos en main.dart
    await tester.pumpWidget(const KuskaHakuApp(initialRoute: '/login'));

    // Verificamos que el widget principal se monte sin errores
    expect(find.byType(KuskaHakuApp), findsOneWidget);
  });
}
