// GENERATED CODE, DO NOT EDIT MANUALLY

import { CoreUserDataCollection } from './CoreUserDataCollection.gen';
import { EatingUserDataCollection } from './EatingUserDataCollection.gen';
import { User } from './User.gen';

export interface IUserDataCollection {
  core: CoreUserDataCollection;
  eating: EatingUserDataCollection;
  reportID: string;
  user: User;
}

export class UserDataCollection implements IUserDataCollection {
  core: CoreUserDataCollection;
  eating: EatingUserDataCollection;
  reportID: string;
  user: User;
  constructor(input: Partial<UserDataCollection> = {}) {
    this.core = input.core || new CoreUserDataCollection();
    this.eating = input.eating || new EatingUserDataCollection();
    this.reportID = input.reportID || '';
    this.user = input.user || new User();
  }
}
