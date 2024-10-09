// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IServiceSettingConfigurationCreationRequestInput {
   value: string;
 belongsToHousehold: string;
 belongsToUser: string;
 notes: string;
 serviceSettingID: string;

}

export class ServiceSettingConfigurationCreationRequestInput implements IServiceSettingConfigurationCreationRequestInput {
   value: string;
 belongsToHousehold: string;
 belongsToUser: string;
 notes: string;
 serviceSettingID: string;
constructor(input: Partial<ServiceSettingConfigurationCreationRequestInput> = {}) {
	 this.value = input.value = '';
 this.belongsToHousehold = input.belongsToHousehold = '';
 this.belongsToUser = input.belongsToUser = '';
 this.notes = input.notes = '';
 this.serviceSettingID = input.serviceSettingID = '';
}
}