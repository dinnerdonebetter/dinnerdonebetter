// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSetting {
  id: string;
  name: string;
  adminsOnly: boolean;
  archivedAt?: string;
  description: string;
  enumeration: string;
  lastUpdatedAt?: string;
  type: string;
  createdAt: string;
  defaultValue?: string;
}

export class ServiceSetting implements IServiceSetting {
  id: string;
  name: string;
  adminsOnly: boolean;
  archivedAt?: string;
  description: string;
  enumeration: string;
  lastUpdatedAt?: string;
  type: string;
  createdAt: string;
  defaultValue?: string;
  constructor(input: Partial<ServiceSetting> = {}) {
    this.id = input.id = '';
    this.name = input.name = '';
    this.adminsOnly = input.adminsOnly = false;
    this.archivedAt = input.archivedAt;
    this.description = input.description = '';
    this.enumeration = input.enumeration = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.type = input.type = '';
    this.createdAt = input.createdAt = '';
    this.defaultValue = input.defaultValue;
  }
}
