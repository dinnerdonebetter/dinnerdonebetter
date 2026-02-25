//
//  TaskListView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct TaskListView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @State private var viewModel: TaskListViewModel

  init(
    mealPlan: Mealplanning_MealPlan,
    tasks: [Mealplanning_MealPlanTask],
    authManager: AuthenticationManager
  ) {
    _viewModel = State(
      initialValue: TaskListViewModel(
        mealPlan: mealPlan,
        tasks: tasks,
        authManager: authManager
      )
    )
  }

  var body: some View {
    ScrollView {
      VStack(alignment: .leading, spacing: DSTheme.Spacing.xl) {
        // Header
        headerSection

        if viewModel.isLoading {
          DSLoadingView("Loading tasks...")
        } else if viewModel.tasks.isEmpty {
          emptyStateView
        } else {
          let readyNowGroups = viewModel.getReadyNowGroups()
          let doLaterGroups = viewModel.getDoLaterGroups()
          let finishedGroups = viewModel.getFinishedGroups()

          if !readyNowGroups.isEmpty {
            tasksSection(
              title: "Ready Now",
              groups: readyNowGroups,
              color: DSTheme.Colors.warning
            )
          }

          if !doLaterGroups.isEmpty {
            laterTasksSection(groups: doLaterGroups)
          }

          if !finishedGroups.isEmpty {
            tasksSection(
              title: "Completed",
              groups: finishedGroups,
              color: DSTheme.Colors.success
            )
          }
        }

        if let errorMessage = viewModel.errorMessage {
          Text(errorMessage)
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.error)
            .padding()
        }
      }
      .dsScreenPadding()
    }
    .navigationTitle(viewModel.mealPlan.notes.isEmpty ? "Tasks" : "Tasks")
    .navigationBarTitleDisplayMode(.large)
    .task {
      await viewModel.loadTasks()
    }
  }

  private var headerSection: some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
      Text(HomeView.formatMealPlanTimeRange(viewModel.mealPlan))
        .font(DSTheme.Typography.body)
        .foregroundColor(DSTheme.Colors.textSecondary)

      Text("\(viewModel.tasks.count) task\(viewModel.tasks.count == 1 ? "" : "s")")
        .font(DSTheme.Typography.caption)
        .foregroundColor(DSTheme.Colors.textSecondary)
    }
  }

  private var emptyStateView: some View {
    DSEmptyState(
      icon: "checklist",
      title: "No tasks",
      message: "Tasks will appear here once the meal plan is finalized."
    )
  }

  private func tasksSection(
    title: String,
    groups: [TaskGroup],
    color: Color
  ) -> some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
      Text(title)
        .font(DSTheme.Typography.label)
        .foregroundColor(color)

      ForEach(groups, id: \.parent.id) { group in
        TaskGroupRow(
          group: group,
          viewModel: viewModel,
          loadedRecipes: viewModel.loadedRecipes,
          loadedPrepTasks: viewModel.loadedPrepTasks
        )
      }
    }
  }

  private func laterTasksSection(groups: [TaskGroup]) -> some View {
    let buckets = bucketByStartTime(groups)

    return VStack(alignment: .leading, spacing: DSTheme.Spacing.md) {
      Text("Later")
        .font(DSTheme.Typography.label)
        .foregroundColor(.secondary)

      ForEach(buckets, id: \.label) { bucket in
        VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
          Text(bucket.label)
            .font(DSTheme.Typography.caption)
            .fontWeight(.semibold)
            .foregroundColor(.orange)

          ForEach(bucket.groups, id: \.parent.id) { group in
            TaskGroupRow(
              group: group,
              viewModel: viewModel,
              loadedRecipes: viewModel.loadedRecipes,
              loadedPrepTasks: viewModel.loadedPrepTasks
            )
          }
        }
      }
    }
  }

  private struct TimeBucket {
    let label: String
    let groups: [TaskGroup]
  }

  private func bucketByStartTime(_ groups: [TaskGroup]) -> [TimeBucket] {
    var bucketMap: [(key: String, groups: [TaskGroup])] = []
    let now = Date()

    for group in groups {
      let label: String
      if let startDate = viewModel.canStartAt(task: group.parent) {
        label = formatCountdownLabel(from: now, to: startDate)
      } else {
        label = "Later"
      }

      if let index = bucketMap.firstIndex(where: { $0.key == label }) {
        bucketMap[index].groups.append(group)
      } else {
        bucketMap.append((key: label, groups: [group]))
      }
    }

    return bucketMap.map { TimeBucket(label: $0.key, groups: $0.groups) }
  }

  private func formatCountdownLabel(from now: Date, to date: Date) -> String {
    let seconds = Int(date.timeIntervalSince(now))
    if seconds <= 0 { return "Now" }

    let hours = seconds / 3600
    let days = hours / 24

    if hours < 1 {
      let minutes = seconds / 60
      return "In \(minutes) min"
    } else if hours < 24 {
      return "In \(hours) hr"
    } else if days == 1 {
      return "In 1 day"
    } else {
      return "In \(days) days"
    }
  }
}

