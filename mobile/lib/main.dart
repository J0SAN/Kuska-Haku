import 'package:flutter/material.dart';

import 'auth/api_client.dart';
import 'auth/login_screen.dart';
import 'driver/presentation/driver_home_screen.dart';
import 'map/map_screen.dart';

void main() {
  runApp(const KuskaHakuApp());
}

class KuskaHakuApp extends StatelessWidget {
  final String? initialRoute;

  const KuskaHakuApp({super.key, this.initialRoute});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Kuska Haku',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF2E7D32)),
        useMaterial3: true,
      ),
      home: const RootScreen(),
    );
  }
}

class RootScreen extends StatefulWidget {
  const RootScreen({super.key});

  @override
  State<RootScreen> createState() => _RootScreenState();
}

class _RootScreenState extends State<RootScreen> {
  final _apiClient = AuthApiClient();
  AuthResult? _session;

  void _logout() => setState(() => _session = null);

  @override
  Widget build(BuildContext context) {
    final session = _session;
    if (session == null) {
      return LoginScreen(
        apiClient: _apiClient,
        onLoginSuccess: (result) => setState(() => _session = result),
      );
    }

    switch (session.user.role) {
      case 'driver':
      case 'conductor':
        return DriverHomeScreen(session: session, onLogout: _logout);
      case 'passenger':
      case 'pasajero':
        return LiveMapScreen(session: session, onLogout: _logout);
      default:
        return _AdminNoticeScreen(onLogout: _logout);
    }
  }
}

class _AdminNoticeScreen extends StatelessWidget {
  final VoidCallback onLogout;
  const _AdminNoticeScreen({required this.onLogout});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Icon(
                Icons.admin_panel_settings,
                size: 56,
                color: Color(0xFF2E7D32),
              ),
              const SizedBox(height: 12),
              const Text(
                'Las cuentas de administrador se gestionan desde el panel web.',
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 16),
              FilledButton(
                onPressed: onLogout,
                child: const Text('Salir'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
