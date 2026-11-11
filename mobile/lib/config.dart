/// Configuración compartida de la aplicación móvil.
/// - Emulador Android: '10.0.2.2:8080'
/// - Dispositivo físico en LAN: la IP de tu PC, por ejemplo '192.168.x.x:8080'
class AppConfig {
  static const String backendHost = '10.0.2.2:8080';

  /// Versión visible en el drawer. Sincronizar a mano con `pubspec.yaml`.
  static const String appVersion = '1.0.0';
}
