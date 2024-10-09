// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserLoginInput {
  totpToken: string;
  username: string;
  password: string;
}

export class UserLoginInput implements IUserLoginInput {
  totpToken: string;
  username: string;
  password: string;
  constructor(input: Partial<UserLoginInput> = {}) {
    this.totpToken = input.totpToken = '';
    this.username = input.username = '';
    this.password = input.password = '';
  }
}
