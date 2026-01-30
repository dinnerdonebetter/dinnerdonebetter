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

        if requiresTOTP {
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
            || (requiresTOTP && totpCode.isEmpty)
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
        totpCode = ""
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
