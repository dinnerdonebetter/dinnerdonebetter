// GENERATED CODE, DO NOT EDIT MANUALLY

import { AuditLogEntry } from './AuditLogEntry.gen';
import { Account } from './Account.gen';
import { AccountInvitation } from './AccountInvitation.gen';
import { ServiceSettingConfiguration } from './ServiceSettingConfiguration.gen';

export interface ICoreUserDataCollection {
  auditLogEntries: object;
  accounts: Account[];
  receivedInvites: AccountInvitation[];
  sentInvites: AccountInvitation[];
  serviceSettingConfigurations: object;
  userAuditLogEntries: AuditLogEntry[];
  userServiceSettingConfigurations: ServiceSettingConfiguration[];
  webhooks: object;
}

export class CoreUserDataCollection implements ICoreUserDataCollection {
  auditLogEntries: object;
  accounts: Account[];
  receivedInvites: AccountInvitation[];
  sentInvites: AccountInvitation[];
  serviceSettingConfigurations: object;
  userAuditLogEntries: AuditLogEntry[];
  userServiceSettingConfigurations: ServiceSettingConfiguration[];
  webhooks: object;
  constructor(input: Partial<CoreUserDataCollection> = {}) {
    this.auditLogEntries = input.auditLogEntries || {};
    this.accounts = input.accounts || [];
    this.receivedInvites = input.receivedInvites || [];
    this.sentInvites = input.sentInvites || [];
    this.serviceSettingConfigurations = input.serviceSettingConfigurations || {};
    this.userAuditLogEntries = input.userAuditLogEntries || [];
    this.userServiceSettingConfigurations = input.userServiceSettingConfigurations || [];
    this.webhooks = input.webhooks || {};
  }
}
