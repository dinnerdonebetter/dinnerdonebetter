//
//  PerformRecipeViewModel.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import Foundation
import GRPCCore
import SwiftProtobuf
import SwiftUI

@Observable
@MainActor
class PerformRecipeViewModel {
  var recipe: Mealplanning_Recipe?
  var prepTasks: [Mealplanning_RecipePrepTask] = []
  var isLoading = false
  var isLoadingPrepTasks = false
  var errorMessage: String?

  // Track which steps are completed (by step ID: "recipeID:stepID")
  var completedSteps: Set<String> = []

  // Track which prep tasks are completed (by prep task ID)
  var completedPrepTasks: Set<String> = []

  // Special step: wash hands (index -1)
  var washHandsCompleted: Bool = false

  // Map from product ID to the step that produces it (recipeID:stepID)
  var productIDToStepKey: [String: String] = [:]

  // Special index for wash hands step
  static let washHandsStepIndex = -1

  private let recipeID: String
  private let authManager: AuthenticationManager

  init(recipeID: String, authManager: AuthenticationManager) {
    self.recipeID = recipeID
    self.authManager = authManager
  }

  func loadRecipe() async {
    isLoading = true
    errorMessage = nil

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "PerformRecipeViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }

      // Get OAuth2 token (will refresh if needed)
      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "PerformRecipeViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      // Create request
      var request = Mealplanning_GetRecipeRequest()
      request.recipeID = recipeID

      // Execute request
      let response = try await clientManager.client.mealPlanning.getRecipe(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      self.recipe = response.result
      buildProductIDToStepIndexMapping()

      // Load prep tasks after recipe is loaded
      await loadPrepTasks()
    } catch {
      errorMessage = "Failed to load recipe: \(error.localizedDescription)"
      print("❌ Error loading recipe: \(error)")
    }

    isLoading = false
  }

  func loadPrepTasks() async {
    isLoadingPrepTasks = true

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "PerformRecipeViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }

      // Get OAuth2 token (will refresh if needed)
      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "PerformRecipeViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      // Create request
      var request = Mealplanning_GetRecipePrepTasksRequest()
      request.recipeID = recipeID

      // Execute request
      let response = try await clientManager.client.mealPlanning.getRecipePrepTasks(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      self.prepTasks = response.results
    } catch {
      print("⚠️ Error loading prep tasks: \(error)")
      // Don't set error message for prep tasks - it's not critical
    }

    isLoadingPrepTasks = false
  }

  // Build a mapping from recipe step product IDs to the step that produces them
  private func buildProductIDToStepIndexMapping() {
    guard let recipe = recipe else { return }

    // Map products from main recipe
    for step in recipe.steps {
      let stepKey = "\(recipe.id):\(step.id)"
      for product in step.products {
        productIDToStepKey[product.id] = stepKey
      }
    }

    // Map products from associated recipes
    for associatedRecipe in recipe.associatedRecipes {
      for step in associatedRecipe.steps {
        let stepKey = "\(associatedRecipe.id):\(step.id)"
        for product in step.products {
          productIDToStepKey[product.id] = stepKey
        }
      }
    }
  }

  // Helper to create a step key from recipeID and stepID
  private func stepKey(recipeID: String, stepID: String) -> String {
    return "\(recipeID):\(stepID)"
  }

  // Check if a step can be checked off (all prerequisites are completed)
  // Supports both old API (stepIndex for main recipe) and new API (recipeID + stepID)
  func canCheckStep(_ stepIndex: Int) -> Bool {
    // Wash hands step can always be checked
    if stepIndex == Self.washHandsStepIndex {
      return true
    }

    guard let recipe = recipe, stepIndex >= 0, stepIndex < recipe.steps.count else {
      return false
    }
    let step = recipe.steps[stepIndex]
    return canCheckStep(recipeID: recipe.id, stepID: step.id)
  }

