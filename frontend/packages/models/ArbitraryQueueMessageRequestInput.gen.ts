// GENERATED CODE, DO NOT EDIT MANUALLY

import { MessageQueueName } from './enums.gen';

export interface IArbitraryQueueMessageRequestInput {
  body: string;
  queueName: MessageQueueName;
}

export class ArbitraryQueueMessageRequestInput implements IArbitraryQueueMessageRequestInput {
  body: string;
  queueName: MessageQueueName;
  constructor(input: Partial<ArbitraryQueueMessageRequestInput> = {}) {
    this.body = input.body || '';
    this.queueName = input.queueName || 'data_changes';
  }
}
