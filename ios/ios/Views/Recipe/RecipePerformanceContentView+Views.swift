//
//  RecipePerformanceContentView+Views.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

// MARK: - Step Card View

struct StepCardView: View {
  let step: Mealplanning_RecipeStep
  let index: Int
  let viewModel: PerformRecipeViewModel
  let formatStepTitle: (Mealplanning_RecipeStep, PerformRecipeViewModel) -> String

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
        StepDetailsView(step: step, viewModel: viewModel, stepIndex: index)
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

struct StepDetailsView: View {
  let step: Mealplanning_RecipeStep
  let viewModel: PerformRecipeViewModel
  let stepIndex: Int

  var body: some View {
    VStack(alignment: .leading, spacing: 8) {
      // Ingredients
      if !step.ingredients.isEmpty {
        StepItemsSectionView(
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
        StepItemsSectionView(
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
        StepItemsSectionView(
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
}

// MARK: - Step Items Section View

struct StepItemsSectionView: View {
  let title: String
  let items: [StepItem]

  var body: some View {
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
