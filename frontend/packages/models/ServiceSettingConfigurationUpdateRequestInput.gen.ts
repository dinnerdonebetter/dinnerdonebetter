// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSettingConfigurationUpdateRequestInput {
  belongsToHousehold: string;
  belongsToUser: string;
  notes: string;
  serviceSettingID: string;
  value: string;
}

export class ServiceSettingConfigurationUpdateRequestInput implements IServiceSettingConfigurationUpdateRequestInput {
  belongsToHousehold: string;
  belongsToUser: string;
  notes: string;
  serviceSettingID: string;
  value: string;
  constructor(input: Partial<ServiceSettingConfigurationUpdateRequestInput> = {}) {
    this.belongsToHousehold = input.belongsToHousehold || '';
    this.belongsToUser = input.belongsToUser || '';
    this.notes = input.notes || '';
    this.serviceSettingID = input.serviceSettingID || '';
    this.value = input.value || '';
  }
}
