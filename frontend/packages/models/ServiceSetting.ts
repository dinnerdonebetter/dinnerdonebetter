// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSetting {
  adminsOnly: boolean;
  archivedAt: string;
  createdAt: string;
  defaultValue: string;
  description: string;
  enumeration: string[];
  id: string;
  lastUpdatedAt: string;
  name: string;
  type: string;
}

export class ServiceSetting implements IServiceSetting {
  adminsOnly: boolean;
  archivedAt: string;
  createdAt: string;
  defaultValue: string;
  description: string;
  enumeration: string[];
  id: string;
  lastUpdatedAt: string;
  name: string;
  type: string;
  constructor(input: Partial<ServiceSetting> = {}) {
    this.adminsOnly = input.adminsOnly || false;
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.defaultValue = input.defaultValue || '';
    this.description = input.description || '';
    this.enumeration = input.enumeration || [];
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.name = input.name || '';
    this.type = input.type || '';
  }
}
