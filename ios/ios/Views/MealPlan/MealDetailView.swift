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
  @State private var showMealCompletedSteps = false
  @State private var customUpNextOrder: [String] = []
  @State private var customForLaterOrder: [String] = []

  let mealID: String
  /// When true, scale editing is hidden because the scale is set by the meal plan.
  let isFromMealPlan: Bool

  init(mealID: String, isFromMealPlan: Bool = false) {
    self.mealID = mealID
    self.isFromMealPlan = isFromMealPlan
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
            loadStepOrderFromUserDefaults(mealID: meal.id)
          }
        }
      }
    }
    .onChange(of: viewModel?.meal?.id) { _, _ in
      if let meal = viewModel?.meal {
        ensureComponentDataLoaded(for: meal)
        loadStepOrderFromUserDefaults(mealID: meal.id)
      }
    }
  }

  // MARK: - Overall Info Section

  private func overallInfoSection(meal: Mealplanning_Meal) -> some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
      // Description
      if !meal.description_p.isEmpty {
        Text(meal.description_p)
          .font(DSTheme.Typography.body)
          .foregroundColor(DSTheme.Colors.textSecondary)
      }

      // Estimated portions
      if meal.hasEstimatedPortions {
        HStack(spacing: DSTheme.Spacing.sm) {
          Image(systemName: "person.2")
            .foregroundColor(DSTheme.Colors.textSecondary)
          Text("Estimated Portions: \(PortionsFormatter.format(meal.estimatedPortions))")
            .font(DSTheme.Typography.body)
            .foregroundColor(DSTheme.Colors.textSecondary)
        }
      }

      // Meal Scale Control (hidden when viewing from meal plan – scale is set by the plan)
      if !isFromMealPlan {
        Divider()
          .padding(.vertical, DSTheme.Spacing.xs)

        HStack(spacing: DSTheme.Spacing.md) {
          Text("Meal Scale:")
            .font(DSTheme.Typography.label)

          HStack(spacing: DSTheme.Spacing.sm) {
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
              .font(DSTheme.Typography.body)
              .foregroundColor(DSTheme.Colors.textSecondary)

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
    }
    .padding(DSTheme.Spacing.lg)
    .frame(maxWidth: .infinity, alignment: .leading)
    .background(DSTheme.Colors.cardBackground)
    .cornerRadius(DSTheme.Radius.lg)
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
    VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
      // Instruments & Vessels
      aggregatedInstrumentsVesselsSection

      // Ingredients
      aggregatedIngredientsSection
    }
  }

  private var mealWashHandsSection: some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
      HStack(alignment: .top, spacing: DSTheme.Spacing.md) {
        Button(
          action: {
            setMealWashHandsCompleted(!mealWashHandsCompleted)
          },
          label: {
            Image(systemName: mealWashHandsCompleted ? "checkmark.circle.fill" : "circle")
              .font(.title2)
              .foregroundColor(
                mealWashHandsCompleted
                  ? DSTheme.Colors.success
                  : DSTheme.Colors.primary
              )
          }
        )
        .buttonStyle(.plain)

        VStack(alignment: .leading, spacing: DSTheme.Spacing.xs) {
          Text("🧼 Wash your hands")
            .font(DSTheme.Typography.title3)
            .foregroundColor(
              mealWashHandsCompleted
                ? DSTheme.Colors.textSecondary
                : DSTheme.Colors.textPrimary
            )
            .italic(mealWashHandsCompleted)

          Text("Complete this once to unlock all component steps.")
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.textSecondary)
        }

        Spacer()
      }
    }
    .padding(DSTheme.Spacing.lg)
    .background(
      mealWashHandsCompleted
        ? DSTheme.Colors.cardBackground
        : Color(.systemBackground)
    )
    .cornerRadius(DSTheme.Radius.lg)
    .overlay(
      RoundedRectangle(cornerRadius: DSTheme.Radius.lg)
        .stroke(
          mealWashHandsCompleted
            ? DSTheme.Colors.success.opacity(0.3)
            : Color.clear,
          lineWidth: 2
        )
    )
  }

  private var mealFlowModeSection: some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
      Text("Step View")
        .font(DSTheme.Typography.label)
        .foregroundColor(DSTheme.Colors.textSecondary)

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
              .font(DSTheme.Typography.title3)
              .foregroundColor(DSTheme.Colors.textPrimary)
            Spacer()
            Image(systemName: isInstrumentsVesselsExpanded ? "chevron.down" : "chevron.right")
              .font(DSTheme.Typography.caption)
              .foregroundColor(DSTheme.Colors.textSecondary)
          }
          .padding(DSTheme.Spacing.lg)
          .background(DSTheme.Colors.cardBackground)
        }
      )
      .buttonStyle(.plain)

      if isInstrumentsVesselsExpanded && !aggregatedItems.isEmpty {
        aggregatedInstrumentsVesselsList(aggregatedItems)
      }
    }
    .background(DSTheme.Colors.cardBackground)
    .cornerRadius(DSTheme.Radius.lg)
  }

  @ViewBuilder
  private func aggregatedInstrumentsVesselsList(_ items: [AggregatedInstrumentVessel]) -> some View
  {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
      ForEach(items, id: \.itemID) { item in
        aggregatedInstrumentVesselRow(item: item)
      }
    }
    .padding(.vertical, DSTheme.Spacing.sm)
    .background(Color(.systemBackground))
  }

  private func aggregatedInstrumentVesselRow(item: AggregatedInstrumentVessel) -> some View {
    HStack(spacing: DSTheme.Spacing.md) {
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
            checkedInstrumentsVessels.contains(item.itemID)
              ? DSTheme.Colors.success
              : DSTheme.Colors.textTertiary
          )
        }
      )
      .buttonStyle(.plain)

      HStack(spacing: DSTheme.Spacing.sm) {
        Image(
          systemName: item.type == .instrument
            ? "wrench.and.screwdriver" : "square.stack.3d.up"
        )
        .font(DSTheme.Typography.caption)
        .foregroundColor(DSTheme.Colors.textSecondary)
        .frame(width: 20)

        HStack {
          Text(item.name)
            .font(DSTheme.Typography.body)
            .foregroundColor(
              checkedInstrumentsVessels.contains(item.itemID)
                ? DSTheme.Colors.textSecondary
                : DSTheme.Colors.textPrimary
            )
            .strikethrough(checkedInstrumentsVessels.contains(item.itemID))

          if let sourceName = item.sourceRecipeName, !sourceName.isEmpty {
            Text("(from \(sourceName))")
              .font(DSTheme.Typography.captionSmall)
              .foregroundColor(DSTheme.Colors.textSecondary)
          }
        }
      }

      Spacer()
    }
    .padding(.horizontal, DSTheme.Spacing.lg)
    .padding(.vertical, DSTheme.Spacing.xs)
  }

  private static func stepOrderKey(mealID: String, group: String) -> String {
    "stepOrder_meal_\(mealID)_\(group)"
  }

  private func loadStepOrderFromUserDefaults(mealID: String) {
    let upNext =
      UserDefaults.standard.stringArray(forKey: Self.stepOrderKey(mealID: mealID, group: "upNext"))
      ?? []
    let forLater =
      UserDefaults.standard.stringArray(
        forKey: Self.stepOrderKey(mealID: mealID, group: "forLater")) ?? []
    customUpNextOrder = upNext
    customForLaterOrder = forLater
  }

  private func saveStepOrderToUserDefaults(mealID: String, upNext: [String], forLater: [String]) {
    UserDefaults.standard.set(upNext, forKey: Self.stepOrderKey(mealID: mealID, group: "upNext"))
    UserDefaults.standard.set(
      forLater, forKey: Self.stepOrderKey(mealID: mealID, group: "forLater"))
  }

  private func applyOrder<T: Identifiable>(_ items: [T], order: [String]) -> [T]
  where T.ID == String {
    guard !order.isEmpty else { return items }
    return items.sorted { firstItem, secondItem in
      let firstIndex = order.firstIndex(of: firstItem.id) ?? Int.max
      let secondIndex = order.firstIndex(of: secondItem.id) ?? Int.max
      return firstIndex < secondIndex
    }
  }

  private func unifiedMealStepsSection(meal: Mealplanning_Meal) -> some View {
    let allSteps = collectUnifiedMealStepsWithMerging(
      meal: meal,
      componentViewModels: componentViewModels,
      loadedRecipes: loadedRecipes,
      baseComponentScales: baseComponentScales,
      mealScale: mealScale,
      formatComponentType: formatComponentType
    )
    let upNext = allSteps.filter { $0.category == .upNext }
    let forLater = allSteps.filter { $0.category == .forLater }

    let orderedUpNext = applyOrder(upNext, order: customUpNextOrder)
    let orderedForLater = applyOrder(forLater, order: customForLaterOrder)

    let focusedGroups = [
      StepFlowGroup(title: "Up Next", color: DSTheme.Colors.warning, items: orderedUpNext),
      StepFlowGroup(title: "Not Yet", color: DSTheme.Colors.primary, items: orderedForLater),
    ]

    return StepFlowSection(
      showCompleted: $showMealCompletedSteps,
      allSteps: allSteps,
      focusedGroups: focusedGroups,
      allStepsTitle: "All Steps",
      emptyMessage: "Loading component steps...",
      onReorder: { groupId, source, destination in
        let items = groupId == "Up Next" ? orderedUpNext : orderedForLater
        var mutable = items
        mutable.move(fromOffsets: source, toOffset: destination)
        let newOrder = mutable.map(\.id)
        if groupId == "Up Next" {
          customUpNextOrder = newOrder
        } else {
          customForLaterOrder = newOrder
        }
        saveStepOrderToUserDefaults(
          mealID: meal.id,
          upNext: customUpNextOrder,
          forLater: customForLaterOrder
        )
      },
      headerContent: {
        Text("Cook Flow")
          .font(DSTheme.Typography.title2)
      },
      allModeLeadingContent: {
        EmptyView()
      },
      focusModeLeadingContent: {
        EmptyView()
      },
      rowContent: { item in
        VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
          let tagStyle = componentTagStyle(for: item.componentIndex)
          Text(item.componentNamesForTag)
            .font(DSTheme.Typography.captionSmall)
            .padding(.horizontal, DSTheme.Spacing.sm)
            .padding(.vertical, DSTheme.Spacing.xxs)
            .background(tagStyle.background)
            .foregroundColor(tagStyle.foreground)
            .cornerRadius(DSTheme.Radius.sm)

          StepCardView(
            step: item.step,
            index: item.stepIndex,
            viewModel: item.viewModel,
            formatStepTitle: { step, _ in
              formatUnifiedStepTitle(step: step, index: item.stepIndex)
            },
            recipeID: item.recipeID,
            isAssociatedRecipeStep: item.isAssociatedRecipeStep,
            associatedRecipeName: item.associatedRecipeName,
            scale: item.isMerged ? 1.0 : item.scale,
            isCompletedOverride: item.isMerged
              ? item.sources.allSatisfy {
                $0.viewModel.isStepCompleted(recipeID: $0.recipeID, stepID: $0.step.id)
              }
              : nil,
            canCheckOverride: item.isMerged
              ? item.sources.allSatisfy { source in
                let prereqs = source.viewModel.getPrerequisiteStepKeys(
                  recipeID: source.recipeID, stepID: source.step.id
                )
                return prereqs.allSatisfy { source.viewModel.completedSteps.contains($0) }
              }
              : nil,
            onToggleOverride: item.isMerged
              ? {
                for source in item.sources {
                  source.viewModel.toggleStep(recipeID: source.recipeID, stepID: source.step.id)
                }
              }
              : nil
          )
        }
      }
    )
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
              .font(DSTheme.Typography.title3)
              .foregroundColor(DSTheme.Colors.textPrimary)
            Spacer()
            Image(systemName: isIngredientsExpanded ? "chevron.down" : "chevron.right")
              .font(DSTheme.Typography.caption)
              .foregroundColor(DSTheme.Colors.textSecondary)
          }
          .padding(DSTheme.Spacing.lg)
          .background(DSTheme.Colors.cardBackground)
        }
      )
      .buttonStyle(.plain)

      if isIngredientsExpanded && !aggregatedIngredients.isEmpty {
        VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
          ForEach(aggregatedIngredients, id: \.ingredientID) { aggregated in
            HStack(spacing: DSTheme.Spacing.md) {
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
                    checkedIngredients.contains(aggregated.ingredientID)
                      ? DSTheme.Colors.success
                      : DSTheme.Colors.textTertiary
                  )
                }
              )
              .buttonStyle(.plain)

              VStack(alignment: .leading, spacing: DSTheme.Spacing.xxs) {
                HStack {
                  if let quantityText = aggregated.quantityText {
                    Text(quantityText)
                      .font(DSTheme.Typography.label)
                      .foregroundColor(DSTheme.Colors.textSecondary)
                  }

                  Text(aggregated.name)
                    .font(DSTheme.Typography.body)
                    .foregroundColor(
                      checkedIngredients.contains(aggregated.ingredientID)
                        ? DSTheme.Colors.textSecondary
                        : DSTheme.Colors.textPrimary
                    )
                    .strikethrough(checkedIngredients.contains(aggregated.ingredientID))

                  if let sourceName = aggregated.sourceRecipeName, !sourceName.isEmpty {
                    Text("(from \(sourceName))")
                      .font(DSTheme.Typography.captionSmall)
                      .foregroundColor(DSTheme.Colors.textSecondary)
                  }
                }

                if !aggregated.quantityNotes.isEmpty {
                  Text(aggregated.quantityNotes)
                    .font(DSTheme.Typography.caption)
                    .foregroundColor(DSTheme.Colors.textSecondary)
                }
              }

              Spacer()
            }
            .padding(.horizontal, DSTheme.Spacing.lg)
            .padding(.vertical, DSTheme.Spacing.xs)
          }
        }
        .padding(.vertical, DSTheme.Spacing.sm)
        .background(Color(.systemBackground))
      }
    }
    .background(DSTheme.Colors.cardBackground)
    .cornerRadius(DSTheme.Radius.lg)
  }

}

