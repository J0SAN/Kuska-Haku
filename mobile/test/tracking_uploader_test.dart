import 'package:flutter_test/flutter_test.dart';
import 'package:kuska_haku_pasajero/driver/tracking/position_buffer.dart';
import 'package:kuska_haku_pasajero/driver/tracking/track_point.dart';
import 'package:kuska_haku_pasajero/driver/tracking/tracking_uploader.dart';

TrackPoint makePoint(int n) => TrackPoint(
      latitude: -7.16,
      longitude: -78.51,
      speedKmh: 20,
      heading: 90,
      timestamp: DateTime.utc(2026, 10, 1, 12, 0, n),
    );

PositionBuffer bufferWith(int n) {
  final b = PositionBuffer();
  for (var i = 0; i < n; i++) {
    b.add(makePoint(i));
  }
  return b;
}

void main() {
  test('sin nada pendiente no llama al envío', () async {
    var calls = 0;
    final u = TrackingUploader(
      buffer: PositionBuffer(),
      send: (_) async {
        calls++;
        return UploadOutcome.sent;
      },
    );
    expect(await u.flushOnce(), isNull);
    expect(calls, 0);
  });

  test('un envío exitoso vacía la cola enviada', () async {
    final b = bufferWith(3);
    final u =
        TrackingUploader(buffer: b, send: (_) async => UploadOutcome.sent);
    expect(await u.flushOnce(), UploadOutcome.sent);
    expect(b.isEmpty, isTrue);
  });

  test('envía en lotes de máximo batchSize', () async {
    final sizes = <int>[];
    final b = bufferWith(250);
    final u = TrackingUploader(
      buffer: b,
      batchSize: 100,
      send: (batch) async {
        sizes.add(batch.length);
        return UploadOutcome.sent;
      },
    );
    while (!b.isEmpty) {
      await u.flushOnce();
    }
    expect(sizes, [100, 100, 50]);
  });

  test('si falla la red CONSERVA los puntos y los reenvía después', () async {
    final b = bufferWith(4);
    var offline = true;
    final u = TrackingUploader(
      buffer: b,
      send: (_) async => offline ? UploadOutcome.retry : UploadOutcome.sent,
    );
    expect(await u.flushOnce(), UploadOutcome.retry);
    expect(b.length, 4);
    offline = false;
    expect(await u.flushOnce(), UploadOutcome.sent);
    expect(b.isEmpty, isTrue);
  });

  test('un lote rechazado (400) se descarta para no bloquear la cola',
      () async {
    final b = bufferWith(2);
    final u = TrackingUploader(
      buffer: b,
      send: (_) async => UploadOutcome.dropBatch,
    );
    await u.flushOnce();
    expect(b.isEmpty, isTrue);
  });

  test('sesión vencida no descarta nada', () async {
    final b = bufferWith(2);
    final u = TrackingUploader(
      buffer: b,
      send: (_) async => UploadOutcome.authExpired,
    );
    await u.flushOnce();
    expect(b.length, 2);
  });

  test('la espera crece exponencialmente con los fallos y tiene tope',
      () async {
    final b = bufferWith(1);
    final u = TrackingUploader(
      buffer: b,
      send: (_) async => UploadOutcome.retry,
    );
    expect(u.nextDelay(null), const Duration(seconds: 3));

    final delays = <Duration>[];
    for (var i = 0; i < 7; i++) {
      final o = await u.flushOnce();
      delays.add(u.nextDelay(o));
    }
    expect(delays.first, const Duration(seconds: 6));
    expect(delays[1], const Duration(seconds: 12));
    expect(delays[2], const Duration(seconds: 24));
    expect(delays.last, const Duration(seconds: 60));
  });

  test('tras recuperar señal con cola pendiente, envía sin espera', () async {
    final b = bufferWith(250);
    final u = TrackingUploader(
      buffer: b,
      batchSize: 100,
      send: (_) async => UploadOutcome.sent,
    );
    final o = await u.flushOnce();
    expect(u.nextDelay(o), Duration.zero);
  });

  test('un éxito reinicia el contador de fallos', () async {
    final b = bufferWith(2);
    var failing = true;
    final u = TrackingUploader(
      buffer: b,
      batchSize: 1,
      send: (_) async => failing ? UploadOutcome.retry : UploadOutcome.sent,
    );
    await u.flushOnce();
    await u.flushOnce();
    expect(u.consecutiveFailures, 2);
    failing = false;
    await u.flushOnce();
    expect(u.consecutiveFailures, 0);
  });
}
