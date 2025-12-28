//
//  HomeViewModel.swift
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
class HomeViewModel {
  // Type alias for convenience
  typealias Timestamp = SwiftProtobuf.Google_Protobuf_Timestamp

  // Static utility function to convert protobuf timestamp to Date
  // Can be used from both HomeViewModel and HomeView
  static func timestampToDate(_ timestamp: SwiftProtobuf.Google_Protobuf_Timestamp) -> Date {
    // Google_Protobuf_Timestamp.seconds is Int64, nanos is Int32
    // TimeInterval is a typealias for Double
    let seconds = TimeInterval(timestamp.seconds)
    let nanos = TimeInterval(timestamp.nanos) / 1_000_000_000.0
    return Date(timeIntervalSince1970: seconds + nanos)
  }
  // Data
  var allMealPlans: [Mealplanning_MealPlan] = []
  var userTasks: [Mealplanning_MealPlanTask] = []
  var groceryLists: [String: [Mealplanning_MealPlanGroceryListItem]] = [:]  // Keyed by meal plan ID

  // Loading states
  var isLoading = false
  var errorMessage: String?

  // Computed properties
  var pendingVoteMealPlans: [Mealplanning_MealPlan] {
    let now = Date()
    return allMealPlans.filter { mealPlan in
      // Meal plan is pending if voting deadline hasn't passed and status indicates voting
      guard mealPlan.status == .awaitingVotes else {
        return false
      }

      let deadline = timestampToDate(mealPlan.votingDeadline)
      return deadline > now
    }
  }

  var upcomingMealPlans: [Mealplanning_MealPlan] {
    let now = Date()
    let fourWeeksFromNow = Calendar.current.date(byAdding: .day, value: 28, to: now) ?? now

    return allMealPlans.filter { mealPlan in
      // Finalized meal plans with events in the next 2 weeks
      guard mealPlan.status == .finalized else {
        return false
      }

      // Check if any event is in the next 2 weeks
      return mealPlan.events.contains { event in
        let eventStart = timestampToDate(event.startsAt)
        return eventStart >= now && eventStart <= fourWeeksFromNow
      }
    }
  }

  var activeGroceryLists: [(mealPlanID: String, items: [Mealplanning_MealPlanGroceryListItem])] {
    return groceryLists.compactMap { mealPlanID, items in
      // Only show grocery lists for finalized meal plans with unacquired items
      let mealPlan = allMealPlans.first { $0.id == mealPlanID }
      guard let mealPlan = mealPlan,
        mealPlan.status == .finalized,
        mealPlan.groceryListInitialized
      else {
        return nil
      }

      // Filter to show only items that aren't fully acquired
      let activeItems = items.filter { item in
        item.status != .acquired
      }

      return activeItems.isEmpty ? nil : (mealPlanID, activeItems)
    }
  }

  private let authManager: AuthenticationManager

  init(authManager: AuthenticationManager) {
    self.authManager = authManager
  }

  func loadData() async {
    isLoading = true
    errorMessage = nil

    do {
      // First fetch meal plans
      let mealPlans = try await fetchMealPlans()
      self.allMealPlans = mealPlans
      
      // Log pending vote meal plans count
      let pendingCount = pendingVoteMealPlans.count
      print("📊 HomeViewModel: Loaded \(pendingCount) meal plan(s) with pending votes")

      // Then fetch tasks (which needs meal plans to be loaded)
      let tasks = try await fetchUserTasks()
      self.userTasks = tasks

      // Fetch grocery lists for finalized meal plans
      await fetchGroceryLists(
        for: mealPlans.filter { $0.status == .finalized && $0.groceryListInitialized }
      )
    } catch {
      errorMessage = "Failed to load data: \(error.localizedDescription)"
      print("❌ Error loading home data: \(error)")
    }

    isLoading = false
  }

