//
//  ContentView.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import SwiftUI

struct ContentView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Environment(DeepLinkHandler.self) private var deepLinkHandler
  @State private var showLogin: Bool = true

  // Invitation data from deep link
  @State private var pendingInvitationID: String = ""
  @State private var pendingInvitationToken: String = ""

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
          ),
          invitationID: pendingInvitationID,
          invitationToken: pendingInvitationToken
        )
      }
    }
    .onChange(of: deepLinkHandler.pendingDestination) { _, newDestination in
      handleDeepLink(newDestination)
    }
    .onAppear {
      // Handle any pending deep link on appear
      handleDeepLink(deepLinkHandler.pendingDestination)
    }
  }

  private func handleDeepLink(_ destination: DeepLinkDestination?) {
    guard let destination = destination else { return }

    switch destination {
    case .acceptInvitation(let invitationID, let token):
      // Store invitation data and navigate to registration
      pendingInvitationID = invitationID
      pendingInvitationToken = token
      showLogin = false  // Show registration view
      deepLinkHandler.clearPendingDestination()

    case .resetPassword(let token):
      // TODO: Navigate to password reset flow
      print("Password reset token: \(token)")
      deepLinkHandler.clearPendingDestination()

    case .verifyEmail(let token):
      // TODO: Handle email verification
      print("Email verification token: \(token)")
      deepLinkHandler.clearPendingDestination()

    case .unknown:
      break
    }
  }
}

#Preview {
  ContentView()
    .environment(AuthenticationManager())
    .environment(DeepLinkHandler())
}
