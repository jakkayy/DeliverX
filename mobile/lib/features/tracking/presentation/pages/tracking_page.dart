import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class TrackingPage extends StatefulWidget {
  final String orderId;
  const TrackingPage({super.key, required this.orderId});

  @override
  State<TrackingPage> createState() => _TrackingPageState();
}

class _TrackingPageState extends State<TrackingPage> {
  GoogleMapController? _mapController;
  WebSocketChannel? _channel;

  static const _initialPosition = LatLng(13.7563, 100.5018);
  final Set<Marker> _markers = {};
  LatLng? _driverPosition;

  @override
  void initState() {
    super.initState();
    _connectWebSocket();
  }

  void _connectWebSocket() {
    const wsBaseUrl = String.fromEnvironment(
      'WS_BASE_URL',
      defaultValue: 'ws://10.0.2.2:8084',
    );

    _channel = WebSocketChannel.connect(
      Uri.parse('$wsBaseUrl/ws/tracking/${widget.orderId}'),
    );

    _channel!.stream.listen((message) {
      // Expected: {"lat": 13.756, "lng": 100.501}
      // TODO: parse JSON and update driver marker
      if (mounted) {
        setState(() {
          _driverPosition = const LatLng(13.756, 100.502);
          _markers.removeWhere((m) => m.markerId.value == 'driver');
          _markers.add(Marker(
            markerId: const MarkerId('driver'),
            position: _driverPosition!,
            icon: BitmapDescriptor.defaultMarkerWithHue(
                BitmapDescriptor.hueGreen),
            infoWindow: const InfoWindow(title: 'Driver'),
          ));
        });
      }
    });
  }

  @override
  void dispose() {
    _channel?.sink.close();
    _mapController?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Live Tracking')),
      body: Stack(
        children: [
          GoogleMap(
            initialCameraPosition: const CameraPosition(
              target: _initialPosition,
              zoom: 14,
            ),
            markers: _markers,
            myLocationEnabled: true,
            myLocationButtonEnabled: true,
            onMapCreated: (c) => _mapController = c,
          ),
          Positioned(
            bottom: 0,
            left: 0,
            right: 0,
            child: Container(
              padding: const EdgeInsets.all(16),
              decoration: const BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
                boxShadow: [BoxShadow(blurRadius: 10, color: Colors.black26)],
              ),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Container(
                    width: 40,
                    height: 4,
                    decoration: BoxDecoration(
                      color: Colors.grey.shade300,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                  const SizedBox(height: 16),
                  const ListTile(
                    leading: CircleAvatar(child: Icon(Icons.two_wheeler)),
                    title: Text('Driver is on the way'),
                    subtitle: Text('Estimated arrival: 10 min'),
                    trailing: Icon(Icons.call, color: Colors.green),
                    contentPadding: EdgeInsets.zero,
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
