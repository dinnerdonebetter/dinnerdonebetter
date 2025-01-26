// GENERATED CODE, DO NOT EDIT MANUALLY

import { AuditLogEntry } from './AuditLogEntry.gen';
import { Household } from './Household.gen';
import { HouseholdInvitation } from './HouseholdInvitation.gen';
import { ServiceSettingConfiguration } from './ServiceSettingConfiguration.gen';

export interface ICoreUserDataCollection {
  auditLogEntries: object;
  households: Household[];
  receivedInvites: HouseholdInvitation[];
  sentInvites: HouseholdInvitation[];
  serviceSettingConfigurations: object;
  userAuditLogEntries: AuditLogEntry[];
  userServiceSettingConfigurations: ServiceSettingConfiguration[];
  webhooks: object;
}

export class CoreUserDataCollection implements ICoreUserDataCollection {
  auditLogEntries: object;
  households: Household[];
  receivedInvites: HouseholdInvitation[];
  sentInvites: HouseholdInvitation[];
  serviceSettingConfigurations: object;
  userAuditLogEntries: AuditLogEntry[];
  userServiceSettingConfigurations: ServiceSettingConfiguration[];
  webhooks: object;
  constructor(input: Partial<CoreUserDataCollection> = {}) {
    this.auditLogEntries = input.auditLogEntries || {};
    this.households = input.households || [];
    this.receivedInvites = input.receivedInvites || [];
    this.sentInvites = input.sentInvites || [];
    this.serviceSettingConfigurations = input.serviceSettingConfigurations || {};
    this.userAuditLogEntries = input.userAuditLogEntries || [];
    this.userServiceSettingConfigurations = input.userServiceSettingConfigurations || [];
    this.webhooks = input.webhooks || {};
  }
}
