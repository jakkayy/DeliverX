package com.grab.notification.config;

import com.google.auth.oauth2.GoogleCredentials;
import com.google.firebase.FirebaseApp;
import com.google.firebase.FirebaseOptions;
import com.google.firebase.messaging.FirebaseMessaging;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.io.ClassPathResource;

import java.io.IOException;
import java.io.InputStream;

@Slf4j
@Configuration
public class FirebaseConfig {

    @Value("${firebase.project-id}")
    private String projectId;

    @Value("${firebase.service-account-path:firebase-service-account.json}")
    private String serviceAccountPath;

    @Bean
    public FirebaseMessaging firebaseMessaging() throws IOException {
        if (FirebaseApp.getApps().isEmpty()) {
            InputStream serviceAccount = loadServiceAccount();
            FirebaseOptions options = FirebaseOptions.builder()
                    .setCredentials(GoogleCredentials.fromStream(serviceAccount))
                    .setProjectId(projectId)
                    .build();
            FirebaseApp.initializeApp(options);
            log.info("Firebase initialized with project: {}", projectId);
        }
        return FirebaseMessaging.getInstance();
    }

    private InputStream loadServiceAccount() throws IOException {
        try {
            return new ClassPathResource(serviceAccountPath).getInputStream();
        } catch (IOException e) {
            log.warn("Service account file not found at classpath, using Application Default Credentials");
            return GoogleCredentials.getApplicationDefault().getClass()
                    .getResourceAsStream(serviceAccountPath);
        }
    }
}
