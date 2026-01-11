//
//  RecipeOptionSelectionView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct RecipeOptionSelectionView: View {
  @Binding var isPresented: Bool
  @State private var selections: [String: [String: [UInt32: UInt32]]] = [:]  // recipeID -> (stepID -> (ingredientIndex -> selectedOptionIndex))

  let recipes: [Mealplanning_Recipe]
  let onSave: ([String: [String: [UInt32: UInt32]]]) -> Void

  var body: some View {
    NavigationView {
      ScrollView {
        VStack(alignment: .leading, spacing: 20) {
          Text("Select Recipe Options")
            .font(.headline)
            .padding(.horizontal)

          Text(
            "Choose your preferred options for each recipe. You can skip this step to use the default options."
          )
          .font(.subheadline)
          .foregroundColor(.secondary)
          .padding(.horizontal)

          ForEach(recipes, id: \.id) { recipe in
            recipeOptionSection(recipe: recipe)
          }
        }
        .padding(.vertical)
      }
      .navigationTitle("Recipe Options")
      .navigationBarTitleDisplayMode(.inline)
      .toolbar {
        ToolbarItem(placement: .cancellationAction) {
          Button("Skip") {
            // Use defaults (optionIndex: 0)
            onSave([:])
            isPresented = false
          }
        }
        ToolbarItem(placement: .confirmationAction) {
          Button("Save") {
            onSave(selections)
            isPresented = false
          }
        }
      }
    }
  }

  private func recipeOptionSection(recipe: Mealplanning_Recipe) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      Text(recipe.name.isEmpty ? "Unnamed Recipe" : recipe.name)
        .font(.headline)
        .padding(.horizontal)

      // Process each step for option groups
      ForEach(Array(recipe.steps.enumerated()), id: \.element.id) { stepIndex, step in
        stepOptionSection(recipe: recipe, step: step, stepIndex: stepIndex)
      }

      // Process associated recipes
      ForEach(recipe.associatedRecipes, id: \.id) { associatedRecipe in
        VStack(alignment: .leading, spacing: 8) {
          Text("From \(associatedRecipe.name.isEmpty ? "Unnamed Recipe" : associatedRecipe.name)")
            .font(.subheadline)
            .foregroundColor(.secondary)
            .padding(.horizontal)

          ForEach(Array(associatedRecipe.steps.enumerated()), id: \.element.id) { stepIndex, step in
            stepOptionSection(recipe: associatedRecipe, step: step, stepIndex: stepIndex)
          }
        }
      }
    }
    .padding(.vertical, 8)
    .background(Color(.systemGray6))
    .cornerRadius(12)
    .padding(.horizontal)
  }

  private func stepOptionSection(
    recipe: Mealplanning_Recipe,
    step: Mealplanning_RecipeStep,
    stepIndex: Int
  ) -> some View {
    // Group ingredients by options (only ingredients have selectable options)
    let (_, ingredientGroups) = groupIngredientsForSelection(
      step.ingredients, stepID: step.id, recipeID: recipe.id)

    guard !ingredientGroups.isEmpty else {
      return AnyView(EmptyView())
    }

    return AnyView(
      VStack(alignment: .leading, spacing: 8) {
        Text("Step \(stepIndex + 1)")
          .font(.subheadline)
          .fontWeight(.semibold)
          .padding(.horizontal)

        // Ingredient option groups (only ingredients have selectable options)
        ForEach(ingredientGroups, id: \.id) { group in
          optionGroupPicker(
            title: "Ingredient",
            options: group.options.map { ($0.ingredient.name, $0.optionIndex) },
            selectedIndex: Binding(
              get: {
                selections[recipe.id]?[step.id]?[group.index]
                  ?? (group.options.first?.optionIndex ?? 0)
              },
              set: { newValue in
                if selections[recipe.id] == nil {
                  selections[recipe.id] = [:]
                }
                if selections[recipe.id]?[step.id] == nil {
                  selections[recipe.id]?[step.id] = [:]
                }
                selections[recipe.id]?[step.id]?[group.index] = newValue
              }
            )
          )
        }
      }
    )
  }

  private func optionGroupPicker(
    title: String,
    options: [(String, UInt32)],
    selectedIndex: Binding<UInt32>
  ) -> some View {
    VStack(alignment: .leading, spacing: 4) {
      Text("\(title):")
        .font(.caption)
        .foregroundColor(.secondary)
        .padding(.horizontal)

      Picker("", selection: selectedIndex) {
        ForEach(options, id: \.1) { name, index in
          Text(name).tag(index)
        }
      }
      .pickerStyle(.segmented)
      .padding(.horizontal)
    }
  }

  // Helper functions to group items (similar to the ones in RecipePerformanceContentView+Helpers)
  private func groupIngredientsForSelection(
    _ ingredients: [Mealplanning_RecipeStepIngredient],
    stepID: String,
    recipeID: String
  ) -> (
    regular: [Mealplanning_RecipeStepIngredient],
    optionGroups: [OptionGroupAggregate]
  ) {
    var regular: [Mealplanning_RecipeStepIngredient] = []
    var optionGroupsByIndex: [UInt32: [Mealplanning_RecipeStepIngredient]] = [:]

    for ingredient in ingredients {
      // Index 0 typically means not in an option group
      if ingredient.index != 0 {
        let index = ingredient.index
        let hasOptions = ingredients.contains { other in
          other.id != ingredient.id && other.index != 0 && other.index == index
        }

        if hasOptions {
          if optionGroupsByIndex[index] == nil {
            optionGroupsByIndex[index] = []
          }
          optionGroupsByIndex[index]?.append(ingredient)
        } else {
          regular.append(ingredient)
        }
      } else {
        regular.append(ingredient)
      }
    }

    var optionGroups: [OptionGroupAggregate] = []
    for (index, groupIngredients) in optionGroupsByIndex {
      let sorted = groupIngredients.sorted { lhs, rhs in
        let lhsIndex = lhs.optionIndex
        let rhsIndex = rhs.optionIndex
        return lhsIndex < rhsIndex
      }

      var options: [IngredientOption] = []
      for ingredient in sorted {
        let optionIndex = ingredient.optionIndex
        let optionID = "\(stepID)-\(index)-\(optionIndex)"

        var aggregated = AggregatedIngredient(
          ingredientID: ingredient.hasIngredient ? ingredient.ingredient.id : ingredient.id,
          name: ingredient.name,
          quantityNotes: ingredient.quantityNotes,
          measurementUnit: ingredient.hasMeasurementUnit ? ingredient.measurementUnit : nil
        )

        if ingredient.hasQuantity {
          aggregated.addQuantity(ingredient.quantity)
        }

        options.append(
          IngredientOption(
            id: optionID,
            ingredient: ingredient,
            optionIndex: optionIndex,
            aggregated: aggregated
          )
        )
      }

      let groupID = "\(recipeID)-\(stepID)-\(index)"
      optionGroups.append(
        OptionGroupAggregate(
          id: groupID,
          recipeID: recipeID,
          stepID: stepID,
          stepIndex: 0,  // Not needed for selection
          index: index,
          options: options,
          selectedOptionIndex: nil,
          sourceRecipeID: nil,
          sourceRecipeName: nil
        )
      )
    }

    optionGroups.sort { $0.index < $1.index }
    return (regular: regular, optionGroups: optionGroups)
  }
}
