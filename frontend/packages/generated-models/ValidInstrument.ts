// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrument {
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  lastUpdatedAt?: string;
  pluralName: string;
  archivedAt?: string;
  createdAt: string;
  description: string;
  slug: string;
  usableForStorage: boolean;
  id: string;
  name: string;
}

export class ValidInstrument implements IValidInstrument {
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  lastUpdatedAt?: string;
  pluralName: string;
  archivedAt?: string;
  createdAt: string;
  description: string;
  slug: string;
  usableForStorage: boolean;
  id: string;
  name: string;
  constructor(input: Partial<ValidInstrument> = {}) {
    this.displayInSummaryLists = input.displayInSummaryLists = false;
    this.iconPath = input.iconPath = '';
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.pluralName = input.pluralName = '';
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.description = input.description = '';
    this.slug = input.slug = '';
    this.usableForStorage = input.usableForStorage = false;
    this.id = input.id = '';
    this.name = input.name = '';
  }
}
