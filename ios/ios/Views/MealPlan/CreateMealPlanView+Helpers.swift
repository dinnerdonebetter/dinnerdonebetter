//
//  CreateMealPlanView+Helpers.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

// MARK: - CreateMealPlanView Helpers

extension CreateMealPlanView {
  // MARK: - Meal Plan Details Section

  func mealPlanDetailsSection(viewModel: CreateMealPlanViewModel) -> some View {
    let bindableViewModel = Bindable(viewModel)
    let validationErrors = viewModel.getValidationErrors()

    return VStack(alignment: .leading, spacing: 16) {
      Text("Meal Plan Details")
        .font(.title2)
        .fontWeight(.bold)

      VStack(alignment: .leading, spacing: 12) {
        HStack {
          Text("Name")
            .font(.headline)
          if validationErrors.hasNameError {
            Image(systemName: "exclamationmark.circle.fill")
              .foregroundColor(.red)
              .font(.caption)
          }
        }
        TextField("Enter meal plan name", text: bindableViewModel.mealPlanName)
          .textFieldStyle(.roundedBorder)
          .focused($focusedField, equals: .mealPlanName)
          .overlay(
            RoundedRectangle(cornerRadius: 5)
              .stroke(validationErrors.hasNameError ? Color.red : Color.clear, lineWidth: 1)
          )
        if validationErrors.hasNameError {
          Text("Meal plan name is required")
            .font(.caption)
            .foregroundColor(.red)
        }
      }

      VStack(alignment: .leading, spacing: 12) {
        HStack {
          Text("Voting Deadline")
            .font(.headline)
          if validationErrors.hasVotingDeadlineError {
            Image(systemName: "exclamationmark.circle.fill")
              .foregroundColor(.red)
              .font(.caption)
          }
        }
        DatePicker(
          "",
          selection: bindableViewModel.votingDeadline,
          displayedComponents: [.date, .hourAndMinute]
        )
        .datePickerStyle(.compact)
        .labelsHidden()
        if validationErrors.hasVotingDeadlineError {
          Text("Voting deadline must be at least 12 hours from now")
            .font(.caption)
            .foregroundColor(.red)
        }
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }

  // MARK: - Meal Plan Details Section (Horizontal for iPad)

  func mealPlanDetailsSectionHorizontal(viewModel: CreateMealPlanViewModel) -> some View {
    let bindableViewModel = Bindable(viewModel)
    let validationErrors = viewModel.getValidationErrors()

    return VStack(alignment: .leading, spacing: 16) {
      Text("Meal Plan Details")
        .font(.title2)
        .fontWeight(.bold)

      HStack(alignment: .top, spacing: 24) {
        VStack(alignment: .leading, spacing: 12) {
          HStack {
            Text("Name")
              .font(.headline)
            if validationErrors.hasNameError {
              Image(systemName: "exclamationmark.circle.fill")
                .foregroundColor(.red)
                .font(.caption)
            }
          }
          TextField("Enter meal plan name", text: bindableViewModel.mealPlanName)
            .textFieldStyle(.roundedBorder)
            .focused($focusedField, equals: .mealPlanName)
            .overlay(
              RoundedRectangle(cornerRadius: 5)
                .stroke(validationErrors.hasNameError ? Color.red : Color.clear, lineWidth: 1)
            )
          if validationErrors.hasNameError {
            Text("Meal plan name is required")
              .font(.caption)
              .foregroundColor(.red)
          }
        }
        .frame(maxWidth: .infinity)

        VStack(alignment: .leading, spacing: 12) {
          HStack {
            Text("Voting Deadline")
              .font(.headline)
            if validationErrors.hasVotingDeadlineError {
              Image(systemName: "exclamationmark.circle.fill")
                .foregroundColor(.red)
                .font(.caption)
            }
          }
          DatePicker(
            "",
            selection: bindableViewModel.votingDeadline,
            displayedComponents: [.date, .hourAndMinute]
          )
          .datePickerStyle(.compact)
          .labelsHidden()
          .frame(maxWidth: .infinity, alignment: .leading)
          if validationErrors.hasVotingDeadlineError {
            Text("Voting deadline must be at least 12 hours from now")
              .font(.caption)
              .foregroundColor(.red)
          }
        }
        .frame(maxWidth: .infinity)
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }

  // MARK: - Events Section

  func eventsSection(viewModel: CreateMealPlanViewModel) -> some View {
    VStack(alignment: .leading, spacing: 16) {
      HStack {
        Text("Events")
          .font(.title2)
          .fontWeight(.bold)
        Spacer()
        Button(
          action: {
            viewModel.addEvent()
          },
          label: {
            HStack {
              Image(systemName: "plus.circle.fill")
              Text("Add Event")
            }
            .font(.subheadline)
            .foregroundColor(.blue)
          })
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

  // MARK: - Error Message

  func errorMessage(_ message: String) -> some View {
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

  func createButton(
    viewModel: CreateMealPlanViewModel,
    recipesForOptionSelection: Binding<[Mealplanning_Recipe]?>
  ) -> some View {
    _ = Bindable(viewModel)
    let canCreate = viewModel.canCreateMealPlan()

    return Button(
      action: {
        // Collect all selected meals from all events
        var allSelectedMeals: [Mealplanning_Meal] = []
        for event in viewModel.events {
          allSelectedMeals.append(contentsOf: event.selectedMeals)
        }

        // Use ViewModel method to check for recipes with options
        let recipesWithOptions = viewModel.collectRecipesWithOptions(from: allSelectedMeals)

        print("🔍 Create button: Found \(recipesWithOptions.count) recipe(s) with options")

        if !recipesWithOptions.isEmpty {
          // Get all unique recipes for the modal
          let allRecipes = viewModel.getAllRecipes(from: allSelectedMeals)
          print("🔍 Create button: getAllRecipes returned \(allRecipes.count) recipe(s)")
          for recipe in allRecipes {
            print("  📝 Recipe: \(recipe.name) (ID: \(recipe.id)), has \(recipe.steps.count) steps")
          }

          // Pass all recipes to the modal - it will filter internally to only show recipes with options
          // This is safer than pre-filtering in case there are ID mismatches
          print(
            "🔍 Create button: Passing \(allRecipes.count) recipe(s) to modal (modal will filter)")
          // Set recipes first, then the sheet will appear automatically via the binding
          recipesForOptionSelection.wrappedValue = allRecipes
        } else {
          print("⚠️ Create button: No recipes with options found, proceeding directly")
          // No options, proceed directly
          Task {
            let success = await viewModel.createMealPlan()
            if success {
              NotificationCenter.default.post(name: .mealPlanCreated, object: nil)
              dismiss()
            }
          }
        }
      },
      label: {
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
          viewModel.isCreating || !canCreate
            ? Color.gray : Color.blue
        )
        .foregroundColor(.white)
        .cornerRadius(10)
      }
    )
    .disabled(viewModel.isCreating || !canCreate)
  }
}
