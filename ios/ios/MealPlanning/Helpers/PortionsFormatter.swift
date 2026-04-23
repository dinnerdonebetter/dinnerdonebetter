//
//  PortionsFormatter.swift
//  ios
//
//  Formats estimated portions for display.
//  - No max: treat min as authoritative (not "limitless")
//  - Same min and max: single authoritative number
//  - Different min and max: show range
//

enum PortionsFormatter {
  /// Formats a portion range for display.
  /// - No max: shows min as the authoritative number (e.g. "4")
  /// - Min equals max: shows single number (e.g. "4")
  /// - Min differs from max: shows range (e.g. "4-6")
  static func format(min: Float, max: Float?) -> String {
    if let max = max {
      if min == max {
        return String(format: "%.1f", min)
      } else {
        return String(format: "%.1f-%.1f", min, max)
      }
    } else {
      return String(format: "%.1f", min)
    }
  }

  /// Formats a scaled portion range for display (e.g. when scaling a meal).
  static func formatScaled(min: Float, max: Float?, scale: Float) -> String {
    let scaledMin = min * scale
    if let max = max {
      let scaledMax = max * scale
      if scaledMin == scaledMax {
        return String(format: "%.1f", scaledMin)
      } else {
        return String(format: "%.1f-%.1f", scaledMin, scaledMax)
      }
    } else {
      return String(format: "%.1f", scaledMin)
    }
  }
}
