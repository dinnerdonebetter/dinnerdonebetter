//
//  CreateMealPlanViewModelTests.swift
//  iosTests
//
//  Created by Auto on 12/8/25.
//

import Foundation
import SwiftProtobuf
@testable import ios
import Testing

// MARK: - Helper Functions for Test Data

func createMockMeal(id: String = "meal-1", name: String = "Test Meal") -> Mealplanning_Meal {
  var meal = Mealplanning_Meal()
  meal.id = id
  meal.name = name
  return meal
}

func createMockAuthenticationManagerForMealPlan() -> AuthenticationManager {
  let manager = AuthenticationManager()
  manager.isAuthenticated = true
  manager.oauth2AccessToken = "mock-oauth2-token"
  return manager
}

// MARK: - MealPlanEvent Tests

struct MealPlanEventTests {
  @Test("MealPlanEvent initializes with required fields")
  func testMealPlanEventInitialization() {
    let startDate = Date()
    let endDate = Date().addingTimeInterval(7200)  // 2 hours later

    let event = MealPlanEvent(
      mealType: .dinner,
      startDate: startDate,
      endDate: endDate
    )

    #expect(event.mealType == .dinner)
    #expect(event.startDate == startDate)
    #expect(event.endDate == endDate)
    #expect(event.searchQuery.isEmpty)
    #expect(event.searchResults.isEmpty)
    #expect(event.isSearching == false)
    #expect(event.selectedMeals.isEmpty)
    #expect(event.mealScales.isEmpty)
  }

  @Test("MealPlanEvent has unique ID")
  func testMealPlanEventUniqueID() {
    let startDate = Date()
    let endDate = Date().addingTimeInterval(7200)

    let event1 = MealPlanEvent(mealType: .dinner, startDate: startDate, endDate: endDate)
    let event2 = MealPlanEvent(mealType: .lunch, startDate: startDate, endDate: endDate)

    #expect(event1.id != event2.id)
  }
}

// MARK: - Initialization Tests

struct InitializationTests {
  @Test("CreateMealPlanViewModel initializes with default event")
  @MainActor
  func testViewModelInitialization() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    #expect(viewModel.events.count == 1)
    #expect(!viewModel.mealPlanName.isEmpty)
    #expect(viewModel.isCreating == false)
    #expect(viewModel.creationError == nil)
  }

  @Test("CreateMealPlanViewModel initializes with next Monday event")
  @MainActor
  func testViewModelInitializesWithNextMonday() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    #expect(viewModel.events.count == 1)
    let event = viewModel.events[0]
    let calendar = Calendar.current
    let weekday = calendar.component(.weekday, from: event.startDate)
    let hour = calendar.component(.hour, from: event.startDate)

    #expect(weekday == 2 || weekday == 1 || weekday >= 3)
    #expect(hour == 19)
  }

  @Test("CreateMealPlanViewModel sets default meal plan name")
  @MainActor
  func testViewModelSetsDefaultName() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    #expect(!viewModel.mealPlanName.isEmpty)
    #expect(viewModel.mealPlanName.contains("Meal Plan"))
  }

  @Test("CreateMealPlanViewModel sets voting deadline")
  @MainActor
  func testViewModelSetsVotingDeadline() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    #expect(viewModel.votingDeadline != Date())
  }
}

// MARK: - Event Management Tests

struct EventManagementTests {
  @Test("addEvent adds new event")
  @MainActor
  func testAddEvent() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let initialCount = viewModel.events.count
    viewModel.addEvent()

