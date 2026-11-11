import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../theme/app_colors.dart';
import 'api_client.dart';
import 'registration.dart';

class RegisterScreen extends StatefulWidget {
  final AuthApiClient apiClient;
  final void Function(AuthResult result) onRegistered;

  const RegisterScreen({
    super.key,
    required this.apiClient,
    required this.onRegistered,
  });

  @override
  State<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends State<RegisterScreen> {
  final _formKey = GlobalKey<FormState>();
  final _dniController = TextEditingController();
  final _checkDigitController = TextEditingController();
  final _firstNamesController = TextEditingController();
  final _lastNamesController = TextEditingController();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmController = TextEditingController();

  bool _acceptTerms = false;
  bool _loading = false;
  bool _obscurePassword = true;
  bool _obscureConfirm = true;
  String? _error;
  bool _showLoginAction = false;

  @override
  void initState() {
    super.initState();
    for (final controller in [
      _dniController,
      _checkDigitController,
      _firstNamesController,
      _lastNamesController,
      _emailController,
      _passwordController,
      _confirmController,
    ]) {
      controller.addListener(_refresh);
    }
  }

  @override
  void dispose() {
    _dniController.dispose();
    _checkDigitController.dispose();
    _firstNamesController.dispose();
    _lastNamesController.dispose();
    _emailController.dispose();
    _passwordController.dispose();
    _confirmController.dispose();
    super.dispose();
  }

  void _refresh() {
    if (mounted) setState(() {});
  }

  bool get _dniAndDigitValid =>
      validateDni(_dniController.text) == null &&
      validateCheckDigit(_dniController.text, _checkDigitController.text) ==
          null;

  bool get _isValid =>
      _dniAndDigitValid &&
      validatePersonName(_firstNamesController.text) == null &&
      validatePersonName(_lastNamesController.text) == null &&
      validateEmail(_emailController.text) == null &&
      validatePassword(_passwordController.text) == null &&
      passwordsMatch(
            _passwordController.text,
            _confirmController.text,
          ) ==
          null &&
      _acceptTerms;

  Future<void> _submit() async {
    if (_loading || !_formKey.currentState!.validate() || !_acceptTerms) {
      return;
    }

    setState(() {
      _loading = true;
      _error = null;
      _showLoginAction = false;
    });

    try {
      final result = await widget.apiClient.register(
        RegisterInput(
          dni: _dniController.text,
          checkDigit: _checkDigitController.text,
          firstNames: _firstNamesController.text,
          lastNames: _lastNamesController.text,
          email: _emailController.text,
          password: _passwordController.text,
          acceptTerms: _acceptTerms,
        ),
      );
      if (!mounted) return;
      widget.onRegistered(result);
      Navigator.of(context).pop();
    } on AuthException catch (e) {
      if (!mounted) return;
      setState(() {
        _error = e.message;
        _showLoginAction = e.statusCode == 409;
      });
    } catch (_) {
      if (!mounted) return;
      setState(() {
        _error = 'No pudimos conectar con el servidor. Revisa tu conexión.';
        _showLoginAction = false;
      });
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  void _goToLogin() {
    Navigator.of(context).pop();
  }

  void _showCheckDigitHelp() {
    showModalBottomSheet<void>(
      context: context,
      showDragHandle: true,
      builder: (context) => const SafeArea(
        child: Padding(
          padding: EdgeInsets.fromLTRB(24, 0, 24, 24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '¿Dónde está el dígito verificador?',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              SizedBox(height: 12),
              Text(
                'En el DNI azul aparece impreso a la derecha del número, '
                'después de un guion. En el DNIe está junto al número, en la '
                'zona de lectura del reverso.',
              ),
              SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    child: _DocumentSample(
                      title: 'DNI azul',
                      number: '46027897 - 4',
                    ),
                  ),
                  SizedBox(width: 12),
                  Expanded(
                    child: _DocumentSample(
                      title: 'DNIe',
                      number: '46027897 4',
                    ),
                  ),
                ],
              ),
              SizedBox(height: 16),
              Text(
                'El dígito no es un dato secreto: solo confirma que copiaste '
                'bien el número de tu documento.',
                style: TextStyle(color: AppColors.inactive),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _showPrivacyPolicy() {
    showModalBottomSheet<void>(
      context: context,
      showDragHandle: true,
      isScrollControlled: true,
      builder: (context) => const SafeArea(
        child: Padding(
          padding: EdgeInsets.fromLTRB(24, 0, 24, 24),
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Política de privacidad',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
                SizedBox(height: 12),
                Text(
                  'Tus datos personales (DNI, nombres, apellidos y correo) se '
                  'usan únicamente para crear y proteger tu cuenta en Kuska '
                  'Haku, conforme a la Ley N.º 29733 de Protección de Datos '
                  'Personales.',
                ),
                SizedBox(height: 12),
                Text(
                  'No compartimos tus datos con terceros ni los usamos para '
                  'publicidad. Puedes solicitar su actualización o eliminación '
                  'escribiendo a soporte.',
                ),
                SizedBox(height: 12),
                Text(
                  'La versión completa de la política se publicará antes del '
                  'lanzamiento en las tiendas de aplicaciones.',
                  style: TextStyle(color: AppColors.inactive),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  InputDecoration _decoration({
    required String label,
    String? hint,
    Widget? suffixIcon,
  }) {
    return InputDecoration(
      labelText: label,
      hintText: hint,
      suffixIcon: suffixIcon,
      border: const OutlineInputBorder(),
      counterText: '',
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        backgroundColor: AppColors.primary,
        foregroundColor: AppColors.onPrimary,
        title: const Text('Crear cuenta'),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24),
          child: Form(
            key: _formKey,
            autovalidateMode: AutovalidateMode.onUserInteraction,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                const Text(
                  'Documento de identidad',
                  style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 8),
                Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Expanded(
                      child: TextFormField(
                        key: const Key('dni-field'),
                        controller: _dniController,
                        keyboardType: TextInputType.number,
                        maxLength: 8,
                        inputFormatters: [
                          FilteringTextInputFormatter.digitsOnly,
                          LengthLimitingTextInputFormatter(8),
                        ],
                        decoration: _decoration(label: 'DNI'),
                        validator: validateDni,
                      ),
                    ),
                    const SizedBox(width: 8),
                    SizedBox(
                      width: 88,
                      child: TextFormField(
                        key: const Key('check-digit-field'),
                        controller: _checkDigitController,
                        maxLength: 1,
                        textCapitalization: TextCapitalization.characters,
                        inputFormatters: [
                          FilteringTextInputFormatter.allow(
                            RegExp(r'[0-9A-Za-z]'),
                          ),
                          LengthLimitingTextInputFormatter(1),
                        ],
                        decoration: _decoration(label: 'Dígito'),
                        validator: (value) => validateCheckDigit(
                          _dniController.text,
                          value,
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.only(top: 4),
                      child: IconButton(
                        tooltip: '¿Dónde está el dígito verificador?',
                        onPressed: _showCheckDigitHelp,
                        icon: const Icon(
                          Icons.info_outline,
                          color: AppColors.primary,
                        ),
                      ),
                    ),
                  ],
                ),
                if (_dniAndDigitValid)
                  Semantics(
                    liveRegion: true,
                    child: const Row(
                      children: [
                        Icon(
                          Icons.check_circle,
                          size: 18,
                          color: AppColors.primary,
                        ),
                        SizedBox(width: 6),
                        Text(
                          'DNI válido',
                          style: TextStyle(color: AppColors.primary),
                        ),
                      ],
                    ),
                  ),
                const SizedBox(height: 16),
                TextFormField(
                  key: const Key('first-names-field'),
                  controller: _firstNamesController,
                  textCapitalization: TextCapitalization.words,
                  autofillHints: const [AutofillHints.givenName],
                  decoration: _decoration(label: 'Nombres'),
                  validator: validatePersonName,
                ),
                const SizedBox(height: 12),
                TextFormField(
                  key: const Key('last-names-field'),
                  controller: _lastNamesController,
                  textCapitalization: TextCapitalization.words,
                  autofillHints: const [AutofillHints.familyName],
                  decoration: _decoration(label: 'Apellidos'),
                  validator: validatePersonName,
                ),
                const SizedBox(height: 12),
                TextFormField(
                  key: const Key('email-field'),
                  controller: _emailController,
                  keyboardType: TextInputType.emailAddress,
                  autocorrect: false,
                  textCapitalization: TextCapitalization.none,
                  autofillHints: const [AutofillHints.email],
                  decoration: _decoration(label: 'Correo electrónico'),
                  validator: validateEmail,
                ),
                const SizedBox(height: 12),
                TextFormField(
                  key: const Key('password-field'),
                  controller: _passwordController,
                  obscureText: _obscurePassword,
                  autofillHints: const [AutofillHints.newPassword],
                  decoration: _decoration(
                    label: 'Contraseña',
                    suffixIcon: IconButton(
                      tooltip: _obscurePassword
                          ? 'Mostrar contraseña'
                          : 'Ocultar contraseña',
                      onPressed: () => setState(
                        () => _obscurePassword = !_obscurePassword,
                      ),
                      icon: Icon(
                        _obscurePassword
                            ? Icons.visibility_off
                            : Icons.visibility,
                      ),
                    ),
                  ),
                  validator: validatePassword,
                ),
                const SizedBox(height: 12),
                TextFormField(
                  key: const Key('confirm-password-field'),
                  controller: _confirmController,
                  obscureText: _obscureConfirm,
                  autofillHints: const [AutofillHints.newPassword],
                  decoration: _decoration(
                    label: 'Confirmar contraseña',
                    suffixIcon: IconButton(
                      tooltip: _obscureConfirm
                          ? 'Mostrar contraseña'
                          : 'Ocultar contraseña',
                      onPressed: () => setState(
                        () => _obscureConfirm = !_obscureConfirm,
                      ),
                      icon: Icon(
                        _obscureConfirm
                            ? Icons.visibility_off
                            : Icons.visibility,
                      ),
                    ),
                  ),
                  validator: (value) => passwordsMatch(
                    _passwordController.text,
                    value,
                  ),
                ),
                const SizedBox(height: 8),
                CheckboxListTile(
                  key: const Key('accept-terms'),
                  value: _acceptTerms,
                  onChanged: (value) => setState(
                    () => _acceptTerms = value ?? false,
                  ),
                  controlAffinity: ListTileControlAffinity.leading,
                  contentPadding: EdgeInsets.zero,
                  title: const Text(
                    'Acepto la política de privacidad (Ley 29733)',
                  ),
                ),
                Align(
                  alignment: Alignment.centerLeft,
                  child: TextButton(
                    onPressed: _showPrivacyPolicy,
                    child: const Text('Ver política de privacidad'),
                  ),
                ),
                if (_error != null) ...[
                  const SizedBox(height: 4),
                  Text(
                    _error!,
                    style: const TextStyle(color: AppColors.error),
                  ),
                  if (_showLoginAction)
                    Align(
                      alignment: Alignment.centerLeft,
                      child: TextButton(
                        onPressed: _goToLogin,
                        child: const Text('Iniciar sesión'),
                      ),
                    ),
                ],
                const SizedBox(height: 12),
                ElevatedButton(
                  key: const Key('submit-button'),
                  onPressed: _isValid && !_loading ? _submit : null,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.primary,
                    foregroundColor: AppColors.onPrimary,
                    disabledBackgroundColor: AppColors.inactive,
                    minimumSize: const Size.fromHeight(48),
                  ),
                  child: _loading
                      ? const SizedBox(
                          height: 20,
                          width: 20,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            color: AppColors.onPrimary,
                          ),
                        )
                      : const Text('Crear cuenta'),
                ),
                const SizedBox(height: 8),
                TextButton(
                  onPressed: _goToLogin,
                  child: const Text('¿Ya tienes cuenta? Inicia sesión'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class _DocumentSample extends StatelessWidget {
  final String title;
  final String number;

  const _DocumentSample({required this.title, required this.number});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        border: Border.all(color: AppColors.primary),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(title, style: const TextStyle(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          Text(
            number,
            style: const TextStyle(
              fontFamily: 'monospace',
              letterSpacing: 1.5,
            ),
          ),
        ],
      ),
    );
  }
}
