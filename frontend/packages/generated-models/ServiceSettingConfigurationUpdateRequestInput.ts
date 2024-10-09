// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IServiceSettingConfigurationUpdateRequestInput {
   value?: string;
 belongsToHousehold?: string;
 belongsToUser?: string;
 notes?: string;
 serviceSettingID?: string;

}

export class ServiceSettingConfigurationUpdateRequestInput implements IServiceSettingConfigurationUpdateRequestInput {
   value?: string;
 belongsToHousehold?: string;
 belongsToUser?: string;
 notes?: string;
 serviceSettingID?: string;
constructor(input: Partial<ServiceSettingConfigurationUpdateRequestInput> = {}) {
	 this.value = input.value;
 this.belongsToHousehold = input.belongsToHousehold;
 this.belongsToUser = input.belongsToUser;
 this.notes = input.notes;
 this.serviceSettingID = input.serviceSettingID;
}
}