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
  @Environment(\.dismiss) private var dismiss
  @State private var viewModel: PerformRecipeViewModel?
  
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
            ScrollView {
              VStack(alignment: .leading, spacing: 16) {
                // Recipe header
                recipeHeader(recipe: recipe, viewModel: viewModel)
                
                // Steps list
                stepsList(recipe: recipe, viewModel: viewModel)
              }
              .padding()
            }
            .navigationTitle("Perform Recipe")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
              ToolbarItem(placement: .navigationBarLeading) {
                Button("Done") {
                  dismiss()
                }
              }
            }
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
  
  // MARK: - Recipe Header
  
  private func recipeHeader(recipe: Mealplanning_Recipe, viewModel: PerformRecipeViewModel) -> some View {
    VStack(alignment: .leading, spacing: 8) {
      Text(recipe.name.isEmpty ? "Unnamed Recipe" : recipe.name)
        .font(.title)
        .fontWeight(.bold)
      
      if !recipe.description_p.isEmpty {
        Text(recipe.description_p)
          .font(.subheadline)
          .foregroundColor(.secondary)
      }
      
      // Progress indicator
      let completedCount = viewModel.completedSteps.count
      let totalSteps = recipe.steps.count
      Text("\(completedCount) of \(totalSteps) steps completed")
        .font(.caption)
        .foregroundColor(.secondary)
        .padding(.top, 4)
    }
    .padding()
    .frame(maxWidth: .infinity, alignment: .leading)
    .background(Color(.systemGray6))
    .cornerRadius(12)
  }
  
  // MARK: - Steps List
  
  private func stepsList(recipe: Mealplanning_Recipe, viewModel: PerformRecipeViewModel) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      Text("Steps")
        .font(.headline)
        .padding(.horizontal, 4)
      
      ForEach(Array(recipe.steps.enumerated()), id: \.element.id) { index, step in
        stepCard(step: step, index: index, viewModel: viewModel)
      }
    }
  }
  
  // MARK: - Step Card
  
  private func stepCard(
    step: Mealplanning_RecipeStep, index: Int, viewModel: PerformRecipeViewModel
  ) -> some View {
    let isCompleted = viewModel.isStepCompleted(index)
    let canCheck = viewModel.canCheckStep(index)
    let prerequisites = viewModel.getPrerequisiteStepIndices(index)
    let hasPrerequisites = !prerequisites.isEmpty
    let allPrerequisitesCompleted = prerequisites.allSatisfy { viewModel.isStepCompleted($0) }
    
    return VStack(alignment: .leading, spacing: 12) {
      // Step header with checkbox
      HStack(alignment: .top, spacing: 12) {
        // Checkbox
        Button(action: {
          viewModel.toggleStep(index)
        }) {
          Image(systemName: isCompleted ? "checkmark.circle.fill" : "circle")
            .font(.title2)
            .foregroundColor(
              canCheck ? (isCompleted ? .green : .blue) : .gray
            )
        }
        .disabled(!canCheck)
        
        // Step number and instructions
        VStack(alignment: .leading, spacing: 4) {
          HStack {
            Text("Step \(Int(step.index) + 1)")
              .font(.headline)
              .foregroundColor(isCompleted ? .secondary : .primary)
            
            if step.optional {
              Text("(Optional)")
                .font(.caption)
                .foregroundColor(.secondary)
            }
          }
          
          if !step.explicitInstructions.isEmpty {
            Text(step.explicitInstructions)
              .font(.body)
              .foregroundColor(isCompleted ? .secondary : .primary)
              .strikethrough(isCompleted)
          }
          
          // Prerequisites warning
          if hasPrerequisites && !allPrerequisitesCompleted {
            HStack(spacing: 4) {
              Image(systemName: "exclamationmark.triangle.fill")
                .font(.caption)
                .foregroundColor(.orange)
              Text("Complete steps \(prerequisites.map { String($0 + 1) }.joined(separator: ", ")) first")
                .font(.caption)
                .foregroundColor(.orange)
            }
            .padding(.top, 4)
          }
        }
        
        Spacer()
      }
      
      // Step details (ingredients, instruments, vessels)
      if !isCompleted || true {  // Show details even when completed
        stepDetails(step: step, viewModel: viewModel, stepIndex: index)
      }
    }
    .padding()
    .background(
      isCompleted ? Color(.systemGray6) : Color(.systemBackground)
    )
    .cornerRadius(12)
    .overlay(
      RoundedRectangle(cornerRadius: 12)
        .stroke(
          isCompleted ? Color.green.opacity(0.3) : Color.clear,
          lineWidth: 2
        )
    )
  }
  
  // MARK: - Step Details
  
  private func stepDetails(
    step: Mealplanning_RecipeStep, viewModel: PerformRecipeViewModel, stepIndex: Int
  ) -> some View {
    VStack(alignment: .leading, spacing: 8) {
      // Ingredients
      if !step.ingredients.isEmpty {
          stepItemsSection(
          title: "Ingredients",
          items: step.ingredients.map { ingredient in
            let isProduct = ingredient.hasRecipeStepProductID
            let productID = isProduct ? ingredient.recipeStepProductID : nil
            let prerequisiteStepIndex = productID.flatMap { viewModel.getStepIndexForProductID($0) }
            let prerequisiteCompleted = prerequisiteStepIndex.map { viewModel.isStepCompleted($0) } ?? true
            
            return StepItem(
              name: ingredient.name,
              isProduct: isProduct,
              prerequisiteStepIndex: prerequisiteStepIndex,
              prerequisiteCompleted: prerequisiteCompleted
            )
          }
        )
      }
      
      // Instruments
      if !step.instruments.isEmpty {
          stepItemsSection(
          title: "Instruments",
          items: step.instruments.map { instrument in
            let isProduct = instrument.hasRecipeStepProductID
            let productID = isProduct ? instrument.recipeStepProductID : nil
            let prerequisiteStepIndex = productID.flatMap { viewModel.getStepIndexForProductID($0) }
            let prerequisiteCompleted = prerequisiteStepIndex.map { viewModel.isStepCompleted($0) } ?? true
            
            return StepItem(
              name: instrument.name,
              isProduct: isProduct,
              prerequisiteStepIndex: prerequisiteStepIndex,
              prerequisiteCompleted: prerequisiteCompleted
            )
          }
        )
      }
      
      // Vessels
      if !step.vessels.isEmpty {
          stepItemsSection(
          title: "Vessels",
          items: step.vessels.map { vessel in
            let isProduct = vessel.hasRecipeStepProductID
            let productID = isProduct ? vessel.recipeStepProductID : nil
            let prerequisiteStepIndex = productID.flatMap { viewModel.getStepIndexForProductID($0) }
            let prerequisiteCompleted = prerequisiteStepIndex.map { viewModel.isStepCompleted($0) } ?? true
            
            return StepItem(
              name: vessel.name,
              isProduct: isProduct,
              prerequisiteStepIndex: prerequisiteStepIndex,
              prerequisiteCompleted: prerequisiteCompleted
            )
          }
        )
      }
      
      // Notes
      if !step.notes.isEmpty {
        Text(step.notes)
          .font(.caption)
          .foregroundColor(.secondary)
          .italic()
          .padding(.top, 4)
      }
    }
    .padding(.leading, 44)  // Align with step content
  }
  
  // MARK: - Step Items Section
  
  private func stepItemsSection(title: String, items: [StepItem]) -> some View {
    VStack(alignment: .leading, spacing: 4) {
      Text(title)
        .font(.subheadline)
        .fontWeight(.semibold)
        .foregroundColor(.secondary)
      
      ForEach(Array(items.enumerated()), id: \.offset) { _, item in
        HStack(spacing: 6) {
          if item.isProduct && !item.prerequisiteCompleted {
            Image(systemName: "clock.fill")
              .font(.caption)
              .foregroundColor(.orange)
          }
          Text(item.name)
            .font(.caption)
            .foregroundColor(
              (item.isProduct && !item.prerequisiteCompleted) ? .orange : .secondary
            )
          if let prerequisiteStepIndex = item.prerequisiteStepIndex {
            Text("(from step \(prerequisiteStepIndex + 1))")
              .font(.caption2)
              .foregroundColor(.secondary)
          }
        }
      }
    }
  }
}

// MARK: - Helper Types

private struct StepItem {
  let name: String
  let isProduct: Bool
  let prerequisiteStepIndex: Int?
  let prerequisiteCompleted: Bool
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

