//
//  HomeView.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import Combine
import SwiftProtobuf
import SwiftUI

struct HomeView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @State private var viewModel: HomeViewModel?

  var body: some View {
    NavigationStack {
      Group {
        if let viewModel = viewModel {
          DSContentState(
            isLoading: viewModel.isLoading,
            loadingMessage: "Loading...",
            error: viewModel.errorMessage,
            onRetry: { await viewModel.loadData() },
            content: {
              ScrollView {
                VStack(spacing: DSTheme.Spacing.xl) {
                  // Welcome & Create CTA
                  heroSection(viewModel: viewModel)

                  // Quick Access Row
                  quickAccessRow

                  // Pending Votes Section
                  if !viewModel.pendingVoteMealPlans.isEmpty {
                    pendingVotesSection(viewModel: viewModel)
                  }

                  // Upcoming Meals Section
                  if !viewModel.upcomingMealPlans.isEmpty {
                    upcomingMealsSection(viewModel: viewModel)
                  }

                  // My Tasks Section
                  if !viewModel.userTasks.isEmpty {
                    myTasksSection(viewModel: viewModel)
                  }

                  // Grocery Lists Section
                  if !viewModel.activeGroceryLists.isEmpty {
                    groceryListsSection(viewModel: viewModel)
                  }

                  // Empty State
                  if viewModel.pendingVoteMealPlans.isEmpty
                    && viewModel.upcomingMealPlans.isEmpty
                    && viewModel.userTasks.isEmpty
                    && viewModel.activeGroceryLists.isEmpty
                  {
                    emptyStateView
                  }
                }
                .dsScreenPadding()
              }
            })
        } else {
          DSInitializingView()
        }
      }
      .navigationTitle("Home")
      .toolbar {
        ToolbarItem(placement: .navigationBarTrailing) {
          NavigationLink(destination: AccountSettingsView()) {
            DSAvatarView(
              name: viewModel?.currentUserDisplayName ?? authManager.username,
              size: .sm
            )
          }
        }
      }
      .refreshable {
        if let viewModel = viewModel {
          await viewModel.loadData()
        }
      }
      .onAppear {
        if viewModel == nil {
          viewModel = HomeViewModel(authManager: authManager)
        }
        if let viewModel = viewModel {
          Task {
            await viewModel.loadData()
          }
        }
      }
      .onReceive(NotificationCenter.default.publisher(for: .mealPlanCreated)) { _ in
        if let viewModel = viewModel {
          Task {
            await viewModel.loadData()
          }
        }
      }
    }
  }

  // MARK: - Hero Section
  private func heroSection(viewModel: HomeViewModel) -> some View {
    VStack(spacing: DSTheme.Spacing.lg) {
      // Greeting
      Text("\(greeting), \(viewModel.currentUserDisplayName)!")
        .font(DSTheme.Typography.largeTitle)
        .foregroundColor(DSTheme.Colors.textPrimary)
        .frame(maxWidth: .infinity, alignment: .leading)

      // Create Meal Plan CTA Card
      NavigationLink(destination: CreateMealPlanView()) {
        HStack(spacing: DSTheme.Spacing.md) {
          ZStack {
            Circle()
              .fill(DSTheme.Colors.primary.opacity(0.15))
              .frame(width: 48, height: 48)

            Image(systemName: "plus")
              .font(.system(size: 20, weight: .semibold))
              .foregroundColor(DSTheme.Colors.primary)
          }

          VStack(alignment: .leading, spacing: DSTheme.Spacing.xxs) {
            Text("Create Meal Plan")
              .font(DSTheme.Typography.label)
              .foregroundColor(DSTheme.Colors.textPrimary)

            Text("Plan your meals for the week")
              .font(DSTheme.Typography.caption)
              .foregroundColor(DSTheme.Colors.textSecondary)
          }

          Spacer()

          Image(systemName: "chevron.right")
            .font(.system(size: 14, weight: .semibold))
            .foregroundColor(DSTheme.Colors.textTertiary)
        }
        .padding(DSTheme.Spacing.lg)
        .background(
          RoundedRectangle(cornerRadius: DSTheme.Radius.lg)
            .fill(DSTheme.Colors.cardBackground)
            .shadow(color: Color.black.opacity(0.06), radius: 8, x: 0, y: 2)
        )
        .overlay(
          RoundedRectangle(cornerRadius: DSTheme.Radius.lg)
            .stroke(DSTheme.Colors.primary.opacity(0.2), lineWidth: 1)
        )
      }
      .buttonStyle(.plain)
    }
  }

  private var greeting: String {
    let hour = Calendar.current.component(.hour, from: Date())
    switch hour {
    case 0..<12:
      return "Good morning"
    case 12..<17:
      return "Good afternoon"
    default:
      return "Good evening"
    }
  }

  // MARK: - Quick Access Row
  private var quickAccessRow: some View {
    HStack(spacing: DSTheme.Spacing.md) {
      QuickAccessButton(
        icon: "book.closed.fill",
        label: "Recipes",
        color: DSTheme.Colors.secondary,
        destination: RecipeListView()
      )

      QuickAccessButton(
        icon: "fork.knife",
        label: "Meals",
        color: DSTheme.Colors.tertiary,
        destination: MealListView()
      )
    }
  }

  // MARK: - Pending Votes Section
  private func pendingVotesSection(viewModel: HomeViewModel) -> some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
      HStack {
        Label("Pending Votes", systemImage: "hand.raised.fill")
          .font(DSTheme.Typography.title2)
          .foregroundColor(DSTheme.Colors.textPrimary)

        Spacer()

        Text("\(viewModel.pendingVoteMealPlans.count)")
          .font(DSTheme.Typography.caption)
          .fontWeight(.semibold)
          .foregroundColor(DSTheme.Colors.warning)
          .padding(.horizontal, DSTheme.Spacing.sm)
          .padding(.vertical, DSTheme.Spacing.xxs)
          .background(DSTheme.Colors.warning.opacity(0.15))
          .cornerRadius(DSTheme.Radius.full)
      }

      ForEach(viewModel.pendingVoteMealPlans, id: \.id) { mealPlan in
        NavigationLink(destination: VoteMealPlanView(mealPlan: mealPlan)) {
          PendingVoteCardContent(
            mealPlan: mealPlan,
            hasVoted: viewModel.hasUserVoted(on: mealPlan),
            timeUntilDeadline: viewModel.timeUntilDeadline(mealPlan.votingDeadline)
          )
        }
        .buttonStyle(.plain)
      }
    }
  }

  // MARK: - Upcoming Meals Section
  private func upcomingMealsSection(viewModel: HomeViewModel) -> some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
      HStack {
        Label("Upcoming Meals", systemImage: "calendar")
          .font(DSTheme.Typography.title2)
          .foregroundColor(DSTheme.Colors.textPrimary)

        Spacer()

        Text(
          "\(viewModel.upcomingMealPlans.count) plan\(viewModel.upcomingMealPlans.count == 1 ? "" : "s")"
        )
        .font(DSTheme.Typography.caption)
        .foregroundColor(DSTheme.Colors.textSecondary)
      }

      ForEach(viewModel.upcomingMealPlans, id: \.id) { mealPlan in
        NavigationLink(
          destination: MealPlanDetailView(
            mealPlan: mealPlan,
            groceryListItems: nil
          )
        ) {
          UpcomingMealCardContent(mealPlan: mealPlan)
        }
        .buttonStyle(.plain)
      }
    }
  }

  // MARK: - My Tasks Section
  private func myTasksSection(viewModel: HomeViewModel) -> some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
      HStack {
        Label("My Tasks", systemImage: "checklist")
          .font(DSTheme.Typography.title2)
          .foregroundColor(DSTheme.Colors.textPrimary)

        Spacer()

        let unfinishedCount = viewModel.userTasks.filter { $0.status != .finished }.count
        if unfinishedCount > 0 {
          Text("\(unfinishedCount) remaining")
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.warning)
        }
      }

      let groupedTasks = Dictionary(grouping: viewModel.userTasks) { task in
        formatTaskDate(task)
      }
      let sortedDates = groupedTasks.keys.sorted()

      ForEach(sortedDates, id: \.self) { date in
        VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
          Text(date)
            .font(DSTheme.Typography.caption)
            .fontWeight(.medium)
            .foregroundColor(DSTheme.Colors.textSecondary)
            .padding(.horizontal, DSTheme.Spacing.xs)

          ForEach(groupedTasks[date] ?? [], id: \.id) { task in
            TaskCard(task: task)
          }
        }
      }
    }
  }

  // MARK: - Grocery Lists Section
  private func groceryListsSection(viewModel: HomeViewModel) -> some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
      HStack {
        Label("Grocery Lists", systemImage: "cart.fill")
          .font(DSTheme.Typography.title2)
          .foregroundColor(DSTheme.Colors.textPrimary)

        Spacer()
      }

      ForEach(viewModel.activeGroceryLists, id: \.mealPlanID) { groceryList in
        if let mealPlan = viewModel.allMealPlans.first(where: { $0.id == groceryList.mealPlanID }) {
          NavigationLink(
            destination: GroceryListView(
              mealPlan: mealPlan,
              items: [],
              authManager: viewModel.authManager
            )
          ) {
            GroceryListCard(
              mealPlan: mealPlan,
              items: groceryList.items
            )
          }
          .buttonStyle(.plain)
        }
      }
    }
  }

  // MARK: - Empty State
  private var emptyStateView: some View {
    DSEmptyState(
      icon: "calendar.badge.plus",
      title: "No Active Meal Plans",
      message: "Create a meal plan to get started!",
      size: .large
    )
  }

  // MARK: - Helper Functions

  private func formatTaskDate(_ task: Mealplanning_MealPlanTask) -> String {
    let formatter = DateFormatter()
    formatter.dateStyle = .medium
    formatter.timeStyle = .none

    let timestamp = task.completedAt.seconds > 0 ? task.completedAt : task.createdAt
    let date = timestampToDate(timestamp)
    return formatter.string(from: date)
  }

  private func timestampToDate(_ timestamp: SwiftProtobuf.Google_Protobuf_Timestamp) -> Date {
    return HomeViewModel.timestampToDate(timestamp)
  }

  static func formatMealPlanTimeRange(_ mealPlan: Mealplanning_MealPlan) -> String {
    guard !mealPlan.events.isEmpty else {
      return ""
    }

    let earliestStart =
      mealPlan.events.map { HomeViewModel.timestampToDate($0.startsAt) }.min() ?? Date()
    let latestEnd = mealPlan.events.map { HomeViewModel.timestampToDate($0.endsAt) }.max() ?? Date()

    let dateFormatter = DateFormatter()
    dateFormatter.dateStyle = .medium
    dateFormatter.timeStyle = .none

    let startString = dateFormatter.string(from: earliestStart)

    let calendar = Calendar.current
    if calendar.isDate(earliestStart, inSameDayAs: latestEnd) {
      let timeFormatter = DateFormatter()
      timeFormatter.dateStyle = .none
      timeFormatter.timeStyle = .short
      return
        "\(startString) • \(timeFormatter.string(from: earliestStart)) - \(timeFormatter.string(from: latestEnd))"
    } else {
      let endString = dateFormatter.string(from: latestEnd)
      return "\(startString) - \(endString)"
    }
  }

  static func mealPlanDisplayTitle(_ mealPlan: Mealplanning_MealPlan, fallback: String) -> String {
    let title = mealPlan.notes.trimmingCharacters(in: .whitespacesAndNewlines)
    guard !title.isEmpty else {
      return fallback
    }

    guard mealPlan.events.count == 1, title.hasPrefix("Meal Plan for ") else {
      return title
    }

    let startDate =
      mealPlan.events
      .map { HomeViewModel.timestampToDate($0.startsAt) }
      .min() ?? Date()
    let formatter = DateFormatter()
    formatter.dateStyle = .medium
    formatter.timeStyle = .none

    return "Meal Plan for \(formatter.string(from: startDate))"
  }
}

