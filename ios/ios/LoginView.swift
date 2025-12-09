//
//  LoginView.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import SwiftUI

struct LoginView: View {
    @Environment(AuthenticationManager.self) private var authManager
    
    @State private var username: String = ""
    @State private var password: String = ""
    @State private var showError: Bool = false
    
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
                
                SecureField("Password", text: $password)
                    .textFieldStyle(.roundedBorder)
                
                if showError {
                    Text("Please enter both username and password")
                        .font(.caption)
                        .foregroundColor(.red)
                }
                
                Button(action: handleLogin) {
                    Text("Sign In")
                        .fontWeight(.semibold)
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.accentColor)
                        .foregroundColor(.white)
                        .cornerRadius(10)
                }
                .padding(.top, 8)
            }
            .padding(.horizontal, 32)
            
            Spacer()
            Spacer()
        }
        .padding()
    }
    
    private func handleLogin() {
        showError = false
        
        if authManager.login(username: username, password: password) {
            // Login successful, state change will trigger view update
            username = ""
            password = ""
        } else {
            showError = true
        }
    }
}

#Preview {
    LoginView()
        .environment(AuthenticationManager())
}
