export class ValidIngredient {
     id: number;
     externalID: string;
     name: string;
     variant: string;
     description: string;
     warning: string;
     iconPath: string;
     containsSoy: boolean;
     containsTreeNut: boolean;
     containsShellfish: boolean;
     containsSesame: boolean;
     containsFish: boolean;
     containsGluten: boolean;
     animalFlesh: boolean;
     animalDerived: boolean;
     volumetric: boolean;
     containsPeanut: boolean;
     containsDairy: boolean;
     containsEgg: boolean;
     containsWheat: boolean;
     createdOn: number;
     lastUpdatedOn?: number;
     archivedOn?: number;

    constructor(
      id: number = 0,
      externalID: string = "",
      name: string = "",
      variant: string = "",
      description: string = "",
      warning: string = "",
      iconPath: string = "",
      containsSoy: boolean = false,
      containsTreeNut: boolean = false,
      containsShellfish: boolean = false,
      containsSesame: boolean = false,
      containsFish: boolean = false,
      containsGluten: boolean = false,
      animalFlesh: boolean = false,
      animalDerived: boolean = false,
      volumetric: boolean = false,
      containsPeanut: boolean = false,
      containsDairy: boolean = false,
      containsEgg: boolean = false,
      containsWheat: boolean = false,
      createdOn: number = 0,
      lastUpdatedOn?: number,
      archivedOn?: number,
    ) {
         this.id = id;
         this.externalID = externalID;
         this.name = name;
         this.variant = variant;
         this.description = description;
         this.warning = warning;
         this.iconPath = iconPath;
         this.containsSoy = containsSoy;
         this.containsTreeNut = containsTreeNut;
         this.containsShellfish = containsShellfish;
         this.containsSesame = containsSesame;
         this.containsFish = containsFish;
         this.containsGluten = containsGluten;
         this.animalFlesh = animalFlesh;
         this.animalDerived = animalDerived;
         this.volumetric = volumetric;
         this.containsPeanut = containsPeanut;
         this.containsDairy = containsDairy;
         this.containsEgg = containsEgg;
         this.containsWheat = containsWheat;
         this.createdOn = createdOn;
         this.lastUpdatedOn = lastUpdatedOn;
         this.archivedOn = archivedOn;
    }
}

export class ValidIngredientList {
     validIngredients: ValidIngredient[];
     totalCount: number;
     page: number;
     limit: number;
     filteredCount: number;

     constructor(
          validIngredients: ValidIngredient[] = [],
          totalCount: number = 0,
          page: number = 0,
          limit: number = 0,
          filteredCount: number = 0,
     ) {
          this.validIngredients = validIngredients;
          this.totalCount = totalCount;
          this.page = page;
          this.limit = limit;
          this.filteredCount = filteredCount;
     }
}