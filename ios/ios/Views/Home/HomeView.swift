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
  @Environment(UserSettingsService.self) private var userSettingsService
  @State private var viewModel: HomeViewModel?

  var body: some View {
    NavigationStack {
      Group {
        if let viewModel = viewModel {
          DSContentState(
            isLoading: viewModel.isLoading,
            loadingMessage: "Loading...",
            error: viewModel.errorMessage,
            errorTitle: viewModel.errorTitle,
            errorIcon: viewModel.errorIcon,
            errorIconColor: viewModel.errorIconColor,
            onRetry: { await viewModel.loadData() },
            content: {
              VStack(spacing: 0) {
                // Sticky header: Greeting
                greetingSection(viewModel: viewModel)
                  .dsScreenPadding()
                  .padding(.bottom, DSTheme.Spacing.md)

                // Sticky: Active Meal Plan
                if let activeMealPlan = viewModel.activeMealPlan {
                  activeMealPlanSection(viewModel: viewModel, mealPlan: activeMealPlan)
                    .dsScreenPadding()
                    .padding(.bottom, DSTheme.Spacing.lg)
                }

                // Sticky: Pending Votes
                if !viewModel.pendingVoteMealPlans.isEmpty {
                  pendingVotesSection(viewModel: viewModel)
                    .dsScreenPadding()
                    .padding(.bottom, DSTheme.Spacing.lg)
                }

                // Sticky: Upcoming Meal Plans (non-finalized)
                if !viewModel.upcomingMealPlans.isEmpty {
                  upcomingMealPlansSection(viewModel: viewModel)
                    .dsScreenPadding()
                    .padding(.bottom, DSTheme.Spacing.lg)
                }

                // Scrollable: Future Meal Plans only (or fill space)
                if !viewModel.futureFinalizedMealPlans.isEmpty {
                  softSeparator

                  ScrollView {
                    futureMealPlansSection(viewModel: viewModel)
                      .dsScreenPadding()
                      .padding(.bottom, DSTheme.Spacing.xl)
                  }
                  .frame(maxHeight: .infinity)
                } else if viewModel.pendingVoteMealPlans.isEmpty
                  && viewModel.activeMealPlan == nil
                  && viewModel.upcomingMealPlans.isEmpty
                  && viewModel.futureFinalizedMealPlans.isEmpty
                {
                  // Empty State
                  emptyStateView
                    .dsScreenPadding()
                    .frame(maxHeight: .infinity)
                } else {
                  // Active/pending/upcoming but no future - fill space so footer stays bottom
                  Spacer()
                    .frame(maxHeight: .infinity)
                }

                // Sticky footer: Create Meal Plan + Recipes/Meals
                VStack(spacing: DSTheme.Spacing.lg) {
                  Divider()

                  createMealPlanSection(viewModel: viewModel)

                  quickAccessRow
                }
                .dsScreenPadding()
                .background(Color(.systemBackground))
              }
            })
        } else {
          DSInitializingView()
        }
      }
      .navigationTitle("")
      .navigationBarTitleDisplayMode(.inline)
      .toolbar {
        ToolbarItem(placement: .navigationBarTrailing) {
          NavigationLink(destination: AccountSettingsView()) {
            DSAvatar(
              name: viewModel?.currentUserDisplayName ?? authManager.username,
              size: .sm,
              imageURL: viewModel.flatMap { homeViewModel in
                guard let user = homeViewModel.currentUser,
                  user.hasAvatar,
                  !user.avatar.storagePath.isEmpty
                else { return nil }
                return APIConfiguration.mediaURL(
                  forStoragePath: user.avatar.storagePath, bucket: "avatars")
              }
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
      .onReceive(NotificationCenter.default.publisher(for: .mealPlanArchived)) { _ in
        if let viewModel = viewModel {
          Task {
            await viewModel.loadData()
          }
        }
      }
    }
  }

  // MARK: - Greeting
  private func greetingSection(viewModel: HomeViewModel) -> some View {
    Text("\(greeting), \(viewModel.currentUserDisplayName)!")
      .font(DSTheme.Typography.largeTitle)
      .foregroundColor(DSTheme.Colors.textPrimary)
      .frame(maxWidth: .infinity, alignment: .leading)
  }

  // MARK: - Create Meal Plan CTA
  private func createMealPlanSection(viewModel: HomeViewModel) -> some View {
    NavigationLink(
      destination: CreateMealPlanWizardView(
        acceptedOccupiedDates: viewModel.acceptedOccupiedDates,
        proposedOccupiedDates: viewModel.proposedOccupiedDates
      )
    ) {
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

  // MARK: - Active Meal Plan Section
  private func activeMealPlanSection(viewModel: HomeViewModel, mealPlan: Mealplanning_MealPlan)
    -> some View
  {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
      Label("Active Meal Plan", systemImage: "star.circle.fill")
        .font(DSTheme.Typography.title2)
        .foregroundColor(DSTheme.Colors.textPrimary)

      NavigationLink(
        destination: MealPlanDetailView(
          mealPlan: mealPlan,
          groceryListItems: nil
        )
      ) {
        UpcomingMealCardContent(mealPlan: mealPlan)
      }
      .buttonStyle(.plain)

      // Task and grocery for active plan only
      activePlanTaskAndGrocery(viewModel: viewModel, mealPlan: mealPlan)
    }
  }

  private func activePlanTaskAndGrocery(
    viewModel: HomeViewModel, mealPlan: Mealplanning_MealPlan
  ) -> some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
      if let taskEntry = viewModel.mealPlansWithTasks.first(where: { $0.mealPlanID == mealPlan.id })
      {
        let summary = taskSummary(
          tasks: taskEntry.tasks, mealPlan: mealPlan, includeDateForFuture: false)
        NavigationLink(
          destination: TaskListView(
            mealPlan: mealPlan,
            tasks: [],
            authManager: viewModel.authManager,
            userSettingsService: userSettingsService
          )
        ) {
          InfoButton(icon: "checklist", text: summary.text, color: summary.color)
        }
        .buttonStyle(.plain)
      }

      if let groceryEntry = viewModel.mealPlansWithGroceryLists.first(where: {
        $0.mealPlanID == mealPlan.id
      }) {
        let neededCount = groceryEntry.items.filter { $0.status == .needs || $0.status == .unknown }
          .count
        NavigationLink(
          destination: GroceryListView(
            mealPlan: mealPlan,
            items: [],
            authManager: viewModel.authManager
          )
        ) {
          InfoButton(
            icon: "cart.fill",
            text: neededCount > 0
              ? "Grocery List (\(neededCount) ingredient\(neededCount == 1 ? "" : "s") needed)"
              : "All ingredients acquired",
            color: neededCount > 0 ? DSTheme.Colors.primary : DSTheme.Colors.success
          )
        }
        .buttonStyle(.plain)
      }
    }
  }

  private var softSeparator: some View {
    VStack(spacing: 0) {
      Spacer()
        .frame(height: DSTheme.Spacing.lg)
      Rectangle()
        .fill(DSTheme.Colors.border.opacity(0.5))
        .frame(height: 1)
        .frame(maxWidth: .infinity)
        .padding(.horizontal, DSTheme.Spacing.xl * 2)
      Spacer()
        .frame(height: DSTheme.Spacing.lg)
    }
  }

  private func futureMealPlansSection(viewModel: HomeViewModel) -> some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
      HStack {
        Label("Future Meal Plans", systemImage: "calendar.badge.clock")
          .font(DSTheme.Typography.title2)
          .foregroundColor(DSTheme.Colors.textPrimary)

        Spacer()

        Text(
          "\(viewModel.futureFinalizedMealPlans.count) plan\(viewModel.futureFinalizedMealPlans.count == 1 ? "" : "s")"
        )
        .font(DSTheme.Typography.caption)
        .foregroundColor(DSTheme.Colors.textSecondary)
      }

      ForEach(viewModel.futureFinalizedMealPlans, id: \.id) { mealPlan in
        futureMealPlanBlock(viewModel: viewModel, mealPlan: mealPlan)
      }
    }
  }

  private func futureMealPlanBlock(
    viewModel: HomeViewModel, mealPlan: Mealplanning_MealPlan
  ) -> some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
      NavigationLink(
        destination: MealPlanDetailView(
          mealPlan: mealPlan,
          groceryListItems: nil
        )
      ) {
        UpcomingMealCardContent(mealPlan: mealPlan)
      }
      .buttonStyle(.plain)

      VStack(alignment: .leading, spacing: DSTheme.Spacing.xs) {
        if let taskEntry = viewModel.mealPlansWithTasks.first(where: {
          $0.mealPlanID == mealPlan.id
        }) {
          let summary = taskSummary(
            tasks: taskEntry.tasks, mealPlan: mealPlan, includeDateForFuture: true)
          NavigationLink(
            destination: TaskListView(
              mealPlan: mealPlan,
              tasks: [],
              authManager: viewModel.authManager,
              userSettingsService: userSettingsService
            )
          ) {
            InfoButton(icon: "checklist", text: summary.text, color: summary.color)
          }
          .buttonStyle(.plain)
        }

        if let groceryEntry = viewModel.mealPlansWithGroceryLists.first(where: {
          $0.mealPlanID == mealPlan.id
        }) {
          let neededCount = groceryEntry.items.filter {
            $0.status == .needs || $0.status == .unknown
          }.count
          NavigationLink(
            destination: GroceryListView(
              mealPlan: mealPlan,
              items: [],
              authManager: viewModel.authManager
            )
          ) {
            InfoButton(
              icon: "cart.fill",
              text: neededCount > 0
                ? "Grocery List (\(neededCount) ingredient\(neededCount == 1 ? "" : "s") needed)"
                : "All ingredients acquired",
              color: neededCount > 0 ? DSTheme.Colors.primary : DSTheme.Colors.success
            )
          }
          .buttonStyle(.plain)
        }
      }
    }
  }

  // MARK: - Upcoming Meal Plans Section
  private func upcomingMealPlansSection(viewModel: HomeViewModel) -> some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
      HStack {
        Label("Upcoming Meal Plans", systemImage: "calendar")
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

  private func taskSummary(
    tasks: [Mealplanning_MealPlanTask], mealPlan: Mealplanning_MealPlan,
    includeDateForFuture: Bool = false
  ) -> (text: String, color: Color) {
    let now = Date()
    let unfinished = tasks.filter { $0.status == .unfinished }

    guard !unfinished.isEmpty else {
      return ("All tasks done", DSTheme.Colors.success)
    }

    var readyCount = 0
    var earliestStart: Date?

    for task in unfinished {
      if let startDate = taskStartDate(task: task, mealPlan: mealPlan) {
        if now >= startDate {
          readyCount += 1
        } else if earliestStart.map({ startDate < $0 }) ?? true {
          earliestStart = startDate
        }
      } else {
        readyCount += 1
      }
    }

    if readyCount > 0 {
      return (
        "Prep Tasks (\(readyCount) ready)",
        DSTheme.Colors.warning
      )
    }

    if let earliest = earliestStart {
      let timeFormatter = DateFormatter()
      timeFormatter.dateStyle = .none
      timeFormatter.timeStyle = .short
      if includeDateForFuture {
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = "M/d/yy"
        return (
          "Next task on \(dateFormatter.string(from: earliest)) at \(timeFormatter.string(from: earliest))",
          DSTheme.Colors.textSecondary
        )
      } else {
        return (
          "Next task at \(timeFormatter.string(from: earliest))",
          DSTheme.Colors.textSecondary
        )
      }
    }

    return ("No tasks ready yet", DSTheme.Colors.textSecondary)
  }

  private func taskStartDate(
    task: Mealplanning_MealPlanTask, mealPlan: Mealplanning_MealPlan
  ) -> Date? {
    guard task.hasRecipePrepTask,
      task.recipePrepTask.hasTimeBufferBeforeRecipeInSeconds,
      task.recipePrepTask.timeBufferBeforeRecipeInSeconds.hasMax,
      task.hasMealPlanOption
    else { return nil }
    let eventID = task.mealPlanOption.belongsToMealPlanEvent
    guard !eventID.isEmpty,
      let event = mealPlan.events.first(where: { $0.id == eventID })
    else { return nil }
    let eventTime = HomeViewModel.timestampToDate(event.startsAt)
    return eventTime.addingTimeInterval(
      -Double(task.recipePrepTask.timeBufferBeforeRecipeInSeconds.max))
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

  /// Compact date range for subtitle when title already conveys meal info (e.g. "Mar 12–14" or "Tue, Thu 7:00 PM").
  static func formatMealPlanTimeRangeCompact(_ mealPlan: Mealplanning_MealPlan) -> String {
    guard !mealPlan.events.isEmpty else {
      return ""
    }

    let eventDates = mealPlan.events.map { HomeViewModel.timestampToDate($0.startsAt) }.sorted()
    let earliestStart = eventDates.first ?? Date()
    let latestEnd = mealPlan.events.map { HomeViewModel.timestampToDate($0.endsAt) }.max() ?? Date()

    let dateFormatter = DateFormatter()
    dateFormatter.dateFormat = "MMM d"

    let calendar = Calendar.current
    if calendar.isDate(earliestStart, inSameDayAs: latestEnd) {
      let timeFormatter = DateFormatter()
      timeFormatter.dateStyle = .none
      timeFormatter.timeStyle = .short
      return
        "\(dateFormatter.string(from: earliestStart)) • \(timeFormatter.string(from: earliestStart))"
    }

    if eventDates.count <= 3 {
      let weekdayFormatter = DateFormatter()
      weekdayFormatter.dateFormat = "EEE"
      let timeFormatter = DateFormatter()
      timeFormatter.dateStyle = .none
      timeFormatter.timeStyle = .short
      let parts = eventDates.map {
        "\(weekdayFormatter.string(from: $0)) \(timeFormatter.string(from: $0))"
      }
      return parts.joined(separator: ", ")
    }

    let startString = dateFormatter.string(from: earliestStart)
    let endString = dateFormatter.string(from: latestEnd)
    return "\(startString)–\(endString)"
  }

  /// Display names from chosen meal options (meal name or recipe names).
  static func chosenMealDisplayNames(from mealPlan: Mealplanning_MealPlan) -> [String] {
    mealPlan.events.compactMap { event in
      guard let chosen = event.options.first(where: { $0.chosen }) else { return nil }
      let meal = chosen.meal
      if !meal.name.isEmpty {
        return meal.name
      }
      let recipeNames = meal.components.compactMap { comp -> String? in
        comp.recipe.name.isEmpty ? nil : comp.recipe.name
      }
      return recipeNames.isEmpty ? nil : recipeNames.joined(separator: ", ")
    }
  }

  /// Whether notes is the auto-generated default from the wizard.
  static func isDefaultMealPlanTitle(_ notes: String, mealPlan: Mealplanning_MealPlan) -> Bool {
    let trimmed = notes.trimmingCharacters(in: .whitespacesAndNewlines)
    guard !trimmed.isEmpty else { return true }

    let dateRange = formatMealPlanTimeRange(mealPlan)
    guard !dateRange.isEmpty else { return false }

    if trimmed == "Dinners \(dateRange)" {
      return true
    }

    if mealPlan.events.count == 1, trimmed.hasPrefix("Meal Plan for ") {
      let startDate =
        mealPlan.events.map { HomeViewModel.timestampToDate($0.startsAt) }.min() ?? Date()
      let formatter = DateFormatter()
      formatter.dateStyle = .medium
      formatter.timeStyle = .none
      return trimmed == "Meal Plan for \(formatter.string(from: startDate))"
    }

    return false
  }

  static func mealPlanDisplayTitle(_ mealPlan: Mealplanning_MealPlan, fallback: String) -> String {
    let title = mealPlan.notes.trimmingCharacters(in: .whitespacesAndNewlines)

    if !title.isEmpty && !Self.isDefaultMealPlanTitle(title, mealPlan: mealPlan) {
      return title
    }

    let names = Self.chosenMealDisplayNames(from: mealPlan)
    if !names.isEmpty {
      return names.joined(separator: " & ")
    }

    guard !mealPlan.events.isEmpty else {
      return fallback
    }

    let earliestStart =
      mealPlan.events.map { HomeViewModel.timestampToDate($0.startsAt) }.min() ?? Date()
    let latestEnd = mealPlan.events.map { HomeViewModel.timestampToDate($0.endsAt) }.max() ?? Date()
    let dateFormatter = DateFormatter()
    dateFormatter.dateFormat = "MMM d"
    let calendar = Calendar.current
    let compactRange: String
    if calendar.isDate(earliestStart, inSameDayAs: latestEnd) {
      compactRange = dateFormatter.string(from: earliestStart)
    } else {
      compactRange =
        "\(dateFormatter.string(from: earliestStart))–\(dateFormatter.string(from: latestEnd))"
    }
    let eventCount = mealPlan.events.count
    let mealType = eventCount == 1 ? "dinner" : "\(eventCount) dinners"
    return "\(mealType.capitalized) • \(compactRange)"
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

  private var usesDerivedTitle: Bool {
    let notes = mealPlan.notes.trimmingCharacters(in: .whitespacesAndNewlines)
    return notes.isEmpty || HomeView.isDefaultMealPlanTitle(notes, mealPlan: mealPlan)
  }

  private var subtitleText: String {
    usesDerivedTitle
      ? HomeView.formatMealPlanTimeRangeCompact(mealPlan)
      : HomeView.formatMealPlanTimeRange(mealPlan)
  }

  var body: some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
      HStack {
        VStack(alignment: .leading, spacing: DSTheme.Spacing.xxs) {
          Text(HomeView.mealPlanDisplayTitle(mealPlan, fallback: "Meal Plan"))
            .font(DSTheme.Typography.label)
            .foregroundColor(DSTheme.Colors.textPrimary)

          if !subtitleText.isEmpty {
            Text(subtitleText)
              .font(DSTheme.Typography.caption)
              .foregroundColor(DSTheme.Colors.textSecondary)
          }
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

              Text(eventDisplayLabel(event))
                .font(DSTheme.Typography.body)
                .foregroundColor(DSTheme.Colors.textPrimary)

              Spacer()

              Text(formatEventDate(event, compact: mealPlan.status == .finalized))
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

  private func eventDisplayLabel(_ event: Mealplanning_MealPlanEvent) -> String {
    guard mealPlan.status == .finalized,
      let chosen = event.options.first(where: { $0.chosen })
    else {
      return MealPlanningUtils.formatMealName(event.mealName)
    }
    let meal = chosen.meal
    if !meal.name.isEmpty {
      return meal.name
    }
    let recipeNames = meal.components.compactMap { comp -> String? in
      comp.recipe.name.isEmpty ? nil : comp.recipe.name
    }
    if !recipeNames.isEmpty {
      return recipeNames.joined(separator: ", ")
    }
    return MealPlanningUtils.formatMealName(event.mealName)
  }

  private func formatEventDate(_ event: Mealplanning_MealPlanEvent, compact: Bool) -> String {
    let date = HomeViewModel.timestampToDate(event.startsAt)
    if compact {
      let weekdayFormatter = DateFormatter()
      weekdayFormatter.dateFormat = "EEE"
      let timeFormatter = DateFormatter()
      timeFormatter.dateStyle = .none
      timeFormatter.timeStyle = .short
      return "\(weekdayFormatter.string(from: date)) \(timeFormatter.string(from: date))"
    }
    let formatter = DateFormatter()
    formatter.dateStyle = .short
    formatter.timeStyle = .short
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

// MARK: - Info Button
struct InfoButton: View {
  let icon: String
  let text: String
  let color: Color

  var body: some View {
    HStack(spacing: DSTheme.Spacing.md) {
      ZStack {
        Circle()
          .fill(color.opacity(0.15))
          .frame(width: 36, height: 36)

        Image(systemName: icon)
          .font(.system(size: 15, weight: .medium))
          .foregroundColor(color)
      }

      Text(text)
        .font(DSTheme.Typography.label)
        .foregroundColor(DSTheme.Colors.textPrimary)

      Spacer()

      Image(systemName: "chevron.right")
        .font(.system(size: 13, weight: .semibold))
        .foregroundColor(DSTheme.Colors.textTertiary)
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

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "John Doe"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  return HomeView()
    .environment(authManager)
}
