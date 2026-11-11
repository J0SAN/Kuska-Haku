import 'package:flutter/material.dart';

import '../../theme/app_colors.dart';
import '../../transport_routes/transport_routes_api.dart';
import '../format.dart';

/// Sheet con los próximos arribos de las unidades de la ruta activa.
Future<void> showRouteStopsSheet(
  BuildContext context, {
  required RouteInfo? route,
  required List<BusEta> etas,
}) {
  return showModalBottomSheet<void>(
    context: context,
    backgroundColor: AppColors.primary,
    shape: const RoundedRectangleBorder(
      borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
    ),
    builder: (ctx) => Padding(
      padding: const EdgeInsets.all(20),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            route == null ? 'Paraderos de la ruta' : 'Ruta: ${route.name}',
            style: const TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: AppColors.onPrimary,
            ),
          ),
          const SizedBox(height: 12),
          if (etas.isEmpty)
            const Padding(
              padding: EdgeInsets.symmetric(vertical: 16),
              child: Text(
                'No hay unidades transmitiendo tiempos de llegada en este momento.',
                style: TextStyle(color: AppColors.onPrimaryMuted),
              ),
            )
          else
            for (final bus in etas)
              ListTile(
                contentPadding: EdgeInsets.zero,
                leading: const Icon(
                  Icons.directions_bus,
                  color: AppColors.accent,
                ),
                title: Text(
                  'Unidad ${bus.busId}',
                  style: const TextStyle(
                    color: AppColors.onPrimary,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                subtitle: Text(
                  bus.offRoute
                      ? 'Fuera de ruta — ETA no disponible'
                      : (bus.nextStop == null
                          ? 'Recorrido completado'
                          : 'Próximo paradero: ${bus.nextStop!.name} '
                              '(${formatEta(bus.nextStop!.etaSeconds ?? 0)})'),
                  style: const TextStyle(color: AppColors.onPrimaryMuted),
                ),
              ),
        ],
      ),
    ),
  );
}
