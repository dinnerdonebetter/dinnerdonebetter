//
//  ContentView.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import SwiftUI

struct ContentView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @State private var showLogin: Bool = true

  var body: some View {
    Group {
      if authManager.isAuthenticated && !authManager.oauth2AccessToken.isEmpty {
        HomeView()
      } else if showLogin {
        LoginView(
          showRegister: Binding(
            get: { false },
            set: { _ in showLogin = false }
          )
        )
      } else {
        RegisterView(
          showLogin: Binding(
            get: { true },
            set: { _ in showLogin = true }
          )
        )
      }
    }
  }
}

#Preview {
  ContentView()
    .environment(AuthenticationManager())
}
