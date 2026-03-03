//
//  AppDelegate.swift
//  ios
//
//  Handles push notification lifecycle (device token registration).
//

import UIKit
import UserNotifications

@available(macOS 15.0, iOS 18.0, watchOS 11.0, tvOS 18.0, visionOS 2.0, *)
class AppDelegate: NSObject, UIApplicationDelegate {
  func application(
    _ application: UIApplication,
    didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]? = nil
  ) -> Bool {
    requestNotificationPermissionAndRegister()
    return true
  }

  func application(
    _ application: UIApplication,
    didRegisterForRemoteNotificationsWithDeviceToken deviceToken: Data
  ) {
    DeviceTokenRegistrationService.shared.reportDeviceToken(deviceToken)
  }

  func application(
    _ application: UIApplication,
    didFailToRegisterForRemoteNotificationsWithError error: Error
  ) {
    print("⚠️ Failed to register for remote notifications: \(error.localizedDescription)")
  }

  private func requestNotificationPermissionAndRegister() {
    UNUserNotificationCenter.current().requestAuthorization(
      options: [.alert, .badge, .sound]
    ) { granted, error in
      if let error {
        print("⚠️ Notification permission error: \(error.localizedDescription)")
        return
      }
      if granted {
        Task { @MainActor in
          UIApplication.shared.registerForRemoteNotifications()
        }
      }
    }
  }
}