  func canCheckStep(recipeID: String, stepID: String) -> Bool {
    // All steps require wash hands to be completed first
    if !washHandsCompleted {
      return false
    }

    guard let recipe = recipe else {
      return false
    }

    // Find the step (could be in main recipe or associated recipe)
    let step: Mealplanning_RecipeStep?
    if recipeID == recipe.id {
      step = recipe.steps.first { $0.id == stepID }
    } else {
      step = recipe.associatedRecipes.first { $0.id == recipeID }?.steps.first { $0.id == stepID }
    }

    guard let step = step else {
      return false
    }

    // Check all ingredients
    for ingredient in step.ingredients where ingredient.hasRecipeStepProductID {
      let productID = ingredient.recipeStepProductID
      if let prerequisiteStepKey = productIDToStepKey[productID] {
        if !completedSteps.contains(prerequisiteStepKey) {
          return false
        }
      }
    }

    // Check all instruments
    for instrument in step.instruments where instrument.hasRecipeStepProductID {
      let productID = instrument.recipeStepProductID
      if let prerequisiteStepKey = productIDToStepKey[productID] {
        if !completedSteps.contains(prerequisiteStepKey) {
          return false
        }
      }
    }

    // Check all vessels
    for vessel in step.vessels where vessel.hasRecipeStepProductID {
      let productID = vessel.recipeStepProductID
      if let prerequisiteStepKey = productIDToStepKey[productID] {
        if !completedSteps.contains(prerequisiteStepKey) {
          return false
        }
      }
    }

    return true
  }

  // Toggle step completion
  // Supports both old API (stepIndex for main recipe) and new API (recipeID + stepID)
  func toggleStep(_ stepIndex: Int) {
    // Handle wash hands step
    if stepIndex == Self.washHandsStepIndex {
      washHandsCompleted.toggle()
      // If unchecking wash hands, uncheck all other steps
      if !washHandsCompleted {
        completedSteps.removeAll()
      }
      return
    }

    guard let recipe = recipe, stepIndex >= 0, stepIndex < recipe.steps.count else {
      return
    }
    let step = recipe.steps[stepIndex]
    toggleStep(recipeID: recipe.id, stepID: step.id)
  }

  func toggleStep(recipeID: String, stepID: String) {
    let currentStepKey = stepKey(recipeID: recipeID, stepID: stepID)

    guard canCheckStep(recipeID: recipeID, stepID: stepID) else {
      return
    }

    if completedSteps.contains(currentStepKey) {
      // When unchecking, also uncheck all dependent steps
      uncheckStepAndDependents(stepKey: currentStepKey)
    } else {
      completedSteps.insert(currentStepKey)
    }
  }

  // Recursively uncheck a step and all steps that depend on it
  private func uncheckStepAndDependents(stepKey: String) {
    guard let recipe = recipe else { return }
    completedSteps.remove(stepKey)

    // Parse recipeID and stepID from stepKey
    let components = stepKey.split(separator: ":", maxSplits: 1)
    guard components.count == 2,
      let stepID = String(components[1]) as String?
    else {
      return
    }
    let recipeID = String(components[0])

    // Find the step to get its products
    let step: Mealplanning_RecipeStep?
    if recipeID == recipe.id {
      step = recipe.steps.first { $0.id == stepID }
    } else {
      step = recipe.associatedRecipes.first { $0.id == recipeID }?.steps.first { $0.id == stepID }
    }

    guard let step = step else { return }

    // Find all steps (in main recipe and associated recipes) that depend on this step's products
    for product in step.products {
      // Check main recipe steps
      processDependentSteps(
        steps: recipe.steps,
        recipeID: recipe.id,
        productID: product.id,
        excludeStepKey: stepKey
      )

      // Check associated recipe steps
      for associatedRecipe in recipe.associatedRecipes {
        processDependentSteps(
          steps: associatedRecipe.steps,
          recipeID: associatedRecipe.id,
          productID: product.id,
          excludeStepKey: stepKey
        )
      }
    }
  }

  // Helper function to check if a step depends on a product
  private func stepDependsOnProduct(_ step: Mealplanning_RecipeStep, productID: String) -> Bool {
    // Check ingredients
    for ingredient in step.ingredients {
      if ingredient.hasRecipeStepProductID && ingredient.recipeStepProductID == productID {
        return true
      }
    }

    // Check instruments
    for instrument in step.instruments {
      if instrument.hasRecipeStepProductID && instrument.recipeStepProductID == productID {
        return true
      }
    }

    // Check vessels
    for vessel in step.vessels {
      if vessel.hasRecipeStepProductID && vessel.recipeStepProductID == productID {
        return true
      }
    }

    return false
  }

