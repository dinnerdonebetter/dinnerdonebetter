export class ValidPreparation {
  id: number;
  name: string;
  description: string;
  icon: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    this.name = '';
    this.description = '';
    this.createdOn = 0;
    this.icon = '';
  }
}

export function validPreparationsAreEqual(
  i1: ValidPreparation,
  i2: ValidPreparation,
): boolean {
  return (
    i1.id === i2.id &&
    i1.name === i2.name &&
    i1.description === i2.description &&
    i1.icon === i2.icon &&
    i1.createdOn === i2.createdOn &&
    i1.updatedOn === i2.updatedOn &&
    i1.archivedOn === i2.archivedOn
  );
}
