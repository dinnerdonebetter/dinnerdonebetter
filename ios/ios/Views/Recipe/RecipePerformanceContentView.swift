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
struct RecipePerformanceContentView: View {  // swiftlint:disable:this type_body_length
  @Binding var checkedIngredients: Set<String>
  @Binding var checkedInstrumentsVessels: Set<String>
  @Binding var isInstrumentsVesselsExpanded: Bool
  @Binding var isIngredientsExpanded: Bool

  let recipe: Mealplanning_Recipe
  let viewModel: PerformRecipeViewModel
  var hideIngredientsAndInstruments: Bool = false
  var mealPlanSelections: [Mealplanning_MealPlanRecipeOptionSelection]?
  var highlightedStepIDs: Set<String>? = nil
  var prepTaskContext: PerformRecipeView.PrepTaskContext? = nil

  // State for option selections (for interactive selection outside meal plan context)
  // Note: Only ingredients have selectable options; instruments and vessels are concrete
  @State private var selectedIngredientOptions: [String: UInt32] = [:]  // optionGroupID -> selectedOptionIndex

  @Environment(AuthenticationManager.self) private var authManager

  var body: some View {
    ScrollView {
      VStack(alignment: .leading, spacing: 16) {
        // Recipe header
        recipeHeader(recipe: recipe, viewModel: viewModel)

        // Associated recipes section
        if !recipe.associatedRecipes.isEmpty {
          associatedRecipesSection(recipe: recipe)
        }

        // Instruments & Vessels section (hidden when embedded in meal view)
        if !hideIngredientsAndInstruments {
          instrumentsVesselsSection(recipe: recipe)
        }

        // Ingredients section (hidden when embedded in meal view)
        if !hideIngredientsAndInstruments {
          ingredientsSection(recipe: recipe)
        }

        // Steps list
        stepsList(recipe: recipe, viewModel: viewModel)
      }
      .padding()
    }
  }

  // MARK: - Associated Recipes Section

  private func associatedRecipesSection(recipe: Mealplanning_Recipe) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      Text("Prerequisite Recipes")
        .font(.headline)
        .padding(.horizontal, 4)

