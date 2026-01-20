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

// MARK: - Meal Plan Event

struct MealPlanEvent: Identifiable {
  let id = UUID()
  var mealType: Mealplanning_MealPlanEventName = .dinner
  var startDate: Date
  var endDate: Date
  var notes: String = ""
  var searchQuery: String = ""
  var searchResults: [Mealplanning_Meal] = []
  var isSearching: Bool = false
  var searchError: String?
  var selectedMeals: [Mealplanning_Meal] = []
  // Meal scales: maps meal ID to scale value (default 1.0)
  var mealScales: [String: Float] = [:]
}

@Observable
@MainActor
// swiftlint:disable:next type_body_length
class CreateMealPlanViewModel {
  // Events - always at least one
  var events: [MealPlanEvent] = []

  // Meal plan creation state
  var mealPlanName: String = ""
  var votingDeadline = Date()
  var isCreating: Bool = false
  var creationError: String?
  var createdMealPlanID: String?

  // Option selections: recipeID -> (stepID -> (ingredientIndex -> selectedOptionIndex))
  // Note: Only ingredients have selectable options; instruments and vessels are concrete
  var recipeOptionSelections: [String: [String: [UInt32: UInt32]]] = [:]

  private let authManager: AuthenticationManager

  init(authManager: AuthenticationManager) {
    self.authManager = authManager
    // Initialize with one blank event - next Monday at 7PM
    let eventStart = nextMondayAt7PM()
    let eventEnd = Calendar.current.date(byAdding: .hour, value: 2, to: eventStart) ?? eventStart
    // Set voting deadline to the preceding Friday at midnight
    self.votingDeadline = precedingFridayAtMidnight(for: eventStart)
    self.events = [
      MealPlanEvent(mealType: .dinner, startDate: eventStart, endDate: eventEnd)
    ]
    // Set default meal plan name
    updateDefaultMealPlanName()
  }

  // MARK: - Helper Functions for Date Calculation

  private func nextMondayAt7PM() -> Date {
    let calendar = Calendar.current
    let now = Date()

    // Get the weekday component (1 = Sunday, 2 = Monday, ..., 7 = Saturday)
    let weekday = calendar.component(.weekday, from: now)

    // Calculate days until next Monday
    // If today is Monday (2), we want next Monday (7 days)
    // If today is Sunday (1), we want the Monday after tomorrow (8 days)
    // If today is Tuesday-Saturday, we want the Monday after that
    let daysUntilMonday: Int
    if weekday == 2 {  // Today is Monday
      daysUntilMonday = 7
    } else if weekday == 1 {  // Today is Sunday
      daysUntilMonday = 8  // Monday after tomorrow
    } else {  // Tuesday through Saturday
      daysUntilMonday = 9 - weekday  // e.g., Tuesday (3) -> 6 days, Wednesday (4) -> 5 days
    }

    // Get next Monday
    guard let nextMonday = calendar.date(byAdding: .day, value: daysUntilMonday, to: now) else {
      // Fallback to tomorrow at 7PM if calculation fails
      return calendar.date(
        bySettingHour: 19, minute: 0, second: 0,
        of: calendar.date(byAdding: .day, value: 1, to: now) ?? now) ?? now
    }

    // Set time to 7PM (19:00)
    return calendar.date(bySettingHour: 19, minute: 0, second: 0, of: nextMonday) ?? nextMonday
  }

  /// Get the preceding Friday at midnight for a given date
  /// - Parameter date: The date to find the preceding Friday for
  /// - Returns: The Friday before the given date at midnight (00:00)
  private func precedingFridayAtMidnight(for date: Date) -> Date {
    let calendar = Calendar.current
    let weekday = calendar.component(.weekday, from: date)

    // Calculate days to subtract to get to the preceding Friday
    // Friday = 6, Saturday = 7, Sunday = 1, Monday = 2, Tuesday = 3, Wednesday = 4, Thursday = 5
    let daysToSubtract: Int
    switch weekday {
    case 1:  // Sunday - Friday is 2 days ago
      daysToSubtract = 2
    case 2:  // Monday - Friday is 3 days ago
      daysToSubtract = 3
    case 3:  // Tuesday - Friday is 4 days ago
      daysToSubtract = 4
    case 4:  // Wednesday - Friday is 5 days ago
      daysToSubtract = 5
    case 5:  // Thursday - Friday is 6 days ago
      daysToSubtract = 6
    case 6:  // Friday - use this Friday (0 days)
      daysToSubtract = 0
    case 7:  // Saturday - Friday is 1 day ago
      daysToSubtract = 1
    default:
      daysToSubtract = 0
    }

    // Get the Friday
    guard let friday = calendar.date(byAdding: .day, value: -daysToSubtract, to: date) else {
      // Fallback to 3 days before the date at midnight
      return calendar.date(
        bySettingHour: 0, minute: 0, second: 0,
        of: calendar.date(byAdding: .day, value: -3, to: date) ?? date) ?? date
    }

    // Set time to midnight (00:00)
    return calendar.date(bySettingHour: 0, minute: 0, second: 0, of: friday) ?? friday
  }