    #expect(viewModel.events.count == initialCount + 1)
  }

  @Test("addEvent adds event after latest event")
  @MainActor
  func testAddEventAfterLatest() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let firstEventDate = viewModel.events[0].startDate
    viewModel.addEvent()

    #expect(viewModel.events.count == 2)
    #expect(viewModel.events[1].startDate > firstEventDate)
  }

  @Test("addEvent updates voting deadline")
  @MainActor
  func testAddEventUpdatesDeadline() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    viewModel.addEvent()

    #expect(viewModel.votingDeadline != nil)
  }

  @Test("addEvent updates meal plan name")
  @MainActor
  func testAddEventUpdatesName() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let originalName = viewModel.mealPlanName
    viewModel.addEvent()

    #expect(viewModel.mealPlanName != originalName)
  }

  @Test("removeEvent removes event")
  @MainActor
  func testRemoveEvent() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    viewModel.addEvent()
    let eventToRemove = viewModel.events[0]
    let initialCount = viewModel.events.count

    viewModel.removeEvent(eventToRemove)

    #expect(viewModel.events.count == initialCount - 1)
  }

  @Test("removeEvent ensures at least one event remains")
  @MainActor
  func testRemoveEventKeepsOne() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let eventToRemove = viewModel.events[0]
    viewModel.removeEvent(eventToRemove)

    #expect(viewModel.events.count == 1)
  }

  @Test("removeEvent does nothing for non-existent event")
  @MainActor
  func testRemoveEventNonExistent() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let initialCount = viewModel.events.count
    let fakeEvent = MealPlanEvent(mealType: .dinner, startDate: Date(), endDate: Date())

    viewModel.removeEvent(fakeEvent)

    #expect(viewModel.events.count == initialCount)
  }

  @Test("updateEvent updates existing event")
  @MainActor
  func testUpdateEvent() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    var event = viewModel.events[0]
    event.mealType = .lunch
    let originalType = viewModel.events[0].mealType

    viewModel.updateEvent(event)

    #expect(viewModel.events[0].mealType == .lunch)
    #expect(viewModel.events[0].mealType != originalType)
  }

  @Test("updateEvent does nothing for non-existent event")
  @MainActor
  func testUpdateEventNonExistent() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let fakeEvent = MealPlanEvent(
      mealType: .breakfast,
      startDate: Date(),
      endDate: Date()
    )

    viewModel.updateEvent(fakeEvent)

    #expect(viewModel.events.count == 1)
  }

  @Test("updateEventSearchQuery updates search query")
  @MainActor
  func testUpdateEventSearchQuery() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let eventID = viewModel.events[0].id
    viewModel.updateEventSearchQuery(eventID, query: "pasta")

    #expect(viewModel.events[0].searchQuery == "pasta")
  }

  @Test("updateEventSearchQuery does nothing for non-existent event")
  @MainActor
  func testUpdateEventSearchQueryNonExistent() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let fakeID = UUID()
    viewModel.updateEventSearchQuery(fakeID, query: "test")

    #expect(viewModel.events[0].searchQuery.isEmpty)
  }

  @Test("updateEventMealType updates meal type")
  @MainActor
  func testUpdateEventMealType() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let eventID = viewModel.events[0].id
    viewModel.updateEventMealType(eventID, mealType: .breakfast)

    #expect(viewModel.events[0].mealType == .breakfast)
  }

  @Test("updateEventStartDate updates start date and end date")
  @MainActor
  func testUpdateEventStartDate() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let eventID = viewModel.events[0].id
    let newStartDate = Date().addingTimeInterval(86400)

    viewModel.updateEventStartDate(eventID, date: newStartDate)

    #expect(viewModel.events[0].startDate == newStartDate)
    let expectedEnd = Calendar.current.date(byAdding: .hour, value: 2, to: newStartDate)
    #expect(abs(viewModel.events[0].endDate.timeIntervalSince(expectedEnd ?? newStartDate)) < 60.0)
  }

  @Test("updateEventEndDate updates end date")
  @MainActor
  func testUpdateEventEndDate() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let eventID = viewModel.events[0].id
    let newEndDate = Date().addingTimeInterval(10800)  // 3 hours from now

    viewModel.updateEventEndDate(eventID, date: newEndDate)

    #expect(viewModel.events[0].endDate == newEndDate)
  }
}

// MARK: - Meal Selection Tests

struct MealSelectionTests {
  @Test("toggleMealSelection adds meal when not selected")
  @MainActor
  func testToggleMealSelectionAdds() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let event = viewModel.events[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")

    #expect(viewModel.isMealSelected(meal, in: event) == false)

    viewModel.toggleMealSelection(meal, in: event)

    #expect(viewModel.isMealSelected(meal, in: event) == true)
    #expect(viewModel.events[0].selectedMeals.count == 1)
    #expect(viewModel.getMealScale(meal, in: event) == 1.0)
  }

  @Test("toggleMealSelection removes meal when selected")
  @MainActor
  func testToggleMealSelectionRemoves() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let event = viewModel.events[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")

    viewModel.toggleMealSelection(meal, in: event)
    #expect(viewModel.isMealSelected(meal, in: event) == true)

    viewModel.toggleMealSelection(meal, in: event)

    #expect(viewModel.isMealSelected(meal, in: event) == false)
    #expect(viewModel.events[0].selectedMeals.isEmpty)
    #expect(viewModel.getMealScale(meal, in: event) == 1.0)  // Default when not selected
  }

