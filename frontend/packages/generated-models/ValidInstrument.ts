// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidInstrument {
   pluralName: string;
 slug: string;
 description: string;
 iconPath: string;
 lastUpdatedAt?: string;
 id: string;
 includeInGeneratedInstructions: boolean;
 name: string;
 usableForStorage: boolean;
 archivedAt?: string;
 createdAt: string;
 displayInSummaryLists: boolean;

}

export class ValidInstrument implements IValidInstrument {
   pluralName: string;
 slug: string;
 description: string;
 iconPath: string;
 lastUpdatedAt?: string;
 id: string;
 includeInGeneratedInstructions: boolean;
 name: string;
 usableForStorage: boolean;
 archivedAt?: string;
 createdAt: string;
 displayInSummaryLists: boolean;
constructor(input: Partial<ValidInstrument> = {}) {
	 this.pluralName = input.pluralName = '';
 this.slug = input.slug = '';
 this.description = input.description = '';
 this.iconPath = input.iconPath = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.id = input.id = '';
 this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
 this.name = input.name = '';
 this.usableForStorage = input.usableForStorage = false;
 this.archivedAt = input.archivedAt;
 this.createdAt = input.createdAt = '';
 this.displayInSummaryLists = input.displayInSummaryLists = false;
}
}