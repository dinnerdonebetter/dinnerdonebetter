//
//  VoteMealPlanViewModel.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import Foundation
import GRPCCore
import SwiftProtobuf
import SwiftUI

/// Represents a ballot for a single event with ranked options
struct EventBallot: Identifiable {
  let id: String // Event ID
  var rankedOptions: [Mealplanning_MealPlanOption] // Ordered list (first = rank 0)
  var isLocked: Bool = false
  
  /// Check if ballot is complete (has all options)
  func isComplete(totalOptions: Int) -> Bool {
    return rankedOptions.count == totalOptions
  }
}

@Observable
@MainActor
class VoteMealPlanViewModel {
  let mealPlan: Mealplanning_MealPlan
  private let authManager: AuthenticationManager
  
  // Ballots for each event (keyed by event ID)
  var ballots: [String: EventBallot] = [:]
  
  // Submission state
  var isSubmitting = false
  var submissionError: String?
  var submissionSuccess = false
  
  init(mealPlan: Mealplanning_MealPlan, authManager: AuthenticationManager) {
    self.mealPlan = mealPlan
    self.authManager = authManager
    
    // Initialize ballots for each event with all options in original order
    for event in mealPlan.events {
      ballots[event.id] = EventBallot(
        id: event.id,
        rankedOptions: event.options
      )
    }
  }
  
  // MARK: - Ballot Management
  
  /// Get ballot for an event
  func getBallot(for eventID: String) -> EventBallot? {
    return ballots[eventID]
  }
  
  /// Reorder options in a ballot (drag and drop)
  func reorderOptions(eventID: String, from source: IndexSet, to destination: Int) {
    guard var ballot = ballots[eventID] else { return }
    
    // Don't allow reordering if locked
    guard !ballot.isLocked else { return }
    
    ballot.rankedOptions.move(fromOffsets: source, toOffset: destination)
    ballots[eventID] = ballot
  }
  
  /// Lock/unlock a ballot for an event
  func toggleLock(eventID: String) {
    guard var ballot = ballots[eventID] else { return }
    
    // Can only lock if ballot is complete
    let event = mealPlan.events.first { $0.id == eventID }
    guard let event = event else { return }
    
    if !ballot.isLocked && !ballot.isComplete(totalOptions: event.options.count) {
      // Cannot lock incomplete ballot
      return
    }
    
    ballot.isLocked.toggle()
    ballots[eventID] = ballot
  }
  
  /// Check if all ballots are complete and locked
  var canSubmit: Bool {
    for event in mealPlan.events {
      guard let ballot = ballots[event.id] else { return false }
      if !ballot.isLocked || !ballot.isComplete(totalOptions: event.options.count) {
        return false
      }
    }
    return true
  }
  
  // MARK: - Vote Submission
  
  /// Submit all votes for all events
  func submitVotes() async -> Bool {
    guard canSubmit else {
      submissionError = "All ballots must be complete and locked before submitting"
      return false
    }
    
    isSubmitting = true
    submissionError = nil
    submissionSuccess = false
    
    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "VoteMealPlanViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }
      
      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "VoteMealPlanViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }
      
      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)
      
      // Submit votes for each event
      for event in mealPlan.events {
        guard let ballot = ballots[event.id] else {
          continue
        }
        
        // Create vote request
        var request = Mealplanning_CreateMealPlanOptionVoteRequest()
        request.mealPlanID = mealPlan.id
        request.mealPlanEventID = event.id
        
        var input = Mealplanning_MealPlanOptionVoteCreationRequestInput()
        
        // Create votes with ranks based on order in rankedOptions
        for (index, option) in ballot.rankedOptions.enumerated() {
          var voteInput = Mealplanning_MealPlanOptionVoteCreationInput()
          voteInput.belongsToMealPlanOption = option.id
          voteInput.rank = UInt32(index) // Rank 0 = first choice, 1 = second choice, etc.
          voteInput.notes = ""
          voteInput.abstain = false
          
          input.votes.append(voteInput)
        }
        
        request.input = input
        
        // Submit votes for this event
        _ = try await clientManager.client.mealPlanning.createMealPlanOptionVote(
          request,
          metadata: metadata,
          options: clientManager.defaultCallOptions
        )
      }
      
      submissionSuccess = true
      isSubmitting = false
      return true
    } catch {
      submissionError = "Failed to submit votes: \(error.localizedDescription)"
      print("❌ Error submitting votes: \(error)")
      isSubmitting = false
      return false
    }
  }
}

