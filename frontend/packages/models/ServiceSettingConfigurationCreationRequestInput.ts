// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSettingConfigurationCreationRequestInput {
  belongsToHousehold: string;
  belongsToUser: string;
  notes: string;
  serviceSettingID: string;
  value: string;
}

export class ServiceSettingConfigurationCreationRequestInput
  implements IServiceSettingConfigurationCreationRequestInput
{
  belongsToHousehold: string;
  belongsToUser: string;
  notes: string;
  serviceSettingID: string;
  value: string;
  constructor(input: Partial<ServiceSettingConfigurationCreationRequestInput> = {}) {
    this.belongsToHousehold = input.belongsToHousehold || '';
    this.belongsToUser = input.belongsToUser || '';
    this.notes = input.notes || '';
    this.serviceSettingID = input.serviceSettingID || '';
    this.value = input.value || '';
  }
}
