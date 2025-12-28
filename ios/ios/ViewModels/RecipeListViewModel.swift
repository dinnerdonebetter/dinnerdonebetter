//
//  RecipeListViewModel.swift
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
class RecipeListViewModel {
  var recipes: [Mealplanning_Recipe] = []
  var isLoading = false
  var errorMessage: String?

  private let authManager: AuthenticationManager

  init(authManager: AuthenticationManager) {
    self.authManager = authManager
  }

  func loadRecipes() async {
    isLoading = true
    errorMessage = nil

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "RecipeListViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }

      // Get OAuth2 token (will refresh if needed)
      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "RecipeListViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      // Create request - empty filter to get all recipes
      let request = Mealplanning_GetRecipesRequest()

      // Execute request
      let response = try await clientManager.client.mealPlanning.getRecipes(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      self.recipes = response.results
    } catch {
      errorMessage = "Failed to load recipes: \(error.localizedDescription)"
      print("❌ Error loading recipes: \(error)")
    }

    isLoading = false
  }
}
