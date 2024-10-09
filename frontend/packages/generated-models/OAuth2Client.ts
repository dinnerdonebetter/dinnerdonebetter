// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IOAuth2Client {
  clientID: string;
  clientSecret: string;
  createdAt: string;
  description: string;
  id: string;
  name: string;
  archivedAt?: string;
}

export class OAuth2Client implements IOAuth2Client {
  clientID: string;
  clientSecret: string;
  createdAt: string;
  description: string;
  id: string;
  name: string;
  archivedAt?: string;
  constructor(input: Partial<OAuth2Client> = {}) {
    this.clientID = input.clientID = '';
    this.clientSecret = input.clientSecret = '';
    this.createdAt = input.createdAt = '';
    this.description = input.description = '';
    this.id = input.id = '';
    this.name = input.name = '';
    this.archivedAt = input.archivedAt;
  }
}
