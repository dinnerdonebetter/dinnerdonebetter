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
                
                SecureField("Password", text: $password)
                    .textFieldStyle(.roundedBorder)
                    .disabled(isLoading)
                
                if !errorMessage.isEmpty {
                    Text(errorMessage)
                        .font(.caption)
                        .foregroundColor(.red)
                        .multilineTextAlignment(.center)
                }
                
                Button(action: { 
                    // Cancel any existing login task
                    loginTask?.cancel()
                    // Create new task and store reference
                    loginTask = Task { await handleLogin() }
                }) {
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
                .disabled(isLoading || username.isEmpty || password.isEmpty)
                .padding(.top, 8)
            }
            .padding(.horizontal, 32)
            
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
        
        let result = await authManager.login(username: username, password: password)
        
        // Check for cancellation before updating UI
        guard !Task.isCancelled else { return }
        
        await MainActor.run {
            isLoading = false
            
            if result.success {
                // Login successful, state change will trigger view update
                username = ""
                password = ""
            } else {
                errorMessage = result.error ?? "Unknown error occurred"
            }
        }
    }
}

#Preview {
    LoginView()
        .environment(AuthenticationManager())
}
