// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IOAuth2ClientCreationResponse {
  clientSecret: string;
  description: string;
  id: string;
  name: string;
  clientID: string;
}

export class OAuth2ClientCreationResponse implements IOAuth2ClientCreationResponse {
  clientSecret: string;
  description: string;
  id: string;
  name: string;
  clientID: string;
  constructor(input: Partial<OAuth2ClientCreationResponse> = {}) {
    this.clientSecret = input.clientSecret = '';
    this.description = input.description = '';
    this.id = input.id = '';
    this.name = input.name = '';
    this.clientID = input.clientID = '';
  }
}
