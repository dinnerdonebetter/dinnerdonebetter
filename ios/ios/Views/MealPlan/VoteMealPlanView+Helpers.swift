//
//  VoteMealPlanView+Helpers.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

// MARK: - Helper Views Extension

extension VoteMealPlanView {
  func rankedOptionsList(
    event: Mealplanning_MealPlanEvent, viewModel: VoteMealPlanViewModel
  ) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      Text("Your Ranking")
        .font(.headline)
        .padding(.horizontal, 4)

      if let ballot = viewModel.getBallot(for: event.id) {
        if ballot.rankedOptions.isEmpty {
          Text("No options available")
            .font(.subheadline)
            .foregroundColor(.secondary)
            .padding()
        } else {
          // Use List for native drag and drop support
          List {
            ForEach(Array(ballot.rankedOptions.enumerated()), id: \.element.id) { index, option in
              rankedOptionCard(option: option, rank: index + 1, event: event)
                .listRowInsets(EdgeInsets(top: 4, leading: 0, bottom: 4, trailing: 0))
                .listRowBackground(Color(.systemBackground))
            }
            .onMove { source, destination in
              if !ballot.isLocked {
                viewModel.reorderOptions(eventID: event.id, from: source, to: destination)
              }
            }
          }
          .listStyle(.plain)
          .frame(minHeight: CGFloat(ballot.rankedOptions.count * 80))
          .scrollContentBackground(.hidden)
          .environment(\.editMode, ballot.isLocked ? .constant(.inactive) : editMode)
        }
      }
    }
  }

  func rankedOptionCard(
    option: Mealplanning_MealPlanOption, rank: Int, event: Mealplanning_MealPlanEvent
  ) -> some View {
    HStack(spacing: 12) {
      // Rank indicator
      ZStack {
        Circle()
          .fill(rankColor(rank))
          .frame(width: 32, height: 32)
        Text("\(rank)")
          .font(.headline)
          .foregroundColor(.white)
      }

      // Option details
      VStack(alignment: .leading, spacing: 4) {
        if option.hasMeal {
          Text(option.meal.name)
            .font(.headline)
          if !option.meal.description_p.isEmpty {
            Text(option.meal.description_p)
              .font(.caption)
              .foregroundColor(.secondary)
              .lineLimit(2)
          }
        } else {
          Text("Option \(rank)")
            .font(.headline)
        }
      }

      Spacer()

      // Drag handle
      Image(systemName: "line.3.horizontal")
        .foregroundColor(.secondary)
    }
    .padding()
    .background(Color(.systemBackground))
    .cornerRadius(8)
    .shadow(color: Color.black.opacity(0.1), radius: 2, x: 0, y: 1)
  }

  func rankColor(_ rank: Int) -> Color {
    switch rank {
    case 1:
      return .green
    case 2:
      return .blue
    case 3:
      return .orange
    default:
      return .gray
    }
  }
}
