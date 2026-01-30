//
//  DSSection.swift
//  ios
//
//  Design System Section Component - Section with header and optional subtitle
//

import SwiftUI

// MARK: - Section Style

enum DSSectionStyle {
  case standard
  case compact
  case card
}

// MARK: - DSSection Component

/// A section container with a styled header and optional subtitle.
///
/// Usage:
/// ```swift
/// DSSection("Members") {
///   ForEach(members) { MemberRow(member: $0) }
/// }
///
/// DSSection("Household", subtitle: "People in your household") {
///   ForEach(members) { MemberRow(member: $0) }
/// }
///
/// DSSection("Tasks", style: .compact) {
///   TaskList()
/// }
/// ```
struct DSSection<Content: View>: View {
  let title: String
  let subtitle: String?
  let style: DSSectionStyle
  let spacing: CGFloat
  let showDivider: Bool
  @ViewBuilder let content: () -> Content

  init(
    _ title: String,
    subtitle: String? = nil,
    style: DSSectionStyle = .standard,
    spacing: CGFloat? = nil,
    showDivider: Bool = false,
    @ViewBuilder content: @escaping () -> Content
  ) {
    self.title = title
    self.subtitle = subtitle
    self.style = style
    self.spacing = spacing ?? (style == .compact ? DSTheme.Spacing.sm : DSTheme.Spacing.md)
    self.showDivider = showDivider
    self.content = content
  }

  var body: some View {
    VStack(alignment: .leading, spacing: spacing) {
      // Header
      DSSectionHeader(title: title, subtitle: subtitle, style: headerStyle)

      if showDivider {
        Divider()
      }

      // Content
      if style == .card {
        DSCard {
          VStack(alignment: .leading, spacing: spacing) {
            content()
          }
          .frame(maxWidth: .infinity, alignment: .leading)
        }
      } else {
        content()
      }
    }
  }

  private var headerStyle: DSSectionHeaderStyle {
    switch style {
    case .compact:
      return .compact
    default:
      return .standard
    }
  }
}

// MARK: - Section Header Style

enum DSSectionHeaderStyle {
  case standard
  case compact
  case large
}

// MARK: - DSSectionHeader Component

/// A standalone section header that can be used independently.
///
/// Usage:
/// ```swift
/// DSSectionHeader(title: "Members")
/// DSSectionHeader(title: "Details", subtitle: "Additional information")
/// ```
struct DSSectionHeader: View {
  let title: String
  let subtitle: String?
  let style: DSSectionHeaderStyle

  init(
    title: String,
    subtitle: String? = nil,
    style: DSSectionHeaderStyle = .standard
  ) {
    self.title = title
    self.subtitle = subtitle
    self.style = style
  }

  var body: some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.xs) {
      Text(title)
        .font(titleFont)
        .foregroundColor(DSTheme.Colors.textPrimary)

      if let subtitle = subtitle, !subtitle.isEmpty {
        Text(subtitle)
          .font(subtitleFont)
          .foregroundColor(DSTheme.Colors.textSecondary)
      }
    }
    .padding(.horizontal, horizontalPadding)
  }

  private var titleFont: Font {
    switch style {
    case .large:
      return DSTheme.Typography.title1
    case .standard:
      return DSTheme.Typography.title2
    case .compact:
      return DSTheme.Typography.title3
    }
  }

  private var subtitleFont: Font {
    switch style {
    case .large:
      return DSTheme.Typography.body
    case .standard, .compact:
      return DSTheme.Typography.bodySmall
    }
  }

  private var horizontalPadding: CGFloat {
    switch style {
    case .compact:
      return DSTheme.Spacing.xxs
    default:
      return DSTheme.Spacing.xs
    }
  }
}

// MARK: - Empty Section Content

/// A helper view for displaying empty state within a section.
struct DSSectionEmptyContent: View {
  let message: String
  let icon: String?

  init(_ message: String, icon: String? = nil) {
    self.message = message
    self.icon = icon
  }

  var body: some View {
    HStack(spacing: DSTheme.Spacing.sm) {
      if let icon = icon {
        Image(systemName: icon)
          .font(.subheadline)
          .foregroundColor(DSTheme.Colors.textSecondary)
      }

      Text(message)
        .font(DSTheme.Typography.body)
        .foregroundColor(DSTheme.Colors.textSecondary)
    }
    .padding(.vertical, DSTheme.Spacing.sm)
    .padding(.horizontal, DSTheme.Spacing.xs)
  }
}

// MARK: - Preview

#Preview("Section Styles") {
  ScrollView {
    VStack(spacing: DSTheme.Spacing.xl) {
      DSSection("Standard Section") {
        Text("Content goes here")
          .dsCard()
      }

      DSSection("With Subtitle", subtitle: "This is additional context") {
        Text("More content")
          .dsCard()
      }

      DSSection("Compact Style", style: .compact) {
        VStack(spacing: DSTheme.Spacing.sm) {
          Text("Item 1")
          Text("Item 2")
          Text("Item 3")
        }
        .dsCard()
      }

      DSSection("Card Style", subtitle: "Content wrapped in a card", style: .card) {
        VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
          Text("Item 1")
          Divider()
          Text("Item 2")
          Divider()
          Text("Item 3")
        }
      }

      DSSection("Empty Section") {
        DSSectionEmptyContent("No items available", icon: "tray")
      }
    }
    .padding()
  }
}
