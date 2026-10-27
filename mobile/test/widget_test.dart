import 'package:flutter_test/flutter_test.dart';
import 'package:kuska_haku_pasajero/main.dart';

void main() {
  testWidgets('Renderiza la pantalla inicial de login correctamente',
      (WidgetTester tester) async {
    await tester.pumpWidget(const KuskaHakuApp());

    expect(find.byType(KuskaHakuApp), findsOneWidget);
    expect(find.text('Kuska Haku'), findsOneWidget);
    expect(find.text('Ingresar'), findsOneWidget);
  });
}
