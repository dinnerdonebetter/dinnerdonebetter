// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrument {
  createdAt: string;
  iconPath: string;
  pluralName: string;
  usableForStorage: boolean;
  archivedAt?: string;
  displayInSummaryLists: boolean;
  id: string;
  includeInGeneratedInstructions: boolean;
  lastUpdatedAt?: string;
  name: string;
  slug: string;
  description: string;
}

export class ValidInstrument implements IValidInstrument {
  createdAt: string;
  iconPath: string;
  pluralName: string;
  usableForStorage: boolean;
  archivedAt?: string;
  displayInSummaryLists: boolean;
  id: string;
  includeInGeneratedInstructions: boolean;
  lastUpdatedAt?: string;
  name: string;
  slug: string;
  description: string;
  constructor(input: Partial<ValidInstrument> = {}) {
    this.createdAt = input.createdAt = '';
    this.iconPath = input.iconPath = '';
    this.pluralName = input.pluralName = '';
    this.usableForStorage = input.usableForStorage = false;
    this.archivedAt = input.archivedAt;
    this.displayInSummaryLists = input.displayInSummaryLists = false;
    this.id = input.id = '';
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
    this.slug = input.slug = '';
    this.description = input.description = '';
  }
}
