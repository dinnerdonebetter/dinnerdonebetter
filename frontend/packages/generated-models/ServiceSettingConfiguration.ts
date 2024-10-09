// GENERATED CODE, DO NOT EDIT MANUALLY

import { ServiceSetting } from './ServiceSetting';

export interface IServiceSettingConfiguration {
  belongsToUser: string;
  createdAt: string;
  value: string;
  serviceSetting: ServiceSetting;
  archivedAt?: string;
  belongsToHousehold: string;
  id: string;
  lastUpdatedAt?: string;
  notes: string;
}

export class ServiceSettingConfiguration implements IServiceSettingConfiguration {
  belongsToUser: string;
  createdAt: string;
  value: string;
  serviceSetting: ServiceSetting;
  archivedAt?: string;
  belongsToHousehold: string;
  id: string;
  lastUpdatedAt?: string;
  notes: string;
  constructor(input: Partial<ServiceSettingConfiguration> = {}) {
    this.belongsToUser = input.belongsToUser = '';
    this.createdAt = input.createdAt = '';
    this.value = input.value = '';
    this.serviceSetting = input.serviceSetting = new ServiceSetting();
    this.archivedAt = input.archivedAt;
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
  }
}
