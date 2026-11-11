import 'package:flutter_test/flutter_test.dart';
import 'package:kuska_haku_pasajero/auth/registration.dart';

void main() {
  group('validateDni', () {
    test('acepta ocho dígitos', () {
      expect(validateDni('12345678'), isNull);
      expect(validateDni(' 12345678 '), isNull);
    });

    test('rechaza longitudes y caracteres inválidos', () {
      expect(validateDni('1234567'), isNotNull);
      expect(validateDni('123456789'), isNotNull);
      expect(validateDni('1234567A'), isNotNull);
      expect(validateDni(''), isNotNull);
      expect(validateDni('１２３４５６７８'), isNotNull);
    });
  });

  group('validateCheckDigit', () {
    // Misma tabla que backend/internal/domain/dni_test.go.
    final cases = <List<dynamic>>[
      ['12345678', '1', true],
      ['12345678', 'E', true],
      ['12345678', 'e', true],
      ['12345678', '2', false],
      ['87654321', '0', true],
      ['87654321', 'D', true],
      ['00000001', '4', true],
      ['11111111', '7', true],
      ['11111111', 'a', true],
      ['99999999', 'I', true],
      ['12345678', '', false],
      ['12345678', '11', false],
    ];

    for (final item in cases) {
      final dni = item[0] as String;
      final check = item[1] as String;
      final valid = item[2] as bool;

      test('$dni con "$check" => ${valid ? 'válido' : 'inválido'}', () {
        final result = validateCheckDigit(dni, check);
        expect(result, valid ? isNull : isNotNull);
      });
    }

    test('no muestra error de dígito si el DNI aún es inválido', () {
      expect(validateCheckDigit('1234', 'X'), isNull);
    });
  });

  group('validatePersonName', () {
    test('acepta letras, ñ, tildes, apóstrofo y guion', () {
      expect(validatePersonName('Roxana Karina'), isNull);
      expect(validatePersonName('Ñahui'), isNull);
      expect(validatePersonName("O'Connor"), isNull);
      expect(validatePersonName('Pérez-Luna'), isNull);
    });

    test('rechaza dígitos, emojis, un carácter y más de 100', () {
      expect(validatePersonName('R2D2'), isNotNull);
      expect(validatePersonName('😀'), isNotNull);
      expect(validatePersonName('A'), isNotNull);
      expect(validatePersonName('a' * 101), isNotNull);
    });
  });

  group('validateEmail', () {
    test('normaliza mayúsculas y espacios', () {
      expect(validateEmail(' Roxana@Correo.COM '), isNull);
    });

    test('acepta TLD largos y rechaza formatos incompletos', () {
      expect(validateEmail('a@b.online'), isNull);
      expect(validateEmail('a@b'), isNotNull);
      expect(validateEmail('sin-arroba'), isNotNull);
    });
  });

  group('contraseñas', () {
    test('exige al menos 8 caracteres', () {
      expect(validatePassword('1234567'), isNotNull);
      expect(validatePassword('12345678'), isNull);
    });

    test('comprueba la confirmación', () {
      expect(passwordsMatch('secreto123', 'secreto123'), isNull);
      expect(passwordsMatch('secreto123', 'otro'), isNotNull);
    });
  });

  group('RegisterInput', () {
    test('toJson normaliza y no envía role ni fullName', () {
      const input = RegisterInput(
        dni: ' 12345678 ',
        checkDigit: 'e',
        firstNames: '  Roxana   Karina ',
        lastNames: 'Delgado Huamán',
        email: ' Roxana@Correo.COM ',
        password: 'secreto123',
        acceptTerms: true,
      );

      final json = input.toJson();

      expect(json['dni'], '12345678');
      expect(json['checkDigit'], 'E');
      expect(json['firstNames'], 'Roxana Karina');
      expect(json['lastNames'], 'Delgado Huamán');
      expect(json['email'], 'roxana@correo.com');
      expect(json['acceptTerms'], isTrue);
      expect(json.containsKey('role'), isFalse);
      expect(json.containsKey('fullName'), isFalse);
    });
  });
}
