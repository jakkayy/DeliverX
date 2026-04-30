import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class CreateOrderPage extends StatefulWidget {
  const CreateOrderPage({super.key});

  @override
  State<CreateOrderPage> createState() => _CreateOrderPageState();
}

class _CreateOrderPageState extends State<CreateOrderPage> {
  final _formKey = GlobalKey<FormState>();
  final _pickupController = TextEditingController();
  final _dropoffController = TextEditingController();
  final _noteController = TextEditingController();
  String _vehicleType = 'MOTORCYCLE';
  bool _isLoading = false;

  @override
  void dispose() {
    _pickupController.dispose();
    _dropoffController.dispose();
    _noteController.dispose();
    super.dispose();
  }

  Future<void> _createOrder() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _isLoading = true);
    // TODO: dispatch CreateOrderBloc event
    await Future.delayed(const Duration(seconds: 1));
    setState(() => _isLoading = false);
    if (mounted) context.go('/home');
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('New Order')),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            Text('Pickup & Dropoff',
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    )),
            const SizedBox(height: 12),
            TextFormField(
              controller: _pickupController,
              decoration: const InputDecoration(
                labelText: 'Pickup Address',
                prefixIcon: Icon(Icons.location_on, color: Colors.green),
              ),
              validator: (v) => v == null || v.isEmpty ? 'Required' : null,
            ),
            const SizedBox(height: 12),
            TextFormField(
              controller: _dropoffController,
              decoration: const InputDecoration(
                labelText: 'Dropoff Address',
                prefixIcon: Icon(Icons.location_on, color: Colors.red),
              ),
              validator: (v) => v == null || v.isEmpty ? 'Required' : null,
            ),
            const SizedBox(height: 24),
            Text('Vehicle Type',
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    )),
            const SizedBox(height: 12),
            SegmentedButton<String>(
              segments: const [
                ButtonSegment(
                    value: 'MOTORCYCLE', label: Text('Motorcycle'), icon: Icon(Icons.two_wheeler)),
                ButtonSegment(
                    value: 'CAR', label: Text('Car'), icon: Icon(Icons.directions_car)),
                ButtonSegment(
                    value: 'VAN', label: Text('Van'), icon: Icon(Icons.airport_shuttle)),
              ],
              selected: {_vehicleType},
              onSelectionChanged: (v) => setState(() => _vehicleType = v.first),
            ),
            const SizedBox(height: 24),
            TextFormField(
              controller: _noteController,
              maxLines: 3,
              decoration: const InputDecoration(
                labelText: 'Note to driver (optional)',
                prefixIcon: Icon(Icons.note),
              ),
            ),
            const SizedBox(height: 32),
            ElevatedButton(
              onPressed: _isLoading ? null : _createOrder,
              child: _isLoading
                  ? const CircularProgressIndicator(color: Colors.white)
                  : const Text('Find Driver'),
            ),
          ],
        ),
      ),
    );
  }
}
