import 'package:flutter/material.dart';

import '../../auth/api_client.dart';
import '../../config.dart';
import '../../theme/app_colors.dart';
import 'passenger_menu.dart';

/// Drawer del pasajero: solo presentación, sin red ni plugins.
///
/// Recibe el usuario, la sección activa y los callbacks por constructor para
/// poder probarse sin montar `GoogleMap`.
class PassengerDrawer extends StatelessWidget {
  final AuthUser user;
  final PassengerMenuItem selected;
  final ValueChanged<PassengerMenuItem> onSelect;
  final VoidCallback onProfileTap;
  final VoidCallback? onLogout;

  const PassengerDrawer({
    super.key,
    required this.user,
    required this.selected,
    required this.onSelect,
    required this.onProfileTap,
    this.onLogout,
  });

  @override
  Widget build(BuildContext context) {
    return Drawer(
      backgroundColor: AppColors.primary,
      child: SafeArea(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            _buildHeader(context),
            const Divider(color: AppColors.divider, height: 1),
            Expanded(
              child: ListView(
                padding: EdgeInsets.zero,
                children: [
                  for (final entry in passengerMenuEntries)
                    _buildItem(context, entry),
                ],
              ),
            ),
            const Divider(color: AppColors.divider, height: 1),
            _buildFooter(context),
          ],
        ),
      ),
    );
  }

  Widget _buildHeader(BuildContext context) {
    return Semantics(
      button: true,
      label: 'Mi perfil',
      child: InkWell(
        onTap: () {
          Navigator.pop(context);
          onProfileTap();
        },
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 18),
          child: Row(
            children: [
              CircleAvatar(
                radius: 26,
                backgroundColor: AppColors.avatar,
                child: Text(
                  userInitial(user.fullName),
                  style: const TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                    color: AppColors.onPrimary,
                  ),
                ),
              ),
              const SizedBox(width: 14),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      user.fullName,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                        color: AppColors.onPrimary,
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      user.email,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                        color: AppColors.onPrimaryMuted,
                        fontSize: 13,
                      ),
                    ),
                  ],
                ),
              ),
              const Icon(
                Icons.chevron_right,
                color: AppColors.onPrimaryMuted,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildItem(BuildContext context, PassengerMenuEntry entry) {
    final isSelected = entry.item == selected;
    return ListTile(
      selected: isSelected,
      selectedTileColor: AppColors.primaryDark,
      selectedColor: AppColors.onPrimary,
      iconColor: AppColors.onPrimaryMuted,
      textColor: AppColors.onPrimary,
      contentPadding: const EdgeInsets.symmetric(horizontal: 20),
      minLeadingWidth: 24,
      leading: Icon(entry.icon, size: 24),
      title: Text(
        entry.label,
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
        style: TextStyle(
          fontSize: 16,
          fontWeight: isSelected ? FontWeight.w600 : FontWeight.normal,
        ),
      ),
      onTap: () {
        Navigator.pop(context);
        onSelect(entry.item);
      },
    );
  }

  Widget _buildFooter(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 16, 20, 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          if (onLogout != null)
            ElevatedButton(
              onPressed: () {
                Navigator.pop(context);
                onLogout!();
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.accent,
                foregroundColor: AppColors.onAccent,
                padding: const EdgeInsets.symmetric(vertical: 14),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              child: const Text(
                'Cerrar sesión',
                style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
              ),
            ),
          const SizedBox(height: 14),
          const Text(
            'Kuska Haku · v${AppConfig.appVersion}',
            textAlign: TextAlign.center,
            style: TextStyle(color: AppColors.onPrimarySubtle, fontSize: 12),
          ),
        ],
      ),
    );
  }
}
