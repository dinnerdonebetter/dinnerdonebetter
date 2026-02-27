//
//  RecipeTimeEstimation.swift
//  ios
//
//  Aggregates step-level time estimates into a recipe-level range.
//  Steps without time use 1-3 minutes (60-180 seconds) as fallback.
//

import SwiftProtobuf

// MARK: - RecipeTimeEstimate

struct RecipeTimeEstimate {
  let minSeconds: UInt32
  let maxSeconds: UInt32
}

// MARK: - RecipeTimeEstimation

enum RecipeTimeEstimation {
  /// Seconds to assume for min when a step has no time estimate.
  static let defaultMinSeconds: UInt32 = 60  // 1 minute

  /// Seconds to assume for max when a step has no time estimate.
  static let defaultMaxSeconds: UInt32 = 180  // 3 minutes

  /// Estimates total recipe time by aggregating step times.
  /// - Parameter steps: Recipe steps (each may have optional estimatedTimeInSeconds).
  /// - Returns: Aggregated min/max in seconds, or nil if steps is empty.
  static func estimate(steps: [Mealplanning_RecipeStep]) -> RecipeTimeEstimate? {
    guard !steps.isEmpty else { return nil }

    var totalMin: UInt32 = 0
    var totalMax: UInt32 = 0

    for step in steps {
      let range = step.estimatedTimeInSeconds
      let (stepMin, stepMax): (UInt32, UInt32)
      if range.hasMin && range.hasMax {
        stepMin = range.min
        stepMax = range.max
      } else if range.hasMin {
        stepMin = range.min
        stepMax = range.min
      } else if range.hasMax {
        stepMin = range.max
        stepMax = range.max
      } else {
        stepMin = defaultMinSeconds
        stepMax = defaultMaxSeconds
      }
      totalMin = totalMin &+ stepMin
      totalMax = totalMax &+ stepMax
    }

    return RecipeTimeEstimate(minSeconds: totalMin, maxSeconds: totalMax)
  }

  /// Formats a time range for display (e.g. "5–15 min", "1 hr 7 min – 2 hr").
  /// Uses hours for values over 60 minutes; caps display at 24 hours to avoid absurd values.
  static func format(minSeconds: UInt32, maxSeconds: UInt32) -> String {
    let maxCapSeconds: UInt32 = 24 * 3600
    let cappedMax = min(maxSeconds, maxCapSeconds)

    // Compact "X–Y min" when both under 60 minutes
    if minSeconds < 3600 && cappedMax < 3600 {
      let minMinutes = Int(minSeconds / 60)
      let maxMinutes = Int(cappedMax / 60)
      if minMinutes == maxMinutes {
        return "\(minMinutes) min"
      }
      return "\(minMinutes)–\(maxMinutes) min"
    }

    // Use hours for larger values
    let minStr = formatDuration(seconds: minSeconds)
    let maxStr = maxSeconds > maxCapSeconds ? "24+ hr" : formatDuration(seconds: maxSeconds)
    if minStr == maxStr {
      return minStr
    }
    return "\(minStr) – \(maxStr)"
  }

  /// Formats a single duration for display (e.g. "5 min", "1 hr 7 min", "2 hr 30 min").
  private static func formatDuration(seconds: UInt32) -> String {
    let totalMinutes = Int(seconds / 60)
    if totalMinutes < 60 {
      if totalMinutes < 1 { return "< 1 min" }
      return "\(totalMinutes) min"
    }
    let hours = totalMinutes / 60
    let mins = totalMinutes % 60
    if mins == 0 {
      return hours == 1 ? "1 hr" : "\(hours) hr"
    }
    return "\(hours) hr \(mins) min"
  }

  /// Formats a step's optional time range for display.
  /// Returns nil when the step has no time estimate (neither min nor max set).
  static func formatStepTime(_ range: Common_OptionalUint32Range) -> String? {
    if range.hasMin && range.hasMax {
      return format(minSeconds: range.min, maxSeconds: range.max)
    }
    if range.hasMin {
      return formatDuration(seconds: range.min)
    }
    if range.hasMax {
      let maxCapSeconds: UInt32 = 24 * 3600
      return range.max > maxCapSeconds ? "24+ hr" : formatDuration(seconds: range.max)
    }
    return nil
  }
}
