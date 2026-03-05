//
//  EntitlementService.swift
//  ios
//
//  Checks RevenueCat entitlements (e.g. "Dinner Done Better Pro").
//

import Foundation
import RevenueCat

enum EntitlementService {
  /// Entitlement identifier as configured in the RevenueCat dashboard.
  static let proEntitlementID = "Dinner Done Better Pro"

  /// Returns whether the user has an active "Dinner Done Better Pro" entitlement.
  /// Returns `false` when RevenueCat is not configured or on error.
  static func isProActive() async -> Bool {
    guard RevenueCatConfiguration.isConfigured else {
      return false
    }

    do {
      let customerInfo = try await Purchases.shared.customerInfo()
      return customerInfo.entitlements.all[Self.proEntitlementID]?.isActive == true
    } catch {
      print("⚠️ EntitlementService: Error checking entitlement: \(error)")
      return false
    }
  }
}