// MARK: - Meal Detail View Extensions

extension MealDetailView {
  func getAggregatedInstrumentsAndVessels() -> [AggregatedInstrumentVessel] {
    var aggregated: [String: AggregatedInstrumentVessel] = [:]

    for (_, recipeData) in loadedRecipes {
      let recipe = recipeData.recipe
      let scale = recipeData.scale

      func processSteps(
        _ steps: [Mealplanning_RecipeStep], sourceRecipeID: String?, sourceRecipeName: String?
      ) {
        for step in steps {
          for instrument in step.instruments where instrument.hasInstrument {
            let validInstrument = instrument.instrument
            if validInstrument.displayInSummaryLists {
              let itemID = validInstrument.id
              if !itemID.isEmpty {
                if aggregated[itemID] == nil {
                  aggregated[itemID] = AggregatedInstrumentVessel(
                    itemID: itemID,
                    name: instrument.name,
                    type: .instrument,
                    sourceRecipeID: sourceRecipeID,
                    sourceRecipeName: sourceRecipeName
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
                    type: .vessel,
                    sourceRecipeID: sourceRecipeID,
                    sourceRecipeName: sourceRecipeName
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

      // Process associated (prerequisite) recipe steps first
      for associatedRecipe in recipe.associatedRecipes {
        processSteps(
          associatedRecipe.steps,
          sourceRecipeID: associatedRecipe.id,
          sourceRecipeName: associatedRecipe.name.isEmpty ? nil : associatedRecipe.name
        )
      }
      // Process main recipe steps
      processSteps(recipe.steps, sourceRecipeID: nil, sourceRecipeName: nil)
    }

    return Array(aggregated.values).sorted { $0.name < $1.name }
  }

  func getAggregatedIngredients() -> [AggregatedIngredient] {
    var aggregated: [String: AggregatedIngredient] = [:]

    for (_, recipeData) in loadedRecipes {
      let recipe = recipeData.recipe
      let scale = recipeData.scale

      func processSteps(
        _ steps: [Mealplanning_RecipeStep], sourceRecipeID: String?, sourceRecipeName: String?
      ) {
        for step in steps {
          for ingredient in step.ingredients where ingredient.hasIngredient {
            let validIngredient = ingredient.ingredient
            let key = validIngredient.id
            if !key.isEmpty {
              if aggregated[key] == nil {
                aggregated[key] = AggregatedIngredient(
                  ingredientID: key,
                  name: ingredient.name,
                  quantityNotes: ingredient.quantityNotes,
                  measurementUnit: ingredient.hasMeasurementUnit ? ingredient.measurementUnit : nil,
                  sourceRecipeID: sourceRecipeID,
                  sourceRecipeName: sourceRecipeName
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

      // Process associated (prerequisite) recipe steps first
      for associatedRecipe in recipe.associatedRecipes {
        processSteps(
          associatedRecipe.steps,
          sourceRecipeID: associatedRecipe.id,
          sourceRecipeName: associatedRecipe.name.isEmpty ? nil : associatedRecipe.name
        )
      }
      // Process main recipe steps
      processSteps(recipe.steps, sourceRecipeID: nil, sourceRecipeName: nil)
    }

    return Array(aggregated.values).sorted { $0.name < $1.name }
  }

  private func componentsSection(meal: Mealplanning_Meal) -> some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
      Text("Components")
        .font(DSTheme.Typography.title2)

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
          sharedShowCompletedSteps: $showMealCompletedSteps,
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
  let sharedShowCompletedSteps: Binding<Bool>?
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
                  .font(DSTheme.Typography.labelSmall)
                  .padding(.horizontal, DSTheme.Spacing.sm)
                  .padding(.vertical, DSTheme.Spacing.xs)
                  .background(DSTheme.Colors.primary.opacity(0.2))
                  .foregroundColor(DSTheme.Colors.primary)
                  .cornerRadius(DSTheme.Radius.sm)
              }

              // Recipe name
              if let recipe = viewModel?.recipe {
                Text(recipe.name.isEmpty ? "Unnamed Recipe" : recipe.name)
                  .font(DSTheme.Typography.label)
                  .foregroundColor(DSTheme.Colors.textPrimary)
              } else {
                Text("Loading...")
                  .font(DSTheme.Typography.body)
                  .foregroundColor(DSTheme.Colors.textSecondary)
              }

              Spacer()

              // Chevron
              Image(systemName: isExpanded ? "chevron.down" : "chevron.right")
                .font(DSTheme.Typography.caption)
                .foregroundColor(DSTheme.Colors.textSecondary)
            }
            .padding(DSTheme.Spacing.lg)
            .background(DSTheme.Colors.cardBackground)
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
              .padding(DSTheme.Spacing.lg)
          } else if let errorMessage = viewModel.errorMessage {
            VStack(spacing: DSTheme.Spacing.sm) {
              Image(systemName: "exclamationmark.triangle")
                .foregroundColor(DSTheme.Colors.warning)
              Text("Error loading recipe: \(errorMessage)")
                .font(DSTheme.Typography.caption)
                .foregroundColor(DSTheme.Colors.textSecondary)
            }
            .frame(maxWidth: .infinity)
            .padding(DSTheme.Spacing.lg)
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
              sharedCompletedStepsVisibility: sharedShowCompletedSteps,
              allowCompletedStepsToggle: true
            )
            .onAppear {
              onRecipeLoaded?(recipe)
            }
          }
        } else {
          ProgressView("Initializing...")
            .frame(maxWidth: .infinity)
            .padding(DSTheme.Spacing.lg)
        }
      }
    }
    .background(DSTheme.Colors.cardBackground)
    .cornerRadius(DSTheme.Radius.md)
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
