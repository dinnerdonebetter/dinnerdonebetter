// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserPermissionsResponse {
   permissions: object;

}

export class UserPermissionsResponse implements IUserPermissionsResponse {
   permissions: object;
constructor(input: Partial<UserPermissionsResponse> = {}) {
	 this.permissions = input.permissions || {};
}
}