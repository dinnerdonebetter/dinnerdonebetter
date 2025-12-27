//
//  ContentView.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import SwiftUI

struct ContentView: View {
  @Environment(AuthenticationManager.self) private var authManager

  var body: some View {
    Group {
      if authManager.isAuthenticated && !authManager.oauth2AccessToken.isEmpty {
      HomeView()
    } else {
      LoginView()
      }
    }
  }
}

#Preview {
  ContentView()
    .environment(AuthenticationManager())
}
