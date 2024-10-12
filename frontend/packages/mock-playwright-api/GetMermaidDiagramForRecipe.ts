// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetMermaidDiagramForRecipeResponseConfig extends ResponseConfig<string> {
		   recipeID: string;
		

		  constructor( recipeID: string, status: number = 200, body?: string) {
		    super();

		 this.recipeID = recipeID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockGetMermaidDiagramForRecipe = (resCfg: MockGetMermaidDiagramForRecipeResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/mermaid`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};