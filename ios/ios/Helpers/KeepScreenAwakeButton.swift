//
//  KeepScreenAwakeButton.swift
//  ios
//
//  A button that toggles keeping the screen awake while cooking.
//  Uses UIApplication.shared.isIdleTimerDisabled to prevent screen dimming.
//

import SwiftUI
import UIKit

/// A button that toggles keeping the device screen awake.
/// When enabled, prevents the screen from dimming or locking while following a recipe.
/// - Parameter inline: When true, renders as a full-width card-style button (e.g. above wash hands).
///   When false, renders as a compact icon for toolbar use.
struct KeepScreenAwakeButton: View {
  @State private var isScreenAwake = false
  var inline: Bool = false

  var body: some View {
    Button {
      isScreenAwake.toggle()
      UIApplication.shared.isIdleTimerDisabled = isScreenAwake
    } label: {
      if inline {
        HStack(spacing: 12) {
          Image(systemName: isScreenAwake ? "sun.max.fill" : "sun.max")
            .font(.title2)
            .foregroundColor(isScreenAwake ? DSTheme.Colors.primary : DSTheme.Colors.textSecondary)
          Text(isScreenAwake ? "Screen awake (tap to allow sleep)" : "Keep screen awake")
            .font(.subheadline)
            .foregroundColor(isScreenAwake ? DSTheme.Colors.primary : .primary)
          Spacer()
        }
        .padding()
        .background(isScreenAwake ? DSTheme.Colors.primary.opacity(0.1) : Color(.systemGray6))
        .cornerRadius(12)
      } else {
        Image(systemName: isScreenAwake ? "sun.max.fill" : "sun.max")
          .foregroundColor(isScreenAwake ? DSTheme.Colors.primary : DSTheme.Colors.textSecondary)
      }
    }
    .buttonStyle(.plain)
    .accessibilityLabel(isScreenAwake ? "Screen awake (tap to allow sleep)" : "Keep screen awake")
    .accessibilityHint("Prevents the screen from dimming while cooking")
    .onDisappear {
      UIApplication.shared.isIdleTimerDisabled = false
    }
  }
}
