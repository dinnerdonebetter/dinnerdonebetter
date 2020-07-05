import NProgress from 'nprogress';
import 'nprogress/nprogress.css';
import { NavigationGuardNext, Route } from 'vue-router';

import router from './router';
import { UserModule } from '@/store/modules/user';
import {AppModule} from "@/store/modules/app";

NProgress.configure({ showSpinner: false });

const whiteList = ['/login', '/register', '/'];

router.beforeEach(async(to: Route, from: Route, next: NavigationGuardNext) => {
  // Start progress bar
  NProgress.start();

  // conditionally wait, since the login page promise doesn't resolve properly
  const waitPeriod = from.path === "/login" ? 1500 : 0; // lol I fucking hate this

  if (!AppModule.frontendDevMode) {
    setTimeout(() => {
      // Determine whether the user has logged in
      if (UserModule.isLoggedIn) {
        if (to.path === '/login') {
          // If is logged in, redirect to the home page
          next({ path: '/' });
          NProgress.done();
        } else if (to.path.startsWith("/admin") && !UserModule.isAdmin) {
          // if isn't an admin, divert away from admin routes
          next({ path: '/' });
          NProgress.done();
        } else {
          next();
        }
      } else {
        // Has no token
        if (whiteList.indexOf(to.path) !== -1) {
          // In the free login whitelist, go directly
          next();
        } else {
          // Other pages that do not have permission to access are redirected to the login page.
          next(`/login?redirect=${to.path}`);
          NProgress.done();
        }
      }
    }, waitPeriod);
  } else {
    next();
  }
});

router.afterEach((to: Route) => {
  // Finish progress bar
  NProgress.done();

  // set page title
  document.title = to.meta.title;
});
