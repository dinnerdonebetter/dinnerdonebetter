//
//  AnalyticsConfiguration.swift
//  ios
//
//  Configuration and factory for EventReporter.
//  Segment write key is read from Info.plist (Secrets.xcconfig -> $(SEGMENT_WRITE_KEY)).
//

import Foundation

enum AnalyticsConfiguration {
  private static let placeholderWriteKey = "your-segment-write-key-here"

  /// Segment write key from Info.plist. Injected via Secrets.xcconfig at build time.
  static var segmentWriteKey: String {
    Bundle.main.infoDictionary?["SegmentWriteKey"] as? String ?? ""
  }

  /// Shared EventReporter instance. Cached so we don't create multiple Segment clients.
  private static var _sharedReporter: (any EventReporter)?
  private static let reporterLock = NSLock()

  /// Returns the shared EventReporter. Uses Segment when write key is valid, else Noop.
  static func provideEventReporter() -> any EventReporter {
    reporterLock.lock()
    defer { reporterLock.unlock() }
    if let existing = _sharedReporter {
      return existing
    }
    let key = segmentWriteKey.trimmingCharacters(in: .whitespaces)
    let reporter: any EventReporter
    if !key.isEmpty, key != placeholderWriteKey {
      reporter = SegmentEventReporter(writeKey: key)
    } else {
      reporter = NoopEventReporter()
    }
    _sharedReporter = reporter
    return reporter
  }
}
