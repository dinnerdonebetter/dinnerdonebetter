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
  var searchResults: [Mealplanning_Recipe] = []
  var isLoading = false
  var isSearching = false
  var errorMessage: String?
  var searchError: String?

  private let authManager: AuthenticationManager
  private var searchTask: Task<Void, Never>?

  init(authManager: AuthenticationManager) {
    self.authManager = authManager
  }

  var displayedRecipes: [Mealplanning_Recipe] {
    // If we have search results, show those; otherwise show all recipes
    return searchResults.isEmpty ? recipes : searchResults
  }

  var isInSearchMode: Bool {
    return !searchResults.isEmpty
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

      // Create request - use "submitted" status to match webapp behavior
      var request = Mealplanning_GetRecipesRequest()
      request.status = "submitted"

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

  func searchRecipes(query: String) {
    // Cancel any existing search task
    searchTask?.cancel()

    let trimmedQuery = query.trimmingCharacters(in: .whitespacesAndNewlines)

    // If query is empty, clear search results
    if trimmedQuery.isEmpty {
      searchResults = []
      searchError = nil
      isSearching = false
      return
    }

    // Debounce: wait 500ms before executing search
    searchTask = Task {
      try? await Task.sleep(nanoseconds: 500_000_000)  // 500ms

      // Check if task was cancelled
      guard !Task.isCancelled else { return }

      await performSearch(query: trimmedQuery)
    }
  }

  private func performSearch(query: String) async {
    isSearching = true
    searchError = nil

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "RecipeListViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }

      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "RecipeListViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      // Create search request
      var request = Mealplanning_SearchForRecipesRequest()
      request.query = query
      request.useSearchService = APIConfiguration.useSearchService

      // Execute search
      let response = try await clientManager.client.mealPlanning.searchForRecipes(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      searchResults = response.results
    } catch {
      searchError = "Failed to search recipes: \(error.localizedDescription)"
      print("❌ Error searching for recipes: \(error)")
      searchResults = []
    }

    isSearching = false
  }
}
