// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IAvatarUpdateInput {
  base64EncodedData: string;
}

export class AvatarUpdateInput implements IAvatarUpdateInput {
  base64EncodedData: string;
  constructor(input: Partial<AvatarUpdateInput> = {}) {
    this.base64EncodedData = input.base64EncodedData || '';
  }
}
