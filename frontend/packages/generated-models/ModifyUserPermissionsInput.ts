// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IModifyUserPermissionsInput {
  newRole: string;
  reason: string;
}

export class ModifyUserPermissionsInput implements IModifyUserPermissionsInput {
  newRole: string;
  reason: string;
  constructor(input: Partial<ModifyUserPermissionsInput> = {}) {
    this.newRole = input.newRole = '';
    this.reason = input.reason = '';
  }
}
