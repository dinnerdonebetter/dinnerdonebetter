// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IOAuth2ClientCreationResponse {
   clientID: string;
 clientSecret: string;
 description: string;
 id: string;
 name: string;

}

export class OAuth2ClientCreationResponse implements IOAuth2ClientCreationResponse {
   clientID: string;
 clientSecret: string;
 description: string;
 id: string;
 name: string;
constructor(input: Partial<OAuth2ClientCreationResponse> = {}) {
	 this.clientID = input.clientID || '';
 this.clientSecret = input.clientSecret || '';
 this.description = input.description || '';
 this.id = input.id || '';
 this.name = input.name || '';
}
}