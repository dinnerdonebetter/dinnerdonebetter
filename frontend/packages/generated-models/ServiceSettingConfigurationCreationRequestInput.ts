// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSettingConfigurationCreationRequestInput {
  serviceSettingID: string;
  value: string;
  belongsToHousehold: string;
  belongsToUser: string;
  notes: string;
}

export class ServiceSettingConfigurationCreationRequestInput
  implements IServiceSettingConfigurationCreationRequestInput
{
  serviceSettingID: string;
  value: string;
  belongsToHousehold: string;
  belongsToUser: string;
  notes: string;
  constructor(input: Partial<ServiceSettingConfigurationCreationRequestInput> = {}) {
    this.serviceSettingID = input.serviceSettingID = '';
    this.value = input.value = '';
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.belongsToUser = input.belongsToUser = '';
    this.notes = input.notes = '';
  }
}