      ForEach(recipe.associatedRecipes, id: \.id) { associatedRecipe in
        AssociatedRecipeCard(recipe: associatedRecipe)
      }
    }
  }

  // MARK: - Recipe Header

  private func formatEventTime(_ date: Date) -> String {
    let formatter = DateFormatter()
    formatter.dateStyle = .medium
    formatter.timeStyle = .short
    return formatter.string(from: date)
  }

  private func recipeHeader(recipe: Mealplanning_Recipe, viewModel: PerformRecipeViewModel)
    -> some View
  {
    VStack(alignment: .leading, spacing: 8) {
      // Show prep task context if available
      if let context = prepTaskContext {
        VStack(alignment: .leading, spacing: 4) {
          if let prepTaskName = context.prepTaskName, !prepTaskName.isEmpty {
            Text(prepTaskName)
              .font(.headline)
              .foregroundColor(.blue)
          }
          
          HStack(spacing: 4) {
            if let recipeName = context.recipeName, !recipeName.isEmpty {
              Text("for \(recipeName)")
                .font(.subheadline)     
                .foregroundColor(.secondary)
            }
            
            if let eventName = context.eventName, let eventTime = context.eventTime {
              Text("• \(eventName) at \(formatEventTime(eventTime))")
                .font(.subheadline)
                .foregroundColor(.secondary)
            }
          }
        }
        .padding(.bottom, 8)
      }
      
      Text(recipe.name.isEmpty ? "Unnamed Recipe" : recipe.name)
        .font(.title)
        .fontWeight(.bold)

      if !recipe.description_p.isEmpty {
        Text(recipe.description_p)
          .font(.subheadline)
          .foregroundColor(.secondary)
      }

      // Progress indicator
      let completedCount = viewModel.completedSteps.count + (viewModel.washHandsCompleted ? 1 : 0)
      let totalSteps = recipe.steps.count + 1  // +1 for wash hands step
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
    let (regularItems, _, _) =
      getAggregatedInstrumentsAndVessels(from: recipe)

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

      if isInstrumentsVesselsExpanded {
        VStack(alignment: .leading, spacing: 8) {
          // Regular items
          if !regularItems.isEmpty {
            ForEach(regularItems, id: \.itemID) { item in
              instrumentVesselRow(item: item)
            }
          }

          // Note: Instruments and vessels option groups are displayed in step details
          // but are not selectable (they're concrete, unchanging things)
          // Option groups for instruments/vessels will appear in the step details view
        }
        .padding(.vertical, 8)
        .background(Color(.systemBackground))
      }
    }
    .background(Color(.systemGray6))
    .cornerRadius(12)
  }

  private func instrumentVesselRow(item: AggregatedInstrumentVessel) -> some View {
    HStack(spacing: 12) {
      // Checkbox
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
          .foregroundColor(checkedInstrumentsVessels.contains(item.itemID) ? .green : .gray)
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

          if let sourceRecipeName = item.sourceRecipeName {
            Text("(from \(sourceRecipeName))")
              .font(.caption2)
              .foregroundColor(.secondary)
          }
        }
      }

      Spacer()
    }
    .padding(.horizontal)
    .padding(.vertical, 4)
  }

  private func filterInstrumentOptionGroups(
    _ groups: [InstrumentOptionGroupAggregate]
  ) -> [InstrumentOptionGroupAggregate] {
    guard let selections = mealPlanSelections else {
      return groups
    }

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
          return InstrumentOptionGroupAggregate(
            id: group.id,
            recipeID: group.recipeID,
            stepID: group.stepID,
            stepIndex: group.stepIndex,
            index: group.index,
            options: selectedOptions,
            selectedOptionIndex: selection.selectedOptionIndex,
            sourceRecipeID: group.sourceRecipeID,
            sourceRecipeName: group.sourceRecipeName
          )
        }
        return nil
      }

      return group
    }
  }

  private func filterVesselOptionGroups(
    _ groups: [VesselOptionGroupAggregate]
  ) -> [VesselOptionGroupAggregate] {
    guard let selections = mealPlanSelections else {
      return groups
    }

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
          return VesselOptionGroupAggregate(
            id: group.id,
            recipeID: group.recipeID,
            stepID: group.stepID,
            stepIndex: group.stepIndex,
            index: group.index,
            options: selectedOptions,
            selectedOptionIndex: selection.selectedOptionIndex,
            sourceRecipeID: group.sourceRecipeID,
            sourceRecipeName: group.sourceRecipeName
          )
        }
        return nil
      }

      return group
    }
  }

  // MARK: - Ingredients Section

  private func ingredientsSection(recipe: Mealplanning_Recipe) -> some View {
    let (regularIngredients, optionGroups) = getAggregatedIngredients(from: recipe)

    // Filter option groups based on meal plan selections or user selections
    let filteredOptionGroups = filterIngredientOptionGroups(optionGroups)

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

      if isIngredientsExpanded {
        VStack(alignment: .leading, spacing: 8) {
          // Regular ingredients
          if !regularIngredients.isEmpty {
            ForEach(regularIngredients, id: \.ingredientID) { aggregated in
              ingredientRow(aggregated: aggregated)
            }
          }

          // Options section (only ingredients have selectable options)
          if !filteredOptionGroups.isEmpty {
            Text("Options")
              .font(.subheadline)
              .fontWeight(.semibold)
              .foregroundColor(.secondary)
              .padding(.top, 8)
              .padding(.horizontal)

            ForEach(filteredOptionGroups) { group in
              InteractiveIngredientOptionGroupView(
                group: group,
                selectedOptionIndex: Binding(
                  get: {
                    selectedIngredientOptions[group.id] ?? (group.options.first?.optionIndex ?? 0)
                  },
                  set: { selectedIngredientOptions[group.id] = $0 }
                )
              )
            }
          }
        }
        .padding(.vertical, 8)
        .background(Color(.systemBackground))
      }
    }
    .background(Color(.systemGray6))
    .cornerRadius(12)
  }

  private func ingredientRow(aggregated: AggregatedIngredient) -> some View {
    HStack(spacing: 12) {
      // Checkbox
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

          if let sourceRecipeName = aggregated.sourceRecipeName {
            Text("(from \(sourceRecipeName))")
              .font(.caption2)
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

  private func filterIngredientOptionGroups(
    _ groups: [OptionGroupAggregate]
  ) -> [OptionGroupAggregate] {
    guard let selections = mealPlanSelections else {
      // No meal plan selections - return all groups for interactive selection
      return groups
    }

    // Filter based on meal plan selections
    return groups.compactMap { group in
      let selection = selections.first { sel in
        sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
          && sel.ingredientIndex == group.index && sel.selectionType == .ingredient
      }

      if let selection = selection {
        // Only show the selected option
        let selectedOptions = group.options.filter {
          $0.optionIndex == selection.selectedOptionIndex
        }
        if !selectedOptions.isEmpty {
          return OptionGroupAggregate(
            id: group.id,
            recipeID: group.recipeID,
            stepID: group.stepID,
            stepIndex: group.stepIndex,
            index: group.index,
            options: selectedOptions,
            selectedOptionIndex: selection.selectedOptionIndex,
            sourceRecipeID: group.sourceRecipeID,
            sourceRecipeName: group.sourceRecipeName
          )
        }
        return nil
      }

      // No selection - show all options
      return group
    }
  }

  // MARK: - Steps List

  private func shouldShowStep(stepID: String) -> Bool {
    // If no highlighted steps specified, show all steps
    guard let highlightedStepIDs = highlightedStepIDs else {
      return true
    }
    // Only show steps that are in the highlighted set
    return highlightedStepIDs.contains(stepID)
  }

  private func stepsList(recipe: Mealplanning_Recipe, viewModel: PerformRecipeViewModel)
    -> some View
  {
    VStack(alignment: .leading, spacing: 12) {
      Text("Steps")
        .font(.headline)
        .padding(.horizontal, 4)

      // Associated recipe steps (prerequisites) - render first
      if !recipe.associatedRecipes.isEmpty {
        ForEach(recipe.associatedRecipes, id: \.id) { associatedRecipe in
          if !associatedRecipe.steps.isEmpty {
            // Header for associated recipe
            associatedRecipeStepsHeader(recipe: associatedRecipe)

            // Wash hands step before first step of this associated recipe
            washHandsStepCard(viewModel: viewModel)

            // Steps from this associated recipe
            ForEach(
              Array(associatedRecipe.steps.enumerated())
                .filter { shouldShowStep(stepID: $0.element.id) },
              id: \.element.id
            ) { index, step in
              StepCardView(
                step: step,
                index: index,
                viewModel: viewModel,
                formatStepTitle: formatStepTitle,
                recipeID: associatedRecipe.id,
                mealPlanSelections: mealPlanSelections,
                isAssociatedRecipeStep: true,
                associatedRecipeName: associatedRecipe.name,
                highlightedStepIDs: highlightedStepIDs
              )
            }

            // Separator after associated recipe steps
            if associatedRecipe.id != recipe.associatedRecipes.last?.id
              || !recipe.steps.isEmpty
            {
              Divider()
                .padding(.vertical, 8)
            }
          }
        }

        // Header for main recipe steps (if both exist)
        if !recipe.steps.isEmpty {
          mainRecipeStepsHeader()
        }
      }

      // Main recipe steps
      if !recipe.steps.isEmpty {
        // Wash hands step before first step of main recipe
        washHandsStepCard(viewModel: viewModel)

        ForEach(
          Array(recipe.steps.enumerated())
            .filter { shouldShowStep(stepID: $0.element.id) },
          id: \.element.id
        ) { index, step in
          StepCardView(
            step: step,
            index: index,
            viewModel: viewModel,
            formatStepTitle: formatStepTitle,
            recipeID: recipe.id,
            mealPlanSelections: mealPlanSelections,
            highlightedStepIDs: highlightedStepIDs
          )
        }
      }
    }
  }

  // MARK: - Associated Recipe Steps Header

  private func associatedRecipeStepsHeader(recipe: Mealplanning_Recipe) -> some View {
    VStack(alignment: .leading, spacing: 4) {
      HStack(spacing: 8) {
        Text("PREREQUISITE:")
          .font(.caption2)
          .fontWeight(.semibold)
          .foregroundColor(.purple)
          .textCase(.uppercase)
          .tracking(0.5)

        NavigationLink(
          destination: {
            PerformRecipeView(recipeID: recipe.id)
              .environment(authManager)
          },
          label: {
            Text(recipe.name)
              .font(.subheadline)
              .fontWeight(.semibold)
              .foregroundColor(.purple)
          }
        )
      }

      if !recipe.description_p.isEmpty {
        Text(recipe.description_p)
          .font(.caption)
          .foregroundColor(.secondary)
      }
    }
    .padding()
    .frame(maxWidth: .infinity, alignment: .leading)
    .background(Color.purple.opacity(0.1))
    .cornerRadius(8)
    .overlay(
      RoundedRectangle(cornerRadius: 8)
        .stroke(Color.purple.opacity(0.3), lineWidth: 2)
    )
    .padding(.vertical, 4)
  }

  // MARK: - Main Recipe Steps Header

  private func mainRecipeStepsHeader() -> some View {
    HStack {
      Text("MAIN RECIPE STEPS")
        .font(.caption2)
        .fontWeight(.semibold)
        .foregroundColor(.blue)
        .textCase(.uppercase)
        .tracking(0.5)
    }
    .padding()
    .frame(maxWidth: .infinity, alignment: .leading)
    .background(Color.blue.opacity(0.1))
    .cornerRadius(8)
    .overlay(
      RoundedRectangle(cornerRadius: 8)
        .stroke(Color.blue.opacity(0.3), lineWidth: 2)
    )
    .padding(.vertical, 4)
  }

  // MARK: - Wash Hands Step Card

  private func washHandsStepCard(viewModel: PerformRecipeViewModel) -> some View {
    let isCompleted = viewModel.isStepCompleted(PerformRecipeViewModel.washHandsStepIndex)
    let canCheck = viewModel.canCheckStep(PerformRecipeViewModel.washHandsStepIndex)

    return VStack(alignment: .leading, spacing: 12) {
      // Step header with checkbox
      HStack(alignment: .top, spacing: 12) {
        // Checkbox
        Button(
          action: {
            viewModel.toggleStep(PerformRecipeViewModel.washHandsStepIndex)
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

        // Step title
        VStack(alignment: .leading, spacing: 4) {
          Text("Wash your hands")
            .font(.headline)
            .foregroundColor(isCompleted ? .secondary : .primary)
            .italic(isCompleted)
        }

        Spacer()
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
