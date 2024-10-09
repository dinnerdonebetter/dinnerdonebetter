// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IPasswordResetToken {
   createdAt: string;
 expiresAt: string;
 id: string;
 lastUpdatedAt?: string;
 token: string;
 archivedAt?: string;
 belongsToUser: string;

}

export class PasswordResetToken implements IPasswordResetToken {
   createdAt: string;
 expiresAt: string;
 id: string;
 lastUpdatedAt?: string;
 token: string;
 archivedAt?: string;
 belongsToUser: string;
constructor(input: Partial<PasswordResetToken> = {}) {
	 this.createdAt = input.createdAt = '';
 this.expiresAt = input.expiresAt = '';
 this.id = input.id = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.token = input.token = '';
 this.archivedAt = input.archivedAt;
 this.belongsToUser = input.belongsToUser = '';
}
}