// MARK: - Task Group Row

struct TaskGroupRow: View {
  let group: TaskGroup
  let viewModel: TaskListViewModel
  let loadedRecipes: [String: Mealplanning_Recipe]
  let loadedPrepTasks: [String: Mealplanning_RecipePrepTask]

  var hasSubtasks: Bool {
    !group.subtasks.isEmpty
  }

  var isExpanded: Bool {
    viewModel.isExpanded(taskID: group.parent.id)
  }

  var body: some View {
    VStack(alignment: .leading, spacing: 0) {
      // Parent task
      TaskRow(
        task: group.parent,
        viewModel: viewModel,
        loadedRecipes: loadedRecipes,
        loadedPrepTasks: loadedPrepTasks,
        isParent: hasSubtasks,
        isExpanded: isExpanded,
        onToggleExpand: {
          viewModel.toggleExpanded(taskID: group.parent.id)
        }
      )

      // Subtasks (shown when expanded)
      if hasSubtasks && isExpanded {
        ForEach(group.subtasks) { subtask in
          SubtaskRow(
            subtask: subtask,
            parentTask: group.parent,
            viewModel: viewModel
          )
          .padding(.leading, 32)
        }
      }
    }
  }
}

// MARK: - Subtask Row

struct SubtaskRow: View {
  let subtask: SubtaskDisplayInfo
  let parentTask: Mealplanning_MealPlanTask
  let viewModel: TaskListViewModel

  private var isCompleted: Bool {
    viewModel.isSubtaskCompleted(parentTaskID: subtask.parentTaskID, stepID: subtask.stepID)
      || parentTask.status == .finished
  }

  var body: some View {
    HStack(alignment: .top, spacing: 12) {
      Spacer()
        .frame(width: 20)

      // Checkbox
      Button {
        viewModel.toggleSubtaskCompletion(
          parentTaskID: subtask.parentTaskID, stepID: subtask.stepID)
      } label: {
        Image(systemName: isCompleted ? "checkmark.circle.fill" : "circle")
          .font(.body)
          .foregroundColor(isCompleted ? .green : .blue)
      }
      .buttonStyle(.plain)
      .disabled(viewModel.isUpdating)

      // Subtask content -- tappable to navigate when a recipe step is available
      if let recipeID = subtask.recipeID, let recipeStepID = subtask.recipeStepID {
        NavigationLink(
          destination: PerformRecipeView(
            recipeID: recipeID,
            highlightedStepIDs: [recipeStepID]
          )
        ) {
          subtaskContent
        }
        .buttonStyle(.plain)
      } else {
        subtaskContent
      }

      Spacer()

      if subtask.recipeID != nil && subtask.recipeStepID != nil {
        Image(systemName: "chevron.right")
          .font(.caption2)
          .foregroundColor(.secondary)
      }
    }
    .padding()
    .background(Color(.systemGray5))
    .cornerRadius(8)
    .opacity(isCompleted ? 0.7 : 1.0)
  }

  private var subtaskContent: some View {
    let label =
      subtask.ingredientNames.isEmpty
      ? subtask.description
      : "\(subtask.description) \(subtask.ingredientNames.joined(separator: ", "))"

    return Text(label)
      .font(.body)
      .foregroundColor(isCompleted ? .secondary : .primary)
      .strikethrough(isCompleted)
  }
}

// MARK: - Task Row

struct TaskRow: View {
  let task: Mealplanning_MealPlanTask
  let viewModel: TaskListViewModel
  let loadedRecipes: [String: Mealplanning_Recipe]
  let loadedPrepTasks: [String: Mealplanning_RecipePrepTask]
  let isParent: Bool
  let isExpanded: Bool
  let onToggleExpand: (() -> Void)?

