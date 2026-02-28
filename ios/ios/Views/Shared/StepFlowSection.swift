//
//  StepFlowSection.swift
//  ios
//
//  Created by Auto on 2/23/26.
//

import SwiftUI

private func focusModeSectionHeader(title: String, color: Color) -> some View {
  HStack(spacing: 12) {
    Rectangle()
      .fill(Color.secondary.opacity(0.5))
      .frame(height: 1)
      .frame(maxWidth: .infinity)
    Text(title)
      .font(.headline)
      .fontWeight(.semibold)
      .foregroundColor(color)
    Rectangle()
      .fill(Color.secondary.opacity(0.5))
      .frame(height: 1)
      .frame(maxWidth: .infinity)
  }
  .padding(.vertical, 8)
  .frame(maxWidth: .infinity)
  .background(Color(uiColor: .systemBackground))
}

struct StepFlowGroup<Item: Identifiable>: Identifiable {
  let id: String
  let title: String
  let color: Color
  let items: [Item]

  init(title: String, color: Color, items: [Item]) {
    self.id = title
    self.title = title
    self.color = color
    self.items = items
  }
}

struct StepFlowSection<
  Item: Identifiable,
  Header: View,
  AllLead: View,
  FocusLead: View,
  Row: View
>: View {
  let showCompleted: Binding<Bool>
  let allSteps: [Item]
  let focusedGroups: [StepFlowGroup<Item>]
  let allowToggle: Bool
  let allStepsTitle: String
  let emptyMessage: String?
  let headerContent: () -> Header
  let allModeLeadingContent: () -> AllLead
  let focusModeLeadingContent: () -> FocusLead
  let rowContent: (Item) -> Row
  let onReorder: ((String, IndexSet, Int) -> Void)?

  init(
    showCompleted: Binding<Bool>,
    allSteps: [Item],
    focusedGroups: [StepFlowGroup<Item>],
    allowToggle: Bool = true,
    allStepsTitle: String = "All Steps",
    emptyMessage: String? = nil,
    onReorder: ((String, IndexSet, Int) -> Void)? = nil,
    @ViewBuilder headerContent: @escaping () -> Header,
    @ViewBuilder allModeLeadingContent: @escaping () -> AllLead,
    @ViewBuilder focusModeLeadingContent: @escaping () -> FocusLead,
    @ViewBuilder rowContent: @escaping (Item) -> Row
  ) {
    self.showCompleted = showCompleted
    self.allSteps = allSteps
    self.focusedGroups = focusedGroups
    self.allowToggle = allowToggle
    self.allStepsTitle = allStepsTitle
    self.emptyMessage = emptyMessage
    self.headerContent = headerContent
    self.allModeLeadingContent = allModeLeadingContent
    self.focusModeLeadingContent = focusModeLeadingContent
    self.rowContent = rowContent
    self.onReorder = onReorder
  }

  var body: some View {
    VStack(alignment: .leading, spacing: 16) {
      HStack {
        headerContent()
        Spacer()

        if allowToggle {
          Button(showCompleted.wrappedValue ? "Focus mode" : "Complete view") {
            withAnimation {
              showCompleted.wrappedValue.toggle()
            }
          }
          .font(.caption)
          .buttonStyle(.bordered)
          .controlSize(.small)
        }
      }

      if allowToggle && showCompleted.wrappedValue {
        VStack(alignment: .leading, spacing: 8) {
          Text(allStepsTitle)
            .font(.subheadline)
            .fontWeight(.semibold)
            .foregroundColor(.secondary)
            .padding(.horizontal, 4)

          allModeLeadingContent()

          ForEach(allSteps) { item in
            rowContent(item)
          }
        }
      } else {
        focusModeLeadingContent()

        LazyVStack(alignment: .leading, spacing: 0, pinnedViews: .sectionHeaders) {
          ForEach(focusedGroups) { group in
            if !group.items.isEmpty {
              Section {
                if let onReorder {
                  List {
                    ForEach(group.items) { item in
                      rowContent(item)
                        .listRowSeparator(.hidden)
                        .listRowInsets(EdgeInsets(top: 4, leading: 0, bottom: 4, trailing: 0))
                    }
                    .onMove { source, destination in
                      onReorder(group.id, source, destination)
                    }
                  }
                  .listStyle(.plain)
                  .environment(\.editMode, .constant(.active))
                  .scrollDisabled(true)
                  .frame(minHeight: CGFloat(group.items.count) * 80)
                } else {
                  ForEach(group.items) { item in
                    rowContent(item)
                  }
                }
              } header: {
                focusModeSectionHeader(title: group.title, color: group.color)
              }
            }
          }
        }
      }

      if let emptyMessage, allSteps.isEmpty {
        Text(emptyMessage)
          .font(.caption)
          .foregroundColor(.secondary)
      }
    }
  }
}
