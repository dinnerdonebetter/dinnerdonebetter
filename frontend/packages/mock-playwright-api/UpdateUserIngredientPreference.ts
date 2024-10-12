// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { UserIngredientPreference } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockUpdateUserIngredientPreferenceResponseConfig extends ResponseConfig<UserIngredientPreference> {
		   userIngredientPreferenceID: string;
		

		  constructor( userIngredientPreferenceID: string, status: number = 200, body?: UserIngredientPreference) {
		    super();

		 this.userIngredientPreferenceID = userIngredientPreferenceID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockUpdateUserIngredientPreference = (resCfg: MockUpdateUserIngredientPreferenceResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/user_ingredient_preferences/${resCfg.userIngredientPreferenceID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};