//
//  VoteMealPlanView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct VoteMealPlanView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Environment(\.dismiss) private var dismiss
  @State private var viewModel: VoteMealPlanViewModel?
  @State private var currentEventIndex = 0
  
  let mealPlan: Mealplanning_MealPlan
  
  init(mealPlan: Mealplanning_MealPlan) {
    self.mealPlan = mealPlan
    // Note: viewModel will be initialized in onAppear to access @Environment authManager
  }
  
  var body: some View {
    NavigationStack {
      Group {
        if let viewModel = viewModel {
          VStack(spacing: 0) {
            // Event indicator
            eventIndicator(viewModel: viewModel)
            
            // Swipeable events
            TabView(selection: $currentEventIndex) {
              ForEach(Array(viewModel.mealPlan.events.enumerated()), id: \.element.id) { index, event in
                eventVotingView(event: event, index: index, viewModel: viewModel)
                  .tag(index)
              }
            }
            .tabViewStyle(.page)
            .indexViewStyle(.page(backgroundDisplayMode: .always))
            
            // Submit button (only shown when all ballots are complete and locked)
            if viewModel.canSubmit {
              submitButton(viewModel: viewModel)
                .padding()
            }
          }
          .navigationTitle(viewModel.mealPlan.notes.isEmpty ? "Vote on Meal Plan" : viewModel.mealPlan.notes)
          .navigationBarTitleDisplayMode(.inline)
          .toolbar {
            ToolbarItem(placement: .navigationBarLeading) {
              Button("Cancel") {
                dismiss()
              }
            }
          }
          .alert("Error", isPresented: .constant(viewModel.submissionError != nil)) {
            Button("OK") {
              viewModel.submissionError = nil
            }
          } message: {
            if let error = viewModel.submissionError {
              Text(error)
            }
          }
          .alert("Success", isPresented: Binding(
            get: { viewModel.submissionSuccess },
            set: { viewModel.submissionSuccess = $0 }
          )) {
            Button("OK") {
              dismiss()
            }
          } message: {
            Text("Your votes have been submitted successfully!")
          }
        } else {
          ProgressView("Initializing...")
            .frame(maxWidth: .infinity, maxHeight: .infinity)
        }
      }
      .onAppear {
        if viewModel == nil {
          viewModel = VoteMealPlanViewModel(mealPlan: mealPlan, authManager: authManager)
        }
      }
    }
  }
  
  // MARK: - Event Indicator
  
  private func eventIndicator(viewModel: VoteMealPlanViewModel) -> some View {
    HStack {
      Text("Event \(currentEventIndex + 1) of \(viewModel.mealPlan.events.count)")
        .font(.subheadline)
        .foregroundColor(.secondary)
      Spacer()
    }
    .padding(.horizontal)
    .padding(.vertical, 8)
    .background(Color(.systemGray6))
  }
  
  // MARK: - Event Voting View
  
  private func eventVotingView(event: Mealplanning_MealPlanEvent, index: Int, viewModel: VoteMealPlanViewModel) -> some View {
    ScrollView {
      VStack(alignment: .leading, spacing: 16) {
        // Event header
        eventHeader(event: event)
        
        // Lock status
        lockStatusView(event: event, viewModel: viewModel)
        
        // Instructions
        instructionsView(event: event)
        
        // Ranked options list (drag and drop)
        rankedOptionsList(event: event, viewModel: viewModel)
        
        // Lock button
        if let ballot = viewModel.getBallot(for: event.id) {
          LockButtonView(event: event, ballot: ballot) {
            viewModel.toggleLock(eventID: event.id)
          }
          .padding(.top, 8)
        }
      }
      .padding()
    }
  }
  
  private func eventHeader(event: Mealplanning_MealPlanEvent) -> some View {
    let date = HomeViewModel.timestampToDate(event.startsAt)
    let formatter = DateFormatter()
    formatter.dateStyle = .medium
    formatter.timeStyle = .short
    
    return VStack(alignment: .leading, spacing: 8) {
      Text(MealPlanningUtils.formatMealName(event.mealName))
        .font(.title2)
        .fontWeight(.bold)
      
      Text(formatter.string(from: date))
        .font(.subheadline)
        .foregroundColor(.secondary)
    }
  }
  
  private func lockStatusView(event: Mealplanning_MealPlanEvent, viewModel: VoteMealPlanViewModel) -> some View {
    if let ballot = viewModel.getBallot(for: event.id) {
      return AnyView(
        HStack {
          if ballot.isLocked {
            Label("Locked", systemImage: "lock.fill")
              .font(.subheadline)
              .foregroundColor(.green)
          } else {
            Label("Unlocked", systemImage: "lock.open.fill")
              .font(.subheadline)
              .foregroundColor(.orange)
          }
          
          Spacer()
          
          if !ballot.isComplete(totalOptions: event.options.count) {
            Text("\(ballot.rankedOptions.count) of \(event.options.count) options ranked")
              .font(.caption)
              .foregroundColor(.secondary)
          }
        }
        .padding()
        .background(Color(.systemGray6))
        .cornerRadius(8)
      )
    } else {
      return AnyView(EmptyView())
    }
  }
  
  private func instructionsView(event: Mealplanning_MealPlanEvent) -> some View {
    VStack(alignment: .leading, spacing: 4) {
      Text("Rank your preferences")
        .font(.headline)
      Text("Drag options to reorder them. Your first choice should be at the top.")
        .font(.subheadline)
        .foregroundColor(.secondary)
    }
    .padding()
    .background(Color(.systemGray6))
    .cornerRadius(8)
  }
  
  @State private var editMode: EditMode = .active
  
  private func rankedOptionsList(event: Mealplanning_MealPlanEvent, viewModel: VoteMealPlanViewModel) -> some View {
    VStack(alignment: .leading, spacing: 12) {
      Text("Your Ranking")
        .font(.headline)
        .padding(.horizontal, 4)
      
      if let ballot = viewModel.getBallot(for: event.id) {
        if ballot.rankedOptions.isEmpty {
          Text("No options available")
            .font(.subheadline)
            .foregroundColor(.secondary)
            .padding()
        } else {
          // Use List for native drag and drop support
          List {
            ForEach(Array(ballot.rankedOptions.enumerated()), id: \.element.id) { index, option in
              rankedOptionCard(option: option, rank: index + 1, event: event)
                .listRowInsets(EdgeInsets(top: 4, leading: 0, bottom: 4, trailing: 0))
                .listRowBackground(Color(.systemBackground))
            }
            .onMove { source, destination in
              if !ballot.isLocked {
                viewModel.reorderOptions(eventID: event.id, from: source, to: destination)
              }
            }
          }
          .listStyle(.plain)
          .frame(minHeight: CGFloat(ballot.rankedOptions.count * 80))
          .scrollContentBackground(.hidden)
          .environment(\.editMode, ballot.isLocked ? .constant(.inactive) : $editMode)
        }
      }
    }
  }
  
  private func rankedOptionCard(option: Mealplanning_MealPlanOption, rank: Int, event: Mealplanning_MealPlanEvent) -> some View {
    HStack(spacing: 12) {
      // Rank indicator
      ZStack {
        Circle()
          .fill(rankColor(rank))
          .frame(width: 32, height: 32)
        Text("\(rank)")
          .font(.headline)
          .foregroundColor(.white)
      }
      
      // Option details
      VStack(alignment: .leading, spacing: 4) {
        if option.hasMeal {
          Text(option.meal.name)
            .font(.headline)
          if !option.meal.description_p.isEmpty {
            Text(option.meal.description_p)
              .font(.caption)
              .foregroundColor(.secondary)
              .lineLimit(2)
          }
        } else {
          Text("Option \(rank)")
            .font(.headline)
        }
      }
      
      Spacer()
      
      // Drag handle
      Image(systemName: "line.3.horizontal")
        .foregroundColor(.secondary)
    }
    .padding()
    .background(Color(.systemBackground))
    .cornerRadius(8)
    .shadow(color: Color.black.opacity(0.1), radius: 2, x: 0, y: 1)
  }
  
  private func rankColor(_ rank: Int) -> Color {
    switch rank {
    case 1:
      return .green
    case 2:
      return .blue
    case 3:
      return .orange
    default:
      return .gray
    }
  }
  
  // MARK: - Submit Button
  
  private func submitButton(viewModel: VoteMealPlanViewModel) -> some View {
    Button(action: {
      Task {
        let success = await viewModel.submitVotes()
        if success {
          // Alert will be shown via viewModel.submissionSuccess
        }
      }
    }, label: {
      HStack {
        if viewModel.isSubmitting {
          ProgressView()
            .progressViewStyle(CircularProgressViewStyle(tint: .white))
        }
        Text(viewModel.isSubmitting ? "Submitting..." : "Submit Votes")
          .fontWeight(.semibold)
      }
      .frame(maxWidth: .infinity)
      .padding()
      .background(viewModel.isSubmitting ? Color.gray : Color.green)
      .foregroundColor(.white)
      .cornerRadius(10)
    })
    .disabled(viewModel.isSubmitting)
  }
}

