/// Formatea una duración en segundos para mostrarla en la tarjeta de ETA.
String formatEta(int seconds) {
  if (seconds < 60) {
    return '$seconds s';
  }
  final min = seconds ~/ 60;
  final sec = seconds % 60;
  return sec == 0 ? '$min min' : '$min min $sec s';
}
