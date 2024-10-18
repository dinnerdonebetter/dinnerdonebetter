// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IFinalizeMealPlansRequest {
  returnCount: boolean;
}

export class FinalizeMealPlansRequest implements IFinalizeMealPlansRequest {
  returnCount: boolean;
  constructor(input: Partial<FinalizeMealPlansRequest> = {}) {
    this.returnCount = input.returnCount || false;
  }
}