  // Helper function to process dependent steps for a given recipe
  private func processDependentSteps(
    steps: [Mealplanning_RecipeStep],
    recipeID: String,
    productID: String,
    excludeStepKey: String
  ) {
    for dependentStep in steps {
      let dependentStepKey = self.stepKey(recipeID: recipeID, stepID: dependentStep.id)
      if dependentStepKey == excludeStepKey {
        continue  // Skip self
      }

      if stepDependsOnProduct(dependentStep, productID: productID)
        && completedSteps.contains(dependentStepKey)
      {
        uncheckStepAndDependents(stepKey: dependentStepKey)
      }
    }
  }

  // Get the step at the given index
  func getStep(_ stepIndex: Int) -> Mealplanning_RecipeStep? {
    guard let recipe = recipe, stepIndex < recipe.steps.count else {
      return nil
    }
    return recipe.steps[stepIndex]
  }

  // Check if a step is completed
  // Supports both old API (stepIndex for main recipe) and new API (recipeID + stepID)
  func isStepCompleted(_ stepIndex: Int) -> Bool {
    if stepIndex == Self.washHandsStepIndex {
      return washHandsCompleted
    }
    guard let recipe = recipe, stepIndex >= 0, stepIndex < recipe.steps.count else {
      return false
    }
    let step = recipe.steps[stepIndex]
    return isStepCompleted(recipeID: recipe.id, stepID: step.id)
  }

  func isStepCompleted(recipeID: String, stepID: String) -> Bool {
    let currentStepKey = stepKey(recipeID: recipeID, stepID: stepID)
    return completedSteps.contains(currentStepKey)
  }

  // Get all prerequisite step indices for a given step (old API for main recipe)
  func getPrerequisiteStepIndices(_ stepIndex: Int) -> [Int] {
    guard let recipe = recipe, stepIndex >= 0, stepIndex < recipe.steps.count else {
      return []
    }
    let step = recipe.steps[stepIndex]
    let stepKeys = getPrerequisiteStepKeys(recipeID: recipe.id, stepID: step.id)

    // Convert step keys back to indices for main recipe steps only
    var prerequisites: Set<Int> = []
    for stepKey in stepKeys {
      let components = stepKey.split(separator: ":", maxSplits: 1)
      if components.count == 2, String(components[0]) == recipe.id {
        if let stepID = String(components[1]) as String?,
          let index = recipe.steps.firstIndex(where: { $0.id == stepID })
        {
          prerequisites.insert(index)
        }
      }
    }
    return Array(prerequisites).sorted()
  }

  // Step categorization
  enum StepCategory {
    case upNext
    case forLater
    case done
  }

  // Categorize a step based on its completion status and dependencies
  func categorizeStep(recipeID: String, stepID: String) -> StepCategory {
    let stepKey = self.stepKey(recipeID: recipeID, stepID: stepID)

    // If step is completed, it's done
    if completedSteps.contains(stepKey) {
      return .done
    }

    // Check if all dependencies are satisfied
    let prerequisiteKeys = getPrerequisiteStepKeys(recipeID: recipeID, stepID: stepID)
    let allDependenciesDone = prerequisiteKeys.allSatisfy { completedSteps.contains($0) }

    if allDependenciesDone {
      return .upNext
    } else {
      return .forLater
    }
  }

