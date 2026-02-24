//
//  MealDetailView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct MealDetailView: View {
  private enum MealFlowMode: String, CaseIterable, Identifiable {
    case unified = "Unified"
    case byComponent = "By Component"

    var id: String { rawValue }
  }

  private static let componentTagColors: [Color] = [
    .indigo, .green, .yellow,
    .purple, .orange, .blue,
    .mint, .red, .teal,
    .pink, .cyan, .brown,
    Color(red: 0.85, green: 0.20, blue: 0.20),
    Color(red: 0.90, green: 0.45, blue: 0.10),
    Color(red: 0.75, green: 0.65, blue: 0.15),
    Color(red: 0.25, green: 0.65, blue: 0.25),
    Color(red: 0.15, green: 0.70, blue: 0.55),
    Color(red: 0.15, green: 0.65, blue: 0.70),
    Color(red: 0.20, green: 0.45, blue: 0.85),
    Color(red: 0.30, green: 0.35, blue: 0.85),
    Color(red: 0.45, green: 0.25, blue: 0.80),
    Color(red: 0.70, green: 0.20, blue: 0.70),
    Color(red: 0.85, green: 0.20, blue: 0.50),
    Color(red: 0.60, green: 0.35, blue: 0.20),
    Color(red: 0.55, green: 0.20, blue: 0.20),
    Color(red: 0.60, green: 0.35, blue: 0.10),
    Color(red: 0.55, green: 0.50, blue: 0.10),
    Color(red: 0.20, green: 0.55, blue: 0.20),
    Color(red: 0.10, green: 0.55, blue: 0.45),
    Color(red: 0.15, green: 0.45, blue: 0.70),
    Color(red: 0.25, green: 0.25, blue: 0.65),
    Color(red: 0.55, green: 0.20, blue: 0.60),
  ]

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
  @State private var mealWashHandsCompleted = false
  @State private var componentViewModels: [String: PerformRecipeViewModel] = [:]
  @State private var mealFlowMode: MealFlowMode = .unified
  @State private var showUnifiedCompletedSteps = false

  let mealID: String

  init(mealID: String) {
    self.mealID = mealID
  }

  var body: some View {
    Group {
      if let viewModel = viewModel {
        DSContentState(
          isLoading: viewModel.isLoading,
          loadingMessage: "Loading meal...",
          error: viewModel.errorMessage,
          onRetry: { await viewModel.loadMeal() },
          content: {
            if let meal = viewModel.meal {
              ScrollView {
                VStack(alignment: .leading, spacing: DSTheme.Spacing.xl) {
                  // Overall Info Section
                  overallInfoSection(meal: meal)

                  // Aggregated Ingredients & Instruments/Vessels (consolidated from all recipes)
                  if !meal.components.isEmpty {
                    aggregatedListsSection
                  }

                  if !meal.components.isEmpty && !mealWashHandsCompleted {
                    mealWashHandsSection
                  }

                  if !meal.components.isEmpty {
                    mealFlowModeSection
                    if mealFlowMode == .unified {
                      unifiedMealStepsSection(meal: meal)
                    } else {
                      componentsSection(meal: meal)
                    }
                  }
                }
                .dsScreenPadding()
              }
            } else {
              DSEmptyState(
                icon: "fork.knife",
                title: "Meal not found",
                message: "This meal could not be loaded."
              )
            }
          })
      } else {
        DSInitializingView()
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
          if let meal = viewModel.meal {
            ensureComponentDataLoaded(for: meal)
          }
        }
      }
    }
    .onChange(of: viewModel?.meal?.id) { _, _ in
      if let meal = viewModel?.meal {
        ensureComponentDataLoaded(for: meal)
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

          Button {
            adjustMealScale(by: -0.25)
          } label: {
            Image(systemName: "minus.circle")
          }
          .buttonStyle(.plain)

          Button {
            adjustMealScale(by: 0.25)
          } label: {
            Image(systemName: "plus.circle")
          }
          .buttonStyle(.plain)
        }
      }

      Slider(
        value: Binding(
          get: { Double(mealScale) },
          set: { setMealScale(Float($0)) }
        ),
        in: 0.25...4.0,
        step: 0.25
      )
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
    mealScale = min(max(newScale, 0.25), 4.0)
    mealScaleText = String(format: "%.2f", mealScale)
    // Update loadedRecipes with new scales
    for (recipeID, baseScale) in baseComponentScales {
      let newScale = baseScale * mealScale
      if let recipeData = loadedRecipes[recipeID] {
        loadedRecipes[recipeID] = (recipe: recipeData.recipe, scale: newScale)
      }
    }
  }

  private func adjustMealScale(by delta: Float) {
    setMealScale(mealScale + delta)
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

  private var mealWashHandsSection: some View {
    VStack(alignment: .leading, spacing: 8) {
      HStack(alignment: .top, spacing: 12) {
        Button(
          action: {
            setMealWashHandsCompleted(!mealWashHandsCompleted)
          },
          label: {
            Image(systemName: mealWashHandsCompleted ? "checkmark.circle.fill" : "circle")
              .font(.title2)
              .foregroundColor(mealWashHandsCompleted ? .green : .blue)
          }
        )
        .buttonStyle(.plain)

        VStack(alignment: .leading, spacing: 4) {
          Text("🧼 Wash your hands")
            .font(.headline)
            .foregroundColor(mealWashHandsCompleted ? .secondary : .primary)
            .italic(mealWashHandsCompleted)

          Text("Complete this once to unlock all component steps.")
            .font(.caption)
            .foregroundColor(.secondary)
        }

        Spacer()
      }
    }
    .padding()
    .background(mealWashHandsCompleted ? Color(.systemGray6) : Color(.systemBackground))
    .cornerRadius(12)
    .overlay(
      RoundedRectangle(cornerRadius: 12)
        .stroke(mealWashHandsCompleted ? Color.green.opacity(0.3) : Color.clear, lineWidth: 2)
    )
  }

  private var mealFlowModeSection: some View {
    VStack(alignment: .leading, spacing: 8) {
      Text("Step View")
        .font(.subheadline)
        .fontWeight(.semibold)
        .foregroundColor(.secondary)

      Picker("Step View", selection: $mealFlowMode) {
        ForEach(MealFlowMode.allCases) { mode in
          Text(mode.rawValue).tag(mode)
        }
      }
      .pickerStyle(.segmented)
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
            Text("Equipment")
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

  private struct UnifiedMealStep: Identifiable {
    enum Category {
      case upNext
      case forLater
      case done
    }

    let componentID: String
    let componentIndex: Int
    let componentName: String
    let step: Mealplanning_RecipeStep
    let stepIndex: Int
    let recipeID: String
    let scale: Float
    let viewModel: PerformRecipeViewModel
    let category: Category

    var id: String {
      "\(componentID):\(recipeID):\(step.id)"
    }
  }

  private func collectUnifiedMealSteps(meal: Mealplanning_Meal) -> [UnifiedMealStep] {
    var result: [UnifiedMealStep] = []

    for (componentIndex, component) in meal.components.enumerated() {
      let componentID = component.recipe.id
      guard let viewModel = componentViewModels[componentID], let recipe = viewModel.recipe else {
        continue
      }

      for (index, step) in recipe.steps.enumerated() {
        let category: UnifiedMealStep.Category
        switch viewModel.categorizeStep(recipeID: recipe.id, stepID: step.id) {
        case .upNext: category = .upNext
        case .forLater: category = .forLater
        case .done: category = .done
        }
        result.append(
          UnifiedMealStep(
            componentID: componentID,
            componentIndex: componentIndex,
            componentName: recipe.name.isEmpty
              ? formatComponentType(component.componentType) : recipe.name,
            step: step,
            stepIndex: index,
            recipeID: recipe.id,
            scale: loadedRecipes[componentID]?.scale ?? (baseComponentScales[componentID] ?? 1.0)
              * mealScale,
            viewModel: viewModel,
            category: category
          ))
      }
    }

    return result
  }

  private func unifiedMealStepsSection(meal: Mealplanning_Meal) -> some View {
    let allSteps = collectUnifiedMealSteps(meal: meal)
    let upNext = allSteps.filter { $0.category == .upNext }
    let forLater = allSteps.filter { $0.category == .forLater }

    return VStack(alignment: .leading, spacing: 12) {
      HStack {
        Text("Cook Flow")
          .font(.title2)
          .fontWeight(.bold)

        Spacer()

        Button(showUnifiedCompletedSteps ? "Focus mode" : "Overview") {
          withAnimation {
            showUnifiedCompletedSteps.toggle()
          }
        }
        .font(.caption)
        .buttonStyle(.bordered)
        .controlSize(.small)
      }

      if showUnifiedCompletedSteps {
        unifiedMealStepGroup(title: "All Steps", color: .secondary, steps: allSteps)
      } else {
        if !upNext.isEmpty {
          unifiedMealStepGroup(title: "Up Next", color: .orange, steps: upNext)
        }
        if !forLater.isEmpty {
          unifiedMealStepGroup(title: "Not Yet", color: .blue, steps: forLater)
        }
      }

      if allSteps.isEmpty {
        Text("Loading component steps...")
          .font(.caption)
          .foregroundColor(.secondary)
      }
    }
  }

  private func unifiedMealStepGroup(
    title: String,
    color: Color,
    steps: [UnifiedMealStep]
  ) -> some View {
    VStack(alignment: .leading, spacing: 8) {
      Text(title)
        .font(.subheadline)
        .fontWeight(.semibold)
        .foregroundColor(color)

      ForEach(steps) { item in
        VStack(alignment: .leading, spacing: 6) {
          let tagStyle = componentTagStyle(for: item.componentIndex)
          Text(item.componentName)
            .font(.caption2)
            .padding(.horizontal, 8)
            .padding(.vertical, 3)
            .background(tagStyle.background)
            .foregroundColor(tagStyle.foreground)
            .cornerRadius(6)

          StepCardView(
            step: item.step,
            index: item.stepIndex,
            viewModel: item.viewModel,
            formatStepTitle: { step, _ in
              formatUnifiedStepTitle(step: step, index: item.stepIndex)
            },
            recipeID: item.recipeID,
            scale: item.scale
          )
        }
      }
    }
  }

  private func formatUnifiedStepTitle(step: Mealplanning_RecipeStep, index: Int) -> String {
    if step.hasPreparation, !step.preparation.name.isEmpty {
      return step.preparation.name
    }
    return "Step \(index + 1)"
  }

  private func componentTagStyle(for componentIndex: Int) -> (background: Color, foreground: Color)
  {
    let color = Self.componentTagColors[componentIndex % Self.componentTagColors.count]
    return (color.opacity(0.16), color)
  }

  private func ensureComponentDataLoaded(for meal: Mealplanning_Meal) {
    for component in meal.components {
      let recipeID = component.recipe.id
      let baseScale = baseComponentScales[recipeID] ?? component.recipeScale

      if let existingViewModel = componentViewModels[recipeID] {
        if let recipe = existingViewModel.recipe, loadedRecipes[recipeID] == nil {
          loadedRecipes[recipeID] = (recipe: recipe, scale: baseScale * mealScale)
        }
        continue
      }

      let recipeViewModel = PerformRecipeViewModel(recipeID: recipeID, authManager: authManager)
      recipeViewModel.washHandsCompleted = mealWashHandsCompleted
      componentViewModels[recipeID] = recipeViewModel

      Task { @MainActor in
        await recipeViewModel.loadRecipe()
        if let recipe = recipeViewModel.recipe {
          loadedRecipes[recipeID] = (recipe: recipe, scale: baseScale * mealScale)
        }
      }
    }
  }

  private func setMealWashHandsCompleted(_ isCompleted: Bool) {
    mealWashHandsCompleted = isCompleted
    for viewModel in componentViewModels.values {
      viewModel.washHandsCompleted = isCompleted
      if !isCompleted {
        viewModel.completedSteps.removeAll()
        viewModel.clearStepCompletionConditionProgress()
      }
    }
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
                  if let quantityText = aggregated.quantityText {
                    Text(quantityText)
                      .font(.subheadline)
                      .fontWeight(.medium)
                      .foregroundColor(.secondary)
                  }

                  Text(aggregated.name)
                    .font(.subheadline)
                    .foregroundColor(
                      checkedIngredients.contains(aggregated.ingredientID) ? .secondary : .primary
                    )
                    .strikethrough(checkedIngredients.contains(aggregated.ingredientID))
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
                let scaledQuantity = DiscreteQuantityScaling.scaled(
                  instrument.quantity, scale: scale)
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
                let scaledQuantity = DiscreteQuantityScaling.scaled(vessel.quantity, scale: scale)
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
          sharedViewModel: componentViewModels[component.recipe.id],
          globalWashHandsCompleted: mealWashHandsCompleted,
          onGlobalWashHandsToggle: { isCompleted in
            setMealWashHandsCompleted(isCompleted)
          },
          onViewModelReady: { recipeViewModel in
            componentViewModels[component.recipe.id] = recipeViewModel
          },
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
  let sharedViewModel: PerformRecipeViewModel?
  let globalWashHandsCompleted: Bool
  let onGlobalWashHandsToggle: ((Bool) -> Void)?
  let onViewModelReady: ((PerformRecipeViewModel) -> Void)?
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
              ),
              showWashHandsStepCard: false,
              sharedWashHandsCompleted: Binding(
                get: { globalWashHandsCompleted },
                set: { newValue in
                  onGlobalWashHandsToggle?(newValue)
                }
              ),
              allowCompletedStepsToggle: true
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
        if let sharedViewModel {
          viewModel = sharedViewModel
          onViewModelReady?(sharedViewModel)
        } else {
          let createdViewModel = PerformRecipeViewModel(
            recipeID: recipeID, authManager: authManager)
          viewModel = createdViewModel
          onViewModelReady?(createdViewModel)
        }
      }

      if let viewModel, viewModel.recipe == nil, !viewModel.isLoading {
        Task {
          await viewModel.loadRecipe()
        }
      }
      viewModel?.washHandsCompleted = globalWashHandsCompleted
      if !globalWashHandsCompleted {
        viewModel?.completedSteps.removeAll()
      }
    }
    .onChange(of: globalWashHandsCompleted) { _, newValue in
      viewModel?.washHandsCompleted = newValue
      if !newValue {
        viewModel?.completedSteps.removeAll()
      }
    }
    .onChange(of: viewModel?.recipe) { _, newRecipe in
      if let recipe = newRecipe {
        onRecipeLoaded?(recipe)
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
