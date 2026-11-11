import 'package:flutter/material.dart';

/// Pide al pasajero el texto de una sugerencia. Devuelve null si cancela.
Future<String?> showPassengerSuggestionDialog(BuildContext context) {
  String draftText = '';
  return showDialog<String>(
    context: context,
    builder: (ctx) => AlertDialog(
      title: const Text('Sugerencia general'),
      content: TextField(
        autofocus: true,
        maxLines: 4,
        maxLength: 1000,
        onChanged: (value) => draftText = value,
        decoration: const InputDecoration(
          hintText:
              'Cuéntanos qué podemos mejorar del servicio. Para opinar sobre un vehículo, escanea su QR.',
          border: OutlineInputBorder(),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(ctx),
          child: const Text('Cancelar'),
        ),
        FilledButton(
          onPressed: () => Navigator.pop(ctx, draftText),
          child: const Text('Enviar'),
        ),
      ],
    ),
  );
}
