// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IChangeLog {
   oldValue: string;
 newValue: string;

}

export class ChangeLog implements IChangeLog {
   oldValue: string;
 newValue: string;
constructor(input: Partial<ChangeLog> = {}) {
	 this.oldValue = input.oldValue = '';
 this.newValue = input.newValue = '';
}
}