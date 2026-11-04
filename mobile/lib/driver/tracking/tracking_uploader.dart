import 'position_buffer.dart';
import 'track_point.dart';

/// Resultado del intento de envío de un lote HTTP al backend:
///   200 -> sent         (procesado correctamente)
///   400 -> dropBatch    (lote malformado; se descarta para no atascar la cola)
///   401 -> authExpired  (token vencido; requiere iniciar sesión nuevamente)
///   409 -> noBus        (conductor sin unidad asignada)
///   5xx / error de red -> retry (conserva el lote y aplica backoff exponencial)
enum UploadOutcome { sent, retry, dropBatch, authExpired, noBus }

typedef BatchSender = Future<UploadOutcome> Function(List<TrackPoint> batch);

class TrackingUploader {
  TrackingUploader({
    required this.buffer,
    required this.send,
    this.batchSize = 100,
    this.baseDelay = const Duration(seconds: 3),
    this.maxDelay = const Duration(seconds: 60),
  });

  final PositionBuffer buffer;
  final BatchSender send;
  final int batchSize;
  final Duration baseDelay;
  final Duration maxDelay;

  int _failures = 0;
  int get consecutiveFailures => _failures;

  /// Intenta enviar un lote. Retorna null si el buffer está vacío.
  Future<UploadOutcome?> flushOnce() async {
    final batch = buffer.peekBatch(batchSize);
    if (batch == null) return null;

    final outcome = await send(batch.points);
    switch (outcome) {
      case UploadOutcome.sent:
      case UploadOutcome.dropBatch:
        buffer.ack(batch.lastSeq);
        _failures = 0;
        break;
      case UploadOutcome.retry:
      case UploadOutcome.noBus:
        _failures++;
        break;
      case UploadOutcome.authExpired:
        break;
    }
    return outcome;
  }

  /// Calcula la espera antes del próximo envío:
  /// - Si el envío fue exitoso y aún quedan puntos en cola: Duration.zero.
  /// - Si hubo fallos: backoff exponencial (6s, 12s, 24s...) con tope en 60s.
  Duration nextDelay(UploadOutcome? last) {
    if (last == UploadOutcome.sent && !buffer.isEmpty) return Duration.zero;
    if (_failures == 0) return baseDelay;
    final exponent = _failures > 5 ? 5 : _failures;
    final delay = baseDelay * (1 << exponent);
    return delay > maxDelay ? maxDelay : delay;
  }
}
