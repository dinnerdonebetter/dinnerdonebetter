//
//  AssociatedRecipeCard.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct AssociatedRecipeCard: View {
  let recipe: Mealplanning_Recipe
  @Environment(AuthenticationManager.self) private var authManager

  var body: some View {
    NavigationLink(
      destination: {
        PerformRecipeView(recipeID: recipe.id)
          .environment(authManager)
      },
      label: {
        DSCard {
          HStack {
            VStack(alignment: .leading, spacing: DSTheme.Spacing.xs) {
              Text(recipe.name.isEmpty ? "Unnamed Recipe" : recipe.name)
                .font(DSTheme.Typography.body)
                .fontWeight(.medium)
                .foregroundColor(DSTheme.Colors.textPrimary)

              if !recipe.description_p.isEmpty {
                Text(recipe.description_p)
                  .font(DSTheme.Typography.caption)
                  .foregroundColor(DSTheme.Colors.textSecondary)
                  .lineLimit(2)
              }
            }

            Spacer()

            Image(systemName: "chevron.right")
              .font(.caption)
              .foregroundColor(DSTheme.Colors.textTertiary)
          }
        }
      }
    )
  }
}
