import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:kuska_haku_pasajero/map/navigation/passenger_menu.dart';

void main() {
  group('userInitial', () {
    test('usa la primera letra del nombre en mayúscula', () {
      expect(userInitial('José'), 'J');
    });

    test('ignora los espacios alrededor', () {
      expect(userInitial('  ana  '), 'A');
    });

    test('devuelve U cuando el nombre está vacío', () {
      expect(userInitial(''), 'U');
      expect(userInitial('   '), 'U');
    });

    test('respeta las tildes del español', () {
      expect(userInitial('ángel'), 'Á');
    });
  });

  group('roleLabel', () {
    test('traduce los roles conocidos al español', () {
      expect(roleLabel('passenger'), 'Pasajero');
      expect(roleLabel('pasajero'), 'Pasajero');
      expect(roleLabel('driver'), 'Conductor');
      expect(roleLabel('conductor'), 'Conductor');
      expect(roleLabel('admin'), 'Administrador');
      expect(roleLabel('administrador'), 'Administrador');
    });

    test('devuelve el rol original si no lo conoce', () {
      expect(roleLabel('supervisor'), 'supervisor');
    });
  });

  group('passengerMenuEntries', () {
    test('mantiene el orden esperado', () {
      expect(
        passengerMenuEntries.map((entry) => entry.item).toList(),
        [
          PassengerMenuItem.liveMap,
          PassengerMenuItem.arrivals,
          PassengerMenuItem.routes,
          PassengerMenuItem.suggestions,
          PassengerMenuItem.help,
        ],
      );
    });

    test('no repite ítems y cubre todos los valores del enum', () {
      final items = passengerMenuEntries.map((entry) => entry.item).toList();
      expect(items.toSet().length, items.length);
      expect(items.toSet(), PassengerMenuItem.values.toSet());
    });

    test('todas las entradas tienen etiqueta e ícono', () {
      for (final entry in passengerMenuEntries) {
        expect(entry.label.trim(), isNotEmpty);
        expect(entry.icon, isA<IconData>());
      }
    });
  });
}
