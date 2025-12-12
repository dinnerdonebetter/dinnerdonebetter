//
//  HomeView.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct HomeView: View {
    @Environment(AuthenticationManager.self) private var authManager
    @State private var viewModel: HomeViewModel?
    
    var body: some View {
        NavigationStack {
            Group {
                if let viewModel = viewModel {
                    if viewModel.isLoading {
                        ProgressView("Loading...")
                            .frame(maxWidth: .infinity, maxHeight: .infinity)
                    } else if let errorMessage = viewModel.errorMessage {
                    VStack(spacing: 16) {
                        Image(systemName: "exclamationmark.triangle")
                            .font(.largeTitle)
                            .foregroundColor(.orange)
                        Text("Error")
                            .font(.headline)
                        Text(errorMessage)
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                            .multilineTextAlignment(.center)
                            .padding(.horizontal)
                        Button("Retry") {
                            Task {
                                await viewModel.loadData()
                            }
                        }
                        .buttonStyle(.borderedProminent)
                    }
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
                    } else {
                        ScrollView {
                            VStack(spacing: 24) {
                                // Header Section
                                headerSection(viewModel: viewModel)
                                
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
                                if viewModel.pendingVoteMealPlans.isEmpty &&
                                   viewModel.upcomingMealPlans.isEmpty &&
                                   viewModel.userTasks.isEmpty &&
                                   viewModel.activeGroceryLists.isEmpty {
                                    emptyStateView
                                }
                                
                                // Sign Out Button
                                signOutButton
                            }
                            .padding()
                        }
                    }
                } else {
                    ProgressView("Initializing...")
                        .frame(maxWidth: .infinity, maxHeight: .infinity)
                }
            }
            .navigationTitle("Home")
            .refreshable {
                if let viewModel = viewModel {
                    await viewModel.loadData()
                }
            }
            .onAppear {
                // Initialize viewModel with the actual authManager from environment
                if viewModel == nil {
                    viewModel = HomeViewModel(authManager: authManager)
                }
                if let viewModel = viewModel {
                    Task {
                        await viewModel.loadData()
                    }
                }
            }
        }
    }
    
    // MARK: - Header Section
    private func headerSection(viewModel: HomeViewModel) -> some View {
        VStack(spacing: 16) {
            Text("Welcome, \(authManager.username)!")
                .font(.largeTitle)
                .fontWeight(.bold)
            
            Button(action: {
                // TODO: Navigate to create meal plan view
                print("Create Meal Plan tapped")
            }) {
                HStack {
                    Image(systemName: "plus.circle.fill")
                    Text("Create New Meal Plan")
                }
                .fontWeight(.semibold)
                .frame(maxWidth: .infinity)
                .padding()
                .background(Color.blue)
                .foregroundColor(.white)
                .cornerRadius(10)
            }
        }
    }
    
    // MARK: - Pending Votes Section
    private func pendingVotesSection(viewModel: HomeViewModel) -> some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Pending Votes")
                .font(.title2)
                .fontWeight(.bold)
                .padding(.horizontal, 4)
            
            ForEach(viewModel.pendingVoteMealPlans, id: \.id) { mealPlan in
                PendingVoteCard(
                    mealPlan: mealPlan,
                    hasVoted: viewModel.hasUserVoted(on: mealPlan),
                    timeUntilDeadline: viewModel.timeUntilDeadline(mealPlan.votingDeadline)
                ) {
                    // TODO: Navigate to voting view
                    print("Vote on meal plan \(mealPlan.id)")
                }
            }
        }
    }
    
    // MARK: - Upcoming Meals Section
    private func upcomingMealsSection(viewModel: HomeViewModel) -> some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Upcoming Meals")
                .font(.title2)
                .fontWeight(.bold)
                .padding(.horizontal, 4)
            
            ForEach(viewModel.upcomingMealPlans, id: \.id) { mealPlan in
                UpcomingMealCard(mealPlan: mealPlan) {
                    // TODO: Navigate to meal plan detail view
                    print("View meal plan \(mealPlan.id)")
                }
            }
        }
    }
    
    // MARK: - My Tasks Section
    private func myTasksSection(viewModel: HomeViewModel) -> some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("My Tasks")
                .font(.title2)
                .fontWeight(.bold)
                .padding(.horizontal, 4)
            
            // Group tasks by date
            let groupedTasks = Dictionary(grouping: viewModel.userTasks) { task in
                formatTaskDate(task)
            }
            let sortedDates = groupedTasks.keys.sorted()
            
            ForEach(sortedDates, id: \.self) { date in
                VStack(alignment: .leading, spacing: 8) {
                    Text(date)
                        .font(.headline)
                        .foregroundColor(.secondary)
                        .padding(.horizontal, 4)
                    
                    ForEach(groupedTasks[date] ?? [], id: \.id) { task in
                        TaskCard(task: task)
                    }
                }
            }
        }
    }
    
    // MARK: - Grocery Lists Section
    private func groceryListsSection(viewModel: HomeViewModel) -> some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Grocery Lists")
                .font(.title2)
                .fontWeight(.bold)
                .padding(.horizontal, 4)
            
            ForEach(viewModel.activeGroceryLists, id: \.mealPlanID) { groceryList in
                if let mealPlan = viewModel.allMealPlans.first(where: { $0.id == groceryList.mealPlanID }) {
                    GroceryListCard(
                        mealPlan: mealPlan,
                        items: groceryList.items
                    ) {
                        // TODO: Navigate to grocery list detail view
                        print("View grocery list for meal plan \(mealPlan.id)")
                    }
                }
            }
        }
    }
    
    // MARK: - Empty State
    private var emptyStateView: some View {
        VStack(spacing: 16) {
            Image(systemName: "calendar.badge.plus")
                .font(.system(size: 60))
                .foregroundColor(.secondary)
            Text("No Active Meal Plans")
                .font(.title2)
                .fontWeight(.semibold)
            Text("Create a meal plan to get started!")
                .font(.subheadline)
                .foregroundColor(.secondary)
        }
        .padding(.vertical, 40)
    }
    
    // MARK: - Sign Out Button
    private var signOutButton: some View {
        Button(action: authManager.logout) {
            Text("Sign Out")
                .fontWeight(.semibold)
                .frame(maxWidth: .infinity)
                .padding()
                .background(Color.red)
                .foregroundColor(.white)
                .cornerRadius(10)
        }
        .padding(.top, 20)
    }
    
    // MARK: - Helper Functions
    private func formatTaskDate(_ task: Mealplanning_MealPlanTask) -> String {
        // Get the date from the task's meal plan option's event
        // For now, use a simple format
        let formatter = DateFormatter()
        formatter.dateStyle = .medium
        formatter.timeStyle = .none
        
        // Try to get date from completedAt or createdAt
        let timestamp = task.completedAt.seconds > 0 ? task.completedAt : task.createdAt
        let date = timestampToDate(timestamp)
        return formatter.string(from: date)
    }
    
    // Reuse the timestampToDate function from HomeViewModel
    private func timestampToDate(_ timestamp: SwiftProtobuf.Google_Protobuf_Timestamp) -> Date {
        return HomeViewModel.timestampToDate(timestamp)
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
            VStack(alignment: .leading, spacing: 8) {
                HStack {
                    Text(mealPlan.notes.isEmpty ? "Meal Plan" : mealPlan.notes)
                        .font(.headline)
                        .foregroundColor(.primary)
                    Spacer()
                    if hasVoted {
                        Image(systemName: "checkmark.circle.fill")
                            .foregroundColor(.green)
                    } else {
                        Image(systemName: "exclamationmark.circle.fill")
                            .foregroundColor(.orange)
                    }
                }
                
                Text(timeUntilDeadline)
                    .font(.subheadline)
                    .foregroundColor(.secondary)
                
                if !hasVoted {
                    Text("Tap to vote")
                        .font(.caption)
                        .foregroundColor(.blue)
                }
            }
            .padding()
            .background(Color(.systemGray6))
            .cornerRadius(10)
        }
        .buttonStyle(.plain)
    }
}