// MARK: - Quick Access Button

struct QuickAccessButton<Destination: View>: View {
  let icon: String
  let label: String
  let color: Color
  let destination: Destination

  var body: some View {
    NavigationLink(destination: destination) {
      VStack(spacing: DSTheme.Spacing.sm) {
        ZStack {
          RoundedRectangle(cornerRadius: DSTheme.Radius.md)
            .fill(color.opacity(0.12))
            .frame(height: 56)

          Image(systemName: icon)
            .font(.system(size: 22))
            .foregroundColor(color)
        }

        Text(label)
          .font(DSTheme.Typography.caption)
          .fontWeight(.medium)
          .foregroundColor(DSTheme.Colors.textPrimary)
      }
      .frame(maxWidth: .infinity)
    }
    .buttonStyle(.plain)
  }
}

// MARK: - Pending Vote Card
struct PendingVoteCard: View {
  let mealPlan: Mealplanning_MealPlan
  let hasVoted: Bool
  let timeUntilDeadline: String
  let onTap: () -> Void

  var body: some View {
    Button(action: onTap) {
      PendingVoteCardContent(
        mealPlan: mealPlan,
        hasVoted: hasVoted,
        timeUntilDeadline: timeUntilDeadline
      )
    }
    .buttonStyle(.plain)
  }
}

