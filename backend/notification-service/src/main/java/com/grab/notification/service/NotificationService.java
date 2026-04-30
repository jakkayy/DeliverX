package com.grab.notification.service;

import com.grab.notification.dto.NotificationResponse;
import com.grab.notification.model.Notification;
import com.grab.notification.repository.NotificationRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Map;
import java.util.UUID;

@Slf4j
@Service
@RequiredArgsConstructor
public class NotificationService {

    private final NotificationRepository notificationRepository;
    private final FcmService fcmService;

    @Transactional
    public void sendAndSave(UUID userId, String fcmToken, String title, String body,
                             Notification.NotificationType type, String referenceId,
                             Map<String, String> data) {
        String messageId = fcmService.sendToToken(fcmToken, title, body, data);

        Notification notification = Notification.builder()
                .userId(userId)
                .title(title)
                .body(body)
                .type(type)
                .referenceId(referenceId)
                .fcmMessageId(messageId)
                .build();

        notificationRepository.save(notification);
        log.info("Notification saved: userId={} type={}", userId, type);
    }

    public Page<NotificationResponse> getNotifications(UUID userId, Pageable pageable) {
        return notificationRepository
                .findByUserIdOrderByCreatedAtDesc(userId, pageable)
                .map(this::toResponse);
    }

    public long getUnreadCount(UUID userId) {
        return notificationRepository.countByUserIdAndIsReadFalse(userId);
    }

    @Transactional
    public void markAllRead(UUID userId) {
        int updated = notificationRepository.markAllReadByUserId(userId);
        log.info("Marked {} notifications as read for userId={}", updated, userId);
    }

    @Transactional
    public void markRead(UUID notificationId, UUID userId) {
        notificationRepository.markReadById(notificationId, userId);
    }

    private NotificationResponse toResponse(Notification n) {
        return NotificationResponse.builder()
                .id(n.getId())
                .title(n.getTitle())
                .body(n.getBody())
                .type(n.getType().name())
                .referenceId(n.getReferenceId())
                .isRead(n.getIsRead())
                .createdAt(n.getCreatedAt())
                .build();
    }
}
