// GENERATED CODE, DO NOT EDIT MANUALLY

import { ServiceSetting } from './ServiceSetting';

export interface IServiceSettingConfiguration {
  belongsToHousehold: string;
  belongsToUser: string;
  id: string;
  notes: string;
  archivedAt?: string;
  createdAt: string;
  lastUpdatedAt?: string;
  serviceSetting: ServiceSetting;
  value: string;
}

export class ServiceSettingConfiguration implements IServiceSettingConfiguration {
  belongsToHousehold: string;
  belongsToUser: string;
  id: string;
  notes: string;
  archivedAt?: string;
  createdAt: string;
  lastUpdatedAt?: string;
  serviceSetting: ServiceSetting;
  value: string;
  constructor(input: Partial<ServiceSettingConfiguration> = {}) {
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.belongsToUser = input.belongsToUser = '';
    this.id = input.id = '';
    this.notes = input.notes = '';
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.serviceSetting = input.serviceSetting = new ServiceSetting();
    this.value = input.value = '';
  }
}
