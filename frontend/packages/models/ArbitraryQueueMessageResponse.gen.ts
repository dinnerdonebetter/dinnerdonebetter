// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IArbitraryQueueMessageResponse {
  success: boolean;
}

export class ArbitraryQueueMessageResponse implements IArbitraryQueueMessageResponse {
  success: boolean;
  constructor(input: Partial<ArbitraryQueueMessageResponse> = {}) {
    this.success = input.success || false;
  }
}