  // Extract recipe ID and step IDs from prep task
  private var recipeID: String? {
    guard task.hasRecipePrepTask else { return nil }
    let prepTask = task.recipePrepTask
    var recipeID = prepTask.belongsToRecipe

    if recipeID.isEmpty && !prepTask.id.isEmpty {
      if let loadedPrepTask = loadedPrepTasks[prepTask.id] {
        recipeID = loadedPrepTask.belongsToRecipe
      }
    }

    if recipeID.isEmpty && task.hasMealPlanOption {
      recipeID = findRecipeIDFromMealOption(task: task, prepTask: prepTask)
    }

    return recipeID.isEmpty ? nil : recipeID
  }

  private var highlightedStepIDs: Set<String>? {
    guard task.hasRecipePrepTask else { return nil }
    let prepTask = task.recipePrepTask
    let stepIDs = prepTask.taskSteps
      .filter { !$0.belongsToRecipeStep.isEmpty }
      .map { $0.belongsToRecipeStep }
    return stepIDs.isEmpty ? nil : Set(stepIDs)
  }

  // Get prep task context information
  private var prepTaskContext: PerformRecipeView.PrepTaskContext? {
    guard task.hasRecipePrepTask else { return nil }
    let prepTask = task.recipePrepTask

    // Get event information from meal plan option
    var eventName: String?
    var eventTime: Date?

    if task.hasMealPlanOption {
      let eventID = task.mealPlanOption.belongsToMealPlanEvent
      if !eventID.isEmpty, let event = viewModel.mealPlan.events.first(where: { $0.id == eventID })
      {
        eventName = MealPlanningUtils.formatMealName(event.mealName)
        eventTime = HomeViewModel.timestampToDate(event.startsAt)
      }
    }

    // Get recipe name
    var recipeName: String?
    if let recipeID = recipeID, let recipe = loadedRecipes[recipeID] {
      recipeName = recipe.name
    }

    // Get prep task name
    let prepTaskName = prepTask.name.isEmpty ? nil : prepTask.name

    return PerformRecipeView.PrepTaskContext(
      prepTaskName: prepTaskName,
      recipeName: recipeName,
      eventName: eventName,
      eventTime: eventTime
    )
  }

