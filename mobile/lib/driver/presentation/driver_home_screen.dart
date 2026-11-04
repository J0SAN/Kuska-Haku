import 'package:flutter/material.dart';

import '../../auth/api_client.dart';
import '../data/driver_api.dart';
import '../tracking/shift_controller.dart';

class DriverHomeScreen extends StatefulWidget {
  final AuthResult session;
  final VoidCallback onLogout;

  const DriverHomeScreen({
    super.key,
    required this.session,
    required this.onLogout,
  });

  @override
  State<DriverHomeScreen> createState() => _DriverHomeScreenState();
}

class _DriverHomeScreenState extends State<DriverHomeScreen> {
  static const _greenColor = Color(0xFF2E7D32);

  late final DriverApi _api;
  late final ShiftController _shift;

  MyBusInfo? _bus;
  String? _busError;
  TodayIncomes? _today;

  final _passengersController = TextEditingController();
  final _amountController = TextEditingController();
  final _notesController = TextEditingController();
  bool _saving = false;

  @override
  void initState() {
    super.initState();
    _api = DriverApi(token: widget.session.token);
    _shift = ShiftController(api: _api);
    _loadData();
  }

  Future<void> _loadData() async {
    try {
      final bus = await _api.getMyBus();
      if (mounted) {
        setState(() {
          _bus = bus;
          _busError = null;
        });
      }
    } on DriverApiException catch (e) {
      if (mounted) setState(() => _busError = e.message);
    } catch (_) {
      if (mounted) {
        setState(() => _busError = 'No se pudo conectar con el servidor');
      }
    }

    try {
      final today = await _api.getTodayIncomes();
      if (mounted) setState(() => _today = today);
    } catch (_) {}
  }

