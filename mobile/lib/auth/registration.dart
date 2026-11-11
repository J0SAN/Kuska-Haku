/// Lógica de validación del formulario de registro, separada del widget para
/// poder probarla sin Flutter. El algoritmo del dígito verificador es el mismo
/// que `backend/internal/domain/dni.go`.
library;

const _dniPattern = r'^[0-9]{8}$';
const _emailPattern = r'^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,24}$';
const _namePattern = r"^[\p{L} '’-]+$";

const _checkDigitWeights = [3, 2, 7, 6, 5, 4, 3, 2];
const _checkDigitNumbers = ['6', '7', '8', '9', '0', '1', '1', '2', '3', '4', '5'];
const _checkDigitLetters = ['K', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J'];

String _collapseSpaces(String value) =>
    value.trim().replaceAll(RegExp(r'\s+'), ' ');

/// Valida el número de DNI: exactamente 8 dígitos ASCII.
String? validateDni(String? value) {
  final dni = (value ?? '').trim();
  if (!RegExp(_dniPattern).hasMatch(dni)) {
    return 'Ingresa tu DNI de 8 dígitos.';
  }
  return null;
}

/// Valida el carácter impreso junto al número. Acepta el dígito numérico o la
/// letra equivalente, en mayúscula o minúscula.
String? validateCheckDigit(String? dni, String? checkDigit) {
  final check = (checkDigit ?? '').trim().toUpperCase();
  if (check.isEmpty) {
    return 'Ingresa el dígito verificador.';
  }
  if (check.length != 1) {
    return 'El dígito verificador es un solo carácter.';
  }

  final number = (dni ?? '').trim();
  if (!RegExp(_dniPattern).hasMatch(number)) {
    return null; // el error del DNI ya se muestra en su campo
  }

  var sum = 0;
  for (var i = 0; i < 8; i++) {
    sum += int.parse(number[i]) * _checkDigitWeights[i];
  }
  var index = 11 - (sum % 11);
  if (index == 11) {
    index = 0;
  }

  if (check == _checkDigitNumbers[index] ||
      check == _checkDigitLetters[index]) {
    return null;
  }
  return 'El dígito verificador no coincide. Revisa el número de tu DNI.';
}

/// Valida nombres o apellidos: 2 a 100 caracteres, solo letras Unicode,
/// espacios, apóstrofos y guiones.
String? validatePersonName(String? value) {
  final name = _collapseSpaces(value ?? '');
  final length = name.runes.length;
  if (length < 2 || length > 100) {
    return 'Ingresa entre 2 y 100 caracteres.';
  }
  if (!RegExp(_namePattern, unicode: true).hasMatch(name)) {
    return 'Solo se permiten letras, espacios, apóstrofos y guiones.';
  }
  return null;
}

/// Normaliza y valida el correo electrónico.
String? validateEmail(String? value) {
  final email = (value ?? '').trim().toLowerCase();
  if (!RegExp(_emailPattern).hasMatch(email)) {
    return 'Ingresa un correo válido.';
  }
  return null;
}

/// La contraseña debe tener al menos 8 caracteres.
String? validatePassword(String? value) {
  if ((value ?? '').length < 8) {
    return 'La contraseña debe tener al menos 8 caracteres.';
  }
  return null;
}

/// Comprueba que la confirmación coincida con la contraseña.
String? passwordsMatch(String? password, String? confirmation) {
  if ((confirmation ?? '') != (password ?? '')) {
    return 'Las contraseñas no coinciden.';
  }
  return null;
}

/// Datos que se envían a `POST /api/v1/auth/register`.
class RegisterInput {
  final String dni;
  final String checkDigit;
  final String firstNames;
  final String lastNames;
  final String email;
  final String password;
  final bool acceptTerms;

  const RegisterInput({
    required this.dni,
    required this.checkDigit,
    required this.firstNames,
    required this.lastNames,
    required this.email,
    required this.password,
    required this.acceptTerms,
  });

  Map<String, dynamic> toJson() => {
        'dni': dni.trim(),
        'checkDigit': checkDigit.trim().toUpperCase(),
        'firstNames': _collapseSpaces(firstNames),
        'lastNames': _collapseSpaces(lastNames),
        'email': email.trim().toLowerCase(),
        'password': password,
        'acceptTerms': acceptTerms,
      };
}
