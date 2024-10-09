// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IPasswordResetResponse {
  successful: boolean;
}

export class PasswordResetResponse implements IPasswordResetResponse {
  successful: boolean;
  constructor(input: Partial<PasswordResetResponse> = {}) {
    this.successful = input.successful || false;
  }
}
