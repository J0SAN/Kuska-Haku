import 'package:flutter/material.dart';

import '../../config.dart';

/// Ayuda breve sobre cómo leer el mapa y el ETA.
Future<void> showPassengerHelpDialog(BuildContext context) {
  return showDialog<void>(
    context: context,
    builder: (ctx) => AlertDialog(
      title: const Text('Ayuda'),
      content: const Text(
        'Kuska Haku muestra en el mapa la posición en vivo de las unidades '
        'de transporte público.\n\n'
        'La línea verde marca el trazado de la ruta activa y cada marcador '
        'es un bus que está transmitiendo su GPS.\n\n'
        'Abre «Próximas llegadas» para ver cuánto falta para que cada '
        'unidad llegue a su siguiente paradero.\n\n'
        'Versión ${AppConfig.appVersion}.',
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(ctx),
          child: const Text('Entendido'),
        ),
      ],
    ),
  );
}
