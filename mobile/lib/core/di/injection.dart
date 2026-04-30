import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:get_it/get_it.dart';

import '../network/api_client.dart';

final getIt = GetIt.instance;

Future<void> configureDependencies() async {
  // Core
  getIt.registerSingleton<Dio>(ApiClient.create());
  getIt.registerSingleton<FlutterSecureStorage>(const FlutterSecureStorage());

  // Repositories and BLoCs are registered here as the app grows
}
