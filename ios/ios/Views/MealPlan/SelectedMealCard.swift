//
//  SelectedMealCard.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

// MARK: - Selected Meal Card

struct SelectedMealCard: View {
  let meal: Mealplanning_Meal
  let scale: Float
  let onScaleChange: (Float) -> Void
  let onRemove: () -> Void
  let isRegularWidth: Bool
  
  @State private var scaleText: String = ""
  @FocusState private var isScaleFocused: Bool

  var body: some View {
    if isRegularWidth {
      // Horizontal layout for iPad
      horizontalLayout
    } else {
      // Vertical layout for iPhone
      verticalLayout
    }
  }
  
  private var horizontalLayout: some View {
    HStack(alignment: .top, spacing: 16) {
      VStack(alignment: .leading, spacing: 4) {
        Text(meal.name.isEmpty ? "Unnamed Meal" : meal.name)
          .font(.headline)

        if !meal.components.isEmpty {
          let recipeNames = meal.components.compactMap { component -> String? in
            component.recipe.name.isEmpty ? nil : component.recipe.name
          }
          if !recipeNames.isEmpty {
            Text(recipeNames.joined(separator: ", "))
              .font(.caption)
              .foregroundColor(.secondary)
          }
        }
      }
      .frame(maxWidth: .infinity, alignment: .leading)

      // Scale control
      VStack(alignment: .leading, spacing: 4) {
        Text("Scale")
          .font(.caption)
          .foregroundColor(.secondary)
        HStack(spacing: 8) {
          TextField("1.0", text: $scaleText)
            .keyboardType(.decimalPad)
            .textFieldStyle(.roundedBorder)
            .frame(width: 100)
            .focused($isScaleFocused)
            .onSubmit {
              updateScaleFromText()
            }
            .onChange(of: isScaleFocused) { _, isFocused in
              if !isFocused {
                updateScaleFromText()
              }
            }
          Text("x")
            .font(.subheadline)
            .foregroundColor(.secondary)
        }
      }
      .frame(maxWidth: 200)

      Button(action: onRemove) {
        Image(systemName: "xmark.circle.fill")
          .foregroundColor(.red)
          .font(.title3)
      }
    }
    .padding()
    .background(Color(.systemBackground))
    .cornerRadius(8)
    .onAppear {
      scaleText = String(format: "%.2f", scale)
    }
    .onChange(of: scale) { _, newValue in
      if !isScaleFocused {
        scaleText = String(format: "%.2f", newValue)
      }
    }
  }
  
  private var verticalLayout: some View {
    VStack(alignment: .leading, spacing: 12) {
      HStack {
        VStack(alignment: .leading, spacing: 4) {
          Text(meal.name.isEmpty ? "Unnamed Meal" : meal.name)
            .font(.headline)

          if !meal.components.isEmpty {
            let recipeNames = meal.components.compactMap { component -> String? in
              component.recipe.name.isEmpty ? nil : component.recipe.name
            }
            if !recipeNames.isEmpty {
              Text(recipeNames.joined(separator: ", "))
                .font(.caption)
                .foregroundColor(.secondary)
            }
          }
        }

        Spacer()

        Button(action: onRemove) {
          Image(systemName: "xmark.circle.fill")
            .foregroundColor(.red)
            .font(.title3)
        }
      }

      // Scale control
      VStack(alignment: .leading, spacing: 4) {
        Text("Scale")
          .font(.caption)
          .foregroundColor(.secondary)
        HStack(spacing: 8) {
          TextField("1.0", text: $scaleText)
            .keyboardType(.decimalPad)
            .textFieldStyle(.roundedBorder)
            .focused($isScaleFocused)
            .onSubmit {
              updateScaleFromText()
            }
            .onChange(of: isScaleFocused) { _, isFocused in
              if !isFocused {
                updateScaleFromText()
              }
            }
          Text("x")
            .font(.subheadline)
            .foregroundColor(.secondary)
        }
      }
    }
    .padding()
    .background(Color(.systemBackground))
    .cornerRadius(8)
    .onAppear {
      scaleText = String(format: "%.2f", scale)
    }
    .onChange(of: scale) { _, newValue in
      if !isScaleFocused {
        scaleText = String(format: "%.2f", newValue)
      }
    }
  }
  
  private func updateScaleFromText() {
    // Parse the text input and validate it's a positive number
    if let parsedValue = Float(scaleText.trimmingCharacters(in: .whitespacesAndNewlines)) {
      if parsedValue > 0 {
        onScaleChange(parsedValue)
        scaleText = String(format: "%.2f", parsedValue)
      } else {
        // Invalid: not positive, reset to current scale
        scaleText = String(format: "%.2f", scale)
      }
    } else {
      // Invalid: not a number, reset to current scale
      scaleText = String(format: "%.2f", scale)
    }
  }
}