  private func fetchMealPlans() async throws -> [Mealplanning_MealPlan] {
    guard let clientManager = try? authManager.getClientManager() else {
      throw NSError(
        domain: "HomeViewModel", code: 1,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
    }

    // Get OAuth2 token (will refresh if needed)
    guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
      throw NSError(
        domain: "HomeViewModel", code: 2,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
    }

    let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

    let response = try await clientManager.client.mealPlanning.getMealPlansForAccount(
      Mealplanning_GetMealPlansForAccountRequest(),
      metadata: metadata,
      options: clientManager.defaultCallOptions
    )

    return response.results
  }

  private func fetchUserTasks() async throws -> [Mealplanning_MealPlanTask] {
    guard let clientManager = try? authManager.getClientManager() else {
      throw NSError(
        domain: "HomeViewModel", code: 1,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
    }

    // Get OAuth2 token (will refresh if needed)
    guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
      throw NSError(
        domain: "HomeViewModel", code: 2,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
    }

    let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

    // We need to get tasks for all meal plans and filter by user
    // For now, let's get tasks from all finalized meal plans
    let finalizedMealPlans = allMealPlans.filter {
      $0.status == .finalized && $0.tasksCreated
    }

    var allTasks: [Mealplanning_MealPlanTask] = []

    for mealPlan in finalizedMealPlans {
      var request = Mealplanning_GetMealPlanTasksRequest()
      request.mealPlanID = mealPlan.id

      do {
        let response = try await clientManager.client.mealPlanning.getMealPlanTasks(
          request,
          metadata: metadata,
          options: clientManager.defaultCallOptions
        )

        // Filter tasks assigned to current user
        let userTasks = response.results.filter { task in
          task.assignedToUser == authManager.userID
        }

        allTasks.append(contentsOf: userTasks)
      } catch {
        print("⚠️ Failed to fetch tasks for meal plan \(mealPlan.id): \(error)")
        // Continue with other meal plans
      }
    }

    return allTasks
  }

  private func fetchGroceryLists(for mealPlans: [Mealplanning_MealPlan]) async {
    guard let clientManager = try? authManager.getClientManager() else {
      return
    }

    // Get OAuth2 token (will refresh if needed)
    guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
      print("⚠️ Failed to get OAuth2 access token for grocery lists")
      return
    }

    let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

    for mealPlan in mealPlans {
      var request = Mealplanning_GetMealPlanGroceryListItemsForMealPlanRequest()
      request.mealPlanID = mealPlan.id

      do {
        let response = try await clientManager.client.mealPlanning
          .getMealPlanGroceryListItemsForMealPlan(
            request,
            metadata: metadata,
            options: clientManager.defaultCallOptions
          )

        await MainActor.run {
          groceryLists[mealPlan.id] = response.results
        }
      } catch {
        print("⚠️ Failed to fetch grocery list for meal plan \(mealPlan.id): \(error)")
      }
    }
  }

  // Instance method that delegates to static method for consistency
  func timestampToDate(_ timestamp: Timestamp) -> Date {
    return Self.timestampToDate(timestamp)
  }

  // Check if user has voted on a meal plan
  func hasUserVoted(on mealPlan: Mealplanning_MealPlan) -> Bool {
    for event in mealPlan.events {
      for option in event.options
      where option.votes.contains(where: { $0.byUser == authManager.userID && !$0.abstain }) {
        return true
      }
    }
    return false
  }

  // Format time until deadline
  func timeUntilDeadline(_ deadline: SwiftProtobuf.Google_Protobuf_Timestamp) -> String {
    let deadlineDate = timestampToDate(deadline)
    let now = Date()

    if deadlineDate <= now {
      return "Deadline passed"
    }

    let components = Calendar.current.dateComponents(
      [.day, .hour, .minute], from: now, to: deadlineDate)

    if let days = components.day, days > 0 {
      return "\(days) day\(days == 1 ? "" : "s") left"
    } else if let hours = components.hour, hours > 0 {
      return "\(hours) hour\(hours == 1 ? "" : "s") left"
    } else if let minutes = components.minute, minutes > 0 {
      return "\(minutes) minute\(minutes == 1 ? "" : "s") left"
    } else {
      return "Less than a minute left"
    }
  }
}
