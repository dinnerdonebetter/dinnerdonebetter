// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSettingConfigurationUpdateRequestInput {
  belongsToUser?: string;
  notes?: string;
  serviceSettingID?: string;
  value?: string;
  belongsToHousehold?: string;
}

export class ServiceSettingConfigurationUpdateRequestInput implements IServiceSettingConfigurationUpdateRequestInput {
  belongsToUser?: string;
  notes?: string;
  serviceSettingID?: string;
  value?: string;
  belongsToHousehold?: string;
  constructor(input: Partial<ServiceSettingConfigurationUpdateRequestInput> = {}) {
    this.belongsToUser = input.belongsToUser;
    this.notes = input.notes;
    this.serviceSettingID = input.serviceSettingID;
    this.value = input.value;
    this.belongsToHousehold = input.belongsToHousehold;
  }
}
