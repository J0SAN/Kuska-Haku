import 'package:flutter/material.dart';

/// Secciones disponibles en el menú lateral del pasajero.
enum PassengerMenuItem { liveMap, arrivals, routes, suggestions, help }

/// Entrada de presentación del menú: etiqueta e ícono de una sección.
class PassengerMenuEntry {
  final PassengerMenuItem item;
  final String label;
  final IconData icon;

  const PassengerMenuEntry({
    required this.item,
    required this.label,
    required this.icon,
  });
}

/// Orden visible en el drawer. Solo incluye funciones que existen hoy.
const List<PassengerMenuEntry> passengerMenuEntries = [
  PassengerMenuEntry(
    item: PassengerMenuItem.liveMap,
    label: 'Mapa en vivo',
    icon: Icons.directions_bus_outlined,
  ),
  PassengerMenuEntry(
    item: PassengerMenuItem.arrivals,
    label: 'Próximas llegadas',
    icon: Icons.schedule,
  ),
  PassengerMenuEntry(
    item: PassengerMenuItem.routes,
    label: 'Rutas',
    icon: Icons.alt_route,
  ),
  PassengerMenuEntry(
    item: PassengerMenuItem.suggestions,
    label: 'Sugerencias generales',
    icon: Icons.rate_review_outlined,
  ),
  PassengerMenuEntry(
    item: PassengerMenuItem.help,
    label: 'Ayuda',
    icon: Icons.info_outline,
  ),
];

/// Inicial del avatar: primera letra del nombre en mayúscula o `U` si no hay.
///
/// `toUpperCase` respeta las tildes del español («ángel» → «Á»).
String userInitial(String fullName) {
  final trimmed = fullName.trim();
  if (trimmed.isEmpty) return 'U';
  return trimmed[0].toUpperCase();
}

/// Etiqueta legible del rol que viene en el JWT.
String roleLabel(String role) {
  switch (role) {
    case 'passenger':
    case 'pasajero':
      return 'Pasajero';
    case 'driver':
    case 'conductor':
      return 'Conductor';
    case 'admin':
    case 'administrador':
      return 'Administrador';
    default:
      return role;
  }
}
