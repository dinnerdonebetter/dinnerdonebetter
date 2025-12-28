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
        StepCardView(
          step: step,
          index: index,
          viewModel: viewModel,
          formatStepTitle: formatStepTitle
        )
      }
    }
  }

}
