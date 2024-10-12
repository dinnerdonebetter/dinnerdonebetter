// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ServiceSetting } from './ServiceSetting';


export interface IServiceSettingConfiguration {
   archivedAt: string;
 belongsToHousehold: string;
 belongsToUser: string;
 createdAt: string;
 id: string;
 lastUpdatedAt: string;
 notes: string;
 serviceSetting: ServiceSetting;
 value: string;

}

export class ServiceSettingConfiguration implements IServiceSettingConfiguration {
   archivedAt: string;
 belongsToHousehold: string;
 belongsToUser: string;
 createdAt: string;
 id: string;
 lastUpdatedAt: string;
 notes: string;
 serviceSetting: ServiceSetting;
 value: string;
constructor(input: Partial<ServiceSettingConfiguration> = {}) {
	 this.archivedAt = input.archivedAt || '';
 this.belongsToHousehold = input.belongsToHousehold || '';
 this.belongsToUser = input.belongsToUser || '';
 this.createdAt = input.createdAt || '';
 this.id = input.id || '';
 this.lastUpdatedAt = input.lastUpdatedAt || '';
 this.notes = input.notes || '';
 this.serviceSetting = input.serviceSetting || new ServiceSetting();
 this.value = input.value || '';
}
}