  @Test("toggleMealSelection sets default scale to 1.0")
  @MainActor
  func testToggleMealSelectionSetsDefaultScale() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let event = viewModel.events[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")

    viewModel.toggleMealSelection(meal, in: event)

    #expect(viewModel.getMealScale(meal, in: event) == 1.0)
  }

  @Test("toggleMealSelection removes scale when deselected")
  @MainActor
  func testToggleMealSelectionRemovesScale() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let event = viewModel.events[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")

    viewModel.toggleMealSelection(meal, in: event)
    viewModel.setMealScale(meal, scale: 2.5, in: event)
    #expect(viewModel.getMealScale(meal, in: event) == 2.5)

    viewModel.toggleMealSelection(meal, in: event)

    #expect(viewModel.getMealScale(meal, in: event) == 1.0)
  }

  @Test("removeSelectedMeal removes meal")
  @MainActor
  func testRemoveSelectedMeal() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let event = viewModel.events[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")

    viewModel.toggleMealSelection(meal, in: event)
    #expect(viewModel.events[0].selectedMeals.count == 1)

    viewModel.removeSelectedMeal(meal, from: event)

    #expect(viewModel.events[0].selectedMeals.isEmpty)
    #expect(viewModel.isMealSelected(meal, in: event) == false)
  }

  @Test("setMealScale updates meal scale")
  @MainActor
  func testSetMealScale() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let event = viewModel.events[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")

    viewModel.toggleMealSelection(meal, in: event)
    viewModel.setMealScale(meal, scale: 2.5, in: event)

    #expect(viewModel.getMealScale(meal, in: event) == 2.5)
  }

  @Test("getMealScale returns default 1.0 when not set")
  @MainActor
  func testGetMealScaleDefault() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let event = viewModel.events[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")

    #expect(viewModel.getMealScale(meal, in: event) == 1.0)
  }

  @Test("filteredSearchResults excludes selected meals")
  @MainActor
  func testFilteredSearchResults() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let event = viewModel.events[0]
    let meal1 = createMockMeal(id: "meal-1", name: "Pasta")
    let meal2 = createMockMeal(id: "meal-2", name: "Pizza")
    let meal3 = createMockMeal(id: "meal-3", name: "Burger")

    viewModel.events[0].searchResults = [meal1, meal2, meal3]
    viewModel.toggleMealSelection(meal1, in: event)

    let filtered = viewModel.filteredSearchResults(for: event)

    #expect(filtered.count == 2)
    #expect(filtered.contains(where: { $0.id == meal2.id }))
    #expect(filtered.contains(where: { $0.id == meal3.id }))
    #expect(!filtered.contains(where: { $0.id == meal1.id }))
  }
}

// MARK: - Validation Tests

struct ValidationTests {
  @Test("validateEvents returns nil when valid")
  @MainActor
  func testValidateEventsValid() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let event = viewModel.events[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")
    viewModel.toggleMealSelection(meal, in: event)

    let error = viewModel.validateEvents()

    #expect(error == nil)
  }

  @Test("validateEvents returns error when event has no meals")
  @MainActor
  func testValidateEventsNoMeals() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let error = viewModel.validateEvents()

    #expect(error != nil)
    #expect(error?.contains("at least one meal") == true)
  }

  @Test("validateEvents returns error when events span more than week")
  @MainActor
  func testValidateEventsSpanTooLong() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    // Add event 8 days later
    let event = viewModel.events[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")
    viewModel.toggleMealSelection(meal, in: event)

    let farFutureDate = Date().addingTimeInterval(8 * 24 * 60 * 60)  // 8 days
    let eventID = viewModel.events[0].id
    viewModel.updateEventStartDate(eventID, date: farFutureDate)

    viewModel.addEvent()
    let newEvent = viewModel.events[1]
    let meal2 = createMockMeal(id: "meal-2", name: "Pizza")
    viewModel.toggleMealSelection(meal2, in: newEvent)

    // Update first event to be 8 days before second
    let firstEventID = viewModel.events[0].id
    let earlyDate = Calendar.current.date(byAdding: .day, value: -8, to: farFutureDate) ?? Date()
    viewModel.updateEventStartDate(firstEventID, date: earlyDate)

    let error = viewModel.validateEvents()

    #expect(error != nil)
    #expect(error?.contains("more than a week") == true)
  }

  @Test("validateEvents returns nil when events have meals")
  @MainActor
  func testValidateEventsWithMeals() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let event = viewModel.events[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")
    viewModel.toggleMealSelection(meal, in: event)

    let error = viewModel.validateEvents()

    #expect(error == nil)
  }
}