// MARK: - Upcoming Meal Card
struct UpcomingMealCard: View {
    let mealPlan: Mealplanning_MealPlan
    let onTap: () -> Void
    
    var body: some View {
        Button(action: onTap) {
            VStack(alignment: .leading, spacing: 8) {
                Text(mealPlan.notes.isEmpty ? "Meal Plan" : mealPlan.notes)
                    .font(.headline)
                    .foregroundColor(.primary)
                
                // Show upcoming events
                ForEach(mealPlan.events.prefix(3), id: \.id) { event in
                    HStack {
                        Text(event.mealName.capitalized)
                            .font(.subheadline)
                        Spacer()
                        Text(formatEventDate(event))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                }
                
                if mealPlan.events.count > 3 {
                    Text("+ \(mealPlan.events.count - 3) more events")
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
            }
            .padding()
            .background(Color(.systemGray6))
            .cornerRadius(10)
        }
        .buttonStyle(.plain)
    }
    
    private func formatEventDate(_ event: Mealplanning_MealPlanEvent) -> String {
        let formatter = DateFormatter()
        formatter.dateStyle = .short
        formatter.timeStyle = .short
        
        let date = HomeViewModel.timestampToDate(event.startsAt)
        return formatter.string(from: date)
    }
}

// MARK: - Task Card
struct TaskCard: View {
    let task: Mealplanning_MealPlanTask
    
    var body: some View {
        HStack {
            VStack(alignment: .leading, spacing: 4) {
                Text(task.creationExplanation)
                    .font(.subheadline)
                    .fontWeight(.medium)
                
                if !task.statusExplanation.isEmpty {
                    Text(task.statusExplanation)
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
            }
            
            Spacer()
            
            // Status indicator
            Circle()
                .fill(task.status.lowercased() == "completed" ? Color.green : Color.orange)
                .frame(width: 12, height: 12)
        }
        .padding()
        .background(Color(.systemGray6))
        .cornerRadius(8)
    }
}

// MARK: - Grocery List Card
struct GroceryListCard: View {
    let mealPlan: Mealplanning_MealPlan
    let items: [Mealplanning_MealPlanGroceryListItem]
    let onTap: () -> Void
    
    // Computed property for items to show
    private var itemsToShow: [Mealplanning_MealPlanGroceryListItem] {
        Array(items.prefix(3))
    }
    
    var body: some View {
        Button(action: onTap) {
            VStack(alignment: .leading, spacing: 8) {
                Text(mealPlan.notes.isEmpty ? "Grocery List" : mealPlan.notes)
                    .font(.headline)
                    .foregroundColor(.primary)
                
                Text("\(items.count) item\(items.count == 1 ? "" : "s") needed")
                    .font(.subheadline)
                    .foregroundColor(.secondary)
                
                // Show first few items
                ForEach(itemsToShow, id: \.id) { item in
                    GroceryItemRow(item: item)
                }
                
                if items.count > 3 {
                    Text("+ \(items.count - 3) more items")
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
            }
            .padding()
            .background(Color(.systemGray6))
            .cornerRadius(10)
        }
        .buttonStyle(.plain)
    }
    
    private func formatQuantity(_ quantity: Common_Float32RangeWithOptionalMax) -> String {
        if quantity.hasMax {
            return "\(quantity.min) - \(quantity.max)"
        } else {
            return "\(quantity.min)+"
        }
    }
}

// MARK: - Grocery Item Row
struct GroceryItemRow: View {
    let item: Mealplanning_MealPlanGroceryListItem
    
    var body: some View {
        HStack {
            Text(item.ingredient.name)
                .font(.caption)
            Spacer()
            if item.hasQuantityNeeded && item.quantityNeeded.hasMax {
                Text(formatQuantity(item.quantityNeeded))
                    .font(.caption)
                    .foregroundColor(.secondary)
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
