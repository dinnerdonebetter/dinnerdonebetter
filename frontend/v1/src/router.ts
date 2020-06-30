import Vue from "vue";
import Router from "vue-router";

import { routes } from "@/constants/routes";

Vue.use(Router);

/*
  redirect:            if set to 'noredirect', no redirect action will be trigger when clicking the breadcrumb
  meta: {
    title: 'title'     the name showed in subMenu and breadcrumb (recommend set)
    icon: 'svg-name'   the icon showed in the sidebar
    breadcrumb: false  if false, the item will be hidden in breadcrumb (default is true)
    hidden: true       if true, this route will not show in the sidebar (default is false)
  }
*/

export default new Router({
  mode: 'history',
  scrollBehavior: (to, from, savedPosition) => {
    if (savedPosition) {
      return savedPosition;
    } else {
      return { x: 0, y: 0 };
    }
  },
  base: process.env.BASE_URL,
  routes: routes,
});
