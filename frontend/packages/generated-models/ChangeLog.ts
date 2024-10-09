// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IChangeLog {
  newValue: string;
  oldValue: string;
}

export class ChangeLog implements IChangeLog {
  newValue: string;
  oldValue: string;
  constructor(input: Partial<ChangeLog> = {}) {
    this.newValue = input.newValue = '';
    this.oldValue = input.oldValue = '';
  }
}
