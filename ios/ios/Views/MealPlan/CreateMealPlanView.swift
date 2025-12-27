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

struct CreateMealPlanView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Environment(\.dismiss) private var dismiss
  @Environment(\.horizontalSizeClass) private var horizontalSizeClass
  @State private var viewModel: CreateMealPlanViewModel?
  @FocusState private var focusedField: Field?
  
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
            createButton(viewModel: viewModel)
          }
          .padding()
          .frame(maxWidth: isRegularWidth ? 800 : .infinity)
          .frame(maxWidth: .infinity)
        }
        .scrollDismissesKeyboard(.interactively)
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

  // MARK: - Meal Plan Details Section

  private func mealPlanDetailsSection(viewModel: CreateMealPlanViewModel) -> some View {
    let bindableViewModel = Bindable(viewModel)
    
    return VStack(alignment: .leading, spacing: 16) {
      Text("Meal Plan Details")
        .font(.title2)
        .fontWeight(.bold)

      VStack(alignment: .leading, spacing: 12) {
        Text("Name")
          .font(.headline)
        TextField("Enter meal plan name", text: bindableViewModel.mealPlanName)
          .textFieldStyle(.roundedBorder)
          .focused($focusedField, equals: .mealPlanName)
      }

      VStack(alignment: .leading, spacing: 12) {
        Text("Voting Deadline")
          .font(.headline)
        DatePicker(
          "",
          selection: bindableViewModel.votingDeadline,
          displayedComponents: [.date, .hourAndMinute]
        )
        .datePickerStyle(.compact)
        .labelsHidden()
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }
  
  // MARK: - Meal Plan Details Section (Horizontal for iPad)
  
  private func mealPlanDetailsSectionHorizontal(viewModel: CreateMealPlanViewModel) -> some View {
    let bindableViewModel = Bindable(viewModel)
    
    return VStack(alignment: .leading, spacing: 16) {
      Text("Meal Plan Details")
        .font(.title2)
        .fontWeight(.bold)

      HStack(alignment: .top, spacing: 24) {
        VStack(alignment: .leading, spacing: 12) {
          Text("Name")
            .font(.headline)
          TextField("Enter meal plan name", text: bindableViewModel.mealPlanName)
            .textFieldStyle(.roundedBorder)
            .focused($focusedField, equals: .mealPlanName)
        }
        .frame(maxWidth: .infinity)

        VStack(alignment: .leading, spacing: 12) {
          Text("Voting Deadline")
            .font(.headline)
          DatePicker(
            "",
            selection: bindableViewModel.votingDeadline,
            displayedComponents: [.date, .hourAndMinute]
          )
          .datePickerStyle(.compact)
          .labelsHidden()
          .frame(maxWidth: .infinity, alignment: .leading)
        }
        .frame(maxWidth: .infinity)
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }

  // MARK: - Events Section

  private func eventsSection(viewModel: CreateMealPlanViewModel) -> some View {
    VStack(alignment: .leading, spacing: 16) {
      HStack {
        Text("Events")
          .font(.title2)
          .fontWeight(.bold)
        Spacer()
        Button(action: {
          viewModel.addEvent()
        }) {
          HStack {
            Image(systemName: "plus.circle.fill")
            Text("Add Event")
          }
          .font(.subheadline)
          .foregroundColor(.blue)
        }
      }

      // Always use vertical stack for events
      ForEach(viewModel.events) { event in
        eventCard(event: event, viewModel: viewModel)
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }

  // MARK: - Event Card

  private func eventCard(event: MealPlanEvent, viewModel: CreateMealPlanViewModel) -> some View {
    VStack(alignment: .leading, spacing: 16) {
      // Event Header
      HStack {
        Text("Event")
          .font(.headline)
        Spacer()
        if viewModel.events.count > 1 {
          Button(action: {
            viewModel.removeEvent(event)
          }) {
            Image(systemName: "trash")
              .foregroundColor(.red)
          }
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
            Picker("Meal Type", selection: Binding(
              get: { event.mealType },
              set: { newValue in
                viewModel.updateEventMealType(event.id, mealType: newValue)
              }
            )) {
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
          Picker("Meal Type", selection: Binding(
            get: { event.mealType },
            set: { newValue in
              viewModel.updateEventMealType(event.id, mealType: newValue)
            }
          )) {
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

  private func searchSection(event: MealPlanEvent, viewModel: CreateMealPlanViewModel) -> some View {
      _ = Bindable(viewModel)
    guard let eventIndex = viewModel.events.firstIndex(where: { $0.id == event.id }) else {
      return AnyView(EmptyView())
    }
    
    return AnyView(
      VStack(alignment: .leading, spacing: 12) {
        Text("Search for Meals")
          .font(.headline)

        HStack {
          TextField("Search meals...", text: Binding(
            get: { event.searchQuery },
            set: { newValue in
              viewModel.updateEventSearchQuery(event.id, query: newValue)
            }
          ))
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
            Button(action: {
              Task {
                await viewModel.searchForMeals(for: event)
              }
            }) {
              Image(systemName: "magnifyingglass")
                .padding(8)
                .background(Color.blue)
                .foregroundColor(.white)
                .cornerRadius(8)
            }
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

  private func selectedMealsSection(event: MealPlanEvent, viewModel: CreateMealPlanViewModel) -> some View {
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

  private func searchResultsSection(event: MealPlanEvent, filteredResults: [Mealplanning_Meal], viewModel: CreateMealPlanViewModel) -> some View {
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

  // MARK: - Error Message

  private func errorMessage(_ message: String) -> some View {
    HStack {
      Image(systemName: "exclamationmark.triangle")
        .foregroundColor(.red)
      Text(message)
        .font(.subheadline)
        .foregroundColor(.red)
    }
    .padding()
    .background(Color.red.opacity(0.1))
    .cornerRadius(8)
  }

  // MARK: - Create Button

  private func createButton(viewModel: CreateMealPlanViewModel) -> some View {
      _ = Bindable(viewModel)
    let hasSelectedMeals = viewModel.events.contains { !$0.selectedMeals.isEmpty }
    
    return Button(action: {
      Task {
        let success = await viewModel.createMealPlan()
        if success {
          // Post notification to refresh home view
          NotificationCenter.default.post(name: .mealPlanCreated, object: nil)
          dismiss()
        }
      }
    }) {
      HStack {
        if viewModel.isCreating {
          ProgressView()
            .progressViewStyle(CircularProgressViewStyle(tint: .white))
        }
        Text(viewModel.isCreating ? "Creating..." : "Create Meal Plan")
          .fontWeight(.semibold)
      }
      .frame(maxWidth: .infinity)
      .padding()
      .background(
        viewModel.isCreating || !hasSelectedMeals
          ? Color.gray : Color.blue
      )
      .foregroundColor(.white)
      .cornerRadius(10)
    }
    .disabled(viewModel.isCreating || !hasSelectedMeals)
  }
}

// MARK: - Meal Search Result Card

struct MealSearchResultCard: View {
  let meal: Mealplanning_Meal
  let isSelected: Bool
  let onTap: () -> Void

  var body: some View {
    HStack {
      VStack(alignment: .leading, spacing: 8) {
        Text(meal.name.isEmpty ? "Unnamed Meal" : meal.name)
          .font(.headline)
          .foregroundColor(.primary)

        if !meal.description_p.isEmpty {
          Text(meal.description_p)
            .font(.subheadline)
            .foregroundColor(.secondary)
            .lineLimit(2)
        }

        // Show recipe names from components
        if !meal.components.isEmpty {
          let recipeNames = meal.components.compactMap { component -> String? in
            component.recipe.name.isEmpty ? nil : component.recipe.name
          }
          if !recipeNames.isEmpty {
            Text(recipeNames.joined(separator: ", "))
              .font(.caption)
              .foregroundColor(.secondary)
              .lineLimit(1)
          }
        }
      }

      Spacer()

      Image(systemName: isSelected ? "checkmark.circle.fill" : "circle")
        .foregroundColor(isSelected ? .blue : .gray)
        .font(.title2)
    }
    .padding()
    .background(isSelected ? Color.blue.opacity(0.1) : Color(.systemBackground))
    .cornerRadius(8)
    .overlay(
      RoundedRectangle(cornerRadius: 8)
        .stroke(isSelected ? Color.blue : Color.clear, lineWidth: 2)
    )
    .contentShape(Rectangle())
    .onTapGesture {
      onTap()
    }
  }
}

// MARK: - Selected Meal Card

struct SelectedMealCard: View {
  let meal: Mealplanning_Meal
  let scale: Float
  let onScaleChange: (Float) -> Void
  let onRemove: () -> Void
  let isRegularWidth: Bool
  
  @State private var scaleText: String = ""
  @FocusState private var isScaleFocused: Bool

  var body: some View {
    if isRegularWidth {
      // Horizontal layout for iPad
      HStack(alignment: .top, spacing: 16) {
        VStack(alignment: .leading, spacing: 4) {
          Text(meal.name.isEmpty ? "Unnamed Meal" : meal.name)
            .font(.headline)

          if !meal.components.isEmpty {
            let recipeNames = meal.components.compactMap { component -> String? in
              component.recipe.name.isEmpty ? nil : component.recipe.name
            }
            if !recipeNames.isEmpty {
              Text(recipeNames.joined(separator: ", "))
                .font(.caption)
                .foregroundColor(.secondary)
            }
          }
        }
        .frame(maxWidth: .infinity, alignment: .leading)

        // Scale control
        VStack(alignment: .leading, spacing: 4) {
          Text("Scale")
            .font(.caption)
            .foregroundColor(.secondary)
          HStack(spacing: 8) {
            TextField("1.0", text: $scaleText)
              .keyboardType(.decimalPad)
              .textFieldStyle(.roundedBorder)
              .frame(width: 100)
              .focused($isScaleFocused)
              .onSubmit {
                updateScaleFromText()
              }
              .onChange(of: isScaleFocused) { wasFocused, isFocused in
                if !isFocused {
                  updateScaleFromText()
                }
              }
            Text("x")
              .font(.subheadline)
              .foregroundColor(.secondary)
          }
        }
        .frame(maxWidth: 200)

        Button(action: onRemove) {
          Image(systemName: "xmark.circle.fill")
            .foregroundColor(.red)
            .font(.title3)
        }
      }
      .padding()
      .background(Color(.systemBackground))
      .cornerRadius(8)
      .onAppear {
        scaleText = String(format: "%.2f", scale)
      }
      .onChange(of: scale) { oldValue, newValue in
        if !isScaleFocused {
          scaleText = String(format: "%.2f", newValue)
        }
      }
    } else {
      // Vertical layout for iPhone
      VStack(alignment: .leading, spacing: 12) {
        HStack {
          VStack(alignment: .leading, spacing: 4) {
            Text(meal.name.isEmpty ? "Unnamed Meal" : meal.name)
              .font(.headline)

            if !meal.components.isEmpty {
              let recipeNames = meal.components.compactMap { component -> String? in
                component.recipe.name.isEmpty ? nil : component.recipe.name
              }
              if !recipeNames.isEmpty {
                Text(recipeNames.joined(separator: ", "))
                  .font(.caption)
                  .foregroundColor(.secondary)
              }
            }
          }

          Spacer()

          Button(action: onRemove) {
            Image(systemName: "xmark.circle.fill")
              .foregroundColor(.red)
              .font(.title3)
          }
        }

        // Scale control
        VStack(alignment: .leading, spacing: 4) {
          Text("Scale")
            .font(.caption)
            .foregroundColor(.secondary)
          HStack(spacing: 8) {
            TextField("1.0", text: $scaleText)
              .keyboardType(.decimalPad)
              .textFieldStyle(.roundedBorder)
              .focused($isScaleFocused)
              .onSubmit {
                updateScaleFromText()
              }
              .onChange(of: isScaleFocused) { wasFocused, isFocused in
                if !isFocused {
                  updateScaleFromText()
                }
              }
            Text("x")
              .font(.subheadline)
              .foregroundColor(.secondary)
          }
        }
      }
      .padding()
      .background(Color(.systemBackground))
      .cornerRadius(8)
      .onAppear {
        scaleText = String(format: "%.2f", scale)
      }
      .onChange(of: scale) { oldValue, newValue in
        if !isScaleFocused {
          scaleText = String(format: "%.2f", newValue)
        }
      }
    }
  }
  
  private func updateScaleFromText() {
    // Parse the text input and validate it's a positive number
    if let parsedValue = Float(scaleText.trimmingCharacters(in: .whitespacesAndNewlines)) {
      if parsedValue > 0 {
        onScaleChange(parsedValue)
        scaleText = String(format: "%.2f", parsedValue)
      } else {
        // Invalid: not positive, reset to current scale
        scaleText = String(format: "%.2f", scale)
      }
    } else {
      // Invalid: not a number, reset to current scale
      scaleText = String(format: "%.2f", scale)
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
