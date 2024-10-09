// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSettingCreationRequestInput {
  description: string;
  enumeration: string;
  name: string;
  type: string;
  adminsOnly: boolean;
  defaultValue?: string;
}

export class ServiceSettingCreationRequestInput implements IServiceSettingCreationRequestInput {
  description: string;
  enumeration: string;
  name: string;
  type: string;
  adminsOnly: boolean;
  defaultValue?: string;
  constructor(input: Partial<ServiceSettingCreationRequestInput> = {}) {
    this.description = input.description = '';
    this.enumeration = input.enumeration = '';
    this.name = input.name = '';
    this.type = input.type = '';
    this.adminsOnly = input.adminsOnly = false;
    this.defaultValue = input.defaultValue;
  }
}