  var body: some View {
    HStack(alignment: .top, spacing: 12) {
      // Disclosure indicator for parent tasks
      if isParent {
        Button {
          onToggleExpand?()
        } label: {
          Image(systemName: isExpanded ? "chevron.down" : "chevron.right")
            .font(.caption)
            .foregroundColor(.secondary)
            .frame(width: 20)
        }
        .buttonStyle(.plain)
      } else {
        // Spacer for alignment when not a parent
        Spacer()
          .frame(width: 20)
      }

      // Checkbox - enabled for parent tasks when all subtasks are completed, or when already finished (to allow unchecking)
      if isParent {
        let allSubtasksCompleted = viewModel.areAllSubtasksCompleted(parentTaskID: task.id)
        let canCheckParent = allSubtasksCompleted || task.status == .finished

        if canCheckParent {
          // Parent task can be checked when all subtasks are done, or unchecked when finished
          Button {
            print("🔘 Clicked parent task checkbox for task \(task.id)")
            Task {
              await viewModel.toggleTaskStatus(task)
            }
          } label: {
            Image(systemName: task.status == .finished ? "checkmark.circle.fill" : "circle")
              .font(.title2)
              .foregroundColor(
                task.status == .finished ? .green : .blue
              )
              .contentShape(Rectangle())
          }
          .buttonStyle(.plain)
          .disabled(viewModel.isUpdating)
          .frame(minWidth: 44, minHeight: 44)  // Ensure tappable area
        } else {
          // Parent task shows status but can't be checked (not all subtasks done yet)
          Image(systemName: task.status == .finished ? "checkmark.circle.fill" : "circle")
            .font(.title2)
            .foregroundColor(.gray)
        }
      } else {
        // Standalone tasks can be checked
        Button {
          Task {
            await viewModel.toggleTaskStatus(task)
          }
        } label: {
          Image(systemName: task.status == .finished ? "checkmark.circle.fill" : "circle")
            .font(.title2)
            .foregroundColor(
              task.status == .finished ? .green : .blue
            )
        }
        .buttonStyle(.plain)
        .disabled(viewModel.isUpdating)
      }

      // Task content - make clickable if it has a recipe prep task
      let hasNavigation = recipeID != nil && highlightedStepIDs != nil
      let isCompleted = task.status == .finished
      let context = prepTaskContext

      Group {
        if hasNavigation, let recipeID = recipeID, let highlightedStepIDs = highlightedStepIDs {
          NavigationLink(
            destination: PerformRecipeView(
              recipeID: recipeID,
              highlightedStepIDs: highlightedStepIDs,
              prepTaskContext: context
            )
          ) {
            taskDescriptionContent(isCompleted: isCompleted)
          }
          .buttonStyle(.plain)
        } else {
          taskDescriptionContent(isCompleted: isCompleted)
        }
      }

      Spacer()
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(8)
    .opacity(task.status == .finished ? 0.7 : 1.0)
  }

  // Get event start time for countdown
  private var eventStartTime: Date? {
    guard task.hasMealPlanOption else { return nil }
    let eventID = task.mealPlanOption.belongsToMealPlanEvent
    guard !eventID.isEmpty,
      let event = viewModel.mealPlan.events.first(where: { $0.id == eventID })
    else {
      return nil
    }
    return HomeViewModel.timestampToDate(event.startsAt)
  }

  private func taskDescriptionContent(isCompleted: Bool) -> some View {
    VStack(alignment: .leading, spacing: 4) {
      // Task description - show prep task name or creation explanation
      if task.hasRecipePrepTask && !task.recipePrepTask.name.isEmpty {
        Text(task.recipePrepTask.name)
          .font(.body)
          .fontWeight(.medium)
          .strikethrough(isCompleted)
          .foregroundColor(
            isCompleted ? .secondary : .primary
          )
      } else {
        Text(task.creationExplanation)
          .font(.body)
          .fontWeight(.medium)
          .strikethrough(isCompleted)
          .foregroundColor(
            isCompleted ? .secondary : .primary
          )
      }

      if !task.statusExplanation.isEmpty {
        Text(task.statusExplanation)
          .font(.caption)
          .foregroundColor(.secondary)
      }

      // Storage info (when task has a recipe prep task with storage data)
      if task.hasRecipePrepTask {
        storageInfoView(prepTask: task.recipePrepTask, isCompleted: isCompleted)
      }

      // Time window for when this task can/should be started
      if !isCompleted, task.hasRecipePrepTask {
        timeWindowView(prepTask: task.recipePrepTask)
      }

      // Countdown timer (only show if task is not completed and we have an event time)
      if !isCompleted, let eventTime = eventStartTime {
        TaskCountdownTimer(dueDate: eventTime)
      }
    }
  }

  @ViewBuilder
  private func storageInfoView(
    prepTask: Mealplanning_RecipePrepTask,
    isCompleted: Bool
  ) -> some View {
    let hasExplicit = !prepTask.explicitStorageInstructions.isEmpty
    let hasType = !prepTask.storageType.isEmpty
    let temp = formatStorageTemperature(prepTask.storageTemperatureInCelsius)
    let hasTemp = !temp.isEmpty

    if hasExplicit || hasType || hasTemp {
      VStack(alignment: .leading, spacing: 2) {
        if hasExplicit {
          Label(prepTask.explicitStorageInstructions, systemImage: "info.circle")
            .font(.caption)
            .foregroundColor(isCompleted ? .secondary : .secondary)
            .strikethrough(isCompleted)
        }
        if hasType || hasTemp {
          HStack(spacing: 6) {
            if hasType {
              Label(prepTask.storageType, systemImage: "archivebox")
                .font(.caption)
                .foregroundColor(isCompleted ? .secondary : .secondary)
                .strikethrough(isCompleted)
            }
            if hasTemp {
              Label(temp, systemImage: "thermometer.medium")
                .font(.caption)
                .foregroundColor(isCompleted ? .secondary : .secondary)
                .strikethrough(isCompleted)
            }
          }
        }
      }
    }
  }

  private func formatStorageTemperature(_ range: Common_OptionalFloat32Range) -> String {
    let hasMin = range.hasMin
    let hasMax = range.hasMax
    if hasMin && hasMax {
      let minVal = Int(range.min.rounded())
      let maxVal = Int(range.max.rounded())
      return minVal == maxVal ? "\(minVal)°C" : "\(minVal)–\(maxVal)°C"
    }
    if hasMax {
      return "below \(Int(range.max.rounded()))°C"
    }
    if hasMin {
      return "above \(Int(range.min.rounded()))°C"
    }
    return ""
  }

  @ViewBuilder
  private func timeWindowView(prepTask: Mealplanning_RecipePrepTask) -> some View {
    if prepTask.hasTimeBufferBeforeRecipeInSeconds,
      prepTask.timeBufferBeforeRecipeInSeconds.min > 0,
      let eventTime = eventStartTime
    {
      let mustStartBy = eventTime.addingTimeInterval(
        -Double(prepTask.timeBufferBeforeRecipeInSeconds.min))

      HStack(spacing: 6) {
        Image(systemName: "calendar")
          .font(.caption2)
        Text("Start by \(formatShortDate(mustStartBy))")
          .font(.caption)
      }
      .foregroundColor(.secondary)
    }
  }

  private func formatShortDate(_ date: Date) -> String {
    let formatter = DateFormatter()
    let calendar = Calendar.current
    if calendar.isDateInToday(date) {
      formatter.dateFormat = "'today at' h:mm a"
    } else if calendar.isDateInTomorrow(date) {
      formatter.dateFormat = "'tomorrow at' h:mm a"
    } else {
      formatter.dateFormat = "EEE h:mm a"
    }
    return formatter.string(from: date)
  }

  // Helper struct for step data
  struct StepInfo {
    let stepID: String
    let description: String
  }

  private func getStepData(task: Mealplanning_MealPlanTask) -> [StepInfo] {
    guard task.hasRecipePrepTask else {
      print("⚠️ Task \(task.id) has no recipePrepTask")
      return []
    }

    let prepTask = task.recipePrepTask
    let recipeID = findRecipeID(for: task, prepTask: prepTask)

    if recipeID.isEmpty {
      print("⚠️ Task \(task.id) prep task has no recipe ID and couldn't find it from meal option")
      return []
    }

    guard let recipe = loadedRecipes[recipeID] else {
      print(
        "⚠️ Recipe \(recipeID) not loaded yet for task \(task.id). Loaded recipes: \(loadedRecipes.keys.joined(separator: ", "))"
      )
      return []
    }

    return buildStepData(prepTask: prepTask, recipe: recipe, recipeID: recipeID)
  }

  private func findRecipeID(
    for task: Mealplanning_MealPlanTask,
    prepTask: Mealplanning_RecipePrepTask
  ) -> String {
    var recipeID = prepTask.belongsToRecipe

    print("🔍 Task \(task.id): prepTask.id = '\(prepTask.id)'")
    print(
      "🔍 Task \(task.id): prepTask.belongsToRecipe = '\(recipeID)' (isEmpty: \(recipeID.isEmpty))")
    print("🔍 Task \(task.id): prepTask.name = '\(prepTask.name)'")
    print("🔍 Task \(task.id): prepTask.taskSteps.count = \(prepTask.taskSteps.count)")

    // If belongsToRecipe is empty, check loaded prep task
    if recipeID.isEmpty && !prepTask.id.isEmpty {
      if let loadedPrepTask = loadedPrepTasks[prepTask.id] {
        recipeID = loadedPrepTask.belongsToRecipe
        print(
          "✅ Found recipe ID \(recipeID) from loaded prep task \(prepTask.id) for task \(task.id)")
      } else {
        print("⚠️ Prep task \(prepTask.id) not found in loadedPrepTasks")
      }
    }

    // If belongsToRecipe is still empty, try to find it from the meal option
    if recipeID.isEmpty && task.hasMealPlanOption {
      recipeID = findRecipeIDFromMealOption(task: task, prepTask: prepTask)
    }

    return recipeID
  }

  private func findRecipeIDFromMealOption(
    task: Mealplanning_MealPlanTask,
    prepTask: Mealplanning_RecipePrepTask
  ) -> String {
    let meal = task.mealPlanOption.meal
    // Try to find the recipe by matching step IDs
    for component in meal.components {
      let candidateRecipeID = component.recipe.id
      // Check if any task step belongs to this recipe's steps
      for taskStep in prepTask.taskSteps where !taskStep.belongsToRecipeStep.isEmpty {
        // Check if this step ID exists in the component's recipe
        if component.recipe.steps.contains(where: { $0.id == taskStep.belongsToRecipeStep }) {
          print("✅ Found recipe ID \(candidateRecipeID) from meal option for task \(task.id)")
          return candidateRecipeID
        }
      }
    }
    return ""
  }

  private func buildStepData(
    prepTask: Mealplanning_RecipePrepTask,
    recipe: Mealplanning_Recipe,
    recipeID: String
  ) -> [StepInfo] {
    var stepData: [StepInfo] = []

    for taskStep in prepTask.taskSteps where !taskStep.belongsToRecipeStep.isEmpty {
      // Try to find step by ID in the loaded recipe
      if let stepIndex = recipe.steps.firstIndex(where: { $0.id == taskStep.belongsToRecipeStep }) {
        let step = recipe.steps[stepIndex]
        let formatted = formatStepTitle(step: step)
        if !formatted.isEmpty {
          stepData.append(StepInfo(stepID: taskStep.id, description: formatted))
        } else {
          print("⚠️ Step \(taskStep.belongsToRecipeStep) formatted to empty string")
        }
      } else {
        print("⚠️ Step \(taskStep.belongsToRecipeStep) not found in recipe \(recipeID)")
      }
    }

    if stepData.isEmpty {
      print(
        "⚠️ No step descriptions found, prep task has \(prepTask.taskSteps.count) task steps"
      )
    }

    return stepData
  }

  private func formatStepTitle(step: Mealplanning_RecipeStep) -> String {
    if step.hasPreparation && !step.preparation.name.isEmpty {
      return step.preparation.name
    }
    return "Step \(Int(step.index) + 1)"
  }
}

// Make MealPlanTask Identifiable
extension Mealplanning_MealPlanTask: Identifiable {
  // Already has id property, so this extension just makes it conform to Identifiable
}

// MARK: - Task Countdown Timer

struct TaskCountdownTimer: View {
  let dueDate: Date
  @State private var timeRemaining: TimeInterval = 0
  @State private var timer: Timer?

