// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSettingConfigurationCreationRequestInput {
  belongsToAccount: string;
  belongsToUser: string;
  notes: string;
  serviceSettingID: string;
  value: string;
}

export class ServiceSettingConfigurationCreationRequestInput
  implements IServiceSettingConfigurationCreationRequestInput
{
  belongsToAccount: string;
  belongsToUser: string;
  notes: string;
  serviceSettingID: string;
  value: string;
  constructor(input: Partial<ServiceSettingConfigurationCreationRequestInput> = {}) {
    this.belongsToAccount = input.belongsToAccount || '';
    this.belongsToUser = input.belongsToUser || '';
    this.notes = input.notes || '';
    this.serviceSettingID = input.serviceSettingID || '';
    this.value = input.value || '';
  }
}
