package com.grab.auth.service;

import com.grab.auth.dto.LoginRequest;
import com.grab.auth.dto.RegisterRequest;
import com.grab.auth.model.User;
import com.grab.auth.repository.UserRepository;
import com.grab.auth.security.JwtUtil;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.data.redis.core.ValueOperations;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.crypto.password.PasswordEncoder;

import java.util.Optional;
import java.util.UUID;

import static org.assertj.core.api.Assertions.*;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class AuthServiceTest {

    @Mock
    private UserRepository userRepository;

    @Mock
    private PasswordEncoder passwordEncoder;

    @Mock
    private JwtUtil jwtUtil;

    @Mock
    private StringRedisTemplate redisTemplate;

    @Mock
    private ValueOperations<String, String> valueOps;

    @InjectMocks
    private AuthService authService;

    private User mockUser;

    @BeforeEach
    void setUp() {
        mockUser = User.builder()
                .id(UUID.randomUUID())
                .name("Test User")
                .phone("0812345678")
                .email("test@example.com")
                .passwordHash("$2a$10$hashedPassword")
                .role(User.Role.CUSTOMER)
                .isActive(true)
                .build();

        when(redisTemplate.opsForValue()).thenReturn(valueOps);
        when(jwtUtil.getRefreshExpirationMs()).thenReturn(604800000L);
        when(jwtUtil.getJwtExpirationMs()).thenReturn(900000L);
        when(jwtUtil.generateAccessToken(any(), any())).thenReturn("mock-access-token");
        when(jwtUtil.generateRefreshToken(any())).thenReturn("mock-refresh-token");
    }

    @Test
    @DisplayName("Register new user successfully")
    void register_success() {
        RegisterRequest request = new RegisterRequest();
        request.setName("Test User");
        request.setPhone("0812345678");
        request.setPassword("password123");
        request.setRole(User.Role.CUSTOMER);

        when(userRepository.existsByPhone("0812345678")).thenReturn(false);
        when(passwordEncoder.encode("password123")).thenReturn("$2a$10$hashedPassword");
        when(userRepository.save(any(User.class))).thenReturn(mockUser);
        doNothing().when(valueOps).set(any(), any(), anyLong(), any());

        var response = authService.register(request);

        assertThat(response.getAccessToken()).isEqualTo("mock-access-token");
        assertThat(response.getUser().getPhone()).isEqualTo("0812345678");
        assertThat(response.getUser().getRole()).isEqualTo("CUSTOMER");
        verify(userRepository).save(any(User.class));
    }

    @Test
    @DisplayName("Register fails when phone already exists")
    void register_phoneAlreadyExists() {
        RegisterRequest request = new RegisterRequest();
        request.setPhone("0812345678");
        request.setPassword("password123");

        when(userRepository.existsByPhone("0812345678")).thenReturn(true);

        assertThatThrownBy(() -> authService.register(request))
                .isInstanceOf(IllegalArgumentException.class)
                .hasMessageContaining("Phone number already registered");
    }

    @Test
    @DisplayName("Login successfully with correct credentials")
    void login_success() {
        LoginRequest request = new LoginRequest();
        request.setPhone("0812345678");
        request.setPassword("password123");

        when(userRepository.findByPhone("0812345678")).thenReturn(Optional.of(mockUser));
        when(passwordEncoder.matches("password123", mockUser.getPasswordHash())).thenReturn(true);
        doNothing().when(valueOps).set(any(), any(), anyLong(), any());

        var response = authService.login(request);

        assertThat(response.getAccessToken()).isEqualTo("mock-access-token");
        assertThat(response.getRefreshToken()).isEqualTo("mock-refresh-token");
    }

    @Test
    @DisplayName("Login fails with wrong password")
    void login_wrongPassword() {
        LoginRequest request = new LoginRequest();
        request.setPhone("0812345678");
        request.setPassword("wrongPassword");

        when(userRepository.findByPhone("0812345678")).thenReturn(Optional.of(mockUser));
        when(passwordEncoder.matches("wrongPassword", mockUser.getPasswordHash())).thenReturn(false);

        assertThatThrownBy(() -> authService.login(request))
                .isInstanceOf(BadCredentialsException.class)
                .hasMessageContaining("Invalid phone or password");
    }

    @Test
    @DisplayName("Login fails when user not found")
    void login_userNotFound() {
        LoginRequest request = new LoginRequest();
        request.setPhone("0899999999");
        request.setPassword("password123");

        when(userRepository.findByPhone("0899999999")).thenReturn(Optional.empty());

        assertThatThrownBy(() -> authService.login(request))
                .isInstanceOf(BadCredentialsException.class);
    }

    @Test
    @DisplayName("Login fails when account is inactive")
    void login_inactiveAccount() {
        mockUser.setIsActive(false);
        LoginRequest request = new LoginRequest();
        request.setPhone("0812345678");
        request.setPassword("password123");

        when(userRepository.findByPhone("0812345678")).thenReturn(Optional.of(mockUser));

        assertThatThrownBy(() -> authService.login(request))
                .isInstanceOf(BadCredentialsException.class)
                .hasMessageContaining("Account is disabled");
    }
}