  // MARK: - Event Management

  func addEvent() {
    // Find the latest event date
    let latestDate = events.map { $0.startDate }.max() ?? Date()
    // Add new event 1 day after the latest, or next Monday if no events
    let newEventDate =
      Calendar.current.date(byAdding: .day, value: 1, to: latestDate) ?? nextMondayAt7PM()
    let eventStart =
      Calendar.current.date(bySettingHour: 19, minute: 0, second: 0, of: newEventDate)
      ?? newEventDate
    let eventEnd = Calendar.current.date(byAdding: .hour, value: 2, to: eventStart) ?? eventStart

    let newEvent = MealPlanEvent(mealType: .dinner, startDate: eventStart, endDate: eventEnd)
    events.append(newEvent)

    // Update voting deadline to be the earliest event's preceding Friday
    updateVotingDeadline()
    // updateDefaultMealPlanName is called by updateVotingDeadline
  }

  private func updateVotingDeadline() {
    // Find the earliest event date
    guard let earliestEvent = events.min(by: { $0.startDate < $1.startDate }) else { return }
    var deadline = precedingFridayAtMidnight(for: earliestEvent.startDate)

    // If the deadline is sooner than 24 hours from now, move it forward by a week
    let now = Date()
    let hoursUntilDeadline = deadline.timeIntervalSince(now) / 3600.0
    if hoursUntilDeadline < 24.0 {
      if let nextWeekDeadline = Calendar.current.date(byAdding: .day, value: 7, to: deadline) {
        deadline = nextWeekDeadline
      }
    }

    votingDeadline = deadline
    // Update default meal plan name when events change
    updateDefaultMealPlanName()
  }

  private func updateDefaultMealPlanName() {
    guard !events.isEmpty else {
      mealPlanName = "Meal Plan"
      return
    }

    // Find earliest start and latest end dates
    let earliestStart = events.map { $0.startDate }.min() ?? Date()
    let latestEnd = events.map { $0.endDate }.max() ?? Date()

    // Format dates
    let dateFormatter = DateFormatter()
    dateFormatter.dateStyle = .medium
    dateFormatter.timeStyle = .short

    let startsAt = dateFormatter.string(from: earliestStart)
    let endsAt = dateFormatter.string(from: latestEnd)

    mealPlanName = "Meal Plan for \(startsAt) to \(endsAt)"
  }

  func removeEvent(_ event: MealPlanEvent) {
    events.removeAll(where: { $0.id == event.id })
    // Ensure at least one event remains
    if events.isEmpty {
      let eventStart = nextMondayAt7PM()
      let eventEnd = Calendar.current.date(byAdding: .hour, value: 2, to: eventStart) ?? eventStart
      events.append(MealPlanEvent(mealType: .dinner, startDate: eventStart, endDate: eventEnd))
    }
    // Update meal plan name when events change
    updateDefaultMealPlanName()
  }

  func updateEvent(_ event: MealPlanEvent) {
    guard let index = events.firstIndex(where: { $0.id == event.id }) else { return }
    events[index] = event
  }

  func updateEventSearchQuery(_ eventID: UUID, query: String) {
    guard let index = events.firstIndex(where: { $0.id == eventID }) else { return }
    events[index].searchQuery = query
  }

  func updateEventMealType(_ eventID: UUID, mealType: Mealplanning_MealPlanEventName) {
    guard let index = events.firstIndex(where: { $0.id == eventID }) else { return }
    events[index].mealType = mealType
  }

  func updateEventStartDate(_ eventID: UUID, date: Date) {
    guard let index = events.firstIndex(where: { $0.id == eventID }) else { return }
    events[index].startDate = date
    // Auto-update end time to be 2 hours after start
    events[index].endDate = Calendar.current.date(byAdding: .hour, value: 2, to: date) ?? date
    // Update voting deadline based on earliest event
    updateVotingDeadline()
  }