// MARK: - Pending Vote Card Content (reusable for NavigationLink)
struct PendingVoteCardContent: View {
  let mealPlan: Mealplanning_MealPlan
  let hasVoted: Bool
  let timeUntilDeadline: String

  var body: some View {
    HStack(spacing: DSTheme.Spacing.md) {
      // Status indicator
      ZStack {
        Circle()
          .fill(
            hasVoted ? DSTheme.Colors.success.opacity(0.15) : DSTheme.Colors.warning.opacity(0.15)
          )
          .frame(width: 44, height: 44)

        Image(systemName: hasVoted ? "checkmark" : "hand.raised.fill")
          .font(.system(size: 18, weight: .medium))
          .foregroundColor(hasVoted ? DSTheme.Colors.success : DSTheme.Colors.warning)
      }

      VStack(alignment: .leading, spacing: DSTheme.Spacing.xxs) {
        Text(HomeView.mealPlanDisplayTitle(mealPlan, fallback: "Meal Plan"))
          .font(DSTheme.Typography.label)
          .foregroundColor(DSTheme.Colors.textPrimary)

        Text(HomeView.formatMealPlanTimeRange(mealPlan))
          .font(DSTheme.Typography.caption)
          .foregroundColor(DSTheme.Colors.textSecondary)

        HStack(spacing: DSTheme.Spacing.xs) {
          Image(systemName: "clock")
            .font(.system(size: 11))
          Text(timeUntilDeadline)
            .font(DSTheme.Typography.caption)
        }
        .foregroundColor(hasVoted ? DSTheme.Colors.textTertiary : DSTheme.Colors.warning)
      }

      Spacer()

      if !hasVoted {
        Text("Vote")
          .font(DSTheme.Typography.caption)
          .fontWeight(.semibold)
          .foregroundColor(.white)
          .padding(.horizontal, DSTheme.Spacing.md)
          .padding(.vertical, DSTheme.Spacing.sm)
          .background(DSTheme.Colors.primary)
          .cornerRadius(DSTheme.Radius.full)
      } else {
        Image(systemName: "chevron.right")
          .font(.system(size: 14, weight: .semibold))
          .foregroundColor(DSTheme.Colors.textTertiary)
      }
    }
    .padding(DSTheme.Spacing.md)
    .background(DSTheme.Colors.cardBackground)
    .cornerRadius(DSTheme.Radius.lg)
    .overlay(
      RoundedRectangle(cornerRadius: DSTheme.Radius.lg)
        .stroke(
          hasVoted ? DSTheme.Colors.border : DSTheme.Colors.warning.opacity(0.3),
          lineWidth: 1
        )
    )
  }
}

