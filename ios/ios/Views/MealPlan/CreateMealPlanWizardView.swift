//
//  CreateMealPlanWizardView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

extension Notification.Name {
  static let mealPlanCreated = Notification.Name("mealPlanCreated")
  static let mealPlanArchived = Notification.Name("mealPlanArchived")
}

struct CreateMealPlanWizardView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Environment(\.dismiss) var dismiss
  @State private var viewModel: CreateMealPlanViewModel?
  @State private var recipesForOptionSelection: [Mealplanning_Recipe]?

  var acceptedOccupiedDates: Set<Date> = []
  var proposedOccupiedDates: Set<Date> = []

  var body: some View {
    Group {
      if let viewModel = viewModel {
        VStack(spacing: 0) {
          stepIndicator(currentStep: viewModel.wizardStep.rawValue, totalSteps: 2)
            .padding()

          ScrollView {
            VStack(spacing: 24) {
              switch viewModel.wizardStep {
              case .weekSelection:
                WeekSelectionStepView(viewModel: viewModel)

              case .mealAssignment:
                MealAssignmentStepView(
                  viewModel: viewModel,
                  recipesForOptionSelection: $recipesForOptionSelection,
                  onDismiss: { dismiss() }
                )
              }
            }
            .padding()
          }

          if let error = viewModel.creationError {
            HStack {
              Image(systemName: "exclamationmark.triangle")
                .foregroundColor(.red)
              Text(error)
                .font(.subheadline)
                .foregroundColor(.red)
            }
            .padding(.horizontal)
          }
        }
        .sheet(
          isPresented: Binding(
            get: { recipesForOptionSelection != nil },
            set: { if !$0 { recipesForOptionSelection = nil } }
          )
        ) {
          if let recipes = recipesForOptionSelection {
            RecipeOptionSelectionView(
              isPresented: Binding(
                get: { recipesForOptionSelection != nil },
                set: { if !$0 { recipesForOptionSelection = nil } }
              ),
              recipes: recipes,
              onSave: { ingredientSelections in
                var finalSelections = ingredientSelections
                if finalSelections.isEmpty {
                  for recipe in recipes {
                    let defaults = viewModel.getDefaultOptionSelections(for: recipe)
                    if !defaults.isEmpty {
                      finalSelections[recipe.id] = defaults
                    }
                  }
                }
                viewModel.setOptionSelections(ingredientSelections: finalSelections)
                Task {
                  let success = await viewModel.createMealPlan()
                  if success {
                    NotificationCenter.default.post(name: .mealPlanCreated, object: nil)
                    recipesForOptionSelection = nil
                    dismiss()
                  }
                }
              }
            )
          }
        }
      } else {
        DSInitializingView()
      }
    }
    .navigationTitle("Plan Dinners")
    .navigationBarTitleDisplayMode(.large)
    .onAppear {
      if viewModel == nil {
        viewModel = CreateMealPlanViewModel(
          authManager: authManager,
          acceptedOccupiedDates: acceptedOccupiedDates,
          proposedOccupiedDates: proposedOccupiedDates
        )
      }
    }
  }

  private func stepIndicator(currentStep: Int, totalSteps: Int) -> some View {
    HStack(spacing: 8) {
      ForEach(1...totalSteps, id: \.self) { step in
        Capsule()
          .fill(step <= currentStep ? Color.blue : Color(.systemGray5))
          .frame(height: 4)
          .frame(maxWidth: .infinity)
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

  return NavigationStack {
    CreateMealPlanWizardView()
      .environment(authManager)
  }
}
