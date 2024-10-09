// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidInstrumentUpdateRequestInput {
  displayInSummaryLists?: boolean;
  iconPath?: string;
  includeInGeneratedInstructions?: boolean;
  name?: string;
  pluralName?: string;
  slug?: string;
  usableForStorage?: boolean;
  description?: string;
}

export class ValidInstrumentUpdateRequestInput implements IValidInstrumentUpdateRequestInput {
  displayInSummaryLists?: boolean;
  iconPath?: string;
  includeInGeneratedInstructions?: boolean;
  name?: string;
  pluralName?: string;
  slug?: string;
  usableForStorage?: boolean;
  description?: string;
  constructor(input: Partial<ValidInstrumentUpdateRequestInput> = {}) {
    this.displayInSummaryLists = input.displayInSummaryLists;
    this.iconPath = input.iconPath;
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions;
    this.name = input.name;
    this.pluralName = input.pluralName;
    this.slug = input.slug;
    this.usableForStorage = input.usableForStorage;
    this.description = input.description;
  }
}
