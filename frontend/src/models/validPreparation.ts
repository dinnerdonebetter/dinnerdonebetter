export class ValidPreparation {
     id: number;
     name: string;
     description: string;
     iconPath: string;
     externalID: string;
     createdOn: number;
     lastUpdatedOn: number;
     archivedOn?: number;

    constructor(
      id: number = 0,
      name: string = '',
      description: string = '',
      iconPath: string = '',
      externalID: string = '',
      createdOn: number = 0,
      lastUpdatedOn: number = 0,
      archivedOn?: number,
    ) {
      this.id = id;
      this.name = name;
      this.description = description;
      this.iconPath = iconPath;
      this.externalID = externalID;
      this.createdOn = createdOn;
      this.lastUpdatedOn = lastUpdatedOn;
      this.archivedOn = archivedOn;
    }
}

export class ValidPreparationList {
     validPreparations: ValidPreparation[];
     totalCount: number;
     page: number;
     limit: number;
     filteredCount: number;

     constructor(
          validPreparations: ValidPreparation[] = [],
          totalCount: number = 0,
          page: number = 0,
          limit: number = 0,
          filteredCount: number = 0,
     ) {
          this.validPreparations = validPreparations;
          this.totalCount = totalCount;
          this.page = page;
          this.limit = limit;
          this.filteredCount = filteredCount;
     }
}