import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:http/http.dart' as http;

import '../../config.dart';
import '../tracking/track_point.dart';
import '../tracking/tracking_uploader.dart';

class DriverApiException implements Exception {
  final String message;
  DriverApiException(this.message);

  @override
  String toString() => message;
}

class MyBusInfo {
  final String id;
  final String busId;
  final String status;
  final String? routeId;
  final String? routeName;

  MyBusInfo({
    required this.id,
    required this.busId,
    required this.status,
    this.routeId,
    this.routeName,
  });

  factory MyBusInfo.fromJson(Map<String, dynamic> json) {
    final bus = json['bus'] as Map<String, dynamic>;
    final route = json['route'] as Map<String, dynamic>?;
    return MyBusInfo(
      id: bus['id'] as String,
      busId: bus['busId'] as String,
      status: (bus['status'] as String?) ?? 'active',
      routeId: (route?['id'] ?? bus['routeId']) as String?,
      routeName: route?['name'] as String?,
    );
  }
}

class IncomeItem {
  final int passengers;
  final double amount;
  final String notes;
  final DateTime createdAt;

  IncomeItem({
    required this.passengers,
    required this.amount,
    required this.notes,
    required this.createdAt,
  });

  factory IncomeItem.fromJson(Map<String, dynamic> json) => IncomeItem(
        passengers: json['passengers'] as int,
        amount: (json['amount'] as num).toDouble(),
        notes: (json['notes'] as String?) ?? '',
        createdAt: DateTime.parse(json['createdAt'] as String).toLocal(),
      );
}

class TodayIncomes {
  final int totalPassengers;
  final double totalAmount;
  final List<IncomeItem> records;

  TodayIncomes({
    required this.totalPassengers,
    required this.totalAmount,
    required this.records,
  });

  factory TodayIncomes.fromJson(Map<String, dynamic> json) => TodayIncomes(
        totalPassengers: json['totalPassengers'] as int,
        totalAmount: (json['totalAmount'] as num).toDouble(),
        records: ((json['records'] as List<dynamic>?) ?? [])
            .map((r) => IncomeItem.fromJson(r as Map<String, dynamic>))
            .toList(),
      );
}

class DriverApi {
  DriverApi({required this.token});
  final String token;

  static const _timeout = Duration(seconds: 10);

  Map<String, String> get _headers => {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      };

  Uri _uri(String path) => Uri.parse('http://${AppConfig.backendHost}$path');

  String _extractError(http.Response res) {
    try {
      final body = jsonDecode(res.body) as Map<String, dynamic>;
      return (body['error'] ?? body['message'] ?? 'Error ${res.statusCode}')
          as String;
    } catch (_) {
      return 'Error ${res.statusCode}';
    }
  }

  Future<MyBusInfo> getMyBus() async {
    final res = await http
        .get(_uri('/api/v1/driver/my-bus'), headers: _headers)
        .timeout(_timeout);
    if (res.statusCode != 200) throw DriverApiException(_extractError(res));
    return MyBusInfo.fromJson(jsonDecode(res.body) as Map<String, dynamic>);
  }

  Future<TodayIncomes> getTodayIncomes() async {
    final res = await http
        .get(_uri('/api/v1/driver/incomes'), headers: _headers)
        .timeout(_timeout);
    if (res.statusCode != 200) throw DriverApiException(_extractError(res));
    return TodayIncomes.fromJson(jsonDecode(res.body) as Map<String, dynamic>);
  }

  Future<void> createIncome({
    required int passengers,
    required double amount,
    required String notes,
  }) async {
    final res = await http
        .post(
          _uri('/api/v1/driver/incomes'),
          headers: _headers,
          body: jsonEncode({
            'passengers': passengers,
            'amount': amount,
            'notes': notes,
          }),
        )
        .timeout(_timeout);
    if (res.statusCode != 201) throw DriverApiException(_extractError(res));
  }

  /// Envía un lote de coordenadas GPS a POST /api/v1/driver/positions.
  Future<UploadOutcome> sendPositions(List<TrackPoint> batch) async {
    try {
      final res = await http
          .post(
            _uri('/api/v1/driver/positions'),
            headers: _headers,
            body: jsonEncode({
              'positions': batch.map((p) => p.toJson()).toList(),
            }),
          )
          .timeout(_timeout);
      switch (res.statusCode) {
        case 200:
          return UploadOutcome.sent;
        case 400:
          return UploadOutcome.dropBatch;
        case 401:
          return UploadOutcome.authExpired;
        case 409:
          return UploadOutcome.noBus;
        default:
          return UploadOutcome.retry;
      }
    } on TimeoutException {
      return UploadOutcome.retry;
    } on SocketException {
      return UploadOutcome.retry;
    } on http.ClientException {
      return UploadOutcome.retry;
    }
  }
}
