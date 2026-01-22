//
//  PerformRecipeView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct PerformRecipeView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @State private var viewModel: PerformRecipeViewModel?
  @State private var isInstrumentsVesselsExpanded = false
  @State private var isIngredientsExpanded = false
  @State private var checkedIngredients: Set<String> = []
  // Track by ValidInstrument/ValidVessel ID
  @State private var checkedInstrumentsVessels: Set<String> = []

  let recipeID: String
  let highlightedStepIDs: Set<String>?
  let prepTaskContext: PrepTaskContext?

  struct PrepTaskContext {
    let prepTaskName: String?
    let recipeName: String?
    let eventName: String?
    let eventTime: Date?
  }

  init(recipeID: String, highlightedStepIDs: Set<String>? = nil, prepTaskContext: PrepTaskContext? = nil) {
    self.recipeID = recipeID
    self.highlightedStepIDs = highlightedStepIDs
    self.prepTaskContext = prepTaskContext
  }

  var body: some View {
    NavigationStack {
      Group {
        if let viewModel = viewModel {
          if viewModel.isLoading {
            ProgressView("Loading recipe...")
              .frame(maxWidth: .infinity, maxHeight: .infinity)
          } else if let errorMessage = viewModel.errorMessage {
            VStack(spacing: 16) {
              Image(systemName: "exclamationmark.triangle")
                .font(.largeTitle)
                .foregroundColor(.orange)
              Text("Error")
                .font(.headline)
              Text(errorMessage)
                .font(.subheadline)
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)
                .padding(.horizontal)
              Button("Retry") {
                Task {
                  await viewModel.loadRecipe()
                }
              }
              .buttonStyle(.borderedProminent)
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
          } else if let recipe = viewModel.recipe {
            RecipePerformanceContentView(
              checkedIngredients: $checkedIngredients,
              checkedInstrumentsVessels: $checkedInstrumentsVessels,
              isInstrumentsVesselsExpanded: $isInstrumentsVesselsExpanded,
              isIngredientsExpanded: $isIngredientsExpanded,
              recipe: recipe,
              viewModel: viewModel,
              highlightedStepIDs: highlightedStepIDs,
              prepTaskContext: prepTaskContext
            )
            .navigationTitle(navigationTitle)
            .navigationBarTitleDisplayMode(.inline)
          } else {
            ProgressView("Loading...")
              .frame(maxWidth: .infinity, maxHeight: .infinity)
          }
        } else {
          ProgressView("Initializing...")
            .frame(maxWidth: .infinity, maxHeight: .infinity)
        }
      }
      .onAppear {
        if viewModel == nil {
          viewModel = PerformRecipeViewModel(recipeID: recipeID, authManager: authManager)
          Task {
            await viewModel?.loadRecipe()
          }
        }
      }
    }
  }
  
  private var navigationTitle: String {
    if let context = prepTaskContext {
      var parts: [String] = []
      
      if let prepTaskName = context.prepTaskName, !prepTaskName.isEmpty {
        parts.append(prepTaskName)
      } else {
        parts.append("Prep Task")
      }
      
      if let recipeName = context.recipeName, !recipeName.isEmpty {
        parts.append("for \(recipeName)")
      }
      
      if let eventName = context.eventName, let eventTime = context.eventTime {
        let formatter = DateFormatter()
        formatter.dateStyle = .none
        formatter.timeStyle = .short
        parts.append("at \(formatter.string(from: eventTime))")
      }
      
      return parts.joined(separator: " ")
    }
    
    return "Perform Recipe"
  }
}

// MARK: - Preview

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "Test User"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  return PerformRecipeView(recipeID: "test-recipe")
    .environment(authManager)
}