  // Get all prerequisite step keys for a given step
  func getPrerequisiteStepKeys(recipeID: String, stepID: String) -> [String] {
    guard let recipe = recipe else { return [] }

    // Find the step
    let step: Mealplanning_RecipeStep?
    if recipeID == recipe.id {
      step = recipe.steps.first { $0.id == stepID }
    } else {
      step = recipe.associatedRecipes.first { $0.id == recipeID }?.steps.first { $0.id == stepID }
    }

    guard let step = step else { return [] }

    var prerequisites: Set<String> = []

    // Check all ingredients
    for ingredient in step.ingredients where ingredient.hasRecipeStepProductID {
      let productID = ingredient.recipeStepProductID
      if let prerequisiteStepKey = productIDToStepKey[productID] {
        prerequisites.insert(prerequisiteStepKey)
      }
    }

    // Check all instruments
    for instrument in step.instruments where instrument.hasRecipeStepProductID {
      let productID = instrument.recipeStepProductID
      if let prerequisiteStepKey = productIDToStepKey[productID] {
        prerequisites.insert(prerequisiteStepKey)
      }
    }

    // Check all vessels
    for vessel in step.vessels where vessel.hasRecipeStepProductID {
      let productID = vessel.recipeStepProductID
      if let prerequisiteStepKey = productIDToStepKey[productID] {
        prerequisites.insert(prerequisiteStepKey)
      }
    }

    return Array(prerequisites).sorted()
  }

  // Get the step key that produces a given product ID
  func getStepKeyForProductID(_ productID: String) -> String? {
    return productIDToStepKey[productID]
  }

  // Get the step index for a product ID (for main recipe only, for display purposes)
  // Returns nil if the product is from an associated recipe
  func getStepIndexForProductID(_ productID: String) -> Int? {
    guard let recipe = recipe,
      let stepKey = productIDToStepKey[productID]
    else {
      return nil
    }

    // Parse step key
    let components = stepKey.split(separator: ":", maxSplits: 1)
    guard components.count == 2,
      String(components[0]) == recipe.id,
      let stepID = String(components[1]) as String?
    else {
      return nil  // Product is from associated recipe
    }

    // Find the step index in main recipe
    return recipe.steps.firstIndex(where: { $0.id == stepID })
  }

  // MARK: - Prep Task Completion

  // Check if a prep task is completed
  func isPrepTaskCompleted(_ prepTaskID: String) -> Bool {
    return completedPrepTasks.contains(prepTaskID)
  }

  // Toggle prep task completion and mark/unmark associated steps
  func togglePrepTask(_ prepTask: Mealplanning_RecipePrepTask) {
    let prepTaskID = prepTask.id

    if completedPrepTasks.contains(prepTaskID) {
      // Uncheck prep task
      completedPrepTasks.remove(prepTaskID)

      // Uncheck all associated steps
      for taskStep in prepTask.taskSteps where !taskStep.belongsToRecipeStep.isEmpty {
        let stepID = taskStep.belongsToRecipeStep
        // Find which recipe this step belongs to
        if let recipe = recipe {
          // Check main recipe steps
          if recipe.steps.contains(where: { $0.id == stepID }) {
            let stepKey = stepKey(recipeID: recipe.id, stepID: stepID)
            uncheckStepAndDependents(stepKey: stepKey)
          } else {
            // Check associated recipe steps
            for associatedRecipe in recipe.associatedRecipes {
              if associatedRecipe.steps.contains(where: { $0.id == stepID }) {
                let stepKey = stepKey(recipeID: associatedRecipe.id, stepID: stepID)
                uncheckStepAndDependents(stepKey: stepKey)
                break
              }
            }
          }
        }
      }
    } else {
      // Check prep task
      completedPrepTasks.insert(prepTaskID)

      // Check all associated steps (allow checking even if prerequisites aren't satisfied,
      // since the user has already done the prep task ahead of time)
      for taskStep in prepTask.taskSteps where !taskStep.belongsToRecipeStep.isEmpty {
        let stepID = taskStep.belongsToRecipeStep
        // Find which recipe this step belongs to
        if let recipe = recipe {
          // Check main recipe steps
          if recipe.steps.contains(where: { $0.id == stepID }) {
            let stepKey = stepKey(recipeID: recipe.id, stepID: stepID)
            // Check the step (don't require prerequisites since prep task is done)
            completedSteps.insert(stepKey)
          } else {
            // Check associated recipe steps
            for associatedRecipe in recipe.associatedRecipes {
              if associatedRecipe.steps.contains(where: { $0.id == stepID }) {
                let stepKey = stepKey(recipeID: associatedRecipe.id, stepID: stepID)
                // Check the step (don't require prerequisites since prep task is done)
                completedSteps.insert(stepKey)
                break
              }
            }
          }
        }
      }
    }
  }
}
