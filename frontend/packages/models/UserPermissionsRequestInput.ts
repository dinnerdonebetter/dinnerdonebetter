// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserPermissionsRequestInput {
   permissions: string[];

}

export class UserPermissionsRequestInput implements IUserPermissionsRequestInput {
   permissions: string[];
constructor(input: Partial<UserPermissionsRequestInput> = {}) {
	 this.permissions = input.permissions || [];
}
}