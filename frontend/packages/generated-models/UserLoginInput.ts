// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserLoginInput {
  username: string;
  password: string;
  totpToken: string;
}

export class UserLoginInput implements IUserLoginInput {
  username: string;
  password: string;
  totpToken: string;
  constructor(input: Partial<UserLoginInput> = {}) {
    this.username = input.username = '';
    this.password = input.password = '';
    this.totpToken = input.totpToken = '';
  }
}
