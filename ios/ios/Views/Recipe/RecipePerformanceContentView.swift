//
//  RecipePerformanceContentView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

/// A reusable view for displaying recipe performance content (ingredients, instruments, vessels, and steps)
/// This can be embedded in PerformRecipeView, Meal views, or any other context where recipe performance is needed
struct RecipePerformanceContentView: View {
  @Binding var checkedIngredients: Set<String>
  @Binding var checkedInstrumentsVessels: Set<String>
  @Binding var isInstrumentsVesselsExpanded: Bool
  @Binding var isIngredientsExpanded: Bool

  let recipe: Mealplanning_Recipe
  let viewModel: PerformRecipeViewModel

  var body: some View {
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
  }

  // MARK: - Recipe Header

  private func recipeHeader(recipe: Mealplanning_Recipe, viewModel: PerformRecipeViewModel)
    -> some View
  {
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
    let aggregatedItems = getAggregatedInstrumentsAndVessels(from: recipe)

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

      if isInstrumentsVesselsExpanded && !aggregatedItems.isEmpty {
        VStack(alignment: .leading, spacing: 8) {
          ForEach(aggregatedItems, id: \.itemID) { item in
            HStack(spacing: 12) {
              // Checkbox
              Button(action: {
                if checkedInstrumentsVessels.contains(item.itemID) {
                  checkedInstrumentsVessels.remove(item.itemID)
                } else {
                  checkedInstrumentsVessels.insert(item.itemID)
                }
              }) {
                Image(
                  systemName: checkedInstrumentsVessels.contains(item.itemID)
                    ? "checkmark.circle.fill" : "circle"
                )
                .font(.title3)
                .foregroundColor(checkedInstrumentsVessels.contains(item.itemID) ? .green : .gray)
              }
              .buttonStyle(.plain)

              HStack(spacing: 8) {
                Image(
                  systemName: item.type == .instrument
                    ? "wrench.and.screwdriver" : "square.stack.3d.up"
                )
                .font(.caption)
                .foregroundColor(.secondary)
                .frame(width: 20)

                HStack {
                  Text(item.name)
                    .font(.subheadline)
                    .foregroundColor(
                      checkedInstrumentsVessels.contains(item.itemID) ? .secondary : .primary
                    )
                    .strikethrough(checkedInstrumentsVessels.contains(item.itemID))

                  if let quantityText = item.quantityText {
                    Text(quantityText)
                      .font(.subheadline)
                      .fontWeight(.medium)
                      .foregroundColor(.secondary)
                  }
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

  // MARK: - Ingredients Section

  private func ingredientsSection(recipe: Mealplanning_Recipe) -> some View {
    let aggregatedIngredients = getAggregatedIngredients(from: recipe)

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

      if isIngredientsExpanded && !aggregatedIngredients.isEmpty {
        VStack(alignment: .leading, spacing: 8) {
          ForEach(aggregatedIngredients, id: \.ingredientID) { aggregated in
            HStack(spacing: 12) {
              // Checkbox
              Button(action: {
                if checkedIngredients.contains(aggregated.ingredientID) {
                  checkedIngredients.remove(aggregated.ingredientID)
                } else {
                  checkedIngredients.insert(aggregated.ingredientID)
                }
              }) {
                Image(
                  systemName: checkedIngredients.contains(aggregated.ingredientID)
                    ? "checkmark.circle.fill" : "circle"
                )
                .font(.title3)
                .foregroundColor(
                  checkedIngredients.contains(aggregated.ingredientID) ? .green : .gray)
              }
              .buttonStyle(.plain)

              VStack(alignment: .leading, spacing: 2) {
                HStack {
                  Text(aggregated.name)
                    .font(.subheadline)
                    .foregroundColor(
                      checkedIngredients.contains(aggregated.ingredientID) ? .secondary : .primary
                    )
                    .strikethrough(checkedIngredients.contains(aggregated.ingredientID))

                  if let quantityText = aggregated.quantityText {
                    Text(quantityText)
                      .font(.subheadline)
                      .fontWeight(.medium)
                      .foregroundColor(.secondary)
                  }
                }

                if !aggregated.quantityNotes.isEmpty {
                  Text(aggregated.quantityNotes)
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

  // MARK: - Steps List

  private func stepsList(recipe: Mealplanning_Recipe, viewModel: PerformRecipeViewModel)
    -> some View
  {
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

        // Step title with preparation and ingredients
        VStack(alignment: .leading, spacing: 4) {
          HStack {
            Text(formatStepTitle(step: step, viewModel: viewModel))
              .font(.headline)
              .foregroundColor(isCompleted ? .secondary : .primary)
              .italic(isCompleted)

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
              Text(
                "Complete steps \(prerequisites.map { String($0 + 1) }.joined(separator: ", ")) first"
              )
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
            let prerequisiteCompleted =
              prerequisiteStepIndex.map { viewModel.isStepCompleted($0) } ?? true

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
            let prerequisiteCompleted =
              prerequisiteStepIndex.map { viewModel.isStepCompleted($0) } ?? true

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
            let prerequisiteCompleted =
              prerequisiteStepIndex.map { viewModel.isStepCompleted($0) } ?? true

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

  // MARK: - Helper Methods

  private func getAggregatedInstrumentsAndVessels(from recipe: Mealplanning_Recipe)
    -> [AggregatedInstrumentVessel]
  {
    var aggregated: [String: AggregatedInstrumentVessel] = [:]

    for step in recipe.steps {
      // Collect instruments (only if it has a ValidInstrument and displayInSummaryLists is true)
      for instrument in step.instruments {
        if instrument.hasInstrument {
          let validInstrument = instrument.instrument
          if validInstrument.displayInSummaryLists {
            let itemID = validInstrument.id
            if !itemID.isEmpty {
              // Initialize or update aggregated item
              if aggregated[itemID] == nil {
                aggregated[itemID] = AggregatedInstrumentVessel(
                  itemID: itemID,
                  name: instrument.name,
                  type: .instrument
                )
              }

              // Aggregate quantities
              if instrument.hasQuantity {
                var current = aggregated[itemID]!
                let quantity = instrument.quantity
                current.addQuantity(quantity)
                aggregated[itemID] = current
              }
            }
          }
        }
      }

      // Collect vessels (only if it has a ValidVessel and displayInSummaryLists is true)
      for vessel in step.vessels {
        if vessel.hasVessel {
          let validVessel = vessel.vessel
          if validVessel.displayInSummaryLists {
            let itemID = validVessel.id
            if !itemID.isEmpty {
              // Initialize or update aggregated item
              if aggregated[itemID] == nil {
                aggregated[itemID] = AggregatedInstrumentVessel(
                  itemID: itemID,
                  name: vessel.name,
                  type: .vessel
                )
              }

              // Aggregate quantities
              if vessel.hasQuantity {
                var current = aggregated[itemID]!
                let quantity = vessel.quantity
                current.addQuantity(quantity)
                aggregated[itemID] = current
              }
            }
          }
        }
      }
    }

    return Array(aggregated.values).sorted { $0.name < $1.name }
  }

  private func getAggregatedIngredients(from recipe: Mealplanning_Recipe) -> [AggregatedIngredient]
  {
    var aggregated: [String: AggregatedIngredient] = [:]

    for step in recipe.steps {
      for ingredient in step.ingredients {
        // Only include if it has a ValidIngredient (not a recipe step product)
        if ingredient.hasIngredient {
          let validIngredient = ingredient.ingredient
          // Use ValidIngredient ID as key to ensure uniqueness
          let key = validIngredient.id
          if !key.isEmpty {
            // Initialize or update aggregated ingredient
            if aggregated[key] == nil {
              aggregated[key] = AggregatedIngredient(
                ingredientID: key,
                name: ingredient.name,
                quantityNotes: ingredient.quantityNotes,
                measurementUnit: ingredient.hasMeasurementUnit ? ingredient.measurementUnit : nil
              )
            }

            // Aggregate quantities
            if ingredient.hasQuantity {
              var current = aggregated[key]!
              let quantity = ingredient.quantity
              current.addQuantity(quantity)
              aggregated[key] = current
            }
          }
        }
      }
    }

    return Array(aggregated.values).sorted { $0.name < $1.name }
  }

  // MARK: - Step Title Formatting

  private func formatStepTitle(step: Mealplanning_RecipeStep, viewModel: PerformRecipeViewModel)
    -> String
  {
    var parts: [String] = []

    // Add preparation name
    if step.hasPreparation && !step.preparation.name.isEmpty {
      parts.append(step.preparation.name)
    }

    // Add ingredient names (only those with ValidIngredient)
    let ingredientNames = step.ingredients
      .filter { $0.hasIngredient }
      .map { $0.name }

    if !ingredientNames.isEmpty {
      parts.append(formatList(ingredientNames))
    }

    // Add instruments (only those with ValidInstrument and displayInSummaryLists)
    let instrumentNames = step.instruments
      .filter { $0.hasInstrument && $0.instrument.displayInSummaryLists }
      .map { $0.name }

    if !instrumentNames.isEmpty {
      parts.append("with \(formatList(instrumentNames))")
    }

    // If no parts, fall back to step number
    if parts.isEmpty {
      return "Step \(Int(step.index) + 1)"
    }

    return parts.joined(separator: " ")
  }

  private func formatList(_ items: [String]) -> String {
    guard !items.isEmpty else { return "" }

    if items.count == 1 {
      return items[0]
    } else if items.count == 2 {
      return "\(items[0]) and \(items[1])"
    } else {
      let allButLast = items.dropLast().joined(separator: ", ")
      return "\(allButLast), and \(items.last!)"
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

private struct AggregatedInstrumentVessel: Identifiable {
  let itemID: String
  let name: String
  let type: InstrumentVesselItem.ItemType

  private var totalMin: UInt32 = 0
  private var totalMax: UInt32? = nil
  private var hasAnyQuantity: Bool = false

  init(itemID: String, name: String, type: InstrumentVesselItem.ItemType) {
    self.itemID = itemID
    self.name = name
    self.type = type
  }

  mutating func addQuantity(_ quantity: Common_Uint32RangeWithOptionalMax) {
    hasAnyQuantity = true
    totalMin += quantity.min

    if quantity.hasMax {
      if let currentMax = totalMax {
        totalMax = currentMax + quantity.max
      } else {
        totalMax = quantity.max
      }
    } else {
      // If any quantity doesn't have a max, the total doesn't have a max
      totalMax = nil
    }
  }

  mutating func addQuantity(_ quantity: Common_Uint16RangeWithOptionalMax) {
    hasAnyQuantity = true
    totalMin += quantity.min

    if quantity.hasMax {
      if let currentMax = totalMax {
        totalMax = currentMax + quantity.max
      } else {
        totalMax = quantity.max
      }
    } else {
      // If any quantity doesn't have a max, the total doesn't have a max
      totalMax = nil
    }
  }

  var quantityText: String? {
    guard hasAnyQuantity else { return nil }

    if let max = totalMax {
      if totalMin == max {
        return "\(totalMin)"
      } else {
        return "\(totalMin) - \(max)"
      }
    } else {
      return "\(totalMin)+"
    }
  }

  var id: String {
    itemID
  }
}

private struct AggregatedIngredient: Identifiable {
  let ingredientID: String
  let name: String
  let quantityNotes: String
  let measurementUnit: Mealplanning_ValidMeasurementUnit?

  private var totalMin: Float = 0
  private var totalMax: Float? = nil
  private var hasAnyQuantity: Bool = false

  init(
    ingredientID: String, name: String, quantityNotes: String,
    measurementUnit: Mealplanning_ValidMeasurementUnit?
  ) {
    self.ingredientID = ingredientID
    self.name = name
    self.quantityNotes = quantityNotes
    self.measurementUnit = measurementUnit
  }

  mutating func addQuantity(_ quantity: Common_Float32RangeWithOptionalMax) {
    hasAnyQuantity = true
    totalMin += quantity.min

    if quantity.hasMax {
      if let currentMax = totalMax {
        totalMax = currentMax + quantity.max
      } else {
        totalMax = quantity.max
      }
    } else {
      // If any quantity doesn't have a max, the total doesn't have a max
      totalMax = nil
    }
  }

  var quantityText: String? {
    guard hasAnyQuantity else { return nil }

    let unitName = measurementUnit?.name ?? ""
    let unit = unitName.isEmpty ? "" : " \(unitName)"

    // Format numbers - use fewer decimals for whole numbers
    let formatMin = totalMin.truncatingRemainder(dividingBy: 1) == 0 ? "%.0f" : "%.2f"
    let formatMax =
      totalMax.map { $0.truncatingRemainder(dividingBy: 1) == 0 ? "%.0f" : "%.2f" } ?? "%.2f"

    if let max = totalMax {
      if totalMin == max {
        return String(format: "\(formatMin)%@", totalMin, unit).trimmingCharacters(in: .whitespaces)
      } else {
        return String(format: "\(formatMin) - \(formatMax)%@", totalMin, max, unit)
          .trimmingCharacters(in: .whitespaces)
      }
    } else {
      return String(format: "\(formatMin)+%@", totalMin, unit).trimmingCharacters(in: .whitespaces)
    }
  }

  var id: String {
    ingredientID
  }
}
