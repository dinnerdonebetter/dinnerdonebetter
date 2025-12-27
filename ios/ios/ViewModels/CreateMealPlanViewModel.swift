//
//  CreateMealPlanViewModel.swift
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
class CreateMealPlanViewModel {
  // Search state
  var searchQuery: String = ""
  var searchResults: [Mealplanning_Meal] = []
  var isSearching: Bool = false
  var searchError: String?

  // Selected meals
  var selectedMeals: [Mealplanning_Meal] = []

  // Meal plan creation state
  var mealPlanName: String = ""
  var votingDeadline: Date = Calendar.current.date(byAdding: .day, value: 7, to: Date()) ?? Date()
  var isCreating: Bool = false
  var creationError: String?
  var createdMealPlanID: String?

  private let authManager: AuthenticationManager

  init(authManager: AuthenticationManager) {
    self.authManager = authManager
  }

  // MARK: - Search Functions

  func searchForMeals() async {
    guard !searchQuery.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
      searchResults = []
      return
    }

    isSearching = true
    searchError = nil

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "CreateMealPlanViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }

      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "CreateMealPlanViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      // Create search request
      var request = Mealplanning_SearchForMealsRequest()
      request.query = searchQuery
      request.useSearchService = false // TODO: configure?

      // Execute search
      let response = try await clientManager.client.mealPlanning.searchForMeals(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      searchResults = response.results
    } catch {
      searchError = "Failed to search meals: \(error.localizedDescription)"
      print("❌ Error searching for meals: \(error)")
      searchResults = []
    }

    isSearching = false
  }

  // MARK: - Selection Functions

  func toggleMealSelection(_ meal: Mealplanning_Meal) {
    if let index = selectedMeals.firstIndex(where: { $0.id == meal.id }) {
      selectedMeals.remove(at: index)
    } else {
      selectedMeals.append(meal)
    }
  }

  func isMealSelected(_ meal: Mealplanning_Meal) -> Bool {
    selectedMeals.contains(where: { $0.id == meal.id })
  }

  func removeSelectedMeal(_ meal: Mealplanning_Meal) {
    selectedMeals.removeAll(where: { $0.id == meal.id })
  }

  // MARK: - Meal Plan Creation

  func createMealPlan() async -> Bool {
    guard !mealPlanName.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
      creationError = "Meal plan name is required"
      return false
    }

    guard !selectedMeals.isEmpty else {
      creationError = "Please select at least one meal"
      return false
    }

    isCreating = true
    creationError = nil

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "CreateMealPlanViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }

      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "CreateMealPlanViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      // Create meal plan request
      var input = Mealplanning_MealPlanCreationRequestInput()
      input.notes = mealPlanName
      input.votingDeadline = dateToTimestamp(votingDeadline)
      input.electionMethod = .schulze  // Default to Schulze method

      // Create a single event for now (can be expanded later)
      // Default to dinner tomorrow
      let tomorrow = Calendar.current.date(byAdding: .day, value: 1, to: Date()) ?? Date()
      let eventStart = Calendar.current.date(bySettingHour: 18, minute: 0, second: 0, of: tomorrow) ?? tomorrow
      let eventEnd = Calendar.current.date(byAdding: .hour, value: 2, to: eventStart) ?? eventStart

      var eventInput = Mealplanning_MealPlanEventCreationRequestInput()
      eventInput.startsAt = dateToTimestamp(eventStart)
      eventInput.endsAt = dateToTimestamp(eventEnd)
      eventInput.mealName = "dinner"
      eventInput.notes = ""

      // Create options for each selected meal
      for meal in selectedMeals {
        var optionInput = Mealplanning_MealPlanOptionCreationRequestInput()
        optionInput.mealID = meal.id
        optionInput.mealScale = 1.0
        optionInput.notes = ""
        eventInput.options.append(optionInput)
      }

      input.events.append(eventInput)

      var request = Mealplanning_CreateMealPlanRequest()
      request.input = input

      // Create the meal plan
      let response = try await clientManager.client.mealPlanning.createMealPlan(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      createdMealPlanID = response.created.id
      isCreating = false
      return true
    } catch {
      creationError = "Failed to create meal plan: \(error.localizedDescription)"
      print("❌ Error creating meal plan: \(error)")
      isCreating = false
      return false
    }
  }

  // MARK: - Helper Functions

  private func dateToTimestamp(_ date: Date) -> SwiftProtobuf.Google_Protobuf_Timestamp {
    var timestamp = SwiftProtobuf.Google_Protobuf_Timestamp()
    timestamp.seconds = Int64(date.timeIntervalSince1970)
    timestamp.nanos = Int32((date.timeIntervalSince1970 - Double(timestamp.seconds)) * 1_000_000_000)
    return timestamp
  }
}

