package com.grab.tracking.websocket;

import com.grab.tracking.service.LocationService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.socket.WebSocketHandler;
import org.springframework.web.reactive.socket.WebSocketMessage;
import org.springframework.web.reactive.socket.WebSocketSession;
import reactor.core.publisher.Mono;

@Slf4j
@Component
@RequiredArgsConstructor
public class TrackingWebSocketHandler implements WebSocketHandler {

    private final LocationService locationService;

    @Override
    public Mono<Void> handle(WebSocketSession session) {
        String path = session.getHandshakeInfo().getUri().getPath();
        String orderId = extractOrderId(path);

        log.info("WebSocket connected: sessionId={} orderId={}", session.getId(), orderId);

        Mono<Void> locationStream = session.send(
                locationService.subscribeToOrderTracking(orderId)
                        .map(session::textMessage)
                        .doOnError(e -> log.error("Error streaming location for order {}: {}", orderId, e.getMessage()))
        );

        Mono<Void> heartbeat = session.receive()
                .filter(msg -> msg.getType() == WebSocketMessage.Type.PING)
                .flatMap(ping -> session.send(Mono.just(session.pongMessage(db -> db))))
                .then();

        return Mono.zip(locationStream, heartbeat)
                .then()
                .doFinally(sig -> log.info("WebSocket disconnected: sessionId={} orderId={}", session.getId(), orderId));
    }

    private String extractOrderId(String path) {
        String[] parts = path.split("/");
        return parts[parts.length - 1];
    }
}