// MARK: - Upcoming Meal Card Content
struct UpcomingMealCardContent: View {
  let mealPlan: Mealplanning_MealPlan

  var body: some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
      HStack {
        VStack(alignment: .leading, spacing: DSTheme.Spacing.xxs) {
          Text(HomeView.mealPlanDisplayTitle(mealPlan, fallback: "Meal Plan"))
            .font(DSTheme.Typography.label)
            .foregroundColor(DSTheme.Colors.textPrimary)

          Text(HomeView.formatMealPlanTimeRange(mealPlan))
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.textSecondary)
        }

        Spacer()

        Image(systemName: "chevron.right")
          .font(.system(size: 14, weight: .semibold))
          .foregroundColor(DSTheme.Colors.textTertiary)
      }

      if !mealPlan.events.isEmpty {
        Divider()

        VStack(spacing: DSTheme.Spacing.xs) {
          ForEach(mealPlan.events.prefix(3), id: \.id) { event in
            HStack {
              Circle()
                .fill(mealColor(for: event.mealName))
                .frame(width: 8, height: 8)

              Text(MealPlanningUtils.formatMealName(event.mealName))
                .font(DSTheme.Typography.body)
                .foregroundColor(DSTheme.Colors.textPrimary)

              Spacer()

              Text(formatEventDate(event))
                .font(DSTheme.Typography.caption)
                .foregroundColor(DSTheme.Colors.textSecondary)
            }
          }
        }

        if mealPlan.events.count > 3 {
          Text("+ \(mealPlan.events.count - 3) more")
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.textTertiary)
        }
      }
    }
    .padding(DSTheme.Spacing.md)
    .background(DSTheme.Colors.cardBackground)
    .cornerRadius(DSTheme.Radius.lg)
    .overlay(
      RoundedRectangle(cornerRadius: DSTheme.Radius.lg)
        .stroke(DSTheme.Colors.border, lineWidth: 1)
    )
  }

  private func formatEventDate(_ event: Mealplanning_MealPlanEvent) -> String {
    let formatter = DateFormatter()
    formatter.dateStyle = .short
    formatter.timeStyle = .short

    let date = HomeViewModel.timestampToDate(event.startsAt)
    return formatter.string(from: date)
  }

  private func mealColor(for mealName: Mealplanning_MealPlanEventName) -> Color {
    switch mealName {
    case .breakfast, .secondBreakfast:
      return .orange
    case .brunch:
      return .yellow
    case .lunch:
      return DSTheme.Colors.primary
    case .supper, .dinner:
      return .purple
    default:
      return DSTheme.Colors.secondary
    }
  }
}

