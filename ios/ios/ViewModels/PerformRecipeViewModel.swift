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
  var isLoading = false
  var errorMessage: String?

  // Track which steps are completed (by step index)
  var completedSteps: Set<Int> = []

  // Map from product ID to the step index that produces it
  var productIDToStepIndex: [String: Int] = [:]

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
    } catch {
      errorMessage = "Failed to load recipe: \(error.localizedDescription)"
      print("❌ Error loading recipe: \(error)")
    }

    isLoading = false
  }

  // Build a mapping from recipe step product IDs to the step index that produces them
  private func buildProductIDToStepIndexMapping() {
    guard let recipe = recipe else { return }
    for (stepIndex, step) in recipe.steps.enumerated() {
      for product in step.products {
        productIDToStepIndex[product.id] = stepIndex
      }
    }
  }

  // Check if a step can be checked off (all prerequisites are completed)
  func canCheckStep(_ stepIndex: Int) -> Bool {
    guard let recipe = recipe, stepIndex < recipe.steps.count else {
      return false
    }

    let step = recipe.steps[stepIndex]

    // Check all ingredients
    for ingredient in step.ingredients where ingredient.hasRecipeStepProductID {
      let productID = ingredient.recipeStepProductID
      if let prerequisiteStepIndex = productIDToStepIndex[productID] {
        if !completedSteps.contains(prerequisiteStepIndex) {
          return false
        }
      }
    }

    // Check all instruments
    for instrument in step.instruments where instrument.hasRecipeStepProductID {
      let productID = instrument.recipeStepProductID
      if let prerequisiteStepIndex = productIDToStepIndex[productID] {
        if !completedSteps.contains(prerequisiteStepIndex) {
          return false
        }
      }
    }

    // Check all vessels
    for vessel in step.vessels where vessel.hasRecipeStepProductID {
      let productID = vessel.recipeStepProductID
      if let prerequisiteStepIndex = productIDToStepIndex[productID] {
        if !completedSteps.contains(prerequisiteStepIndex) {
          return false
        }
      }
    }

    return true
  }

  // Toggle step completion
  func toggleStep(_ stepIndex: Int) {
    guard canCheckStep(stepIndex) else {
      return
    }

    if completedSteps.contains(stepIndex) {
      // When unchecking, also uncheck all dependent steps
      uncheckStepAndDependents(stepIndex)
    } else {
      completedSteps.insert(stepIndex)
    }
  }

  // Recursively uncheck a step and all steps that depend on it
  private func uncheckStepAndDependents(_ stepIndex: Int) {
    guard let recipe = recipe else { return }
    completedSteps.remove(stepIndex)

    // Find all steps that depend on this step's products
    let step = recipe.steps[stepIndex]
    for product in step.products {
      // Find all steps that use this product
      for (dependentStepIndex, dependentStep) in recipe.steps.enumerated() {
        if dependentStepIndex <= stepIndex {
          continue  // Only check later steps
        }

        var dependsOnProduct = false

        // Check if any ingredient uses this product
        for ingredient in dependentStep.ingredients {
          if ingredient.hasRecipeStepProductID && ingredient.recipeStepProductID == product.id {
            dependsOnProduct = true
            break
          }
        }

        // Check if any instrument uses this product
        if !dependsOnProduct {
          for instrument in dependentStep.instruments {
            if instrument.hasRecipeStepProductID && instrument.recipeStepProductID == product.id {
              dependsOnProduct = true
              break
            }
          }
        }

        // Check if any vessel uses this product
        if !dependsOnProduct {
          for vessel in dependentStep.vessels {
            if vessel.hasRecipeStepProductID && vessel.recipeStepProductID == product.id {
              dependsOnProduct = true
              break
            }
          }
        }

        if dependsOnProduct && completedSteps.contains(dependentStepIndex) {
          uncheckStepAndDependents(dependentStepIndex)
        }
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
  func isStepCompleted(_ stepIndex: Int) -> Bool {
    return completedSteps.contains(stepIndex)
  }

  // Get all prerequisite step indices for a given step
  func getPrerequisiteStepIndices(_ stepIndex: Int) -> [Int] {
    guard let recipe = recipe, stepIndex < recipe.steps.count else {
      return []
    }

    var prerequisites: Set<Int> = []
    let step = recipe.steps[stepIndex]

    // Check all ingredients
    for ingredient in step.ingredients where ingredient.hasRecipeStepProductID {
      let productID = ingredient.recipeStepProductID
      if let prerequisiteStepIndex = productIDToStepIndex[productID] {
        prerequisites.insert(prerequisiteStepIndex)
      }
    }

    // Check all instruments
    for instrument in step.instruments where instrument.hasRecipeStepProductID {
      let productID = instrument.recipeStepProductID
      if let prerequisiteStepIndex = productIDToStepIndex[productID] {
        prerequisites.insert(prerequisiteStepIndex)
      }
    }

    // Check all vessels
    for vessel in step.vessels where vessel.hasRecipeStepProductID {
      let productID = vessel.recipeStepProductID
      if let prerequisiteStepIndex = productIDToStepIndex[productID] {
        prerequisites.insert(prerequisiteStepIndex)
      }
    }

    return Array(prerequisites).sorted()
  }

  // Get the step index that produces a given product ID
  func getStepIndexForProductID(_ productID: String) -> Int? {
    return productIDToStepIndex[productID]
  }
}