  var body: some View {
    HStack(spacing: 6) {
      Image(systemName: "clock.fill")
        .font(.caption2)
      Text(formattedDisplayText)
        .font(.caption)
        .fontWeight(.semibold)
        .monospacedDigit()
    }
    .foregroundColor(timeRemainingColor)
    .onAppear {
      updateTimeRemaining()
      startTimer()
    }
    .onDisappear {
      stopTimer()
    }
  }

  private var formattedDate: String {
    let formatter = DateFormatter()
    formatter.dateStyle = .medium
    formatter.timeStyle = .short
    return formatter.string(from: dueDate)
  }

  private var formattedTimeRemaining: String {
    if timeRemaining <= 0 {
      return "Overdue"
    }

    let days = Int(timeRemaining) / 86400
    let hours = (Int(timeRemaining) % 86400) / 3600
    let minutes = (Int(timeRemaining) % 3600) / 60
    let seconds = Int(timeRemaining) % 60

    if days > 0 {
      return String(format: "%dd %02d:%02d:%02d", days, hours, minutes, seconds)
    } else if hours > 0 {
      return String(format: "%02d:%02d:%02d", hours, minutes, seconds)
    } else {
      return String(format: "%02d:%02d", minutes, seconds)
    }
  }

  private var formattedDisplayText: String {
    let dateString = formattedDate
    let countdownString = formattedTimeRemaining
    return "by \(dateString) (\(countdownString))"
  }