// MARK: - Upcoming Meal Card (legacy wrapper)
struct UpcomingMealCard: View {
  let mealPlan: Mealplanning_MealPlan
  let onTap: () -> Void

  var body: some View {
    UpcomingMealCardContent(mealPlan: mealPlan)
  }
}

// MARK: - Task Card
struct TaskCard: View {
  let task: Mealplanning_MealPlanTask

  var body: some View {
    HStack(spacing: DSTheme.Spacing.md) {
      ZStack {
        Circle()
          .fill(
            task.status == .finished
              ? DSTheme.Colors.success.opacity(0.15) : DSTheme.Colors.warning.opacity(0.15)
          )
          .frame(width: 36, height: 36)

        Image(systemName: task.status == .finished ? "checkmark" : "circle")
          .font(.system(size: 14, weight: .medium))
          .foregroundColor(
            task.status == .finished ? DSTheme.Colors.success : DSTheme.Colors.warning)
      }

      VStack(alignment: .leading, spacing: DSTheme.Spacing.xxs) {
        Text(task.creationExplanation)
          .font(DSTheme.Typography.body)
          .foregroundColor(DSTheme.Colors.textPrimary)

        if !task.statusExplanation.isEmpty {
          Text(task.statusExplanation)
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.textSecondary)
        }
      }

      Spacer()
    }
    .padding(DSTheme.Spacing.md)
    .background(DSTheme.Colors.cardBackground)
    .cornerRadius(DSTheme.Radius.md)
    .overlay(
      RoundedRectangle(cornerRadius: DSTheme.Radius.md)
        .stroke(DSTheme.Colors.border, lineWidth: 1)
    )
  }
}

