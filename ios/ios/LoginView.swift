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
                
                if requiresTOTP {
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
                
                Button(action: { 
                    // Cancel any existing login task
                    loginTask?.cancel()
                    // Create new task and store reference
                    loginTask = Task { await handleLogin() }
                }, label: {
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
                })
                .disabled(isLoading || username.isEmpty || password.isEmpty || (requiresTOTP && totpCode.isEmpty))
                .accessibilityIdentifier("signInButton")
                .padding(.top, 8)
            }
            .padding(.horizontal, 32)
            .animation(.easeInOut(duration: 0.3), value: requiresTOTP)
            
            Spacer()
            Spacer()
        }
        .padding()
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
        let result = await authManager.login(username: username, password: password, totpToken: totpToken)
        
        // Check for cancellation before updating UI
        guard !Task.isCancelled else { return }
        
        await MainActor.run {
            isLoading = false
            
            if result.success {
                // Login successful, state change will trigger view update
                username = ""
                password = ""
                totpCode = ""
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
