import 'package:flutter/material.dart';

import '../../auth/api_client.dart';
import '../../theme/app_colors.dart';
import 'passenger_menu.dart';

/// Abre el sheet de perfil del pasajero.
Future<void> showPassengerProfileSheet(BuildContext context, AuthUser user) {
  return showModalBottomSheet<void>(
    context: context,
    backgroundColor: AppColors.primary,
    shape: const RoundedRectangleBorder(
      borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
    ),
    builder: (_) => PassengerProfileSheet(user: user),
  );
}

/// Contenido del perfil: legible sobre el fondo verde del sheet.
class PassengerProfileSheet extends StatelessWidget {
  final AuthUser user;

  const PassengerProfileSheet({super.key, required this.user});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          CircleAvatar(
            radius: 36,
            backgroundColor: AppColors.avatar,
            child: Text(
              userInitial(user.fullName),
              style: const TextStyle(
                fontSize: 30,
                fontWeight: FontWeight.bold,
                color: AppColors.onPrimary,
              ),
            ),
          ),
          const SizedBox(height: 12),
          Text(
            user.fullName,
            textAlign: TextAlign.center,
            style: const TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: AppColors.onPrimary,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            user.email,
            textAlign: TextAlign.center,
            style: const TextStyle(
              color: AppColors.onPrimaryMuted,
              fontSize: 14,
            ),
          ),
          const SizedBox(height: 16),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
            decoration: BoxDecoration(
              color: AppColors.primaryDark,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Icon(
                  Icons.verified_user,
                  color: AppColors.accent,
                  size: 18,
                ),
                const SizedBox(width: 8),
                Text(
                  'Rol activo: ${roleLabel(user.role)}',
                  style: const TextStyle(
                    color: AppColors.onPrimary,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 20),
        ],
      ),
    );
  }
}
