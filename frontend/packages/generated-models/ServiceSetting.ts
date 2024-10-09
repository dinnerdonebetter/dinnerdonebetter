// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSetting {
  adminsOnly: boolean;
  archivedAt?: string;
  description: string;
  enumeration: string;
  name: string;
  type: string;
  createdAt: string;
  defaultValue?: string;
  id: string;
  lastUpdatedAt?: string;
}

export class ServiceSetting implements IServiceSetting {
  adminsOnly: boolean;
  archivedAt?: string;
  description: string;
  enumeration: string;
  name: string;
  type: string;
  createdAt: string;
  defaultValue?: string;
  id: string;
  lastUpdatedAt?: string;
  constructor(input: Partial<ServiceSetting> = {}) {
    this.adminsOnly = input.adminsOnly = false;
    this.archivedAt = input.archivedAt;
    this.description = input.description = '';
    this.enumeration = input.enumeration = '';
    this.name = input.name = '';
    this.type = input.type = '';
    this.createdAt = input.createdAt = '';
    this.defaultValue = input.defaultValue;
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
  }
}
