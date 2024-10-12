// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserRegistrationInput {
   acceptedPrivacyPolicy: boolean;
 acceptedTOS: boolean;
 birthday: string;
 emailAddress: string;
 firstName: string;
 householdName: string;
 invitationID: string;
 invitationToken: string;
 lastName: string;
 password: string;
 username: string;

}

export class UserRegistrationInput implements IUserRegistrationInput {
   acceptedPrivacyPolicy: boolean;
 acceptedTOS: boolean;
 birthday: string;
 emailAddress: string;
 firstName: string;
 householdName: string;
 invitationID: string;
 invitationToken: string;
 lastName: string;
 password: string;
 username: string;
constructor(input: Partial<UserRegistrationInput> = {}) {
	 this.acceptedPrivacyPolicy = input.acceptedPrivacyPolicy || false;
 this.acceptedTOS = input.acceptedTOS || false;
 this.birthday = input.birthday || '';
 this.emailAddress = input.emailAddress || '';
 this.firstName = input.firstName || '';
 this.householdName = input.householdName || '';
 this.invitationID = input.invitationID || '';
 this.invitationToken = input.invitationToken || '';
 this.lastName = input.lastName || '';
 this.password = input.password || '';
 this.username = input.username || '';
}
}