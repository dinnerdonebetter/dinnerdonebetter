// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserRegistrationInput {
   username: string;
 acceptedPrivacyPolicy: boolean;
 acceptedTOS: boolean;
 firstName: string;
 householdName: string;
 invitationID: string;
 password: string;
 birthday?: string;
 emailAddress: string;
 invitationToken: string;
 lastName: string;

}

export class UserRegistrationInput implements IUserRegistrationInput {
   username: string;
 acceptedPrivacyPolicy: boolean;
 acceptedTOS: boolean;
 firstName: string;
 householdName: string;
 invitationID: string;
 password: string;
 birthday?: string;
 emailAddress: string;
 invitationToken: string;
 lastName: string;
constructor(input: Partial<UserRegistrationInput> = {}) {
	 this.username = input.username = '';
 this.acceptedPrivacyPolicy = input.acceptedPrivacyPolicy = false;
 this.acceptedTOS = input.acceptedTOS = false;
 this.firstName = input.firstName = '';
 this.householdName = input.householdName = '';
 this.invitationID = input.invitationID = '';
 this.password = input.password = '';
 this.birthday = input.birthday;
 this.emailAddress = input.emailAddress = '';
 this.invitationToken = input.invitationToken = '';
 this.lastName = input.lastName = '';
}
}