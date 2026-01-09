//
//  CreateMealPlanView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

extension Notification.Name {
  static let mealPlanCreated = Notification.Name("mealPlanCreated")
}

// swiftlint:disable:next type_body_length
struct CreateMealPlanView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Environment(\.dismiss) var dismiss
  @Environment(\.horizontalSizeClass) private var horizontalSizeClass
  @State private var viewModel: CreateMealPlanViewModel?
  @FocusState var focusedField: Field?
  @State private var showOptionSelectionModal = false
  @State private var recipesForOptionSelection: [Mealplanning_Recipe] = []

  enum Field: Hashable {
    case mealPlanName
    case searchQuery(UUID)  // Per-event search
  }

  private var isRegularWidth: Bool {
    horizontalSizeClass == .regular
  }

  var body: some View {
    Group {
      if let viewModel = viewModel {
        ScrollView(.vertical, showsIndicators: true) {
          VStack(spacing: 24) {
            // Meal Plan Details Section - responsive layout
            if isRegularWidth {
              mealPlanDetailsSectionHorizontal(viewModel: viewModel)
            } else {
              mealPlanDetailsSection(viewModel: viewModel)
            }

            // Events Section
            eventsSection(viewModel: viewModel)

            // Error Messages
            if let error = viewModel.creationError {
              errorMessage(error)
            }

            // Create Button
            createButton(
              viewModel: viewModel,
              showOptionSelectionModal: $showOptionSelectionModal,
              recipesForOptionSelection: $recipesForOptionSelection
            )
          }
          .padding()
          .frame(maxWidth: isRegularWidth ? 800 : .infinity)
          .frame(maxWidth: .infinity)
        }
        .scrollDismissesKeyboard(.interactively)
        .sheet(isPresented: $showOptionSelectionModal) {
          RecipeOptionSelectionView(
            isPresented: $showOptionSelectionModal,
            recipes: recipesForOptionSelection,
            onSave: { ingredientSelections in
              viewModel?.setOptionSelections(ingredientSelections: ingredientSelections)
              // Continue with meal plan creation
              Task {
                if let viewModel = viewModel {
                  let success = await viewModel.createMealPlan()
                  if success {
                    NotificationCenter.default.post(name: .mealPlanCreated, object: nil)
                    dismiss()
                  }
                }
              }
            }
          )
        }
      } else {
        ProgressView("Initializing...")
          .frame(maxWidth: .infinity, maxHeight: .infinity)
      }
    }
    .navigationTitle("Create Meal Plan")
    .navigationBarTitleDisplayMode(.large)
    .onAppear {
      if viewModel == nil {
        viewModel = CreateMealPlanViewModel(authManager: authManager)
      }
      // Ensure no field is focused on view load to prevent keyboard from appearing
      focusedField = nil
    }
  }

  // MARK: - Event Card

  func eventCard(event: MealPlanEvent, viewModel: CreateMealPlanViewModel) -> some View {
    VStack(alignment: .leading, spacing: 16) {
      // Event Header
      HStack {
        Text("Event")
          .font(.headline)
        Spacer()
        if viewModel.events.count > 1 {
          Button(
            action: {
              viewModel.removeEvent(event)
            },
            label: {
              Image(systemName: "trash")
                .foregroundColor(.red)
            })
        }
      }

      if isRegularWidth {
        // Horizontal layout for iPad
        HStack(alignment: .top, spacing: 20) {
          // Meal Type
          VStack(alignment: .leading, spacing: 8) {
            Text("Meal Type")
              .font(.subheadline)
              .foregroundColor(.secondary)
            Picker(
              "Meal Type",
              selection: Binding(
                get: { event.mealType },
                set: { newValue in
                  viewModel.updateEventMealType(event.id, mealType: newValue)
                }
              )
            ) {
              Text("Breakfast").tag(Mealplanning_MealPlanEventName.breakfast)
              Text("Second Breakfast").tag(Mealplanning_MealPlanEventName.secondBreakfast)
              Text("Brunch").tag(Mealplanning_MealPlanEventName.brunch)
              Text("Lunch").tag(Mealplanning_MealPlanEventName.lunch)
              Text("Supper").tag(Mealplanning_MealPlanEventName.supper)
              Text("Dinner").tag(Mealplanning_MealPlanEventName.dinner)
            }
            .pickerStyle(.menu)
          }
          .frame(maxWidth: .infinity)

          // Start Time
          VStack(alignment: .leading, spacing: 8) {
            Text("Start Time")
              .font(.subheadline)
              .foregroundColor(.secondary)
            DatePicker(
              "Start Time",
              selection: Binding(
                get: { event.startDate },
                set: { newValue in
                  viewModel.updateEventStartDate(event.id, date: newValue)
                }
              ),
              displayedComponents: [.date, .hourAndMinute]
            )
            .datePickerStyle(.compact)
            .labelsHidden()
          }
          .frame(maxWidth: .infinity)

          // End Time
          VStack(alignment: .leading, spacing: 8) {
            Text("End Time")
              .font(.subheadline)
              .foregroundColor(.secondary)
            DatePicker(
              "End Time",
              selection: Binding(
                get: { event.endDate },
                set: { newValue in
                  viewModel.updateEventEndDate(event.id, date: newValue)
                }
              ),
              displayedComponents: [.date, .hourAndMinute]
            )
            .datePickerStyle(.compact)
            .labelsHidden()
          }
          .frame(maxWidth: .infinity)
        }
      } else {
        // Vertical layout for iPhone
        // Meal Type
        VStack(alignment: .leading, spacing: 8) {
          Text("Meal Type")
            .font(.subheadline)
            .foregroundColor(.secondary)
          Picker(
            "Meal Type",
            selection: Binding(
              get: { event.mealType },
              set: { newValue in
                viewModel.updateEventMealType(event.id, mealType: newValue)
              }
            )
          ) {
            Text("Breakfast").tag(Mealplanning_MealPlanEventName.breakfast)
            Text("Second Breakfast").tag(Mealplanning_MealPlanEventName.secondBreakfast)
            Text("Brunch").tag(Mealplanning_MealPlanEventName.brunch)
            Text("Lunch").tag(Mealplanning_MealPlanEventName.lunch)
            Text("Supper").tag(Mealplanning_MealPlanEventName.supper)
            Text("Dinner").tag(Mealplanning_MealPlanEventName.dinner)
          }
          .pickerStyle(.menu)
        }

        // Date and Time
        VStack(alignment: .leading, spacing: 12) {
          Text("Start Time")
            .font(.subheadline)
            .foregroundColor(.secondary)
          DatePicker(
            "Start Time",
            selection: Binding(
              get: { event.startDate },
              set: { newValue in
                viewModel.updateEventStartDate(event.id, date: newValue)
              }
            ),
            displayedComponents: [.date, .hourAndMinute]
          )
          .datePickerStyle(.compact)
        }

        VStack(alignment: .leading, spacing: 12) {
          Text("End Time")
            .font(.subheadline)
            .foregroundColor(.secondary)
          DatePicker(
            "End Time",
            selection: Binding(
              get: { event.endDate },
              set: { newValue in
                viewModel.updateEventEndDate(event.id, date: newValue)
              }
            ),
            displayedComponents: [.date, .hourAndMinute]
          )
          .datePickerStyle(.compact)
        }
      }

      // Search Section
      searchSection(event: event, viewModel: viewModel)

      // Selected Meals
      if !event.selectedMeals.isEmpty {
        selectedMealsSection(event: event, viewModel: viewModel)
      }

      // Search Results (filtered to exclude selected meals)
      let filteredResults = viewModel.filteredSearchResults(for: event)
      if !filteredResults.isEmpty {
        searchResultsSection(event: event, filteredResults: filteredResults, viewModel: viewModel)
      }
    }
    .padding()
    .background(Color(.systemBackground))
    .cornerRadius(10)
    .overlay(
      RoundedRectangle(cornerRadius: 10)
        .stroke(Color.blue.opacity(0.3), lineWidth: 1)
    )
  }

  // MARK: - Search Section (per event)

  func searchSection(event: MealPlanEvent, viewModel: CreateMealPlanViewModel) -> some View {
    _ = Bindable(viewModel)
    guard let eventIndex = viewModel.events.firstIndex(where: { $0.id == event.id }) else {
      return AnyView(EmptyView())
    }

    return AnyView(
      VStack(alignment: .leading, spacing: 12) {
        Text("Search for Meals")
          .font(.headline)

        HStack {
          TextField(
            "Search meals...",
            text: Binding(
              get: { event.searchQuery },
              set: { newValue in
                viewModel.updateEventSearchQuery(event.id, query: newValue)
              }
            )
          )
          .textFieldStyle(.roundedBorder)
          .autocorrectionDisabled()
          .textInputAutocapitalization(.never)
          .submitLabel(.search)
          .focused($focusedField, equals: .searchQuery(event.id))
          .onSubmit {
            Task {
              await viewModel.searchForMeals(for: event)
            }
          }

          if viewModel.events[eventIndex].isSearching {
            ProgressView()
              .padding(.leading, 8)
          } else {
            Button(
              action: {
                Task {
                  await viewModel.searchForMeals(for: event)
                }
              },
              label: {
                Image(systemName: "magnifyingglass")
                  .padding(8)
                  .background(Color.blue)
                  .foregroundColor(.white)
                  .cornerRadius(8)
              })
          }
        }

        if let searchError = viewModel.events[eventIndex].searchError {
          Text(searchError)
            .font(.caption)
            .foregroundColor(.red)
        }
      }
    )
  }

  // MARK: - Selected Meals Section (per event)

  func selectedMealsSection(event: MealPlanEvent, viewModel: CreateMealPlanViewModel) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      Text("Selected Meals (\(event.selectedMeals.count))")
        .font(.headline)

      ForEach(event.selectedMeals, id: \.id) { meal in
        SelectedMealCard(
          meal: meal,
          scale: viewModel.getMealScale(meal, in: event),
          onScaleChange: { newScale in
            viewModel.setMealScale(meal, scale: newScale, in: event)
          },
          onRemove: {
            viewModel.removeSelectedMeal(meal, from: event)
          },
          isRegularWidth: isRegularWidth
        )
      }
    }
  }

  // MARK: - Search Results Section (per event)

  func searchResultsSection(
    event: MealPlanEvent, filteredResults: [Mealplanning_Meal], viewModel: CreateMealPlanViewModel
  ) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      Text("Search Results (\(filteredResults.count))")
        .font(.headline)

      ForEach(filteredResults, id: \.id) { meal in
        MealSearchResultCard(
          meal: meal,
          isSelected: viewModel.isMealSelected(meal, in: event),
          onTap: {
            viewModel.toggleMealSelection(meal, in: event)
          }
        )
      }
    }
  }

}

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "John Doe"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  return CreateMealPlanView()
    .environment(authManager)
}
