// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { UserNotification } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockUpdateUserNotificationResponseConfig extends ResponseConfig<UserNotification> {
		   userNotificationID: string;
		

		  constructor( userNotificationID: string, status: number = 200, body?: UserNotification) {
		    super();

		 this.userNotificationID = userNotificationID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockUpdateUserNotification = (resCfg: MockUpdateUserNotificationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/user_notifications/${resCfg.userNotificationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PATCH', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};