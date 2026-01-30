//
//  MealSearchResultCard.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

// MARK: - Meal Search Result Card

struct MealSearchResultCard: View {
  let meal: Mealplanning_Meal
  let isSelected: Bool
  let onTap: () -> Void

  var body: some View {
    DSCard(style: isSelected ? .selected : .outlined, action: onTap) {
      HStack {
        VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
          Text(meal.name.isEmpty ? "Unnamed Meal" : meal.name)
            .font(DSTheme.Typography.label)
            .foregroundColor(DSTheme.Colors.textPrimary)

          if !meal.description_p.isEmpty {
            Text(meal.description_p)
              .font(DSTheme.Typography.body)
              .foregroundColor(DSTheme.Colors.textSecondary)
              .lineLimit(2)
          }

          // Show recipe names from components
          if !meal.components.isEmpty {
            let recipeNames = meal.components.compactMap { component -> String? in
              component.recipe.name.isEmpty ? nil : component.recipe.name
            }
            if !recipeNames.isEmpty {
              Text(recipeNames.joined(separator: ", "))
                .font(DSTheme.Typography.caption)
                .foregroundColor(DSTheme.Colors.textSecondary)
                .lineLimit(1)
            }
          }
        }

        Spacer()

        Image(systemName: isSelected ? "checkmark.circle.fill" : "circle")
          .foregroundColor(isSelected ? DSTheme.Colors.primary : DSTheme.Colors.textTertiary)
          .font(.title2)
      }
    }
  }
}