  Future<void> _submitIncome() async {
    final passengers = int.tryParse(_passengersController.text.trim()) ?? 0;
    final amount =
        double.tryParse(_amountController.text.trim().replaceAll(',', '.')) ??
            0;

    setState(() => _saving = true);
    String feedback = 'Registro guardado';
    try {
      await _api.createIncome(
        passengers: passengers,
        amount: amount,
        notes: _notesController.text.trim(),
      );
      _passengersController.clear();
      _amountController.clear();
      _notesController.clear();
      await _loadData();
    } on DriverApiException catch (e) {
      feedback = e.message;
    } catch (_) {
      feedback = 'No se pudo guardar el registro';
    }
    if (!mounted) return;
    setState(() => _saving = false);
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(feedback)),
    );
  }

  Future<void> _handleLogout() async {
    await _shift.stop();
    widget.onLogout();
  }

  @override
  void dispose() {
    _shift.dispose();
    _passengersController.dispose();
    _amountController.dispose();
    _notesController.dispose();
    super.dispose();
  }

  String _formatPen(double value) => 'S/ ${value.toStringAsFixed(2)}';

  String _formatTime(DateTime d) => '${d.hour.toString().padLeft(2, '0')}:'
      '${d.minute.toString().padLeft(2, '0')}:'
      '${d.second.toString().padLeft(2, '0')}';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.session.user.fullName),
        backgroundColor: _greenColor,
        foregroundColor: Colors.white,
        actions: [
          IconButton(
            tooltip: 'Cerrar sesión',
            icon: const Icon(Icons.logout),
            onPressed: _handleLogout,
          ),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          _buildBusCard(),
          const SizedBox(height: 12),
          _buildShiftCard(),
          const SizedBox(height: 12),
          _buildIncomeFormCard(),
          const SizedBox(height: 12),
          _buildSummaryCard(),
        ],
      ),
    );
  }

  Widget _buildBusCard() {
    final bus = _bus;
    return Card(
      child: ListTile(
        leading: const Icon(Icons.directions_bus, color: _greenColor, size: 32),
        title: Text(
          bus == null ? 'Sin bus asignado' : 'Unidad ${bus.busId}',
          style: const TextStyle(fontWeight: FontWeight.bold),
        ),
        subtitle: Text(
          bus != null
              ? 'Ruta: ${bus.routeName ?? 'Sin ruta'}'
              : (_busError ?? 'Cargando...'),
        ),
        trailing: IconButton(
          tooltip: 'Actualizar',
          icon: const Icon(Icons.refresh),
          onPressed: _loadData,
        ),
      ),
    );
  }

  Widget _buildShiftCard() {
    return ListenableBuilder(
      listenable: _shift,
      builder: (context, _) {
        final isRunning = _shift.running;
        return Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(
                      isRunning ? Icons.gps_fixed : Icons.gps_off,
                      color: isRunning ? _greenColor : Colors.grey,
                    ),
                    const SizedBox(width: 8),
                    Text(
                      isRunning ? 'Turno activo' : 'Turno detenido',
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 8),
                Text('Puntos pendientes de envío: ${_shift.pending}'),
                Text(
                  _shift.lastSentAt == null
                      ? 'Aún no se envió ningún punto'
                      : 'Último envío: ${_formatTime(_shift.lastSentAt!)}',
                ),
                if (_shift.buffer.dropped > 0)
                  Text(
                    'Se descartaron ${_shift.buffer.dropped} puntos antiguos por falta prolongada de conexión',
                    style: const TextStyle(color: Colors.orange),
                  ),
                if (_shift.gpsMessage != null)
                  Padding(
                    padding: const EdgeInsets.only(top: 8),
                    child: Text(
                      _shift.gpsMessage!,
                      style: const TextStyle(color: Colors.red),
                    ),
                  ),
                if (_shift.uploadMessage != null)
                  Padding(
                    padding: const EdgeInsets.only(top: 8),
                    child: Text(
                      _shift.uploadMessage!,
                      style: const TextStyle(color: Colors.orange),
                    ),
                  ),
                if (_shift.sessionExpired)
                  Padding(
                    padding: const EdgeInsets.only(top: 8),
                    child: FilledButton(
                      onPressed: widget.onLogout,
                      child: const Text('Volver a ingresar'),
                    ),
                  ),
                const SizedBox(height: 12),
                SizedBox(
                  width: double.infinity,
                  child: FilledButton.icon(
                    style: FilledButton.styleFrom(
                      backgroundColor:
                          isRunning ? Colors.red.shade700 : _greenColor,
                      padding: const EdgeInsets.symmetric(vertical: 14),
                    ),
                    icon: Icon(isRunning ? Icons.stop : Icons.play_arrow),
                    label: Text(
                      isRunning ? 'Finalizar turno' : 'Iniciar turno',
                    ),
                    onPressed: _bus == null && !isRunning
                        ? null
                        : () => isRunning ? _shift.stop() : _shift.start(),
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildIncomeFormCard() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Text(
              'Registrar pasajeros e ingreso',
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _passengersController,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(
                labelText: 'Pasajeros que subieron',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 10),
            TextField(
              controller: _amountController,
              keyboardType:
                  const TextInputType.numberWithOptions(decimal: true),
              decoration: const InputDecoration(
                labelText: 'Monto recaudado (S/)',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 10),
            TextField(
              controller: _notesController,
              maxLength: 300,
              decoration: const InputDecoration(
                labelText: 'Notas (opcional)',
                border: OutlineInputBorder(),
              ),
            ),
            FilledButton(
              style: FilledButton.styleFrom(backgroundColor: _greenColor),
              onPressed: _saving || _bus == null ? null : _submitIncome,
              child: _saving
                  ? const SizedBox(
                      height: 18,
                      width: 18,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        color: Colors.white,
                      ),
                    )
                  : const Text('Registrar'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSummaryCard() {
    final today = _today;
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Resumen de hoy',
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
            ),
            const SizedBox(height: 8),
            if (today == null)
              const Text('Sin datos todavía')
            else ...[
              Text(
                '${today.totalPassengers} pasajeros  ·  ${_formatPen(today.totalAmount)}',
                style: const TextStyle(
                  fontSize: 20,
                  color: _greenColor,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              if (today.records.isEmpty)
                const Text('Todavía no registraste ingresos hoy.'),
              for (final r in today.records)
                ListTile(
                  dense: true,
                  contentPadding: EdgeInsets.zero,
                  leading: const Icon(Icons.receipt_long, size: 20),
                  title: Text(
                    '${r.passengers} pasajeros  ·  ${_formatPen(r.amount)}',
                  ),
                  subtitle: Text(
                    r.notes.isEmpty
                        ? _formatTime(r.createdAt)
                        : '${_formatTime(r.createdAt)}  ·  ${r.notes}',
                  ),
                ),
            ],
          ],
        ),
      ),
    );
  }
}
