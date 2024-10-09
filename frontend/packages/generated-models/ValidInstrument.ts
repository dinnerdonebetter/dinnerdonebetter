// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrument {
  archivedAt?: string;
  createdAt: string;
  id: string;
  lastUpdatedAt?: string;
  usableForStorage: boolean;
  slug: string;
  description: string;
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  name: string;
  pluralName: string;
}

export class ValidInstrument implements IValidInstrument {
  archivedAt?: string;
  createdAt: string;
  id: string;
  lastUpdatedAt?: string;
  usableForStorage: boolean;
  slug: string;
  description: string;
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  name: string;
  pluralName: string;
  constructor(input: Partial<ValidInstrument> = {}) {
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.usableForStorage = input.usableForStorage = false;
    this.slug = input.slug = '';
    this.description = input.description = '';
    this.displayInSummaryLists = input.displayInSummaryLists = false;
    this.iconPath = input.iconPath = '';
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
    this.name = input.name = '';
    this.pluralName = input.pluralName = '';
  }
}
