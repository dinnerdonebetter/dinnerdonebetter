//
//  Item.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import Foundation
import SwiftData

@Model
final class Item {
    var timestamp: Date
    
    init(timestamp: Date) {
        self.timestamp = timestamp
    }
}
