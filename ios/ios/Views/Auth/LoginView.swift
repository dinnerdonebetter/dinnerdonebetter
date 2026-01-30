//
//  LoginView.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import SwiftUI

struct LoginView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Binding var showRegister: Bool

  @State private var username: String = "admin_user"
  @State private var password: String = "admin_pass"
  @State private var totpCode: String = ""
  @State private var requiresTOTP: Bool = false
  @State private var errorMessage: String = ""
  @State private var isLoading: Bool = false
  @State private var loginTask: Task<Void, Never>?

  // Temporary dev feature: always show TOTP and auto-generate
  @State private var alwaysShowTOTP: Bool = true
  private let totpSecret: String =
    "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
  @State private var totpUpdateTask: Task<Void, Never>?

  var body: some View {
    VStack(spacing: DSTheme.Spacing.xl) {
      Spacer()

      // App Title
      VStack(spacing: DSTheme.Spacing.sm) {
        Text("Dinner Done Better")
          .font(DSTheme.Typography.largeTitle)
          .foregroundColor(DSTheme.Colors.textPrimary)

        Text("Sign in to continue")
          .font(DSTheme.Typography.body)
          .foregroundColor(DSTheme.Colors.textSecondary)
      }

      Spacer()

      // Login Form
      VStack(spacing: DSTheme.Spacing.lg) {
        DSTextField(
          "Username",
          text: $username,
          isDisabled: isLoading
        )
        .accessibilityIdentifier("usernameTextField")

        DSTextField(
          "Password",
          text: $password,
          type: .password,
          isDisabled: isLoading
        )
        .accessibilityIdentifier("passwordTextField")

        if requiresTOTP || alwaysShowTOTP {
          DSTextField(
            "2FA Code",
            text: $totpCode,
            type: .number,
            isDisabled: isLoading
          )
          .accessibilityIdentifier("totpTextField")
        }

        if !errorMessage.isEmpty {
          Text(errorMessage)
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.error)
            .multilineTextAlignment(.center)
            .accessibilityIdentifier("errorMessage")
        }

        DSButton(
          isLoading ? "Signing In..." : "Sign In",
          fullWidth: true,
          isLoading: isLoading,
          isDisabled: username.isEmpty || password.isEmpty
            || ((requiresTOTP || alwaysShowTOTP) && totpCode.isEmpty)
        ) {
          loginTask?.cancel()
          loginTask = Task { await handleLogin() }
        }
        .accessibilityIdentifier("signInButton")
        .padding(.top, DSTheme.Spacing.sm)

        // Navigation to register
        Button {
          showRegister = true
        } label: {
          Text("Don't have an account? Sign up")
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.primary)
        }
        .padding(.top, DSTheme.Spacing.sm)
      }
      .padding(.horizontal, DSTheme.Spacing.xxl)
      .animation(DSTheme.Animation.normal, value: requiresTOTP)

      Spacer()
      Spacer()
    }
    .dsScreenPadding()
    .onAppear {
      startTOTPTimer()
    }
    .onDisappear {
      stopTOTPTimer()
    }
  }

  // MARK: - TOTP Generation

  private func updateTOTPCode() {
    guard !totpSecret.isEmpty else {
      totpCode = ""
      return
    }

    if let code = TOTPGenerator.generate(secret: totpSecret) {
      totpCode = code
    } else {
      totpCode = ""
    }
  }

  private func startTOTPTimer() {
    if alwaysShowTOTP {
      updateTOTPCode()
    }

    totpUpdateTask?.cancel()

    guard alwaysShowTOTP else { return }

    totpUpdateTask = Task {
      while !Task.isCancelled {
        try? await Task.sleep(nanoseconds: 1_000_000_000)
        guard !Task.isCancelled else { break }
        await MainActor.run {
          updateTOTPCode()
        }
      }
    }
  }

  private func stopTOTPTimer() {
    totpUpdateTask?.cancel()
    totpUpdateTask = nil
  }

  private func handleLogin() async {
    guard !Task.isCancelled else { return }

    await MainActor.run {
      errorMessage = ""
      isLoading = true
    }

    let totpToken = totpCode.isEmpty ? nil : totpCode
    let result = await authManager.login(
      username: username, password: password, totpToken: totpToken)

    guard !Task.isCancelled else { return }

    await MainActor.run {
      isLoading = false

      if result.success {
        username = ""
        password = ""
        if !alwaysShowTOTP {
          totpCode = ""
        }
        requiresTOTP = false
      } else {
        if result.requiresTOTP {
          requiresTOTP = true
          errorMessage = result.error ?? "Please enter your 2FA code."
        } else {
          errorMessage = result.error ?? "Unknown error occurred"
        }
      }
    }
  }
}

#Preview {
  LoginView(showRegister: .constant(true))
    .environment(AuthenticationManager())
}
