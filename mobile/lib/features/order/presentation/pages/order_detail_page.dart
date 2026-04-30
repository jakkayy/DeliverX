import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class OrderDetailPage extends StatelessWidget {
  final String orderId;
  const OrderDetailPage({super.key, required this.orderId});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Order #${orderId.substring(0, 8)}')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          _StatusCard(),
          const SizedBox(height: 16),
          _LocationCard(),
          const SizedBox(height: 16),
          _DriverCard(orderId: orderId),
          const SizedBox(height: 24),
          ElevatedButton.icon(
            onPressed: () => context.push('/tracking/$orderId'),
            icon: const Icon(Icons.map),
            label: const Text('Track on Map'),
          ),
        ],
      ),
    );
  }
}

class _StatusCard extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
              decoration: BoxDecoration(
                color: Colors.orange.shade100,
                borderRadius: BorderRadius.circular(20),
              ),
              child: Text('IN_TRANSIT',
                  style: TextStyle(
                      color: Colors.orange.shade800,
                      fontWeight: FontWeight.bold)),
            ),
            const Spacer(),
            Text('Est. 15 min', style: Theme.of(context).textTheme.bodyMedium),
          ],
        ),
      ),
    );
  }
}

class _LocationCard extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return const Card(
      child: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          children: [
            ListTile(
              leading: Icon(Icons.location_on, color: Colors.green),
              title: Text('Pickup'),
              subtitle: Text('123 Sukhumvit Rd, Bangkok'),
              contentPadding: EdgeInsets.zero,
            ),
            Divider(),
            ListTile(
              leading: Icon(Icons.location_on, color: Colors.red),
              title: Text('Dropoff'),
              subtitle: Text('456 Silom Rd, Bangkok'),
              contentPadding: EdgeInsets.zero,
            ),
          ],
        ),
      ),
    );
  }
}

class _DriverCard extends StatelessWidget {
  final String orderId;
  const _DriverCard({required this.orderId});

  @override
  Widget build(BuildContext context) {
    return const Card(
      child: ListTile(
        leading: CircleAvatar(child: Icon(Icons.person)),
        title: Text('Prasert Nakorn'),
        subtitle: Text('Honda PCX • กข-1234 • ⭐ 4.85'),
        trailing: Icon(Icons.call),
      ),
    );
  }
}
