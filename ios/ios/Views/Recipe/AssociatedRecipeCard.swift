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
        HStack {
          VStack(alignment: .leading, spacing: 4) {
            Text(recipe.name.isEmpty ? "Unnamed Recipe" : recipe.name)
              .font(.subheadline)
              .fontWeight(.medium)
              .foregroundColor(.primary)

            if !recipe.description_p.isEmpty {
              Text(recipe.description_p)
                .font(.caption)
                .foregroundColor(.secondary)
                .lineLimit(2)
            }
          }

          Spacer()

          Image(systemName: "chevron.right")
            .font(.caption)
            .foregroundColor(.secondary)
        }
        .padding()
        .background(Color(.systemGray6))
        .cornerRadius(8)
      }
    )
  }
}
