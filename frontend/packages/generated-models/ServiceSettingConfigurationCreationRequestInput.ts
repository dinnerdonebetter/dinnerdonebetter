// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSettingConfigurationCreationRequestInput {
  notes: string;
  serviceSettingID: string;
  value: string;
  belongsToHousehold: string;
  belongsToUser: string;
}

export class ServiceSettingConfigurationCreationRequestInput
  implements IServiceSettingConfigurationCreationRequestInput
{
  notes: string;
  serviceSettingID: string;
  value: string;
  belongsToHousehold: string;
  belongsToUser: string;
  constructor(input: Partial<ServiceSettingConfigurationCreationRequestInput> = {}) {
    this.notes = input.notes = '';
    this.serviceSettingID = input.serviceSettingID = '';
    this.value = input.value = '';
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.belongsToUser = input.belongsToUser = '';
  }
}
