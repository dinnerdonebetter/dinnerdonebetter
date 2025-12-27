//
//  CreateMealPlanView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct CreateMealPlanView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Environment(\.dismiss) private var dismiss
  @State private var viewModel: CreateMealPlanViewModel?
  @FocusState private var focusedField: Field?
  
  enum Field {
    case mealPlanName
    case searchQuery
  }

  var body: some View {
    Group {
      if let viewModel = viewModel {
        ScrollView(.vertical, showsIndicators: true) {
          VStack(spacing: 24) {
            // Meal Plan Details Section
            mealPlanDetailsSection(viewModel: viewModel)

            // Search Section
            searchSection(viewModel: viewModel)

            // Selected Meals Section
            if !viewModel.selectedMeals.isEmpty {
              selectedMealsSection(viewModel: viewModel)
            }

            // Search Results Section
            if !viewModel.searchResults.isEmpty {
              searchResultsSection(viewModel: viewModel)
            }

            // Error Messages
            if let error = viewModel.creationError {
              errorMessage(error)
            }

            // Create Button
            createButton(viewModel: viewModel)
          }
          .padding()
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
          "Voting Deadline",
          selection: bindableViewModel.votingDeadline,
          displayedComponents: [.date, .hourAndMinute]
        )
        .datePickerStyle(.compact)
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }

  // MARK: - Search Section

  private func searchSection(viewModel: CreateMealPlanViewModel) -> some View {
    let bindableViewModel = Bindable(viewModel)
    
    return VStack(alignment: .leading, spacing: 16) {
      Text("Search for Meals")
        .font(.title2)
        .fontWeight(.bold)

      HStack {
        TextField("Search meals...", text: bindableViewModel.searchQuery)
          .textFieldStyle(.roundedBorder)
          .autocorrectionDisabled()
          .textInputAutocapitalization(.never)
          .submitLabel(.search)
          .focused($focusedField, equals: .searchQuery)
          .onSubmit {
            Task {
              await viewModel.searchForMeals()
            }
          }

        if viewModel.isSearching {
          ProgressView()
            .padding(.leading, 8)
        } else {
          Button(action: {
            Task {
              await viewModel.searchForMeals()
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

      if let searchError = viewModel.searchError {
        Text(searchError)
          .font(.caption)
          .foregroundColor(.red)
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }

  // MARK: - Selected Meals Section

  private func selectedMealsSection(viewModel: CreateMealPlanViewModel) -> some View {
    VStack(alignment: .leading, spacing: 16) {
      Text("Selected Meals (\(viewModel.selectedMeals.count))")
        .font(.title2)
        .fontWeight(.bold)

      ForEach(viewModel.selectedMeals, id: \.id) { meal in
        SelectedMealCard(
          meal: meal,
          onRemove: {
            viewModel.removeSelectedMeal(meal)
          }
        )
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }

  // MARK: - Search Results Section

  private func searchResultsSection(viewModel: CreateMealPlanViewModel) -> some View {
    VStack(alignment: .leading, spacing: 16) {
      Text("Search Results (\(viewModel.searchResults.count))")
        .font(.title2)
        .fontWeight(.bold)

      ForEach(viewModel.searchResults, id: \.id) { meal in
        MealSearchResultCard(
          meal: meal,
          isSelected: viewModel.isMealSelected(meal),
          onTap: {
            viewModel.toggleMealSelection(meal)
          }
        )
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
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
    Button(action: {
      Task {
        let success = await viewModel.createMealPlan()
        if success {
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
        viewModel.isCreating || viewModel.selectedMeals.isEmpty
          ? Color.gray : Color.blue
      )
      .foregroundColor(.white)
      .cornerRadius(10)
    }
    .disabled(viewModel.isCreating || viewModel.selectedMeals.isEmpty)
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
  let onRemove: () -> Void

  var body: some View {
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
    .padding()
    .background(Color(.systemBackground))
    .cornerRadius(8)
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

