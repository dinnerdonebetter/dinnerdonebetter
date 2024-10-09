// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSettingCreationRequestInput {
  name: string;
  type: string;
  adminsOnly: boolean;
  defaultValue?: string;
  description: string;
  enumeration: string;
}

export class ServiceSettingCreationRequestInput implements IServiceSettingCreationRequestInput {
  name: string;
  type: string;
  adminsOnly: boolean;
  defaultValue?: string;
  description: string;
  enumeration: string;
  constructor(input: Partial<ServiceSettingCreationRequestInput> = {}) {
    this.name = input.name = '';
    this.type = input.type = '';
    this.adminsOnly = input.adminsOnly = false;
    this.defaultValue = input.defaultValue;
    this.description = input.description = '';
    this.enumeration = input.enumeration = '';
  }
}
