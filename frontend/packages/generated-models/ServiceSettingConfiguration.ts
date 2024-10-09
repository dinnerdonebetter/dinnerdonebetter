// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ServiceSetting } from './ServiceSetting';


export interface IServiceSettingConfiguration {
   belongsToUser: string;
 createdAt: string;
 notes: string;
 value: string;
 archivedAt?: string;
 belongsToHousehold: string;
 id: string;
 lastUpdatedAt?: string;
 serviceSetting: ServiceSetting;

}

export class ServiceSettingConfiguration implements IServiceSettingConfiguration {
   belongsToUser: string;
 createdAt: string;
 notes: string;
 value: string;
 archivedAt?: string;
 belongsToHousehold: string;
 id: string;
 lastUpdatedAt?: string;
 serviceSetting: ServiceSetting;
constructor(input: Partial<ServiceSettingConfiguration> = {}) {
	 this.belongsToUser = input.belongsToUser = '';
 this.createdAt = input.createdAt = '';
 this.notes = input.notes = '';
 this.value = input.value = '';
 this.archivedAt = input.archivedAt;
 this.belongsToHousehold = input.belongsToHousehold = '';
 this.id = input.id = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.serviceSetting = input.serviceSetting = new ServiceSetting();
}
}