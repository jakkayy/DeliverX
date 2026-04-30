package com.grab.user.service;

import com.grab.user.dto.UpdateProfileRequest;
import com.grab.user.exception.ResourceNotFoundException;
import com.grab.user.model.User;
import com.grab.user.repository.UserRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Optional;
import java.util.UUID;

import static org.assertj.core.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UserServiceTest {

    @Mock
    private UserRepository userRepository;

    @InjectMocks
    private UserService userService;

    private User mockUser;
    private UUID userId;

    @BeforeEach
    void setUp() {
        userId = UUID.randomUUID();
        mockUser = User.builder()
                .id(userId)
                .name("Somchai")
                .phone("0812345678")
                .email("somchai@example.com")
                .role(User.Role.CUSTOMER)
                .isActive(true)
                .build();
    }

    @Test
    @DisplayName("Get profile successfully")
    void getProfile_success() {
        when(userRepository.findById(userId)).thenReturn(Optional.of(mockUser));

        var response = userService.getProfile(userId);

        assertThat(response.getId()).isEqualTo(userId);
        assertThat(response.getName()).isEqualTo("Somchai");
        assertThat(response.getPhone()).isEqualTo("0812345678");
    }

    @Test
    @DisplayName("Get profile throws when user not found")
    void getProfile_notFound() {
        when(userRepository.findById(userId)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> userService.getProfile(userId))
                .isInstanceOf(ResourceNotFoundException.class);
    }

    @Test
    @DisplayName("Update profile name successfully")
    void updateProfile_success() {
        UpdateProfileRequest request = new UpdateProfileRequest();
        request.setName("New Name");

        when(userRepository.findById(userId)).thenReturn(Optional.of(mockUser));
        when(userRepository.save(mockUser)).thenReturn(mockUser);

        var response = userService.updateProfile(userId, request);

        assertThat(response.getName()).isEqualTo("New Name");
        verify(userRepository).save(mockUser);
    }

    @Test
    @DisplayName("Deactivate account successfully")
    void deactivateAccount_success() {
        when(userRepository.findById(userId)).thenReturn(Optional.of(mockUser));
        when(userRepository.save(mockUser)).thenReturn(mockUser);

        userService.deactivateAccount(userId);

        assertThat(mockUser.getIsActive()).isFalse();
        verify(userRepository).save(mockUser);
    }
}
