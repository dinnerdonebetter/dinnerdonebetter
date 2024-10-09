// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IServiceSetting {
  archivedAt: string;
  defaultValue: string;
  enumeration: string[];
  name: string;
  type: string;
  adminsOnly: boolean;
  description: string;
  id: string;
  lastUpdatedAt: string;
  createdAt: string;
}

export class ServiceSetting implements IServiceSetting {
  archivedAt: string;
  defaultValue: string;
  enumeration: string[];
  name: string;
  type: string;
  adminsOnly: boolean;
  description: string;
  id: string;
  lastUpdatedAt: string;
  createdAt: string;
  constructor(input: Partial<ServiceSetting> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.defaultValue = input.defaultValue || '';
    this.enumeration = input.enumeration || [];
    this.name = input.name || '';
    this.type = input.type || '';
    this.adminsOnly = input.adminsOnly || false;
    this.description = input.description || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.createdAt = input.createdAt || '';
  }
}
