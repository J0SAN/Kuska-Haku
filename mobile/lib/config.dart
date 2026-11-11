/// Configuración compartida de la aplicación móvil.
///
/// El host y el uso de TLS se eligen al compilar, sin tocar el código:
/// - Emulador Android (por defecto): `flutter run` → http://10.0.2.2:8080
/// - Dispositivo físico en la LAN:
///   `flutter run --dart-define=BACKEND_HOST=192.168.x.x:8080`
/// - Backend desplegado (HTTPS/WSS):
///   `flutter build apk --dart-define=BACKEND_HOST=<servicio>.onrender.com --dart-define=BACKEND_TLS=true`
class AppConfig {
  static const String backendHost = String.fromEnvironment(
    'BACKEND_HOST',
    defaultValue: '10.0.2.2:8080',
  );

  /// true en el backend desplegado: Render termina TLS y redirige HTTP a HTTPS.
  static const bool backendTls = bool.fromEnvironment('BACKEND_TLS');

  static const String httpScheme = backendTls ? 'https' : 'http';
  static const String wsScheme = backendTls ? 'wss' : 'ws';

  /// Versión visible en el drawer. Sincronizar a mano con `pubspec.yaml`.
  static const String appVersion = '1.0.0';
}
