// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidVessel } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockUpdateValidVesselResponseConfig extends ResponseConfig<ValidVessel> {
		   validVesselID: string;
		

		  constructor( validVesselID: string, status: number = 200, body?: ValidVessel) {
		    super();

		 this.validVesselID = validVesselID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockUpdateValidVessel = (resCfg: MockUpdateValidVesselResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_vessels/${resCfg.validVesselID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};