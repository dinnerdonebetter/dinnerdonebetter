//
//  MealPlanDetailView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct MealPlanDetailView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Environment(\.dismiss) private var dismiss
  let mealPlan: Mealplanning_MealPlan
  let groceryListItems: [Mealplanning_MealPlanGroceryListItem]?
  @State private var taskCount: Int?
  @State private var groceryListItemCount: Int?
  @State private var showArchiveConfirmation = false
  @State private var isArchiving = false
  @State private var archiveError: String?

  var body: some View {
    ScrollView {
      VStack(alignment: .leading, spacing: DSTheme.Spacing.xl) {
        // Header with title, status, and range inline
        headerSection

        // Voting deadline (if not finalized) - shown in a card if needed
        if mealPlan.status == .awaitingVotes {
          votingDeadlineCard
        }

        // Events
        eventsSection

        // Grocery List Link (if finalized and has items)
        if mealPlan.status == .finalized {
          groceryListSection
        }

        // Tasks Link (if finalized and tasks created)
        if mealPlan.status == .finalized && mealPlan.tasksCreated {
          tasksSection
        }
      }
      .dsScreenPadding()
    }
    .navigationBarTitleDisplayMode(.inline)
    .toolbar {
      ToolbarItem(placement: .primaryAction) {
        Button {
          showArchiveConfirmation = true
        } label: {
          Label("Archive", systemImage: "archivebox")
        }
      }
    }
    .alert("Archive Meal Plan", isPresented: $showArchiveConfirmation) {
      Button("Archive", role: .destructive) {
        Task { await archiveMealPlan() }
      }
      Button("Cancel", role: .cancel) {}
    } message: {
      Text("This will remove the meal plan from your active plans. You can still access archived plans later.")
    }
    .alert("Archive Failed", isPresented: .constant(archiveError != nil)) {
      Button("OK") { archiveError = nil }
    } message: {
      if let error = archiveError {
        Text(error)
      }
    }
    .disabled(isArchiving)
    .task {
      await loadTaskCount()
      await loadGroceryListItemCount()
    }
  }

  private func archiveMealPlan() async {
    isArchiving = true
    archiveError = nil

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        archiveError = "Failed to connect"
        isArchiving = false
        return
      }

      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        archiveError = "Please sign in again"
        isArchiving = false
        return
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      var request = Mealplanning_ArchiveMealPlanRequest()
      request.mealPlanID = mealPlan.id

      _ = try await clientManager.client.mealPlanning.archiveMealPlan(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      NotificationCenter.default.post(name: .mealPlanArchived, object: nil)
      dismiss()
    } catch {
      archiveError = error.localizedDescription
    }

    isArchiving = false
  }

  private func loadTaskCount() async {
    guard mealPlan.status == .finalized && mealPlan.tasksCreated else {
      return
    }

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        return
      }

      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        return
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      var request = Mealplanning_GetMealPlanTasksRequest()
      request.mealPlanID = mealPlan.id

      let response = try await clientManager.client.mealPlanning.getMealPlanTasks(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      // Deduplicate by ID and get count
      var seenIDs: Set<String> = []
      let uniqueCount = response.results.filter { task in
        if seenIDs.contains(task.id) {
          return false
        }
        seenIDs.insert(task.id)
        return true
      }.count

      taskCount = uniqueCount
    } catch {
      print("⚠️ Failed to fetch task count: \(error)")
      taskCount = 0
    }
  }

  private func loadGroceryListItemCount() async {
    guard mealPlan.status == .finalized else {
      return
    }

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        return
      }

      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        return
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      var request = Mealplanning_GetMealPlanGroceryListItemsForMealPlanRequest()
      request.mealPlanID = mealPlan.id

      let response = try await clientManager.client.mealPlanning
        .getMealPlanGroceryListItemsForMealPlan(
          request,
          metadata: metadata,
          options: clientManager.defaultCallOptions
        )

      groceryListItemCount = response.results.count
    } catch {
      print("⚠️ Failed to fetch grocery list item count: \(error)")
      groceryListItemCount = 0
    }
  }

  private var headerSection: some View {
    VStack(alignment: .leading, spacing: 8) {
      // Title with status badge and range inline
      HStack(alignment: .firstTextBaseline, spacing: 8) {
        Text(mealPlan.notes.isEmpty ? "Meal Plan" : mealPlan.notes)
          .font(.title2)
          .fontWeight(.bold)

        statusBadge

        Spacer()
      }

      // Time range on second line
      Text(HomeView.formatMealPlanTimeRange(mealPlan))
        .font(.subheadline)
        .foregroundColor(.secondary)
    }
  }

  private var votingDeadlineCard: some View {
    votingDeadlineView
      .padding()
      .background(Color(.systemGray6))
      .cornerRadius(10)
  }

  private var statusBadge: some View {
    let (text, color) = statusInfo(mealPlan.status)
    return Text(text)
      .font(.caption)
      .fontWeight(.semibold)
      .padding(.horizontal, 8)
      .padding(.vertical, 4)
      .background(color.opacity(0.2))
      .foregroundColor(color)
      .cornerRadius(6)
  }

  private var votingDeadlineView: some View {
    let deadline = HomeViewModel.timestampToDate(mealPlan.votingDeadline)
    let formatter = DateFormatter()
    formatter.dateStyle = .medium
    formatter.timeStyle = .short

    return HStack {
      Image(systemName: "clock")
        .foregroundColor(.secondary)
      Text("Voting deadline: \(formatter.string(from: deadline))")
        .font(.subheadline)
        .foregroundColor(.secondary)
    }
  }

  private func statusInfo(_ status: Mealplanning_MealPlanStatus) -> (String, Color) {
    switch status {
    case .awaitingVotes:
      return ("Awaiting Votes", .orange)
    case .finalized:
      return ("Finalized", .green)
    default:
      return ("Unknown", .gray)
    }
  }

  private var upcomingEvents: [Mealplanning_MealPlanEvent] {
    let now = Date()
    return mealPlan.events.filter { HomeViewModel.timestampToDate($0.startsAt) >= now }
  }

  private var pastEvents: [Mealplanning_MealPlanEvent] {
    let now = Date()
    return mealPlan.events.filter { HomeViewModel.timestampToDate($0.startsAt) < now }
  }

  private var eventsSection: some View {
    VStack(alignment: .leading, spacing: 12) {
      if !upcomingEvents.isEmpty {
        Text("Upcoming")
          .font(.title2)
          .fontWeight(.bold)

        ForEach(upcomingEvents, id: \.id) { event in
          EventCard(event: event)
        }
      }

      if !pastEvents.isEmpty {
        Text("Past")
          .font(.title2)
          .fontWeight(.bold)
          .foregroundColor(.secondary)

        ForEach(pastEvents, id: \.id) { event in
          EventCard(event: event)
            .opacity(0.7)
        }
      }
    }
  }

  private var groceryListSection: some View {
    let count = groceryListItemCount ?? 0
    let hasItems = count > 0

    return Group {
      if hasItems {
        NavigationLink(
          destination: GroceryListView(
            mealPlan: mealPlan,
            items: [],  // Always start with empty array, GroceryListView will fetch fresh data
            authManager: authManager
          )
        ) {
          groceryListCardContent(count: count)
        }
        .buttonStyle(.plain)
      } else {
        groceryListCardContent(count: count)
          .opacity(0.6)
      }
    }
  }

  private func groceryListCardContent(count: Int) -> some View {
    HStack {
      Image(systemName: "cart")
        .foregroundColor(.blue)
      Text("View Grocery List\(count > 0 ? " (\(count))" : "")")
        .font(.headline)
      Spacer()
      if count > 0 {
        Image(systemName: "chevron.right")
          .foregroundColor(.secondary)
          .font(.caption)
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }

  private var tasksSection: some View {
    let count = taskCount ?? 0
    let hasTasks = count > 0

    return Group {
      if hasTasks {
        NavigationLink(
          destination: TaskListView(
            mealPlan: mealPlan,
            tasks: [],  // Always start with empty array, TaskListView will fetch fresh data
            authManager: authManager
          )
        ) {
          tasksCardContent(count: count)
        }
        .buttonStyle(.plain)
      } else {
        tasksCardContent(count: count)
          .opacity(0.6)
      }
    }
  }

  private func tasksCardContent(count: Int) -> some View {
    HStack {
      Image(systemName: "checklist")
        .foregroundColor(.blue)
      Text("View Tasks (\(count))")
        .font(.headline)
      Spacer()
      if count > 0 {
        Image(systemName: "chevron.right")
          .foregroundColor(.secondary)
          .font(.caption)
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }
}

// MARK: - Event Card

struct EventCard: View {
  let event: Mealplanning_MealPlanEvent

  private func eventTimeRange(event: Mealplanning_MealPlanEvent) -> some View {
    let startDate = HomeViewModel.timestampToDate(event.startsAt)
    let formatter = DateFormatter()
    formatter.dateStyle = .medium
    formatter.timeStyle = .short

    return Text(formatter.string(from: startDate))
      .font(.caption)
      .foregroundColor(.secondary)
  }

  var body: some View {
    VStack(alignment: .leading, spacing: 12) {
      // Event header
      HStack {
        VStack(alignment: .leading, spacing: 4) {
          Text(MealPlanningUtils.formatMealName(event.mealName))
            .font(.headline)

          eventTimeRange(event: event)
        }

        Spacer()
      }

      // Selected meals
      if !event.options.isEmpty {
        VStack(alignment: .leading, spacing: 8) {
          Text("Meals")
            .font(.subheadline)
            .fontWeight(.medium)
            .foregroundColor(.secondary)

          ForEach(event.options.filter { $0.chosen }, id: \.id) { option in
            NavigationLink(
              destination: MealDetailView(mealID: option.meal.id, isFromMealPlan: true)
            ) {
              MealOptionCard(option: option, isChosen: true)
            }
            .buttonStyle(.plain)
          }
        }
      } else {
        Text("No meals selected")
          .font(.subheadline)
          .foregroundColor(.secondary)
      }
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(10)
  }
}

// MARK: - Meal Option Card

struct MealOptionCard: View {
  let option: Mealplanning_MealPlanOption
  var isChosen: Bool = true

  var body: some View {
    HStack {
      if isChosen {
        Image(systemName: "checkmark.circle.fill")
          .foregroundColor(.green)
          .font(.caption)
      }

      VStack(alignment: .leading, spacing: 4) {
        HStack(spacing: 6) {
          Text(option.meal.name.isEmpty ? "Unnamed Meal" : option.meal.name)
            .font(.subheadline)
            .fontWeight(isChosen ? .semibold : .regular)

          // Tiebroken indicator
          if isChosen && option.tieBroken {
            HStack(spacing: 2) {
              Image(systemName: "dice")
                .font(.caption2)
              Text("Tiebroken")
                .font(.caption2)
            }
            .foregroundColor(.orange)
            .padding(.horizontal, 4)
            .padding(.vertical, 2)
            .background(Color.orange.opacity(0.2))
            .cornerRadius(4)
          }
        }

        if !option.meal.components.isEmpty {
          let recipeNames = option.meal.components.compactMap { component -> String? in
            component.recipe.name.isEmpty ? nil : component.recipe.name
          }
          if !recipeNames.isEmpty {
            Text(recipeNames.joined(separator: ", "))
              .font(.caption)
              .foregroundColor(.secondary)
              .lineLimit(2)
          }
        }

        if option.mealScale != 1.0 {
          Text("Scale: \(String(format: "%.1f", option.mealScale))x")
            .font(.caption2)
            .foregroundColor(.secondary)
        }
      }

      Spacer()

      // Show chevron only for chosen meals (indicating they're clickable)
      if isChosen {
        Image(systemName: "chevron.right")
          .foregroundColor(.secondary)
          .font(.caption)
      }
    }
    .padding(.vertical, 4)
    .opacity(isChosen ? 1.0 : 0.6)
    .contentShape(Rectangle())
  }
}

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "John Doe"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  // Create a sample meal plan
  var mealPlan = Mealplanning_MealPlan()
  mealPlan.id = "mealplan123"
  mealPlan.notes = "Sample Meal Plan"
  mealPlan.status = .finalized
  mealPlan.groceryListInitialized = true
  mealPlan.tasksCreated = true

  var event = Mealplanning_MealPlanEvent()
  event.id = "event123"
  event.mealName = .dinner
  var startTime = Google_Protobuf_Timestamp()
  startTime.seconds = Int64(Date().timeIntervalSince1970)
  event.startsAt = startTime
  var endTime = Google_Protobuf_Timestamp()
  endTime.seconds = Int64(Date().timeIntervalSince1970 + 7200)
  event.endsAt = endTime

  var option = Mealplanning_MealPlanOption()
  option.id = "option123"
  option.chosen = true
  var meal = Mealplanning_Meal()
  meal.name = "Chicken Dinner"
  option.meal = meal
  option.mealScale = 1.0
  event.options = [option]

  mealPlan.events = [event]

  return NavigationView {
    MealPlanDetailView(mealPlan: mealPlan, groceryListItems: nil)
  }
  .environment(authManager)
}