  func updateEventEndDate(_ eventID: UUID, date: Date) {
    guard let index = events.firstIndex(where: { $0.id == eventID }) else { return }
    events[index].endDate = date
    // Update meal plan name when end date changes
    updateDefaultMealPlanName()
  }

  func updateEventNotes(_ eventID: UUID, notes: String) {
    guard let index = events.firstIndex(where: { $0.id == eventID }) else { return }
    events[index].notes = notes
  }

  // MARK: - Search Functions (per event)

  func searchForMeals(for event: MealPlanEvent) async {
    guard !event.searchQuery.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
      if let index = events.firstIndex(where: { $0.id == event.id }) {
        events[index].searchResults = []
      }
      return
    }

    guard let index = events.firstIndex(where: { $0.id == event.id }) else { return }

    events[index].isSearching = true
    events[index].searchError = nil

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
      request.query = event.searchQuery
      request.useSearchService = false  // disabled for local testing

      // Execute search
      let response = try await clientManager.client.mealPlanning.searchForMeals(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      events[index].searchResults = response.results
    } catch {
      events[index].searchError = "Failed to search meals: \(error.localizedDescription)"
      print("❌ Error searching for meals: \(error)")
      events[index].searchResults = []
    }

    events[index].isSearching = false
  }

  // MARK: - Selection Functions (per event)

  func toggleMealSelection(_ meal: Mealplanning_Meal, in event: MealPlanEvent) {
    guard let index = events.firstIndex(where: { $0.id == event.id }) else { return }

    if let mealIndex = events[index].selectedMeals.firstIndex(where: { $0.id == meal.id }) {
      events[index].selectedMeals.remove(at: mealIndex)
      // Remove scale when meal is deselected
      events[index].mealScales.removeValue(forKey: meal.id)
    } else {
      events[index].selectedMeals.append(meal)
      // Set default scale to 1.0 when meal is selected
      events[index].mealScales[meal.id] = 1.0
    }
  }

  func isMealSelected(_ meal: Mealplanning_Meal, in event: MealPlanEvent) -> Bool {
    guard let index = events.firstIndex(where: { $0.id == event.id }) else { return false }
    return events[index].selectedMeals.contains(where: { $0.id == meal.id })
  }

  func removeSelectedMeal(_ meal: Mealplanning_Meal, from event: MealPlanEvent) {
    guard let index = events.firstIndex(where: { $0.id == event.id }) else { return }
    events[index].selectedMeals.removeAll(where: { $0.id == meal.id })
    // Remove scale when meal is removed
    events[index].mealScales.removeValue(forKey: meal.id)
    // Meal will automatically reappear in search results since we filter them
  }

  func getMealScale(_ meal: Mealplanning_Meal, in event: MealPlanEvent) -> Float {
    guard let index = events.firstIndex(where: { $0.id == event.id }) else { return 1.0 }
    return events[index].mealScales[meal.id] ?? 1.0
  }

  func setMealScale(_ meal: Mealplanning_Meal, scale: Float, in event: MealPlanEvent) {
    guard let index = events.firstIndex(where: { $0.id == event.id }) else { return }
    events[index].mealScales[meal.id] = scale
  }

  /// Get filtered search results for an event (excluding already selected meals)
  func filteredSearchResults(for event: MealPlanEvent) -> [Mealplanning_Meal] {
    guard let index = events.firstIndex(where: { $0.id == event.id }) else { return [] }
    let selectedMealIDs = Set(events[index].selectedMeals.map { $0.id })
    return events[index].searchResults.filter { !selectedMealIDs.contains($0.id) }
  }

  // MARK: - Validation

  struct ValidationErrors {
    var hasNameError: Bool = false
    var hasVotingDeadlineError: Bool = false
    var eventTimelineErrors: [UUID: Bool] = [:]  // Event ID -> has error
  }

  func getValidationErrors() -> ValidationErrors {
    var errors = ValidationErrors()
    let now = Date()
    let twelveHoursInSeconds: TimeInterval = 12 * 60 * 60

    // Check meal plan name
    if mealPlanName.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty {
      errors.hasNameError = true
    }

    // Check voting deadline (must be at least 12 hours from now)
    let timeUntilDeadline = votingDeadline.timeIntervalSince(now)
    if timeUntilDeadline < twelveHoursInSeconds {
      errors.hasVotingDeadlineError = true
    }

    // Check each event's start and end times (must be at least 12 hours from now)
    for event in events {
      let timeUntilStart = event.startDate.timeIntervalSince(now)
      let timeUntilEnd = event.endDate.timeIntervalSince(now)
      if timeUntilStart < twelveHoursInSeconds || timeUntilEnd < twelveHoursInSeconds {
        errors.eventTimelineErrors[event.id] = true
      }
    }

    return errors
  }

  func canCreateMealPlan() -> Bool {
    let errors = getValidationErrors()
    return !errors.hasNameError
      && !errors.hasVotingDeadlineError
      && errors.eventTimelineErrors.isEmpty
      && events.contains { !$0.selectedMeals.isEmpty }
  }

  func validateEvents() -> String? {
    // Check that all events have at least one meal selected
    for event in events where event.selectedMeals.isEmpty {
      return "Each event must have at least one meal selected"
    }

    // Check that events don't span more than a week
    let sortedEvents = events.sorted(by: { $0.startDate < $1.startDate })
    guard let firstEvent = sortedEvents.first,
      let lastEvent = sortedEvents.last
    else {
      return "At least one event is required"
    }

    let daysBetween =
      Calendar.current.dateComponents([.day], from: firstEvent.startDate, to: lastEvent.startDate)
      .day ?? 0
    if daysBetween > 7 {
      return "Meal plan events cannot span more than a week"
    }

    return nil
  }

  // MARK: - Meal Plan Creation

  func createMealPlan() async -> Bool {
    guard !mealPlanName.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
      creationError = "Meal plan name is required"
      return false
    }

    // Validate events
    if let validationError = validateEvents() {
      creationError = validationError
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

      // Create events from the events array
      for event in events {
        var eventInput = Mealplanning_MealPlanEventCreationRequestInput()
        eventInput.startsAt = dateToTimestamp(event.startDate)
        eventInput.endsAt = dateToTimestamp(event.endDate)
        eventInput.mealName = event.mealType
        eventInput.notes = event.notes

        // Create options for each selected meal in this event
        for meal in event.selectedMeals {
          var optionInput = Mealplanning_MealPlanOptionCreationRequestInput()
          optionInput.mealID = meal.id
          optionInput.mealScale = event.mealScales[meal.id] ?? 1.0
          optionInput.notes = ""

          // Add option selections for this meal's recipes
          for component in meal.components {
            let recipeID = component.recipe.id
            // Add ingredient option selections (only ingredients have selectable options)
            if let stepSelections = recipeOptionSelections[recipeID] {
              for (stepID, indexSelections) in stepSelections {
                for (ingredientIndex, selectedOptionIndex) in indexSelections {
                  var selectionInput =
                    Mealplanning_MealPlanRecipeOptionSelectionCreationRequestInput()
                  selectionInput.recipeID = recipeID
                  selectionInput.recipeStepID = stepID
                  selectionInput.ingredientIndex = ingredientIndex
                  selectionInput.selectedOptionIndex = selectedOptionIndex
                  selectionInput.selectionType = .ingredient
                  optionInput.selections.append(selectionInput)
                }
              }
            }
          }

          eventInput.options.append(optionInput)
        }

        input.events.append(eventInput)
      }

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

  // MARK: - Option Selection Helpers

  func setOptionSelections(
    ingredientSelections: [String: [String: [UInt32: UInt32]]]
  ) {
    recipeOptionSelections = ingredientSelections
  }

  /// Collects all unique recipes (including associated recipes) from selected meals that have option groups
  /// - Parameter meals: Array of meals to check
  /// - Returns: Set of recipe IDs that have option groups
  func collectRecipesWithOptions(from meals: [Mealplanning_Meal]) -> Set<String> {
    var recipeIDsWithOptions: Set<String> = []

    // Collect all recipes from selected meals
    var allRecipes: [Mealplanning_Recipe] = []
    for meal in meals {
      for component in meal.components {
        allRecipes.append(component.recipe)
        // Also include associated recipes
        allRecipes.append(contentsOf: component.recipe.associatedRecipes)
      }
    }

    // Check if any recipes have option groups
    for recipe in allRecipes {
      // Debug: Check if recipe has steps loaded
      if recipe.steps.isEmpty {
        print("⚠️ Recipe '\(recipe.name)' (ID: \(recipe.id)) has no steps loaded - may need to fetch full recipe details")
        continue
      }
      
      let hasOptions = recipe.steps.contains { step in
        // Check for ingredient option groups (only ingredients have selectable options)
        // Group ingredients by their index to find option groups
        var indexGroups: [UInt32: [Mealplanning_RecipeStepIngredient]] = [:]
        
        for ingredient in step.ingredients {
          // Index 0 typically means not in an option group
          if ingredient.index != 0 {
            let index = ingredient.index
            if indexGroups[index] == nil {
              indexGroups[index] = []
            }
            indexGroups[index]?.append(ingredient)
          }
        }
        
        // Check if any index group has multiple ingredients (i.e., has options)
        let hasIngredientOptions = indexGroups.values.contains { $0.count > 1 }
        if hasIngredientOptions {
          print("  📝 Step \(step.id) in recipe '\(recipe.name)' has option groups: \(indexGroups.filter { $0.value.count > 1 }.keys.sorted())")
        }
        return hasIngredientOptions
      }

      if hasOptions {
        recipeIDsWithOptions.insert(recipe.id)
        print("✅ Found recipe with options: \(recipe.name) (ID: \(recipe.id))")
      } else {
        print("ℹ️ Recipe '\(recipe.name)' (ID: \(recipe.id)) has no option groups")
      }
    }
    
    print("📊 collectRecipesWithOptions: Found \(recipeIDsWithOptions.count) recipe(s) with options out of \(allRecipes.count) total recipe(s)")

    return recipeIDsWithOptions
  }

  /// Gets all unique recipes (including associated recipes) from selected meals
  /// - Parameter meals: Array of meals to extract recipes from
  /// - Returns: Array of unique recipes
  func getAllRecipes(from meals: [Mealplanning_Meal]) -> [Mealplanning_Recipe] {
    var allRecipes: [Mealplanning_Recipe] = []
    for meal in meals {
      for component in meal.components {
        allRecipes.append(component.recipe)
        // Also include associated recipes
        allRecipes.append(contentsOf: component.recipe.associatedRecipes)
      }
    }

    // Return unique recipes (by ID)
    var seenIDs: Set<String> = []
    return allRecipes.filter { recipe in
      guard !seenIDs.contains(recipe.id) else {
        return false
      }
      seenIDs.insert(recipe.id)
      return true
    }
  }

  /// Checks if a recipe has unselected options
  /// - Parameter recipeID: The recipe ID to check
  /// - Returns: True if the recipe has option groups that haven't been selected
  func hasUnselectedOptions(for recipeID: String) -> Bool {
    // If recipe has no selections at all, it might have unselected options
    // This is a simple check - in practice, you'd want to compare against
    // the actual option groups in the recipe
    return false  // For now, we assume all options are selected or use defaults
  }

  /// Gets default option selections for a recipe (optionIndex: 0 for all option groups)
  /// - Parameter recipe: The recipe to get defaults for
  /// - Returns: A dictionary mapping stepID -> (ingredientIndex -> selectedOptionIndex)
  func getDefaultOptionSelections(for recipe: Mealplanning_Recipe) -> [String: [UInt32: UInt32]] {
    var defaults: [String: [UInt32: UInt32]] = [:]

    // Check all steps for ingredient option groups
    for step in recipe.steps {
      // Group ingredients by index to find option groups
      var optionGroupsByIndex: [UInt32: [Mealplanning_RecipeStepIngredient]] = [:]

      for ingredient in step.ingredients where ingredient.index != 0 {
        // Index 0 typically means not in an option group
        if optionGroupsByIndex[ingredient.index] == nil {
          optionGroupsByIndex[ingredient.index] = []
        }
        optionGroupsByIndex[ingredient.index]?.append(ingredient)
      }

      // For each option group, set default to optionIndex 0 (first option)
      for (index, groupIngredients) in optionGroupsByIndex where groupIngredients.count > 1 {
        // Only create a default if there are multiple options (more than one ingredient with same index)
        // Find the ingredient with optionIndex 0 (default)
        if let defaultIngredient = groupIngredients.first(where: { $0.optionIndex == 0 }) {
          if defaults[step.id] == nil {
            defaults[step.id] = [:]
          }
          defaults[step.id]?[index] = 0  // Default to optionIndex 0
        }
      }
    }

    return defaults
  }

  // MARK: - Helper Functions

  private func dateToTimestamp(_ date: Date) -> SwiftProtobuf.Google_Protobuf_Timestamp {
    var timestamp = SwiftProtobuf.Google_Protobuf_Timestamp()
    timestamp.seconds = Int64(date.timeIntervalSince1970)
    timestamp.nanos = Int32(
      (date.timeIntervalSince1970 - Double(timestamp.seconds)) * 1_000_000_000)
    return timestamp
  }

}
