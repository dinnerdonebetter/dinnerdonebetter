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

// MARK: - Initialization Tests

struct InitializationTests {
  @Test("CreateMealPlanViewModel initializes with wizard state")
  @MainActor
  func testViewModelInitialization() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    #expect(viewModel.wizardStep == .weekSelection)
    #expect(viewModel.selectedWeekOffset == 0)
    #expect(viewModel.selectedDates.isEmpty)
    #expect(viewModel.dayMeals.isEmpty)
    #expect(viewModel.isCreating == false)
    #expect(viewModel.creationError == nil)
  }

  @Test("CreateMealPlanViewModel has default meal plan name")
  @MainActor
  func testViewModelSetsDefaultName() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    #expect(!viewModel.mealPlanName.isEmpty)
    #expect(viewModel.mealPlanName.contains("Meal Plan"))
  }
}

// MARK: - Date Selection Tests

struct DateSelectionTests {
  @Test("toggleDateSelection adds date when not selected")
  @MainActor
  func testToggleDateSelectionAdds() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    viewModel.goToNextWeek()  // Use next week so dates are planable (not past/today)

    let date = viewModel.displayedWeekDays[0]
    viewModel.toggleDateSelection(date)

    #expect(!viewModel.selectedDates.isEmpty)
  }

  @Test("toggleDateSelection removes date when selected")
  @MainActor
  func testToggleDateSelectionRemoves() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    viewModel.goToNextWeek()  // Use next week so dates are planable (not past/today)

    let date = viewModel.displayedWeekDays[0]
    #expect(viewModel.isDateSelected(date) == false)

    viewModel.toggleDateSelection(date)
    #expect(viewModel.isDateSelected(date) == true)

    viewModel.toggleDateSelection(date)
    #expect(viewModel.isDateSelected(date) == false)
  }

  @Test("setDateRangeSelection selects range of days")
  @MainActor
  func testSetDateRangeSelection() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    viewModel.goToNextWeek()  // Use next week so dates are planable (not past/today)
    let days = viewModel.displayedWeekDays
    #expect(days.count >= 3)

    viewModel.setDateRangeSelection(from: days[0], to: days[2])

    #expect(viewModel.selectedDates.count == 3)
  }
}

// MARK: - Meal Assignment Tests

struct MealAssignmentTests {
  @Test("assignMeal adds meal to date")
  @MainActor
  func testAssignMeal() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    viewModel.goToNextWeek()  // Use next week so dates are planable (not past/today)
    let days = viewModel.displayedWeekDays
    viewModel.setDateRangeSelection(from: days[0], to: days[0])

    let date = viewModel.selectedDates[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")

    viewModel.assignMeal(meal, to: date)

    #expect(viewModel.mealForDate(date)?.id == meal.id)
  }

  @Test("removeMeal removes meal from date")
  @MainActor
  func testRemoveMeal() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    viewModel.goToNextWeek()  // Use next week so dates are planable (not past/today)
    let days = viewModel.displayedWeekDays
    viewModel.setDateRangeSelection(from: days[0], to: days[0])

    let date = viewModel.selectedDates[0]
    let meal = createMockMeal(id: "meal-1", name: "Pasta")
    viewModel.assignMeal(meal, to: date)
    #expect(viewModel.mealForDate(date) != nil)

    viewModel.removeMeal(from: date)

    #expect(viewModel.mealForDate(date) == nil)
  }

  @Test("allDaysHaveMeals returns false when any day missing meal")
  @MainActor
  func testAllDaysHaveMealsPartial() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    viewModel.goToNextWeek()  // Use next week so dates are planable (not past/today)
    let days = viewModel.displayedWeekDays
    viewModel.setDateRangeSelection(from: days[0], to: days[1])

    let dates = viewModel.selectedDates
    #expect(!dates.isEmpty)
    viewModel.assignMeal(createMockMeal(id: "m1", name: "M1"), to: dates[0])

    #expect(viewModel.allDaysHaveMeals == false)
  }

  @Test("allDaysHaveMeals returns true when all days have meals")
  @MainActor
  func testAllDaysHaveMealsComplete() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    viewModel.goToNextWeek()  // Use next week so dates are planable (not past/today)
    let days = viewModel.displayedWeekDays
    viewModel.setDateRangeSelection(from: days[0], to: days[1])

    for (index, date) in viewModel.selectedDates.enumerated() {
      let meal = createMockMeal(id: "meal-\(index)", name: "Meal \(index)")
      viewModel.assignMeal(meal, to: date)
    }

    #expect(viewModel.allDaysHaveMeals == true)
  }
}