// MARK: - Grocery List Card
struct GroceryListCard: View {
  let mealPlan: Mealplanning_MealPlan
  let items: [Mealplanning_MealPlanGroceryListItem]

  private var itemsToShow: [Mealplanning_MealPlanGroceryListItem] {
    Array(items.prefix(3))
  }

  var body: some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
      HStack {
        VStack(alignment: .leading, spacing: DSTheme.Spacing.xxs) {
          Text(HomeView.mealPlanDisplayTitle(mealPlan, fallback: "Grocery List"))
            .font(DSTheme.Typography.label)
            .foregroundColor(DSTheme.Colors.textPrimary)

          Text(HomeView.formatMealPlanTimeRange(mealPlan))
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.textSecondary)
        }

        Spacer()

        Text("\(items.count)")
          .font(DSTheme.Typography.caption)
          .fontWeight(.semibold)
          .foregroundColor(DSTheme.Colors.primary)
          .padding(.horizontal, DSTheme.Spacing.sm)
          .padding(.vertical, DSTheme.Spacing.xxs)
          .background(DSTheme.Colors.primary.opacity(0.1))
          .cornerRadius(DSTheme.Radius.full)
      }

      if !itemsToShow.isEmpty {
        Divider()

        ForEach(itemsToShow, id: \.id) { item in
          GroceryItemRow(item: item)
        }

        if items.count > 3 {
          Text("+ \(items.count - 3) more items")
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.textTertiary)
        }
      }
    }
    .padding(DSTheme.Spacing.md)
    .background(DSTheme.Colors.cardBackground)
    .cornerRadius(DSTheme.Radius.lg)
    .overlay(
      RoundedRectangle(cornerRadius: DSTheme.Radius.lg)
        .stroke(DSTheme.Colors.border, lineWidth: 1)
    )
  }
}

// MARK: - Grocery Item Row
struct GroceryItemRow: View {
  let item: Mealplanning_MealPlanGroceryListItem

  var body: some View {
    HStack {
      Circle()
        .stroke(DSTheme.Colors.textTertiary, lineWidth: 1.5)
        .frame(width: 16, height: 16)

      Text(item.ingredient.name)
        .font(DSTheme.Typography.body)
        .foregroundColor(DSTheme.Colors.textPrimary)

      Spacer()

      if item.hasQuantityNeeded && item.quantityNeeded.hasMax {
        Text(formatQuantity(item.quantityNeeded))
          .font(DSTheme.Typography.caption)
          .foregroundColor(DSTheme.Colors.textSecondary)
      }
    }
  }

  private func formatQuantity(_ quantity: Common_Float32RangeWithOptionalMax) -> String {
    if quantity.hasMax {
      return "\(quantity.min) - \(quantity.max)"
    } else {
      return "\(quantity.min)+"
    }
  }
}

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "John Doe"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  return HomeView()
    .environment(authManager)
}
