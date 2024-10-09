// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSetting {
  description: string;
  name: string;
  archivedAt?: string;
  createdAt: string;
  enumeration: string;
  id: string;
  lastUpdatedAt?: string;
  type: string;
  adminsOnly: boolean;
  defaultValue?: string;
}

export class ServiceSetting implements IServiceSetting {
  description: string;
  name: string;
  archivedAt?: string;
  createdAt: string;
  enumeration: string;
  id: string;
  lastUpdatedAt?: string;
  type: string;
  adminsOnly: boolean;
  defaultValue?: string;
  constructor(input: Partial<ServiceSetting> = {}) {
    this.description = input.description = '';
    this.name = input.name = '';
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.enumeration = input.enumeration = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.type = input.type = '';
    this.adminsOnly = input.adminsOnly = false;
    this.defaultValue = input.defaultValue;
  }
}
