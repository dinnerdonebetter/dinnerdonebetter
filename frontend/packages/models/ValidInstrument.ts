// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrument {
  iconPath: string;
  id: string;
  name: string;
  usableForStorage: boolean;
  archivedAt: string;
  createdAt: string;
  description: string;
  displayInSummaryLists: boolean;
  includeInGeneratedInstructions: boolean;
  lastUpdatedAt: string;
  pluralName: string;
  slug: string;
}

export class ValidInstrument implements IValidInstrument {
  iconPath: string;
  id: string;
  name: string;
  usableForStorage: boolean;
  archivedAt: string;
  createdAt: string;
  description: string;
  displayInSummaryLists: boolean;
  includeInGeneratedInstructions: boolean;
  lastUpdatedAt: string;
  pluralName: string;
  slug: string;
  constructor(input: Partial<ValidInstrument> = {}) {
    this.iconPath = input.iconPath || '';
    this.id = input.id || '';
    this.name = input.name || '';
    this.usableForStorage = input.usableForStorage || false;
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.description = input.description || '';
    this.displayInSummaryLists = input.displayInSummaryLists || false;
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions || false;
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.pluralName = input.pluralName || '';
    this.slug = input.slug || '';
  }
}
