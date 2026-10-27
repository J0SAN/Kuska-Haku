import 'package:flutter/material.dart';

import 'auth/api_client.dart';
import 'auth/login_screen.dart';
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

  @override
  Widget build(BuildContext context) {
    final session = _session;
    if (session == null) {
      return LoginScreen(
        apiClient: _apiClient,
        onLoginSuccess: (result) => setState(() => _session = result),
      );
    }
    return LiveMapScreen(
      session: session,
      onLogout: () => setState(() => _session = null),
    );
  }
}
