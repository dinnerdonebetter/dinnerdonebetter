// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSettingConfigurationUpdateRequestInput {
  belongsToAccount: string;
  belongsToUser: string;
  notes: string;
  serviceSettingID: string;
  value: string;
}

export class ServiceSettingConfigurationUpdateRequestInput implements IServiceSettingConfigurationUpdateRequestInput {
  belongsToAccount: string;
  belongsToUser: string;
  notes: string;
  serviceSettingID: string;
  value: string;
  constructor(input: Partial<ServiceSettingConfigurationUpdateRequestInput> = {}) {
    this.belongsToAccount = input.belongsToAccount || '';
    this.belongsToUser = input.belongsToUser || '';
    this.notes = input.notes || '';
    this.serviceSettingID = input.serviceSettingID || '';
    this.value = input.value || '';
  }
}
