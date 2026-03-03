//
//  iosApp.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import SwiftUI

@main
struct IOSApp: App {
  @UIApplicationDelegateAdaptor(AppDelegate.self) private var appDelegate
  @State private var eventReporterService = EventReporterService()
  @State private var authManager = AuthenticationManager()
  @State private var deepLinkHandler = DeepLinkHandler()

  var body: some Scene {
    WindowGroup {
      ContentView()
        .environment(eventReporterService)
        .environment(authManager)
        .environment(deepLinkHandler)
        .onAppear {
          DeviceTokenRegistrationService.shared.configure(authManager: authManager)
        }
        .onOpenURL { url in
          print("Received Universal Link: \(url)")
          deepLinkHandler.handleURL(url)
        }
    }
  }
}
