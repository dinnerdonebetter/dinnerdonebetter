// GENERATED CODE, DO NOT EDIT MANUALLY

import { ServiceSetting } from './ServiceSetting';

export interface IServiceSettingConfiguration {
  serviceSetting: ServiceSetting;
  belongsToUser: string;
  notes: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  value: string;
  archivedAt: string;
  belongsToHousehold: string;
}

export class ServiceSettingConfiguration implements IServiceSettingConfiguration {
  serviceSetting: ServiceSetting;
  belongsToUser: string;
  notes: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  value: string;
  archivedAt: string;
  belongsToHousehold: string;
  constructor(input: Partial<ServiceSettingConfiguration> = {}) {
    this.serviceSetting = input.serviceSetting || new ServiceSetting();
    this.belongsToUser = input.belongsToUser || '';
    this.notes = input.notes || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.value = input.value || '';
    this.archivedAt = input.archivedAt || '';
    this.belongsToHousehold = input.belongsToHousehold || '';
  }
}
