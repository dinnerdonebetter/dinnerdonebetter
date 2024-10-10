// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrument {
  archivedAt: string;
  createdAt: string;
  description: string;
  displayInSummaryLists: boolean;
  iconPath: string;
  id: string;
  includeInGeneratedInstructions: boolean;
  lastUpdatedAt: string;
  name: string;
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
}

export class ValidInstrument implements IValidInstrument {
  archivedAt: string;
  createdAt: string;
  description: string;
  displayInSummaryLists: boolean;
  iconPath: string;
  id: string;
  includeInGeneratedInstructions: boolean;
  lastUpdatedAt: string;
  name: string;
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
  constructor(input: Partial<ValidInstrument> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.description = input.description || '';
    this.displayInSummaryLists = input.displayInSummaryLists || false;
    this.iconPath = input.iconPath || '';
    this.id = input.id || '';
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions || false;
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.name = input.name || '';
    this.pluralName = input.pluralName || '';
    this.slug = input.slug || '';
    this.usableForStorage = input.usableForStorage || false;
  }
}
