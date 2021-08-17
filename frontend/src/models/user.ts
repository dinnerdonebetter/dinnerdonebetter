export class User {
     id: number;
     username: string;
     avatar: string;
     externalID: string;
     reputation: string;
     reputationExplanation: string;
     passwordLastChangedOn?: number;
     serviceRoles: string[];
     requiresPasswordChange: boolean;
     createdOn: number;
     lastUpdatedOn?: number;
     archivedOn?: number;

    constructor(
      id: number = 0,
      username: string = "",
      avatar: string = "",
      externalID: string = "",
      reputation: string = "",
      reputationExplanation: string = "",
      passwordLastChangedOn?: number,
      serviceRoles: string[] = [],
      requiresPasswordChange: boolean = false,
      createdOn: number = 0,
      lastUpdatedOn?: number,
      archivedOn?: number,
    ) {
         this.id = id;
         this.username = username;
         this.avatar = avatar;
         this.externalID = externalID;
         this.reputation = reputation;
         this.reputationExplanation = reputationExplanation;
         this.passwordLastChangedOn = passwordLastChangedOn;
         this.serviceRoles = serviceRoles;
         this.requiresPasswordChange = requiresPasswordChange;
         this.createdOn = createdOn;
         this.lastUpdatedOn = lastUpdatedOn;
         this.archivedOn = archivedOn;
    }
}

export class UserList {
     users: User[];
     totalCount: number;
     page: number;
     limit: number;
     filteredCount: number;

     constructor(
          users: User[] = [],
          totalCount: number = 0,
          page: number = 0,
          limit: number = 0,
          filteredCount: number = 0,
     ) {
          this.users = users;
          this.totalCount = totalCount;
          this.page = page;
          this.limit = limit;
          this.filteredCount = filteredCount;
     }
}