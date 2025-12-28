//
//  LoginView.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import SwiftUI

struct LoginView: View {
  @Environment(AuthenticationManager.self) private var authManager

  @State private var username: String = "admin_user"
  @State private var password: String = "admin_pass"
  @State private var totpCode: String = ""
  @State private var requiresTOTP: Bool = false
  @State private var errorMessage: String = ""
  @State private var isLoading: Bool = false
  @State private var loginTask: Task<Void, Never>?

  // Temporary dev feature: always show TOTP and auto-generate
  // Set alwaysShowTOTP to false to disable this feature
  // The TOTP secret is hardcoded for development
  @State private var alwaysShowTOTP: Bool = true  // Set to false to disable
  // Hardcoded TOTP secret (base32) for dev
  private let totpSecret: String =
    "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
  @State private var totpUpdateTask: Task<Void, Never>?

  var body: some View {
    VStack(spacing: 20) {
      Spacer()

      // App Title
      Text("Dinner Done Better")
        .font(.largeTitle)
        .fontWeight(.bold)

      Text("Sign in to continue")
        .font(.subheadline)
        .foregroundColor(.secondary)

      Spacer()

      // Login Form
      VStack(spacing: 16) {
        TextField("Username", text: $username)
          .textFieldStyle(.roundedBorder)
          .textInputAutocapitalization(.never)
          .autocorrectionDisabled()
          .disabled(isLoading)
          .accessibilityIdentifier("usernameTextField")

        SecureField("Password", text: $password)
          .textFieldStyle(.roundedBorder)
          .disabled(isLoading)
          .accessibilityIdentifier("passwordTextField")

        if requiresTOTP || alwaysShowTOTP {
          TextField("2FA Code", text: $totpCode)
            .textFieldStyle(.roundedBorder)
            .keyboardType(.numberPad)
            .disabled(isLoading)
            .accessibilityIdentifier("totpTextField")
        }

        if !errorMessage.isEmpty {
          Text(errorMessage)
            .font(.caption)
            .foregroundColor(.red)
            .multilineTextAlignment(.center)
            .accessibilityIdentifier("errorMessage")
        }

        Button(
          action: {
            // Cancel any existing login task
            loginTask?.cancel()
            // Create new task and store reference
            loginTask = Task { await handleLogin() }
          },
          label: {
            HStack {
              if isLoading {
                ProgressView()
                  .progressViewStyle(CircularProgressViewStyle(tint: .white))
                  .scaleEffect(0.8)
              }
              Text(isLoading ? "Signing In..." : "Sign In")
                .fontWeight(.semibold)
            }
            .frame(maxWidth: .infinity)
            .padding()
            .background(isLoading ? Color.gray : Color.accentColor)
            .foregroundColor(.white)
            .cornerRadius(10)
          }
        )
        .disabled(
          isLoading || username.isEmpty || password.isEmpty
            || ((requiresTOTP || alwaysShowTOTP) && totpCode.isEmpty)
        )
        .accessibilityIdentifier("signInButton")
        .padding(.top, 8)
      }
      .padding(.horizontal, 32)
      .animation(.easeInOut(duration: 0.3), value: requiresTOTP)

      Spacer()
      Spacer()
    }
    .padding()
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
    // Update immediately
    if alwaysShowTOTP {
      updateTOTPCode()
    }

    // Cancel any existing task
    totpUpdateTask?.cancel()

    // Update every second to refresh the code when it changes
    guard alwaysShowTOTP else { return }

    totpUpdateTask = Task {
      while !Task.isCancelled {
        try? await Task.sleep(nanoseconds: 1_000_000_000)  // 1 second
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
    // Check for cancellation at the start
    guard !Task.isCancelled else { return }

    await MainActor.run {
      errorMessage = ""
      isLoading = true
    }

    // Use TOTP code if it's been entered, otherwise pass nil
    let totpToken = totpCode.isEmpty ? nil : totpCode
    let result = await authManager.login(
      username: username, password: password, totpToken: totpToken)

    // Check for cancellation before updating UI
    guard !Task.isCancelled else { return }

    await MainActor.run {
      isLoading = false

      if result.success {
        // Login successful, state change will trigger view update
        username = ""
        password = ""
        // Only clear TOTP code if not using always-show feature
        if !alwaysShowTOTP {
          totpCode = ""
        }
        requiresTOTP = false
      } else {
        // Check if TOTP is required
        if result.requiresTOTP {
          requiresTOTP = true
          errorMessage = result.error ?? "Please enter your 2FA code."
        } else {
          // If TOTP was required but login failed, keep the TOTP field visible
          // but show the error message
          errorMessage = result.error ?? "Unknown error occurred"
        }
      }
    }
  }
}

#Preview {
  LoginView()
    .environment(AuthenticationManager())
}
