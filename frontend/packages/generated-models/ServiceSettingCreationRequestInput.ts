// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSettingCreationRequestInput {
  adminsOnly: boolean;
  defaultValue?: string;
  description: string;
  enumeration: string;
  name: string;
  type: string;
}

export class ServiceSettingCreationRequestInput implements IServiceSettingCreationRequestInput {
  adminsOnly: boolean;
  defaultValue?: string;
  description: string;
  enumeration: string;
  name: string;
  type: string;
  constructor(input: Partial<ServiceSettingCreationRequestInput> = {}) {
    this.adminsOnly = input.adminsOnly = false;
    this.defaultValue = input.defaultValue;
    this.description = input.description = '';
    this.enumeration = input.enumeration = '';
    this.name = input.name = '';
    this.type = input.type = '';
  }
}
