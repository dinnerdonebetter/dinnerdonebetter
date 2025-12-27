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
    HStack {
      VStack(alignment: .leading, spacing: 8) {
        Text(meal.name.isEmpty ? "Unnamed Meal" : meal.name)
          .font(.headline)
          .foregroundColor(.primary)

        if !meal.description_p.isEmpty {
          Text(meal.description_p)
            .font(.subheadline)
            .foregroundColor(.secondary)
            .lineLimit(2)
        }

        // Show recipe names from components
        if !meal.components.isEmpty {
          let recipeNames = meal.components.compactMap { component -> String? in
            component.recipe.name.isEmpty ? nil : component.recipe.name
          }
          if !recipeNames.isEmpty {
            Text(recipeNames.joined(separator: ", "))
              .font(.caption)
              .foregroundColor(.secondary)
              .lineLimit(1)
          }
        }
      }

      Spacer()

      Image(systemName: isSelected ? "checkmark.circle.fill" : "circle")
        .foregroundColor(isSelected ? .blue : .gray)
        .font(.title2)
    }
    .padding()
    .background(isSelected ? Color.blue.opacity(0.1) : Color(.systemBackground))
    .cornerRadius(8)
    .overlay(
      RoundedRectangle(cornerRadius: 8)
        .stroke(isSelected ? Color.blue : Color.clear, lineWidth: 2)
    )
    .contentShape(Rectangle())
    .onTapGesture {
      onTap()
    }
  }
}

