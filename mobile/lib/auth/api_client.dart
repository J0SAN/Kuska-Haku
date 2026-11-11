import 'dart:convert';
import 'package:http/http.dart' as http;

import '../config.dart';
import 'registration.dart';

class AuthUser {
  final String id;
  final String email;
  final String fullName;
  final String role;

  AuthUser({
    required this.id,
    required this.email,
    required this.fullName,
    required this.role,
  });

  factory AuthUser.fromJson(Map<String, dynamic> json) => AuthUser(
        id: json['id'] as String,
        email: json['email'] as String,
        fullName: json['fullName'] as String,
        role: json['role'] as String,
      );
}

class AuthResult {
  final String token;
  final AuthUser user;

  AuthResult({required this.token, required this.user});
}

class AuthException implements Exception {
  final String message;
  final int? statusCode;

  AuthException(this.message, {this.statusCode});

  @override
  String toString() => message;
}

class AuthApiClient {
  final String backendHost;
  final http.Client _client;

  AuthApiClient({this.backendHost = AppConfig.backendHost, http.Client? client})
      : _client = client ?? http.Client();

  Uri _uri(String path) => Uri.parse('http://$backendHost$path');

  Future<AuthResult> login(String email, String password) async {
    final response = await _post('/api/v1/auth/login', {
      'email': email.trim(),
      'password': password,
    });

    if (response.statusCode != 200) {
      _throwFromResponse(response, fallback: 'Credenciales inválidas');
    }
    return _parseAuth(response);
  }

  Future<AuthResult> register(RegisterInput input) async {
    final response = await _post('/api/v1/auth/register', input.toJson());

    if (response.statusCode != 201) {
      _throwFromResponse(response, fallback: 'No se pudo crear la cuenta');
    }
    return _parseAuth(response);
  }

  Future<http.Response> _post(String path, Map<String, dynamic> body) async {
    try {
      return await _client.post(
        _uri(path),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode(body),
      );
    } catch (_) {
      throw AuthException(
        'No pudimos conectar con el servidor. Revisa tu conexión.',
      );
    }
  }

  AuthResult _parseAuth(http.Response response) {
    final body = jsonDecode(response.body) as Map<String, dynamic>;
    return AuthResult(
      token: body['token'] as String,
      user: AuthUser.fromJson(body['user'] as Map<String, dynamic>),
    );
  }

  Never _throwFromResponse(
    http.Response response, {
    required String fallback,
  }) {
    var message = fallback;
    try {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      message = (body['message'] ?? body['error'] ?? fallback) as String;
    } catch (_) {
      if (response.body.trim().isNotEmpty) {
        message = response.body.trim();
      }
    }
    throw AuthException(message, statusCode: response.statusCode);
  }
}
