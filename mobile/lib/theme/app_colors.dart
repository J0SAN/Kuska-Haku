import 'package:flutter/material.dart';

/// Paleta única de Kuska Haku.
///
/// Toda pantalla o widget nuevo debe tomar sus colores de aquí para que el
/// panel, el login, el mapa y el drawer mantengan la misma identidad visual.
class AppColors {
  AppColors._();

  static const Color primary = Color(0xFF2E7D32);
  static const Color primaryDark = Color(0xFF1B5E20);
  static const Color background = Color(0xFFF1F8F2);
  static const Color accent = Color(0xFFAEEA00);
  static const Color onAccent = Colors.black;

  static const Color onPrimary = Colors.white;
  static const Color onPrimaryMuted = Colors.white70;
  static const Color onPrimarySubtle = Colors.white54;
  static const Color divider = Colors.white24;

  static const Color avatar = Color(0xFFC2185B);
  static const Color warning = Colors.orange;
  static const Color error = Colors.red;
  static const Color inactive = Colors.grey;
}
