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

  init(recipeID: String) {
    self.recipeID = recipeID
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
              viewModel: viewModel
            )
            .navigationTitle("Perform Recipe")
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
