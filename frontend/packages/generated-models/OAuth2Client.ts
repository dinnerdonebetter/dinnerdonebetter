// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IOAuth2Client {
  description: string;
  id: string;
  name: string;
  archivedAt?: string;
  clientID: string;
  clientSecret: string;
  createdAt: string;
}

export class OAuth2Client implements IOAuth2Client {
  description: string;
  id: string;
  name: string;
  archivedAt?: string;
  clientID: string;
  clientSecret: string;
  createdAt: string;
  constructor(input: Partial<OAuth2Client> = {}) {
    this.description = input.description = '';
    this.id = input.id = '';
    this.name = input.name = '';
    this.archivedAt = input.archivedAt;
    this.clientID = input.clientID = '';
    this.clientSecret = input.clientSecret = '';
    this.createdAt = input.createdAt = '';
  }
}
