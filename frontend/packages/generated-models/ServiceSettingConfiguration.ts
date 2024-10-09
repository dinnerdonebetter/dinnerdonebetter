// GENERATED CODE, DO NOT EDIT MANUALLY

import { ServiceSetting } from './ServiceSetting';

export interface IServiceSettingConfiguration {
  archivedAt?: string;
  belongsToHousehold: string;
  id: string;
  notes: string;
  serviceSetting: ServiceSetting;
  value: string;
  belongsToUser: string;
  createdAt: string;
  lastUpdatedAt?: string;
}

export class ServiceSettingConfiguration implements IServiceSettingConfiguration {
  archivedAt?: string;
  belongsToHousehold: string;
  id: string;
  notes: string;
  serviceSetting: ServiceSetting;
  value: string;
  belongsToUser: string;
  createdAt: string;
  lastUpdatedAt?: string;
  constructor(input: Partial<ServiceSettingConfiguration> = {}) {
    this.archivedAt = input.archivedAt;
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.id = input.id = '';
    this.notes = input.notes = '';
    this.serviceSetting = input.serviceSetting = new ServiceSetting();
    this.value = input.value = '';
    this.belongsToUser = input.belongsToUser = '';
    this.createdAt = input.createdAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
  }
}
