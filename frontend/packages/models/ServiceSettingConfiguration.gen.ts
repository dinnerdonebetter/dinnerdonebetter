// GENERATED CODE, DO NOT EDIT MANUALLY

import { ServiceSetting } from './ServiceSetting.gen';

export interface IServiceSettingConfiguration {
  archivedAt: string;
  belongsToAccount: string;
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
  belongsToAccount: string;
  belongsToUser: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  notes: string;
  serviceSetting: ServiceSetting;
  value: string;
  constructor(input: Partial<ServiceSettingConfiguration> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToAccount = input.belongsToAccount || '';
    this.belongsToUser = input.belongsToUser || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.serviceSetting = input.serviceSetting || new ServiceSetting();
    this.value = input.value || '';
  }
}
