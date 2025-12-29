//
//  MealDetailViewModel.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import Foundation
import GRPCCore
import SwiftProtobuf
import SwiftUI

@Observable
@MainActor
class MealDetailViewModel {
  var meal: Mealplanning_Meal?
  var isLoading = false
  var errorMessage: String?

  private let mealID: String
  private let authManager: AuthenticationManager

  init(mealID: String, authManager: AuthenticationManager) {
    self.mealID = mealID
    self.authManager = authManager
  }

  func loadMeal() async {
    isLoading = true
    errorMessage = nil

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "MealDetailViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }

      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "MealDetailViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      var request = Mealplanning_GetMealRequest()
      request.mealID = mealID

      let response = try await clientManager.client.mealPlanning.getMeal(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      self.meal = response.result
    } catch {
      errorMessage = "Failed to load meal: \(error.localizedDescription)"
      print("❌ Error loading meal: \(error)")
    }

    isLoading = false
  }
}
