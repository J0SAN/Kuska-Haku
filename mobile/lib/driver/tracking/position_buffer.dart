import 'track_point.dart';

/// Lote de puntos listo para enviar junto con el número de secuencia del último
/// punto, utilizado para confirmar únicamente el tramo recibido por el servidor.
class BufferedBatch {
  final List<TrackPoint> points;
  final int lastSeq;
  const BufferedBatch({required this.points, required this.lastSeq});
}

/// Cola circular en memoria de puntos GPS pendientes de envío.
/// Con capacidad de 2000 puntos (1 punto cada 3 s), soporta ~100 minutos sin señal.
class PositionBuffer {
  PositionBuffer({this.capacity = 2000}) : assert(capacity > 0);

  final int capacity;
  final List<_Entry> _entries = [];
  int _nextSeq = 1;
  int _dropped = 0;

  int get length => _entries.length;
  bool get isEmpty => _entries.isEmpty;

  /// Cantidad de puntos antiguos descartados por desborde de capacidad.
  int get dropped => _dropped;

  void add(TrackPoint point) {
    _entries.add(_Entry(_nextSeq++, point));
    while (_entries.length > capacity) {
      _entries.removeAt(0);
      _dropped++;
    }
  }

  /// Devuelve hasta [max] puntos del más antiguo al más reciente sin extraerlos.
  BufferedBatch? peekBatch(int max) {
    if (_entries.isEmpty) return null;
    final slice = _entries.take(max).toList(growable: false);
    return BufferedBatch(
      points: slice.map((e) => e.point).toList(growable: false),
      lastSeq: slice.last.seq,
    );
  }

  /// Elimina de la cola todos los puntos con secuencia <= [lastSeq].
  void ack(int lastSeq) {
    _entries.removeWhere((e) => e.seq <= lastSeq);
  }
}

class _Entry {
  final int seq;
  final TrackPoint point;
  _Entry(this.seq, this.point);
}