// MARK: - Lock Button View

struct LockButtonView: View {
  let event: Mealplanning_MealPlanEvent
  let ballot: EventBallot
  let onToggle: () -> Void
  
  var body: some View {
    Button(action: onToggle) {
      HStack {
        Image(systemName: ballot.isLocked ? "lock.fill" : "lock.open.fill")
        Text(ballot.isLocked ? "Unlock Ballot" : "Lock Ballot")
      }
      .font(.headline)
      .frame(maxWidth: .infinity)
      .padding()
      .background(ballot.isLocked ? Color.orange : Color.blue)
      .foregroundColor(.white)
      .cornerRadius(10)
    }
    .disabled(!ballot.isComplete(totalOptions: event.options.count))
  }
}

#Preview {
  // This would need a mock meal plan for preview
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "Test User"
  authManager.userID = "user123"
  authManager.accountID = "account123"
  
  // Create a mock meal plan
  var mealPlan = Mealplanning_MealPlan()
  mealPlan.id = "test-plan"
  mealPlan.notes = "Test Meal Plan"
  
  var event = Mealplanning_MealPlanEvent()
  event.id = "test-event"
  event.mealName = .dinner
  
  var option1 = Mealplanning_MealPlanOption()
  option1.id = "option1"
  var meal1 = Mealplanning_Meal()
  meal1.name = "Pasta"
  option1.meal = meal1
  
  var option2 = Mealplanning_MealPlanOption()
  option2.id = "option2"
  var meal2 = Mealplanning_Meal()
  meal2.name = "Pizza"
  option2.meal = meal2
  
  event.options = [option1, option2]
  mealPlan.events = [event]
  
  return VoteMealPlanView(mealPlan: mealPlan)
    .environment(authManager)
}

