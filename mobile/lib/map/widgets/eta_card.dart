import 'package:flutter/material.dart';

import '../../theme/app_colors.dart';
import '../../transport_routes/transport_routes_api.dart';
import '../format.dart';

/// Tarjeta inferior con las próximas llegadas de las unidades de la ruta.
class EtaCard extends StatelessWidget {
  final List<BusEta> etas;

  const EtaCard({super.key, required this.etas});

  @override
  Widget build(BuildContext context) {
    if (etas.isEmpty) {
      return const SizedBox.shrink();
    }

    return Card(
      elevation: 4,
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Próximas llegadas',
              style: TextStyle(
                fontWeight: FontWeight.bold,
                color: AppColors.primary,
              ),
            ),
            const SizedBox(height: 6),
            for (final bus in etas) _EtaRow(bus: bus),
          ],
        ),
      ),
    );
  }
}

class _EtaRow extends StatelessWidget {
  final BusEta bus;

  const _EtaRow({required this.bus});

  @override
  Widget build(BuildContext context) {
    if (bus.offRoute) {
      return Padding(
        padding: const EdgeInsets.symmetric(vertical: 2),
        child: Row(
          children: [
            const Icon(
              Icons.alt_route,
              size: 18,
              color: AppColors.warning,
            ),
            const SizedBox(width: 8),
            Expanded(
              child: Text('${bus.busId}: fuera de ruta, ETA no disponible'),
            ),
          ],
        ),
      );
    }

    final next = bus.nextStop;
    final text = next == null
        ? '${bus.busId}: recorrido completado'
        : '${bus.busId} llega a ${next.name} en '
            '${formatEta(next.etaSeconds ?? 0)} '
            '(${next.distanceMeters.round()} m)';
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        children: [
          const Icon(
            Icons.directions_bus,
            size: 18,
            color: AppColors.primary,
          ),
          const SizedBox(width: 8),
          Expanded(child: Text(text)),
        ],
      ),
    );
  }
}
