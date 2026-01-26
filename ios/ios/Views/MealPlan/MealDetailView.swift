//
//  MealDetailView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct MealDetailView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @State private var viewModel: MealDetailViewModel?
  @State private var loadedRecipes: [String: (recipe: Mealplanning_Recipe, scale: Float)] = [:]
  @State private var baseComponentScales: [String: Float] = [:]  // Maps component recipe ID to base scale from meal
  @State private var mealScale: Float = 1.0  // Meal-level scale multiplier
  @State private var mealScaleText: String = "1.0"
  @FocusState private var isMealScaleFocused: Bool
  @State private var isInstrumentsVesselsExpanded = false
  @State private var isIngredientsExpanded = false
  @State private var checkedIngredients: Set<String> = []
  @State private var checkedInstrumentsVessels: Set<String> = []

  let mealID: String

  init(mealID: String) {
    self.mealID = mealID
  }

  var body: some View {
    Group {
      if let viewModel = viewModel {
        if viewModel.isLoading {
          ProgressView("Loading meal...")
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
                await viewModel.loadMeal()
              }
            }
            .buttonStyle(.borderedProminent)
          }
          .frame(maxWidth: .infinity, maxHeight: .infinity)
        } else if let meal = viewModel.meal {
          ScrollView {
            VStack(alignment: .leading, spacing: 20) {
              // Overall Info Section
              overallInfoSection(meal: meal)

              // Aggregated Ingredients & Instruments/Vessels
              if !loadedRecipes.isEmpty {
                aggregatedListsSection
              }

              // Components
              if !meal.components.isEmpty {
                componentsSection(meal: meal)
              }
            }
            .padding()
          }
        } else {
          Text("Meal not found")
            .foregroundColor(.secondary)
            .frame(maxWidth: .infinity, maxHeight: .infinity)
        }
      } else {
        ProgressView("Initializing...")
          .frame(maxWidth: .infinity, maxHeight: .infinity)
      }
    }
    .navigationTitle(viewModel?.meal?.name ?? "Meal")
    .navigationBarTitleDisplayMode(.large)
    .onAppear {
      if viewModel == nil {
        viewModel = MealDetailViewModel(mealID: mealID, authManager: authManager)
      }
      if let viewModel = viewModel {
        Task {
          await viewModel.loadMeal()
          // Initialize base scales from meal components after meal loads
          if let meal = viewModel.meal, baseComponentScales.isEmpty {
            for component in meal.components {
              baseComponentScales[component.recipe.id] = component.recipeScale
            }
          }
        }
      }
    }
  }

  // MARK: - Overall Info Section

  private func overallInfoSection(meal: Mealplanning_Meal) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      // Description
      if !meal.description_p.isEmpty {
        Text(meal.description_p)
          .font(.subheadline)
          .foregroundColor(.secondary)
      }

      // Estimated portions
      if meal.hasEstimatedPortions {
        HStack(spacing: 8) {
          Image(systemName: "person.2")
            .foregroundColor(.secondary)
          Text("Estimated Portions: \(formatPortions(meal.estimatedPortions))")
            .font(.subheadline)
            .foregroundColor(.secondary)
        }
      }

      // Meal Scale Control
      Divider()
        .padding(.vertical, 4)

      HStack(spacing: 12) {
        Text("Meal Scale:")
          .font(.subheadline)
          .fontWeight(.medium)

        HStack(spacing: 8) {
          TextField("1.0", text: $mealScaleText)
            .keyboardType(.decimalPad)
            .textFieldStyle(.roundedBorder)
            .frame(width: 80)
            .focused($isMealScaleFocused)
            .onSubmit {
              updateMealScaleFromText()
            }
            .onChange(of: isMealScaleFocused) { _, isFocused in
              if !isFocused {
                updateMealScaleFromText()
              }
            }
            .onChange(of: mealScaleText) { _, newValue in
              // Filter to only allow numbers and a single decimal point
              var filtered = newValue.filter { $0.isNumber || $0 == "." }
              // Ensure only one decimal point
              let parts = filtered.split(separator: ".", omittingEmptySubsequences: false)
              if parts.count > 2 {
                filtered = parts[0] + "." + parts.dropFirst().joined()
              }
              if filtered != newValue {
                mealScaleText = filtered
              }
            }

          Text("x")
            .font(.subheadline)
            .foregroundColor(.secondary)

          // Quick scale buttons
          HStack(spacing: 4) {
            Button("0.5x") {
              setMealScale(0.5)
            }
            .buttonStyle(.bordered)
            .controlSize(.small)

            Button("1x") {
              setMealScale(1.0)
            }
            .buttonStyle(.bordered)
            .controlSize(.small)

            Button("2x") {
              setMealScale(2.0)
            }
            .buttonStyle(.bordered)
            .controlSize(.small)
          }
        }
      }
    }
    .padding()
    .frame(maxWidth: .infinity, alignment: .leading)
    .background(Color(.systemGray6))
    .cornerRadius(12)
  }

  private func updateMealScaleFromText() {
    if let scale = Float(mealScaleText), scale > 0 {
      setMealScale(scale)
    } else {
      mealScaleText = String(format: "%.2f", mealScale)
    }
  }

  private func setMealScale(_ newScale: Float) {
    mealScale = newScale
    mealScaleText = String(format: "%.2f", mealScale)
    // Update loadedRecipes with new scales
    for (recipeID, baseScale) in baseComponentScales {
      let newScale = baseScale * mealScale
      if let recipeData = loadedRecipes[recipeID] {
        loadedRecipes[recipeID] = (recipe: recipeData.recipe, scale: newScale)
      }
    }
  }

  private func formatPortions(_ range: Common_Float32RangeWithOptionalMax) -> String {
    if range.hasMax {
      if range.min == range.max {
        return String(format: "%.1f", range.min)
      } else {
        return String(format: "%.1f-%.1f", range.min, range.max)
      }
    } else {
      return String(format: "%.1f+", range.min)
    }
  }

  private func formatComponentType(_ type: Mealplanning_MealComponentType) -> String {
    switch type {
    case .amuseBouche:
      return "Amuse Bouche"
    case .appetizer:
      return "Appetizer"
    case .soup:
      return "Soup"
    case .main:
      return "Main"
    case .salad:
      return "Salad"
    case .beverage:
      return "Beverage"
    case .side:
      return "Side"
    case .dessert:
      return "Dessert"
    default:
      return ""
    }
  }

  private var aggregatedListsSection: some View {
    VStack(alignment: .leading, spacing: 12) {
      // Instruments & Vessels
      aggregatedInstrumentsVesselsSection

      // Ingredients
      aggregatedIngredientsSection
    }
  }

  private var aggregatedInstrumentsVesselsSection: some View {
    let aggregatedItems = getAggregatedInstrumentsAndVessels()

    return VStack(alignment: .leading, spacing: 0) {
      Button(
        action: {
          withAnimation {
            isInstrumentsVesselsExpanded.toggle()
          }
        },
        label: {
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
      )
      .buttonStyle(.plain)

      if isInstrumentsVesselsExpanded && !aggregatedItems.isEmpty {
        VStack(alignment: .leading, spacing: 8) {
          ForEach(aggregatedItems, id: \.itemID) { item in
            HStack(spacing: 12) {
              Button(
                action: {
                  if checkedInstrumentsVessels.contains(item.itemID) {
                    checkedInstrumentsVessels.remove(item.itemID)
                  } else {
                    checkedInstrumentsVessels.insert(item.itemID)
                  }
                },
                label: {
                  Image(
                    systemName: checkedInstrumentsVessels.contains(item.itemID)
                      ? "checkmark.circle.fill" : "circle"
                  )
                  .font(.title3)
                  .foregroundColor(
                    checkedInstrumentsVessels.contains(item.itemID) ? .green : .gray
                  )
                }
              )
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

  private var aggregatedIngredientsSection: some View {
    let aggregatedIngredients = getAggregatedIngredients()

    return VStack(alignment: .leading, spacing: 0) {
      Button(
        action: {
          withAnimation {
            isIngredientsExpanded.toggle()
          }
        },
        label: {
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
      )
      .buttonStyle(.plain)

      if isIngredientsExpanded && !aggregatedIngredients.isEmpty {
        VStack(alignment: .leading, spacing: 8) {
          ForEach(aggregatedIngredients, id: \.ingredientID) { aggregated in
            HStack(spacing: 12) {
              Button(
                action: {
                  if checkedIngredients.contains(aggregated.ingredientID) {
                    checkedIngredients.remove(aggregated.ingredientID)
                  } else {
                    checkedIngredients.insert(aggregated.ingredientID)
                  }
                },
                label: {
                  Image(
                    systemName: checkedIngredients.contains(aggregated.ingredientID)
                      ? "checkmark.circle.fill" : "circle"
                  )
                  .font(.title3)
                  .foregroundColor(
                    checkedIngredients.contains(aggregated.ingredientID) ? .green : .gray
                  )
                }
              )
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

}

// MARK: - Meal Detail View Extensions

extension MealDetailView {
  func getAggregatedInstrumentsAndVessels() -> [AggregatedInstrumentVessel] {
    var aggregated: [String: AggregatedInstrumentVessel] = [:]

    for (_, recipeData) in loadedRecipes {
      let recipe = recipeData.recipe
      let scale = recipeData.scale

      for step in recipe.steps {
        for instrument in step.instruments where instrument.hasInstrument {
          let validInstrument = instrument.instrument
          if validInstrument.displayInSummaryLists {
            let itemID = validInstrument.id
            if !itemID.isEmpty {
              if aggregated[itemID] == nil {
                aggregated[itemID] = AggregatedInstrumentVessel(
                  itemID: itemID,
                  name: instrument.name,
                  type: .instrument
                )
              }

              if instrument.hasQuantity, var current = aggregated[itemID] {
                var scaledQuantity = instrument.quantity
                scaledQuantity.min = UInt32(Float(scaledQuantity.min) * scale)
                if scaledQuantity.hasMax {
                  scaledQuantity.max = UInt32(Float(scaledQuantity.max) * scale)
                }
                current.addQuantity(scaledQuantity)
                aggregated[itemID] = current
              }
            }
          }
        }

        for vessel in step.vessels where vessel.hasVessel {
          let validVessel = vessel.vessel
          if validVessel.displayInSummaryLists {
            let itemID = validVessel.id
            if !itemID.isEmpty {
              if aggregated[itemID] == nil {
                aggregated[itemID] = AggregatedInstrumentVessel(
                  itemID: itemID,
                  name: vessel.name,
                  type: .vessel
                )
              }

              if vessel.hasQuantity, var current = aggregated[itemID] {
                var scaledQuantity = vessel.quantity
                scaledQuantity.min = UInt32(Float(scaledQuantity.min) * scale)
                if scaledQuantity.hasMax {
                  scaledQuantity.max = UInt32(Float(scaledQuantity.max) * scale)
                }
                current.addQuantity(scaledQuantity)
                aggregated[itemID] = current
              }
            }
          }
        }
      }
    }

    return Array(aggregated.values).sorted { $0.name < $1.name }
  }

  func getAggregatedIngredients() -> [AggregatedIngredient] {
    var aggregated: [String: AggregatedIngredient] = [:]

    for (_, recipeData) in loadedRecipes {
      let recipe = recipeData.recipe
      let scale = recipeData.scale

      for step in recipe.steps {
        for ingredient in step.ingredients where ingredient.hasIngredient {
          let validIngredient = ingredient.ingredient
          let key = validIngredient.id
          if !key.isEmpty {
            if aggregated[key] == nil {
              aggregated[key] = AggregatedIngredient(
                ingredientID: key,
                name: ingredient.name,
                quantityNotes: ingredient.quantityNotes,
                measurementUnit: ingredient.hasMeasurementUnit ? ingredient.measurementUnit : nil
              )
            }

            if ingredient.hasQuantity, var current = aggregated[key] {
              var scaledQuantity = ingredient.quantity
              scaledQuantity.min *= scale
              if scaledQuantity.hasMax {
                scaledQuantity.max *= scale
              }
              current.addQuantity(scaledQuantity)
              aggregated[key] = current
            }
          }
        }
      }
    }

    return Array(aggregated.values).sorted { $0.name < $1.name }
  }

  private func componentsSection(meal: Mealplanning_Meal) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      Text("Components")
        .font(.title2)
        .fontWeight(.bold)

      ForEach(meal.components, id: \.recipe.id) { component in
        EmbeddedRecipeView(
          recipeID: component.recipe.id,
          recipeScale: Binding(
            get: {
              let baseScale = baseComponentScales[component.recipe.id] ?? component.recipeScale
              return baseScale * mealScale
            },
            set: { _ in
              // Scale is controlled at meal level, so we don't allow individual changes
            }
          ),
          componentType: component.componentType,
          onRecipeLoaded: { recipe in
            Task { @MainActor in
              let baseScale = baseComponentScales[component.recipe.id] ?? component.recipeScale
              let currentScale = baseScale * mealScale
              loadedRecipes[component.recipe.id] = (recipe: recipe, scale: currentScale)
            }
          }
        )
      }
    }
    .onAppear {
      // Initialize base scales from meal components if not already set
      if baseComponentScales.isEmpty {
        for component in meal.components {
          baseComponentScales[component.recipe.id] = component.recipeScale
        }
      }
    }
    .onChange(of: mealScale) { _, _ in
      // Update loadedRecipes when meal scale changes
      for (recipeID, baseScale) in baseComponentScales {
        let newScale = baseScale * mealScale
        if let recipeData = loadedRecipes[recipeID] {
          loadedRecipes[recipeID] = (recipe: recipeData.recipe, scale: newScale)
        }
      }
    }
  }
}

// MARK: - Embedded Recipe View

struct EmbeddedRecipeView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @State private var viewModel: PerformRecipeViewModel?
  @State private var isInstrumentsVesselsExpanded = false
  @State private var isIngredientsExpanded = false
  @State private var checkedIngredients: Set<String> = []
  @State private var checkedInstrumentsVessels: Set<String> = []
  @State private var isExpanded = false

  let recipeID: String
  @Binding var recipeScale: Float
  let componentType: Mealplanning_MealComponentType
  let onRecipeLoaded: ((Mealplanning_Recipe) -> Void)?

  var body: some View {
    VStack(alignment: .leading, spacing: 0) {
      // Collapsible header
      VStack(alignment: .leading, spacing: 0) {
        Button(
          action: {
            withAnimation {
              isExpanded.toggle()
            }
          },
          label: {
            HStack {
              // Component type badge
              if componentType != .unspecified {
                Text(formatComponentType(componentType))
                  .font(.caption)
                  .fontWeight(.semibold)
                  .padding(.horizontal, 8)
                  .padding(.vertical, 4)
                  .background(Color.blue.opacity(0.2))
                  .foregroundColor(.blue)
                  .cornerRadius(6)
              }

              // Recipe name
              if let recipe = viewModel?.recipe {
                Text(recipe.name.isEmpty ? "Unnamed Recipe" : recipe.name)
                  .font(.subheadline)
                  .fontWeight(.medium)
                  .foregroundColor(.primary)
              } else {
                Text("Loading...")
                  .font(.subheadline)
                  .foregroundColor(.secondary)
              }

              Spacer()

              // Chevron
              Image(systemName: isExpanded ? "chevron.down" : "chevron.right")
                .font(.caption)
                .foregroundColor(.secondary)
            }
            .padding()
            .background(Color(.systemGray6))
          }
        )
        .buttonStyle(.plain)
      }

      // Recipe content (collapsible)
      if isExpanded {
        if let viewModel = viewModel {
          if viewModel.isLoading {
            ProgressView("Loading recipe...")
              .frame(maxWidth: .infinity)
              .padding()
          } else if let errorMessage = viewModel.errorMessage {
            VStack(spacing: 8) {
              Image(systemName: "exclamationmark.triangle")
                .foregroundColor(.orange)
              Text("Error loading recipe: \(errorMessage)")
                .font(.caption)
                .foregroundColor(.secondary)
            }
            .frame(maxWidth: .infinity)
            .padding()
          } else if let recipe = viewModel.recipe {
            RecipePerformanceContentView(
              checkedIngredients: $checkedIngredients,
              checkedInstrumentsVessels: $checkedInstrumentsVessels,
              isInstrumentsVesselsExpanded: $isInstrumentsVesselsExpanded,
              isIngredientsExpanded: $isIngredientsExpanded,
              recipe: recipe,
              viewModel: viewModel,
              hideIngredientsAndInstruments: true,
              externalScale: Binding(
                get: { recipeScale },
                set: { newValue in
                  if let scale = newValue {
                    recipeScale = scale
                  }
                }
              )
            )
            .onAppear {
              onRecipeLoaded?(recipe)
            }
          }
        } else {
          ProgressView("Initializing...")
            .frame(maxWidth: .infinity)
            .padding()
        }
      }
    }
    .background(Color(.systemGray6))
    .cornerRadius(10)
    .onAppear {
      if viewModel == nil {
        viewModel = PerformRecipeViewModel(recipeID: recipeID, authManager: authManager)
        Task {
          await viewModel?.loadRecipe()
        }
      }
    }
  }

  private func formatComponentType(_ type: Mealplanning_MealComponentType) -> String {
    switch type {
    case .amuseBouche:
      return "Amuse Bouche"
    case .appetizer:
      return "Appetizer"
    case .soup:
      return "Soup"
    case .main:
      return "Main"
    case .salad:
      return "Salad"
    case .beverage:
      return "Beverage"
    case .side:
      return "Side"
    case .dessert:
      return "Dessert"
    default:
      return ""
    }
  }
}

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "John Doe"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  return NavigationView {
    MealDetailView(mealID: "meal123")
  }
  .environment(authManager)
}
