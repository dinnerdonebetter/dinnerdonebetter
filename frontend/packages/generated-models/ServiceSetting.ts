// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IServiceSetting {
   id: string;
 adminsOnly: boolean;
 createdAt: string;
 defaultValue?: string;
 enumeration: string;
 type: string;
 archivedAt?: string;
 description: string;
 lastUpdatedAt?: string;
 name: string;

}

export class ServiceSetting implements IServiceSetting {
   id: string;
 adminsOnly: boolean;
 createdAt: string;
 defaultValue?: string;
 enumeration: string;
 type: string;
 archivedAt?: string;
 description: string;
 lastUpdatedAt?: string;
 name: string;
constructor(input: Partial<ServiceSetting> = {}) {
	 this.id = input.id = '';
 this.adminsOnly = input.adminsOnly = false;
 this.createdAt = input.createdAt = '';
 this.defaultValue = input.defaultValue;
 this.enumeration = input.enumeration = '';
 this.type = input.type = '';
 this.archivedAt = input.archivedAt;
 this.description = input.description = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.name = input.name = '';
}
}