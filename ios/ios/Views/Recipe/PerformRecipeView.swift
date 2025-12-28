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
  @State private var isInstrumentsVesselsExpanded = false
  @State private var isIngredientsExpanded = false
  
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
                
                // Instruments & Vessels section
                instrumentsVesselsSection(recipe: recipe)
                
                // Ingredients section
                ingredientsSection(recipe: recipe)
                
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
  
  // MARK: - Instruments & Vessels Section
  
  private func instrumentsVesselsSection(recipe: Mealplanning_Recipe) -> some View {
    let allInstrumentsVessels = getAllInstrumentsAndVessels(from: recipe)
    
    return VStack(alignment: .leading, spacing: 0) {
      Button(action: {
        withAnimation {
          isInstrumentsVesselsExpanded.toggle()
        }
      }) {
        HStack {
          Text("Instruments & Vessels")
            .font(.headline)
            .foregroundColor(.primary)
          Spacer()
          Image(systemName: isInstrumentsVesselsExpanded ? "chevron.down" : "chevron.right")
            .font(.caption)
            .foregroundColor(.secondary)
        }
        .padding()
        .background(Color(.systemGray6))
      }
      .buttonStyle(.plain)
      
      if isInstrumentsVesselsExpanded && !allInstrumentsVessels.isEmpty {
        VStack(alignment: .leading, spacing: 8) {
          ForEach(allInstrumentsVessels, id: \.id) { item in
            HStack(spacing: 8) {
              Image(systemName: item.type == .instrument ? "wrench.and.screwdriver" : "square.stack.3d.up")
                .font(.caption)
                .foregroundColor(.secondary)
                .frame(width: 20)
              Text(item.name)
                .font(.subheadline)
                .foregroundColor(.primary)
              Spacer()
            }
            .padding(.horizontal)
            .padding(.vertical, 4)
          }
        }
        .padding(.vertical, 8)
        .background(Color(.systemBackground))
      }
    }
    .background(Color(.systemGray6))
    .cornerRadius(12)
  }
  
  // MARK: - Ingredients Section
  
  private func ingredientsSection(recipe: Mealplanning_Recipe) -> some View {
    let allIngredients = getAllIngredients(from: recipe)
    
    return VStack(alignment: .leading, spacing: 0) {
      Button(action: {
        withAnimation {
          isIngredientsExpanded.toggle()
        }
      }) {
        HStack {
          Text("Ingredients")
            .font(.headline)
            .foregroundColor(.primary)
          Spacer()
          Image(systemName: isIngredientsExpanded ? "chevron.down" : "chevron.right")
            .font(.caption)
            .foregroundColor(.secondary)
        }
        .padding()
        .background(Color(.systemGray6))
      }
      .buttonStyle(.plain)
      
      if isIngredientsExpanded && !allIngredients.isEmpty {
        VStack(alignment: .leading, spacing: 8) {
          ForEach(allIngredients, id: \.id) { ingredient in
            HStack(spacing: 8) {
              Image(systemName: "leaf")
                .font(.caption)
                .foregroundColor(.secondary)
                .frame(width: 20)
              VStack(alignment: .leading, spacing: 2) {
                Text(ingredient.name)
                  .font(.subheadline)
                  .foregroundColor(.primary)
                if !ingredient.quantityNotes.isEmpty {
                  Text(ingredient.quantityNotes)
                    .font(.caption)
                    .foregroundColor(.secondary)
                }
              }
              Spacer()
            }
            .padding(.horizontal)
            .padding(.vertical, 4)
          }
        }
        .padding(.vertical, 8)
        .background(Color(.systemBackground))
      }
    }
    .background(Color(.systemGray6))
    .cornerRadius(12)
  }
  
  // MARK: - Helper Methods
  
  private func getAllInstrumentsAndVessels(from recipe: Mealplanning_Recipe) -> [InstrumentVesselItem] {
    var items: [String: InstrumentVesselItem] = [:]
    
    for step in recipe.steps {
      // Collect instruments (only if it has a ValidInstrument and displayInSummaryLists is true)
      for instrument in step.instruments {
        // Only include if it has a ValidInstrument (not a recipe step product)
        // and displayInSummaryLists is true
        if instrument.hasInstrument {
          let validInstrument = instrument.instrument
          if validInstrument.displayInSummaryLists && items[validInstrument.id] == nil {
            items[validInstrument.id] = InstrumentVesselItem(
              id: validInstrument.id,
              name: instrument.name,
              type: .instrument
            )
          }
        }
      }
      
      // Collect vessels (only if it has a ValidVessel and displayInSummaryLists is true)
      for vessel in step.vessels {
        // Only include if it has a ValidVessel (not a recipe step product)
        // and displayInSummaryLists is true
        if vessel.hasVessel {
          let validVessel = vessel.vessel
          if validVessel.displayInSummaryLists && items[validVessel.id] == nil {
            items[validVessel.id] = InstrumentVesselItem(
              id: validVessel.id,
              name: vessel.name,
              type: .vessel
            )
          }
        }
      }
    }
    
    return Array(items.values).sorted { $0.name < $1.name }
  }
  
  private func getAllIngredients(from recipe: Mealplanning_Recipe) -> [Mealplanning_RecipeStepIngredient] {
    var ingredients: [String: Mealplanning_RecipeStepIngredient] = [:]
    
    for step in recipe.steps {
      for ingredient in step.ingredients {
        // Only include if it has a ValidIngredient (not a recipe step product)
        // and displayInSummaryLists is true
        if ingredient.hasIngredient {
          let validIngredient = ingredient.ingredient
          // Use ValidIngredient ID as key to ensure uniqueness
          let key = validIngredient.id
          if !key.isEmpty && ingredients[key] == nil {
            ingredients[key] = ingredient
          }
        }
      }
    }
    
    return Array(ingredients.values).sorted { $0.name < $1.name }
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

private struct InstrumentVesselItem: Identifiable {
  let id: String
  let name: String
  let type: ItemType
  
  enum ItemType {
    case instrument
    case vessel
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

