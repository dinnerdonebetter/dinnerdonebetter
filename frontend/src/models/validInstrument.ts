export class ValidInstrument {
     id: number;
     name: string;
     variant: string;
     description: string;
     iconPath: string;
     externalID: string;
     createdOn: number;
     lastUpdatedOn: number;
     archivedOn?: number;

    constructor(
      id: number = 0,
      name: string = '',
      variant: string = '',
      description: string = '',
      iconPath: string = '',
      externalID: string = '',
      createdOn: number = 0,
      lastUpdatedOn: number = 0,
      archivedOn?: number,
    ) {
      this.id = id;
      this.name = name;
      this.variant = variant;
      this.description = description;
      this.iconPath = iconPath;
      this.externalID = externalID;
      this.createdOn = createdOn;
      this.lastUpdatedOn = lastUpdatedOn;
      this.archivedOn = archivedOn;
    }
}

export class ValidInstrumentList {
     validInstruments: ValidInstrument[];
     totalCount: number;
     page: number;
     limit: number;
     filteredCount: number;

     constructor(
          validInstruments: ValidInstrument[] = [],
          totalCount: number = 0,
          page: number = 0,
          limit: number = 0,
          filteredCount: number = 0,
     ) {
          this.validInstruments = validInstruments;
          this.totalCount = totalCount;
          this.page = page;
          this.limit = limit;
          this.filteredCount = filteredCount;
     }
}