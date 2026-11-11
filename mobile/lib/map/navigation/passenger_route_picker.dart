import 'package:flutter/material.dart';

import '../../theme/app_colors.dart';
import '../../transport_routes/transport_routes_api.dart';

/// Abre el selector de rutas. Devuelve la ruta elegida o null si se cierra.
Future<RouteInfo?> showPassengerRoutePicker(
  BuildContext context, {
  required Future<List<RouteInfo>> routes,
  required String? currentRouteId,
}) {
  return showModalBottomSheet<RouteInfo>(
    context: context,
    backgroundColor: AppColors.primary,
    shape: const RoundedRectangleBorder(
      borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
    ),
    builder: (ctx) => SafeArea(
      child: FutureBuilder<List<RouteInfo>>(
        future: routes,
        builder: (context, snapshot) {
          if (snapshot.connectionState != ConnectionState.done) {
            return const Padding(
              padding: EdgeInsets.all(32),
              child: Center(
                child: CircularProgressIndicator(color: AppColors.onPrimary),
              ),
            );
          }

          if (snapshot.hasError) {
            final error = snapshot.error;
            final message = error is RoutesApiException
                ? error.message
                : 'No se pudieron cargar las rutas';
            return _RoutePickerMessage(message: message);
          }

          final list = snapshot.data ?? const <RouteInfo>[];
          if (list.isEmpty) {
            return const _RoutePickerMessage(
              message: 'No hay rutas activas por el momento.',
            );
          }

          return Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Padding(
                padding: EdgeInsets.fromLTRB(20, 20, 20, 8),
                child: Text(
                  'Elige una ruta',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: AppColors.onPrimary,
                  ),
                ),
              ),
              Flexible(
                child: ListView(
                  shrinkWrap: true,
                  children: [
                    for (final route in list)
                      ListTile(
                        contentPadding:
                            const EdgeInsets.symmetric(horizontal: 20),
                        leading: Icon(
                          route.id == currentRouteId
                              ? Icons.check_circle
                              : Icons.alt_route,
                          color: route.id == currentRouteId
                              ? AppColors.accent
                              : AppColors.onPrimaryMuted,
                        ),
                        title: Text(
                          route.name,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                          style: const TextStyle(
                            color: AppColors.onPrimary,
                            fontSize: 16,
                          ),
                        ),
                        onTap: () => Navigator.pop(ctx, route),
                      ),
                  ],
                ),
              ),
              const SizedBox(height: 12),
            ],
          );
        },
      ),
    ),
  );
}

class _RoutePickerMessage extends StatelessWidget {
  final String message;

  const _RoutePickerMessage({required this.message});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(24),
      child: Text(
        message,
        style: const TextStyle(
          color: AppColors.onPrimaryMuted,
          fontSize: 15,
        ),
      ),
    );
  }
}
