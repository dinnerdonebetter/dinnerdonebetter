//
//  iosApp.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import SwiftUI

@main
struct IOSApp: App {
  @State private var authManager = AuthenticationManager()

  var body: some Scene {
    WindowGroup {
      ContentView()
        .environment(authManager)
    }
  }
}
