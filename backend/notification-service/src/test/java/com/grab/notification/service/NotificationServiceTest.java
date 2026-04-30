package com.grab.notification.service;

import com.grab.notification.model.Notification;
import com.grab.notification.repository.NotificationRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Map;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class NotificationServiceTest {

    @Mock private NotificationRepository notificationRepository;
    @Mock private FcmService fcmService;
    @InjectMocks private NotificationService notificationService;

    private UUID userId;

    @BeforeEach
    void setUp() {
        userId = UUID.randomUUID();
    }

    @Test
    @DisplayName("Send and save notification successfully")
    void sendAndSave_success() {
        when(fcmService.sendToToken(anyString(), anyString(), anyString(), any()))
                .thenReturn("msg-id-123");
        when(notificationRepository.save(any(Notification.class)))
                .thenAnswer(i -> i.getArgument(0));

        notificationService.sendAndSave(
                userId, "fcm-token-abc",
                "Test Title", "Test Body",
                Notification.NotificationType.ORDER_CREATED,
                UUID.randomUUID().toString(),
                Map.of("key", "value")
        );

        verify(fcmService).sendToToken(eq("fcm-token-abc"), eq("Test Title"), eq("Test Body"), any());
        verify(notificationRepository).save(any(Notification.class));
    }

    @Test
    @DisplayName("Send and save still saves even when FCM token is null")
    void sendAndSave_nullToken_stillSaves() {
        when(fcmService.sendToToken(isNull(), anyString(), anyString(), any())).thenReturn(null);
        when(notificationRepository.save(any(Notification.class)))
                .thenAnswer(i -> i.getArgument(0));

        notificationService.sendAndSave(
                userId, null,
                "Title", "Body",
                Notification.NotificationType.GENERAL, null, null
        );

        verify(notificationRepository).save(any(Notification.class));
    }

    @Test
    @DisplayName("Get unread count returns correct value")
    void getUnreadCount() {
        when(notificationRepository.countByUserIdAndIsReadFalse(userId)).thenReturn(5L);

        long count = notificationService.getUnreadCount(userId);

        assertThat(count).isEqualTo(5L);
    }

    @Test
    @DisplayName("Mark all read calls repository")
    void markAllRead() {
        when(notificationRepository.markAllReadByUserId(userId)).thenReturn(3);

        notificationService.markAllRead(userId);

        verify(notificationRepository).markAllReadByUserId(userId);
    }
}