  private var timeRemainingColor: Color {
    if timeRemaining <= 0 {
      return .red
    } else if timeRemaining < 3600 {  // Less than 1 hour
      return .red
    } else if timeRemaining < 86400 {  // Less than 1 day
      return .orange
    } else {
      return .secondary
    }
  }

  private func updateTimeRemaining() {
    let now = Date()
    timeRemaining = max(0, dueDate.timeIntervalSince(now))
  }

  private func startTimer() {
    timer = Timer.scheduledTimer(withTimeInterval: 1.0, repeats: true) { _ in
      updateTimeRemaining()
    }
    if let timer = timer {
      RunLoop.main.add(timer, forMode: .common)
    }
  }

  private func stopTimer() {
    timer?.invalidate()
    timer = nil
  }
}

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "John Doe"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  // Create a sample meal plan and tasks
  var mealPlan = Mealplanning_MealPlan()
  mealPlan.id = "mealplan123"
  mealPlan.notes = "Sample Meal Plan"

  var task1 = Mealplanning_MealPlanTask()
  task1.id = "task1"
  task1.creationExplanation = "Preheat oven to 350°F"
  task1.status = .unfinished

  var task2 = Mealplanning_MealPlanTask()
  task2.id = "task2"
  task2.creationExplanation = "Chop vegetables"
  task2.status = .finished

  return NavigationView {
    TaskListView(
      mealPlan: mealPlan,
      tasks: [task1, task2],
      authManager: authManager
    )
  }
  .environment(authManager)
}
