// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IOAuth2ClientCreationResponse {
  description: string;
  id: string;
  name: string;
  clientID: string;
  clientSecret: string;
}

export class OAuth2ClientCreationResponse implements IOAuth2ClientCreationResponse {
  description: string;
  id: string;
  name: string;
  clientID: string;
  clientSecret: string;
  constructor(input: Partial<OAuth2ClientCreationResponse> = {}) {
    this.description = input.description = '';
    this.id = input.id = '';
    this.name = input.name = '';
    this.clientID = input.clientID = '';
    this.clientSecret = input.clientSecret = '';
  }
}
