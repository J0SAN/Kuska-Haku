import 'package:flutter_test/flutter_test.dart';
import 'package:kuska_haku_pasajero/driver/tracking/position_buffer.dart';
import 'package:kuska_haku_pasajero/driver/tracking/track_point.dart';

TrackPoint makePoint(int n) => TrackPoint(
      latitude: -7.16 + n * 0.0001,
      longitude: -78.51,
      speedKmh: 20,
      heading: 90,
      timestamp: DateTime.utc(2026, 10, 1, 12, 0, n),
    );

void main() {
  test('devuelve los puntos del más viejo al más nuevo sin sacarlos', () {
    final b = PositionBuffer();
    for (var i = 0; i < 5; i++) {
      b.add(makePoint(i));
    }
    final batch = b.peekBatch(3)!;
    expect(batch.points.map((p) => p.timestamp.second), [0, 1, 2]);
    expect(b.length, 5);
  });

  test('ack elimina solo lo confirmado', () {
    final b = PositionBuffer();
    for (var i = 0; i < 5; i++) {
      b.add(makePoint(i));
    }
    b.ack(b.peekBatch(3)!.lastSeq);
    expect(b.length, 2);
    expect(b.peekBatch(10)!.points.first.timestamp.second, 3);
  });

  test(
      'puntos que llegan mientras hay un envío en vuelo NO se pierden en el ack',
      () {
    final b = PositionBuffer();
    b.add(makePoint(0));
    b.add(makePoint(1));
    final inFlight = b.peekBatch(100)!;
    b.add(makePoint(2));
    b.ack(inFlight.lastSeq);
    expect(b.length, 1);
    expect(b.peekBatch(10)!.points.single.timestamp.second, 2);
  });

  test('al desbordar descarta los más viejos y los cuenta', () {
    final b = PositionBuffer(capacity: 3);
    for (var i = 0; i < 5; i++) {
      b.add(makePoint(i));
    }
    expect(b.length, 3);
    expect(b.dropped, 2);
    expect(b.peekBatch(10)!.points.first.timestamp.second, 2);
  });

  test('un ack sobre un lote cuyos puntos ya se descartaron no borra de más',
      () {
    final b = PositionBuffer(capacity: 3);
    for (var i = 0; i < 3; i++) {
      b.add(makePoint(i));
    }
    final inFlight = b.peekBatch(100)!;
    b.add(makePoint(3));
    b.add(makePoint(4));
    b.ack(inFlight.lastSeq);
    expect(b.length, 2);
    expect(b.peekBatch(10)!.points.map((p) => p.timestamp.second), [3, 4]);
  });

  test('el JSON usa UTC y no lleva busId ni routeId', () {
    final json = makePoint(7).toJson();
    expect(json['timestamp'], endsWith('Z'));
    expect(json.containsKey('busId'), isFalse);
    expect(json.containsKey('bus_id'), isFalse);
    expect(json.containsKey('routeId'), isFalse);
    expect(json.containsKey('route_id'), isFalse);
  });
}
