import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:geolocator/geolocator.dart';

import '../data/driver_api.dart';
import 'position_buffer.dart';
import 'track_point.dart';
import 'tracking_uploader.dart';

/// Gestiona el turno del conductor: permisos de GPS, lectura continua en primer
/// plano, almacenamiento en buffer y despacho por lotes HTTP.
class ShiftController extends ChangeNotifier {
  ShiftController({required DriverApi api, PositionBuffer? positionBuffer})
      : _api = api,
        buffer = positionBuffer ?? PositionBuffer() {
    uploader = TrackingUploader(buffer: buffer, send: _api.sendPositions);
  }

  final DriverApi _api;
  final PositionBuffer buffer;
  late final TrackingUploader uploader;

  static const _captureInterval = Duration(seconds: 3);
  static const _tolerance = Duration(milliseconds: 500);

  bool running = false;
  bool sessionExpired = false;
  String? gpsMessage;
  String? uploadMessage;
  DateTime? lastSentAt;
  TrackPoint? lastPoint;

  int get pending => buffer.length;

  StreamSubscription<Position>? _sub;
  Timer? _timer;
  bool _flushing = false;
  bool _disposed = false;
  DateTime? _lastCaptured;

  void _notify() {
    if (!_disposed) notifyListeners();
  }

  Future<String?> _ensurePermission() async {
    if (!await Geolocator.isLocationServiceEnabled()) {
      return 'Activa el GPS del teléfono para iniciar el turno.';
    }
    var permission = await Geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
    }
    if (permission == LocationPermission.denied) {
      return 'Sin permiso de ubicación no se puede compartir tu posición.';
    }
    if (permission == LocationPermission.deniedForever) {
      return 'El permiso de ubicación está bloqueado. Habilítalo en Ajustes del teléfono.';
    }
    return null;
  }

  LocationSettings _settings() {
    if (defaultTargetPlatform == TargetPlatform.android) {
      return AndroidSettings(
        accuracy: LocationAccuracy.high,
        distanceFilter: 0,
        intervalDuration: _captureInterval,
        foregroundNotificationConfig: const ForegroundNotificationConfig(
          notificationTitle: 'Kuska Haku',
          notificationText: 'Turno activo: compartiendo ubicación de la unidad',
          enableWakeLock: true,
          setOngoing: true,
        ),
      );
    }
    if (defaultTargetPlatform == TargetPlatform.iOS ||
        defaultTargetPlatform == TargetPlatform.macOS) {
      return AppleSettings(
        accuracy: LocationAccuracy.high,
        activityType: ActivityType.automotiveNavigation,
        distanceFilter: 0,
        pauseLocationUpdatesAutomatically: false,
        showBackgroundLocationIndicator: true,
        allowBackgroundLocationUpdates: true,
      );
    }
    return const LocationSettings(
      accuracy: LocationAccuracy.high,
      distanceFilter: 0,
    );
  }

  Future<void> start() async {
    if (running) return;

    final problem = await _ensurePermission();
    if (problem != null) {
      gpsMessage = problem;
      _notify();
      return;
    }

    gpsMessage = null;
    uploadMessage = null;
    sessionExpired = false;
    running = true;

    _sub = Geolocator.getPositionStream(locationSettings: _settings()).listen(
      _onPosition,
      onError: (Object e) {
        gpsMessage = 'Error del GPS: $e';
        _notify();
      },
    );
    _schedule(const Duration(seconds: 3));
    _notify();
  }

  Future<void> stop() async {
    if (!running) return;
    running = false;
    await _sub?.cancel();
    _sub = null;
    _timer?.cancel();
    _timer = null;

    while (_flushing) {
      await Future<void>.delayed(const Duration(milliseconds: 100));
    }
    if (!sessionExpired) {
      for (var i = 0; i < 3 && !buffer.isEmpty; i++) {
        final outcome = await uploader.flushOnce();
        if (outcome != UploadOutcome.sent) break;
        lastSentAt = DateTime.now();
      }
    }
    uploadMessage = buffer.isEmpty
        ? null
        : 'Quedaron ${buffer.length} puntos sin enviar por falta de conexión.';
    _notify();
  }

  void _onPosition(Position p) {
    final now = DateTime.now();
    final last = _lastCaptured;
    if (last != null && now.difference(last) < _captureInterval - _tolerance) {
      return;
    }
    _lastCaptured = now;

    final speedKmh = p.speed.isNaN || p.speed < 0 ? 0.0 : p.speed * 3.6;
    final heading = p.heading.isNaN || p.heading < 0 ? 0.0 : p.heading % 360;

    final point = TrackPoint(
      latitude: p.latitude,
      longitude: p.longitude,
      speedKmh: speedKmh,
      heading: heading,
      timestamp: now.toUtc(),
    );
    buffer.add(point);
    lastPoint = point;
    _notify();
  }

  void _schedule(Duration delay) {
    _timer?.cancel();
    if (!running || sessionExpired) return;
    _timer = Timer(delay, _flushLoop);
  }

  Future<void> _flushLoop() async {
    if (!running || _flushing) return;
    _flushing = true;
    UploadOutcome? outcome;
    try {
      outcome = await uploader.flushOnce();
    } finally {
      _flushing = false;
    }
    _apply(outcome);
    if (outcome == UploadOutcome.authExpired) {
      sessionExpired = true;
      _notify();
      return;
    }
    _schedule(uploader.nextDelay(outcome));
    _notify();
  }

  void _apply(UploadOutcome? outcome) {
    switch (outcome) {
      case UploadOutcome.sent:
        lastSentAt = DateTime.now();
        uploadMessage = null;
        break;
      case UploadOutcome.retry:
        uploadMessage =
            'Sin conexión: $pending puntos en espera. Se reintentará automáticamente.';
        break;
      case UploadOutcome.noBus:
        uploadMessage =
            'No tienes un bus asignado. Solicita al administrador que te asigne uno.';
        break;
      case UploadOutcome.dropBatch:
        uploadMessage =
            'El servidor rechazó un lote inválido y fue descartado.';
        break;
      case UploadOutcome.authExpired:
        uploadMessage = 'Tu sesión expiró. Vuelve a ingresar.';
        break;
      case null:
        break;
    }
  }

  @override
  void dispose() {
    _disposed = true;
    _sub?.cancel();
    _timer?.cancel();
    super.dispose();
  }
}