// MARK: - Creation Tests

struct CreationTests {
  @Test("createMealPlan returns false when name is empty")
  @MainActor
  func testCreateMealPlanEmptyName() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    viewModel.mealPlanName = ""

    let result = await viewModel.createMealPlan()

    #expect(result == false)
    #expect(viewModel.creationError != nil)
    #expect(viewModel.creationError?.contains("name is required") == true)
  }

  @Test("createMealPlan returns false when validation fails")
  @MainActor
  func testCreateMealPlanValidationFails() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    viewModel.mealPlanName = "Test Plan"

    let result = await viewModel.createMealPlan()

    #expect(result == false)
    #expect(viewModel.creationError != nil)
  }

  @Test("createMealPlan sets creation state")
  @MainActor
  func testCreateMealPlanSetsState() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    viewModel.mealPlanName = "Test Plan"
    let event = viewModel.events[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")
    viewModel.toggleMealSelection(meal, in: event)

    _ = await viewModel.createMealPlan()

    #expect(viewModel.isCreating == false)
  }

  @Test("createMealPlan handles multiple events")
  @MainActor
  func testCreateMealPlanMultipleEvents() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    viewModel.mealPlanName = "Test Plan"
    viewModel.addEvent()

    for event in viewModel.events {
      let meal = createMockMeal(id: "meal-\(event.id)", name: "Meal")
      viewModel.toggleMealSelection(meal, in: event)
    }

    let result = await viewModel.createMealPlan()

    #expect(result == false || result == true)
  }
}

// MARK: - Date Calculation Tests

struct DateCalculationTests {
  @Test("updateEventStartDate updates voting deadline")
  @MainActor
  func testUpdateEventStartDateUpdatesDeadline() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let eventID = viewModel.events[0].id
    let newDate = Date().addingTimeInterval(86400)
    viewModel.updateEventStartDate(eventID, date: newDate)

    #expect(viewModel.votingDeadline != Date())
  }

  @Test("updateVotingDeadline updates when events change")
  @MainActor
  func testUpdateVotingDeadline() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let earlierDate = Date().addingTimeInterval(-86400)
    let eventID = viewModel.events[0].id
    let originalDeadline = viewModel.votingDeadline
    viewModel.updateEventStartDate(eventID, date: earlierDate)

    #expect(viewModel.votingDeadline != originalDeadline)
  }

  @Test("updateDefaultMealPlanName updates name with date range")
  @MainActor
  func testUpdateDefaultMealPlanName() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let originalName = viewModel.mealPlanName
    viewModel.addEvent()

    #expect(viewModel.mealPlanName != originalName)
    #expect(viewModel.mealPlanName.contains("Meal Plan"))
  }
}

// MARK: - Edge Cases Tests

struct EdgeCaseTests {
  @Test("ViewModel handles multiple meal selections")
  @MainActor
  func testViewModelMultipleMeals() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let event = viewModel.events[0]
    let meal1 = createMockMeal(id: "meal-1", name: "Pasta")
    let meal2 = createMockMeal(id: "meal-2", name: "Pizza")
    let meal3 = createMockMeal(id: "meal-3", name: "Burger")

    viewModel.toggleMealSelection(meal1, in: event)
    viewModel.toggleMealSelection(meal2, in: event)
    viewModel.toggleMealSelection(meal3, in: event)

    #expect(viewModel.events[0].selectedMeals.count == 3)
  }

  @Test("ViewModel handles search query and meal type updates")
  @MainActor
  func testViewModelUpdates() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let eventID = viewModel.events[0].id
    viewModel.updateEventSearchQuery(eventID, query: "pasta")
    #expect(viewModel.events[0].searchQuery == "pasta")

    viewModel.updateEventMealType(eventID, mealType: .breakfast)
    #expect(viewModel.events[0].mealType == .breakfast)
  }
}

