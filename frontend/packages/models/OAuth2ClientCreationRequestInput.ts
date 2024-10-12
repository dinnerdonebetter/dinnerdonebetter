// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IOAuth2ClientCreationRequestInput {
   description: string;
 name: string;

}

export class OAuth2ClientCreationRequestInput implements IOAuth2ClientCreationRequestInput {
   description: string;
 name: string;
constructor(input: Partial<OAuth2ClientCreationRequestInput> = {}) {
	 this.description = input.description || '';
 this.name = input.name || '';
}
}