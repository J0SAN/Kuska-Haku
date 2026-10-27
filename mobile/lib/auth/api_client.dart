import 'dart:convert';
import 'package:http/http.dart' as http;

import '../config.dart';

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
  AuthException(this.message);

  @override
  String toString() => message;
}

class AuthApiClient {
  final String backendHost;

  AuthApiClient({this.backendHost = AppConfig.backendHost});

  Uri _uri(String path) => Uri.parse('http://$backendHost$path');

  Future<AuthResult> login(String email, String password) async {
    final response = await http.post(
      _uri('/api/v1/auth/login'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'email': email, 'password': password}),
    );

    if (response.statusCode != 200) {
      String errorMessage = 'Credenciales inválidas';
      try {
        final body = jsonDecode(response.body) as Map<String, dynamic>;
        errorMessage =
            (body['error'] ?? body['message'] ?? errorMessage) as String;
      } catch (_) {
        if (response.body.trim().isNotEmpty) {
          errorMessage = response.body.trim();
        }
      }
      throw AuthException(errorMessage);
    }

    final body = jsonDecode(response.body) as Map<String, dynamic>;
    return AuthResult(
      token: body['token'] as String,
      user: AuthUser.fromJson(body['user'] as Map<String, dynamic>),
    );
  }
}
