//
//  RecipePerformanceContentView+Views.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

// swiftlint:disable file_length

import SwiftProtobuf
import SwiftUI

// MARK: - Step Card View

struct StepCardView: View {
  let step: Mealplanning_RecipeStep
  let index: Int
  let viewModel: PerformRecipeViewModel
  let formatStepTitle: (Mealplanning_RecipeStep, PerformRecipeViewModel) -> String
  let recipeID: String
  var mealPlanSelections: [Mealplanning_MealPlanRecipeOptionSelection]?

  var body: some View {
    let isCompleted = viewModel.isStepCompleted(index)
    let canCheck = viewModel.canCheckStep(index)
    let prerequisites = viewModel.getPrerequisiteStepIndices(index)
    let hasPrerequisites = !prerequisites.isEmpty
    let allPrerequisitesCompleted = prerequisites.allSatisfy { viewModel.isStepCompleted($0) }

    return VStack(alignment: .leading, spacing: 12) {
      // Step header with checkbox
      HStack(alignment: .top, spacing: 12) {
        // Checkbox
        Button(
          action: {
            viewModel.toggleStep(index)
          },
          label: {
            Image(systemName: isCompleted ? "checkmark.circle.fill" : "circle")
              .font(.title2)
              .foregroundColor(
                canCheck ? (isCompleted ? .green : .blue) : .gray
              )
          }
        )
        .disabled(!canCheck)

        // Step title with preparation and ingredients
        VStack(alignment: .leading, spacing: 4) {
          HStack {
            Text(formatStepTitle(step, viewModel))
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

          // Wash hands requirement warning (show first if blocking)
          if !viewModel.washHandsCompleted && !canCheck {
            HStack(spacing: 4) {
              Image(systemName: "exclamationmark.triangle.fill")
                .font(.caption)
                .foregroundColor(.orange)
              Text("Wash your hands first")
                .font(.caption)
                .foregroundColor(.orange)
            }
            .padding(.top, 4)
          }

          // Prerequisites warning (only show if wash hands is done)
          if viewModel.washHandsCompleted && hasPrerequisites && !allPrerequisitesCompleted {
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
        StepDetailsView(
          step: step,
          viewModel: viewModel,
          stepIndex: index,
          recipeID: recipeID,
          mealPlanSelections: mealPlanSelections
        )
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
}

// MARK: - Step Details View

// swiftlint:disable:next type_body_length
struct StepDetailsView: View {
  let step: Mealplanning_RecipeStep
  let viewModel: PerformRecipeViewModel
  let stepIndex: Int
  let recipeID: String
  var mealPlanSelections: [Mealplanning_MealPlanRecipeOptionSelection]?

  var body: some View {
    VStack(alignment: .leading, spacing: 8) {
      // Ingredients
      if !step.ingredients.isEmpty {
        let (regular, optionGroups) = groupIngredientsForStep(
          step.ingredients,
          stepID: step.id,
          stepIndex: stepIndex,
          recipeID: recipeID
        )

        StepItemsSectionView(
          title: "Ingredients",
          items: regular.map { ingredient in
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
          },
          ingredientOptionGroups: filterOptionGroups(optionGroups, for: .ingredient)
        )
      }

      // Instruments
      if !step.instruments.isEmpty {
        let (regular, optionGroups) = groupInstrumentsForStep(
          step.instruments,
          stepID: step.id,
          stepIndex: stepIndex,
          recipeID: recipeID
        )

        StepItemsSectionView(
          title: "Instruments",
          items: regular.map { instrument in
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
          },
          instrumentOptionGroups: filterInstrumentOptionGroups(optionGroups)
        )
      }

      // Vessels
      if !step.vessels.isEmpty {
        let (regular, optionGroups) = groupVesselsForStep(
          step.vessels,
          stepID: step.id,
          stepIndex: stepIndex,
          recipeID: recipeID
        )

        StepItemsSectionView(
          title: "Vessels",
          items: regular.map { vessel in
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
          },
          vesselOptionGroups: filterVesselOptionGroups(optionGroups)
        )
      }

      // Products
      if !step.products.isEmpty {
        StepProductsSectionView(products: step.products)
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

  // Helper functions to group items for a single step
  private func groupIngredientsForStep(
    _ ingredients: [Mealplanning_RecipeStepIngredient],
    stepID: String,
    stepIndex: Int,
    recipeID: String
  ) -> (
    regular: [Mealplanning_RecipeStepIngredient],
    optionGroups: [OptionGroupAggregate]
  ) {
    // Use the extension function from RecipePerformanceContentView
    // We need to access it through a temporary view extension
    var regular: [Mealplanning_RecipeStepIngredient] = []
    var optionGroupsByIndex: [UInt32: [Mealplanning_RecipeStepIngredient]] = [:]

    for ingredient in ingredients {
      if ingredient.hasIngredientIndex {
        let index = ingredient.ingredientIndex
        let hasOptions = ingredients.contains { other in
          other.id != ingredient.id && other.hasIngredientIndex && other.ingredientIndex == index
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
        let lhsIndex = lhs.hasOptionIndex ? lhs.optionIndex : 0
        let rhsIndex = rhs.hasOptionIndex ? rhs.optionIndex : 0
        return lhsIndex < rhsIndex
      }

      var options: [IngredientOption] = []
      for ingredient in sorted {
        let optionIndex = ingredient.hasOptionIndex ? ingredient.optionIndex : 0
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
          stepIndex: stepIndex,
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

  private func groupInstrumentsForStep(
    _ instruments: [Mealplanning_RecipeStepInstrument],
    stepID: String,
    stepIndex: Int,
    recipeID: String
  ) -> (
    regular: [Mealplanning_RecipeStepInstrument],
    optionGroups: [InstrumentOptionGroupAggregate]
  ) {
    var regular: [Mealplanning_RecipeStepInstrument] = []
    var optionGroupsByIndex: [UInt32: [Mealplanning_RecipeStepInstrument]] = [:]

    for instrument in instruments {
      if instrument.hasIndex {
        let index = instrument.index
        let hasOptions = instruments.contains { other in
          other.id != instrument.id && other.hasIndex && other.index == index
        }

        if hasOptions {
          if optionGroupsByIndex[index] == nil {
            optionGroupsByIndex[index] = []
          }
          optionGroupsByIndex[index]?.append(instrument)
        } else {
          regular.append(instrument)
        }
      } else {
        regular.append(instrument)
      }
    }

    var optionGroups: [InstrumentOptionGroupAggregate] = []
    for (index, groupInstruments) in optionGroupsByIndex {
      let sorted = groupInstruments.sorted { lhs, rhs in
        let lhsIndex = lhs.hasOptionIndex ? lhs.optionIndex : 0
        let rhsIndex = rhs.hasOptionIndex ? rhs.optionIndex : 0
        return lhsIndex < rhsIndex
      }

      var options: [InstrumentOption] = []
      for instrument in sorted {
        let optionIndex = instrument.hasOptionIndex ? instrument.optionIndex : 0
        let optionID = "\(stepID)-\(index)-\(optionIndex)"

        var aggregated = AggregatedInstrumentVessel(
          itemID: instrument.hasInstrument ? instrument.instrument.id : instrument.id,
          name: instrument.name,
          type: .instrument
        )

        if instrument.hasQuantity {
          aggregated.addQuantity(instrument.quantity)
        }

        options.append(
          InstrumentOption(
            id: optionID,
            instrument: instrument,
            optionIndex: optionIndex,
            aggregated: aggregated
          )
        )
      }

      let groupID = "\(recipeID)-\(stepID)-\(index)"
      optionGroups.append(
        InstrumentOptionGroupAggregate(
          id: groupID,
          recipeID: recipeID,
          stepID: stepID,
          stepIndex: stepIndex,
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

  private func groupVesselsForStep(
    _ vessels: [Mealplanning_RecipeStepVessel],
    stepID: String,
    stepIndex: Int,
    recipeID: String
  ) -> (
    regular: [Mealplanning_RecipeStepVessel],
    optionGroups: [VesselOptionGroupAggregate]
  ) {
    var regular: [Mealplanning_RecipeStepVessel] = []
    var optionGroupsByIndex: [UInt32: [Mealplanning_RecipeStepVessel]] = [:]

    for vessel in vessels {
      if vessel.hasIndex {
        let index = vessel.index
        let hasOptions = vessels.contains { other in
          other.id != vessel.id && other.hasIndex && other.index == index
        }

        if hasOptions {
          if optionGroupsByIndex[index] == nil {
            optionGroupsByIndex[index] = []
          }
          optionGroupsByIndex[index]?.append(vessel)
        } else {
          regular.append(vessel)
        }
      } else {
        regular.append(vessel)
      }
    }

    var optionGroups: [VesselOptionGroupAggregate] = []
    for (index, groupVessels) in optionGroupsByIndex {
      let sorted = groupVessels.sorted { lhs, rhs in
        let lhsIndex = lhs.hasOptionIndex ? lhs.optionIndex : 0
        let rhsIndex = rhs.hasOptionIndex ? rhs.optionIndex : 0
        return lhsIndex < rhsIndex
      }

      var options: [VesselOption] = []
      for vessel in sorted {
        let optionIndex = vessel.hasOptionIndex ? vessel.optionIndex : 0
        let optionID = "\(stepID)-\(index)-\(optionIndex)"

        var aggregated = AggregatedInstrumentVessel(
          itemID: vessel.hasVessel ? vessel.vessel.id : vessel.id,
          name: vessel.name,
          type: .vessel
        )

        if vessel.hasQuantity {
          aggregated.addQuantity(vessel.quantity)
        }

        options.append(
          VesselOption(
            id: optionID,
            vessel: vessel,
            optionIndex: optionIndex,
            aggregated: aggregated
          )
        )
      }

      let groupID = "\(recipeID)-\(stepID)-\(index)"
      optionGroups.append(
        VesselOptionGroupAggregate(
          id: groupID,
          recipeID: recipeID,
          stepID: stepID,
          stepIndex: stepIndex,
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

  // Filter option groups based on meal plan selections
  private func filterOptionGroups(
    _ groups: [OptionGroupAggregate],
    for selectionType: Mealplanning_MealPlanRecipeOptionSelectionType
  ) -> [OptionGroupAggregate] {
    guard let selections = mealPlanSelections else { return groups }

    return groups.compactMap { group in
      // Find selection for this group
      let selection = selections.first { sel in
        sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
          && sel.ingredientIndex == group.index && sel.selectionType == selectionType
      }

      if let selection = selection {
        // Only show the selected option
        let selectedOptions = group.options.filter {
          $0.optionIndex == selection.selectedOptionIndex
        }
        if !selectedOptions.isEmpty {
          var filteredGroup = group
          filteredGroup.options = selectedOptions
          filteredGroup.selectedOptionIndex = selection.selectedOptionIndex
          return filteredGroup
        }
        return nil
      }

      // No selection - show all options
      return group
    }
  }

  private func filterInstrumentOptionGroups(
    _ groups: [InstrumentOptionGroupAggregate]
  ) -> [InstrumentOptionGroupAggregate] {
    guard let selections = mealPlanSelections else { return groups }

    return groups.compactMap { group in
      let selection = selections.first { sel in
        sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
          && sel.selectionType == .instrument
      }

      if let selection = selection {
        let selectedOptions = group.options.filter {
          $0.optionIndex == selection.selectedOptionIndex
        }
        if !selectedOptions.isEmpty {
          var filteredGroup = group
          filteredGroup.options = selectedOptions
          filteredGroup.selectedOptionIndex = selection.selectedOptionIndex
          return filteredGroup
        }
        return nil
      }

      return group
    }
  }

  private func filterVesselOptionGroups(
    _ groups: [VesselOptionGroupAggregate]
  ) -> [VesselOptionGroupAggregate] {
    guard let selections = mealPlanSelections else { return groups }

    return groups.compactMap { group in
      let selection = selections.first { sel in
        sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
          && sel.selectionType == .vessel
      }

      if let selection = selection {
        let selectedOptions = group.options.filter {
          $0.optionIndex == selection.selectedOptionIndex
        }
        if !selectedOptions.isEmpty {
          var filteredGroup = group
          filteredGroup.options = selectedOptions
          filteredGroup.selectedOptionIndex = selection.selectedOptionIndex
          return filteredGroup
        }
        return nil
      }

      return group
    }
  }
}

// MARK: - Step Items Section View

struct StepItemsSectionView: View {
  let title: String
  let items: [StepItem]
  var ingredientOptionGroups: [OptionGroupAggregate] = []
  var instrumentOptionGroups: [InstrumentOptionGroupAggregate] = []
  var vesselOptionGroups: [VesselOptionGroupAggregate] = []

  var body: some View {
    VStack(alignment: .leading, spacing: 4) {
      Text(title)
        .font(.subheadline)
        .fontWeight(.semibold)
        .foregroundColor(.secondary)

      // Regular items
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

      // Ingredient option groups
      ForEach(ingredientOptionGroups) { group in
        OptionGroupView(group: group)
      }

      // Instrument option groups
      ForEach(instrumentOptionGroups) { group in
        InstrumentOptionGroupView(group: group)
      }

      // Vessel option groups
      ForEach(vesselOptionGroups) { group in
        VesselOptionGroupView(group: group)
      }
    }
  }
}

// MARK: - Option Group Views

struct OptionGroupView: View {
  let group: OptionGroupAggregate

  var body: some View {
    VStack(alignment: .leading, spacing: 4) {
      Text("one of:")
        .font(.caption)
        .foregroundColor(.secondary)
        .padding(.leading, 16)

      ForEach(group.options) { option in
        HStack(spacing: 6) {
          Text(option.ingredient.name)
            .font(.caption)
            .foregroundColor(.secondary)

          if let quantityText = option.aggregated.quantityText {
            Text(quantityText)
              .font(.caption)
              .foregroundColor(.secondary)
          }
        }
        .padding(.leading, 16)
      }
    }
  }
}

struct InstrumentOptionGroupView: View {
  let group: InstrumentOptionGroupAggregate

  var body: some View {
    VStack(alignment: .leading, spacing: 4) {
      Text("one of:")
        .font(.caption)
        .foregroundColor(.secondary)
        .padding(.leading, 16)

      ForEach(group.options) { option in
        HStack(spacing: 6) {
          Text(option.instrument.name)
            .font(.caption)
            .foregroundColor(.secondary)

          if let quantityText = option.aggregated.quantityText {
            Text(quantityText)
              .font(.caption)
              .foregroundColor(.secondary)
          }
        }
        .padding(.leading, 16)
      }
    }
  }
}

struct VesselOptionGroupView: View {
  let group: VesselOptionGroupAggregate

  var body: some View {
    VStack(alignment: .leading, spacing: 4) {
      Text("one of:")
        .font(.caption)
        .foregroundColor(.secondary)
        .padding(.leading, 16)

      ForEach(group.options) { option in
        HStack(spacing: 6) {
          Text(option.vessel.name)
            .font(.caption)
            .foregroundColor(.secondary)

          if let quantityText = option.aggregated.quantityText {
            Text(quantityText)
              .font(.caption)
              .foregroundColor(.secondary)
          }
        }
        .padding(.leading, 16)
      }
    }
  }
}

// MARK: - Step Products Section View

struct StepProductsSectionView: View {
  let products: [Mealplanning_RecipeStepProduct]

  var body: some View {
    VStack(alignment: .leading, spacing: 4) {
      Text("Products")
        .font(.subheadline)
        .fontWeight(.semibold)
        .foregroundColor(.secondary)

      ForEach(Array(products.enumerated()), id: \.offset) { _, product in
        HStack(spacing: 6) {
          Text(formatProductQuantity(product))
            .font(.caption)
            .foregroundColor(.secondary)
        }
      }
    }
  }

  private func formatProductQuantity(_ product: Mealplanning_RecipeStepProduct) -> String {
    // Check if product is discrete (has ItemQuantity set)
    let isDiscrete =
      product.hasItemQuantity && (product.itemQuantity.hasMin || product.itemQuantity.hasMax)

    if isDiscrete {
      // Discrete product: Format as "4 patties (4 oz each)"
      var itemQtyStr = ""
      if product.itemQuantity.hasMin {
        let min = product.itemQuantity.min
        if product.itemQuantity.hasMax {
          let max = product.itemQuantity.max
          if min == max {
            itemQtyStr = formatQuantity(min)
          } else {
            itemQtyStr = "\(formatQuantity(min))-\(formatQuantity(max))"
          }
        } else {
          itemQtyStr = formatQuantity(min)
        }
      }

      var measurementQtyStr = ""
      if product.hasMeasurementQuantity && product.measurementQuantity.hasMin {
        let min = product.measurementQuantity.min
        if product.measurementQuantity.hasMax {
          let max = product.measurementQuantity.max
          if min == max {
            measurementQtyStr = formatQuantity(min)
          } else {
            measurementQtyStr = "\(formatQuantity(min))-\(formatQuantity(max))"
          }
        } else {
          measurementQtyStr = formatQuantity(min)
        }
      }

      let unitName = product.hasMeasurementUnit ? product.measurementUnit.name : ""

      if !itemQtyStr.isEmpty && !measurementQtyStr.isEmpty && !unitName.isEmpty {
        // Format: "4 patties (4 oz each)"
        return "\(itemQtyStr) \(product.name) (\(measurementQtyStr) \(unitName) each)"
      } else if !itemQtyStr.isEmpty {
        // Fallback: just show count and name if measurement is missing
        return "\(itemQtyStr) \(product.name)"
      } else {
        // Fallback: just show name if quantities are missing
        return product.name
      }
    } else if product.hasMeasurementQuantity && product.measurementQuantity.hasMin {
      // Continuous product: Format as "product name: 16 oz"
      let min = product.measurementQuantity.min
      var qtyStr = formatQuantity(min)

      if product.measurementQuantity.hasMax {
        let max = product.measurementQuantity.max
        if min != max {
          qtyStr = "\(qtyStr)-\(formatQuantity(max))"
        }
      }

      let unitName = product.hasMeasurementUnit ? product.measurementUnit.name : ""
      if !unitName.isEmpty {
        return "\(product.name): \(qtyStr) \(unitName)"
      } else {
        return "\(product.name): \(qtyStr)"
      }
    }

    // Fallback: just show name if no quantities
    return product.name
  }

  private func formatQuantity(_ qty: Float) -> String {
    // Format numbers - use fewer decimals for whole numbers
    if qty.truncatingRemainder(dividingBy: 1) == 0 {
      return String(format: "%.0f", qty)
    } else {
      return String(format: "%.2f", qty)
    }
  }
}
