// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IServiceSettingCreationRequestInput {
   defaultValue?: string;
 description: string;
 enumeration: string;
 name: string;
 type: string;
 adminsOnly: boolean;

}

export class ServiceSettingCreationRequestInput implements IServiceSettingCreationRequestInput {
   defaultValue?: string;
 description: string;
 enumeration: string;
 name: string;
 type: string;
 adminsOnly: boolean;
constructor(input: Partial<ServiceSettingCreationRequestInput> = {}) {
	 this.defaultValue = input.defaultValue;
 this.description = input.description = '';
 this.enumeration = input.enumeration = '';
 this.name = input.name = '';
 this.type = input.type = '';
 this.adminsOnly = input.adminsOnly = false;
}
}