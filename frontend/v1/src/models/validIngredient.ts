export class ValidIngredient {
  id: number;
  name: string;
  variant: string;
  description: string;
  warning: string;
  containsEgg: boolean;
  containsDairy: boolean;
  containsPeanut: boolean;
  containsTreeNut: boolean;
  containsSoy: boolean;
  containsWheat: boolean;
  containsShellfish: boolean;
  containsSesame: boolean;
  containsFish: boolean;
  containsGluten: boolean;
  animalFlesh: boolean;
  animalDerived: boolean;
  measurableByVolume: boolean;
  icon: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.name = '';
    this.id = 0;
    this.variant = '';
    this.description = '';
    this.warning = '';
    this.containsEgg = false;
    this.containsDairy = false;
    this.containsPeanut = false;
    this.containsTreeNut = false;
    this.containsSoy = false;
    this.containsWheat = false;
    this.containsShellfish = false;
    this.containsSesame = false;
    this.containsFish = false;
    this.containsGluten = false;
    this.animalFlesh = false;
    this.animalDerived = false;
    this.measurableByVolume = false;
    this.createdOn = 0;
    this.icon = '';
  }
}

export function validIngredientsAreEqual(
  i1: ValidIngredient,
  i2: ValidIngredient,
): boolean {
  return (
    i1.id === i2.id &&
    i1.name === i2.name &&
    i1.variant === i2.variant &&
    i1.description === i2.description &&
    i1.warning === i2.warning &&    i1.containsDairy === i2.containsDairy &&
    i1.containsPeanut === i2.containsPeanut &&
    i1.containsTreeNut === i2.containsTreeNut &&
    i1.containsSoy === i2.containsSoy &&
    i1.containsWheat === i2.containsWheat &&
    i1.containsShellfish === i2.containsShellfish &&
    i1.containsSesame === i2.containsSesame &&
    i1.containsFish === i2.containsFish &&
    i1.containsGluten === i2.containsGluten &&
    i1.animalFlesh === i2.animalFlesh &&
    i1.animalDerived === i2.animalDerived &&
    i1.measurableByVolume === i2.measurableByVolume &&
    i1.icon === i2.icon &&
    i1.createdOn === i2.createdOn &&
    i1.updatedOn === i2.updatedOn &&
    i1.archivedOn === i2.archivedOn
  );
}
