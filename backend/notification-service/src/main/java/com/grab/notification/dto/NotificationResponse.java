package com.grab.notification.dto;

import lombok.Builder;
import lombok.Data;

import java.time.LocalDateTime;
import java.util.UUID;

@Data
@Builder
public class NotificationResponse {
    private UUID id;
    private String title;
    private String body;
    private String type;
    private String referenceId;
    private Boolean isRead;
    private LocalDateTime createdAt;
}
