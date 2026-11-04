/// Punto GPS capturado en el teléfono del conductor.
/// No incluye busId ni routeId por diseño de seguridad: el backend determina
/// la unidad asignada exclusivamente a partir del JWT del conductor.
class TrackPoint {
  final double latitude;
  final double longitude;
  final double speedKmh;
  final double heading;
  final DateTime timestamp;

  const TrackPoint({
    required this.latitude,
    required this.longitude,
    required this.speedKmh,
    required this.heading,
    required this.timestamp,
  });

  Map<String, dynamic> toJson() => {
        'latitude': latitude,
        'longitude': longitude,
        'speedKmh': speedKmh,
        'heading': heading,
        'timestamp': timestamp.toUtc().toIso8601String(),
      };
}