// MARK: - Search Tests

struct SearchTests {
  @Test("filteredSearchResults excludes assigned meals")
  @MainActor
  func testFilteredSearchResults() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    viewModel.goToNextWeek()  // Use next week so dates are planable (not past/today)
    let days = viewModel.displayedWeekDays
    viewModel.setDateRangeSelection(from: days[0], to: days[0])

    let meal1 = createMockMeal(id: "meal-1", name: "Pasta")
    let meal2 = createMockMeal(id: "meal-2", name: "Pizza")
    let meal3 = createMockMeal(id: "meal-3", name: "Burger")
    viewModel.searchResults = [meal1, meal2, meal3]
    viewModel.assignMeal(meal1, to: viewModel.selectedDates[0])

    let filtered = viewModel.filteredSearchResults(for: viewModel.selectedDates[0])

    #expect(filtered.count == 2)
    #expect(filtered.contains(where: { $0.id == meal2.id }))
    #expect(filtered.contains(where: { $0.id == meal3.id }))
    #expect(!filtered.contains(where: { $0.id == meal1.id }))
  }

  @Test("meal forId finds meal in search results")
  @MainActor
  func testMealForIdFromSearch() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    let meal = createMockMeal(id: "meal-1", name: "Pasta")
    viewModel.searchResults = [meal]

    let found = viewModel.meal(forId: "meal-1")

    #expect(found?.id == "meal-1")
  }

  @Test("meal forId finds meal in dayMeals")
  @MainActor
  func testMealForIdFromDayMeals() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    viewModel.goToNextWeek()  // Use next week so dates are planable (not past/today)
    let days = viewModel.displayedWeekDays
    viewModel.setDateRangeSelection(from: days[0], to: days[0])
    let meal = createMockMeal(id: "meal-1", name: "Pasta")
    viewModel.assignMeal(meal, to: viewModel.selectedDates[0])

    let found = viewModel.meal(forId: "meal-1")

    #expect(found?.id == "meal-1")
  }
}

// MARK: - Validation Tests

struct ValidationTests {
  @Test("canCreateMealPlan returns false when no dates selected")
  @MainActor
  func testCanCreateNoDates() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    #expect(viewModel.canCreateMealPlan() == false)
  }

  @Test("canCreateMealPlan returns false when days lack meals")
  @MainActor
  func testCanCreateMissingMeals() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    viewModel.goToNextWeek()  // Use next week so dates are planable (not past/today)
    let days = viewModel.displayedWeekDays
    viewModel.setDateRangeSelection(from: days[0], to: days[0])

    #expect(viewModel.canCreateMealPlan() == false)
  }

  @Test("canCreateMealPlan returns true when all days have meals")
  @MainActor
  func testCanCreateValid() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    viewModel.goToNextWeek()  // Use next week so dates are planable (not past/today)
    let days = viewModel.displayedWeekDays
    viewModel.setDateRangeSelection(from: days[0], to: days[1])

    for (index, date) in viewModel.selectedDates.enumerated() {
      let meal = createMockMeal(id: "meal-\(index)", name: "Meal \(index)")
      viewModel.assignMeal(meal, to: date)
    }

    #expect(viewModel.canCreateMealPlan() == true)
  }
}

// MARK: - Creation Tests

struct CreationTests {
  @Test("createMealPlan returns false when canCreate is false")
  @MainActor
  func testCreateMealPlanInvalid() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    let result = await viewModel.createMealPlan()

    #expect(result == false)
    #expect(viewModel.creationError != nil)
  }
}

// MARK: - Week Navigation Tests

struct WeekNavigationTests {
  @Test("goToNextWeek increments offset")
  @MainActor
  func testGoToNextWeek() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    #expect(viewModel.selectedWeekOffset == 0)
    viewModel.goToNextWeek()
    #expect(viewModel.selectedWeekOffset == 1)
  }

  @Test("goToPreviousWeek decrements offset")
  @MainActor
  func testGoToPreviousWeek() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)
    viewModel.goToNextWeek()

    viewModel.goToPreviousWeek()
    #expect(viewModel.selectedWeekOffset == 0)
  }

  @Test("goToPreviousWeek does not go below zero")
  @MainActor
  func testGoToPreviousWeekMinZero() async {
    let authManager = createMockAuthenticationManagerForMealPlan()
    let viewModel = CreateMealPlanViewModel(authManager: authManager)

    viewModel.goToPreviousWeek()
    #expect(viewModel.selectedWeekOffset == 0)
  }
}
