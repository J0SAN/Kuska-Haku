import 'package:flutter_test/flutter_test.dart';
import 'package:kuska_haku_pasajero/map/format.dart';

void main() {
  group('formatEta', () {
    test('muestra segundos cuando dura menos de un minuto', () {
      expect(formatEta(0), '0 s');
      expect(formatEta(59), '59 s');
    });

    test('muestra minutos exactos', () {
      expect(formatEta(60), '1 min');
      expect(formatEta(120), '2 min');
    });

    test('muestra minutos y segundos', () {
      expect(formatEta(90), '1 min 30 s');
      expect(formatEta(605), '10 min 5 s');
    });
  });
}
