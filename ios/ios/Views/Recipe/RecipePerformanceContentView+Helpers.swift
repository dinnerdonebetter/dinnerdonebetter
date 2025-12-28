//
//  RecipePerformanceContentView+Helpers.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

// MARK: - Helper Types

struct StepItem {
  let name: String
  let isProduct: Bool
  let prerequisiteStepIndex: Int?
  let prerequisiteCompleted: Bool
}

struct InstrumentVesselItem: Identifiable {
  let id: String
  let name: String
  let type: ItemType

  enum ItemType {
    case instrument
    case vessel
  }
}

struct AggregatedInstrumentVessel: Identifiable {
  let itemID: String
  let name: String
  let type: InstrumentVesselItem.ItemType

  private var totalMin: UInt32 = 0
  private var totalMax: UInt32?
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

struct AggregatedIngredient: Identifiable {
  let ingredientID: String
  let name: String
  let quantityNotes: String
  let measurementUnit: Mealplanning_ValidMeasurementUnit?

  private var totalMin: Float = 0
  private var totalMax: Float?
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

// MARK: - Helper Functions

extension RecipePerformanceContentView {
  func getAggregatedInstrumentsAndVessels(from recipe: Mealplanning_Recipe)
    -> [AggregatedInstrumentVessel]
  {
    var aggregated: [String: AggregatedInstrumentVessel] = [:]

    for step in recipe.steps {
      // Collect instruments (only if it has a ValidInstrument and displayInSummaryLists is true)
      for instrument in step.instruments where instrument.hasInstrument {
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
            if instrument.hasQuantity, var current = aggregated[itemID] {
              let quantity = instrument.quantity
              current.addQuantity(quantity)
              aggregated[itemID] = current
            }
          }
        }
      }

      // Collect vessels (only if it has a ValidVessel and displayInSummaryLists is true)
      for vessel in step.vessels where vessel.hasVessel {
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
            if vessel.hasQuantity, var current = aggregated[itemID] {
              let quantity = vessel.quantity
              current.addQuantity(quantity)
              aggregated[itemID] = current
            }
          }
        }
      }
    }

    return Array(aggregated.values).sorted { $0.name < $1.name }
  }

  func getAggregatedIngredients(from recipe: Mealplanning_Recipe) -> [AggregatedIngredient] {
    var aggregated: [String: AggregatedIngredient] = [:]

    for step in recipe.steps {
      for ingredient in step.ingredients where ingredient.hasIngredient {
        // Only include if it has a ValidIngredient (not a recipe step product)
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
          if ingredient.hasQuantity, var current = aggregated[key] {
            let quantity = ingredient.quantity
            current.addQuantity(quantity)
            aggregated[key] = current
          }
        }
      }
    }

    return Array(aggregated.values).sorted { $0.name < $1.name }
  }

  func formatStepTitle(step: Mealplanning_RecipeStep, viewModel: PerformRecipeViewModel) -> String {
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

  func formatList(_ items: [String]) -> String {
    guard !items.isEmpty else { return "" }

    if items.count == 1 {
      return items[0]
    } else if items.count == 2 {
      return "\(items[0]) and \(items[1])"
    } else {
      let allButLast = items.dropLast().joined(separator: ", ")
      if let last = items.last {
        return "\(allButLast), and \(last)"
      }
      return allButLast
    }
  }
}
