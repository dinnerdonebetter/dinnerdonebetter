import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";

import Dashboard from "./views/admin/Dashboard.vue";
import Login from "./views/app/Login.vue";
import Register from "./views/app/Register.vue";

import AdminUsers from "./views/admin/tables/Users.vue";
import AdminValidInstrument from "./views/admin/editors/ValidInstrument.vue";
import AdminValidIngredient from "./views/admin/editors/ValidIngredient.vue";
import AdminValidPreparation from "./views/admin/editors/ValidPreparation.vue";
import AdminValidIngredients from "./views/admin/tables/ValidIngredients.vue";
import AdminValidInstruments from "./views/admin/tables/ValidInstruments.vue";
import AdminValidPreparations from "./views/admin/tables/ValidPreparations.vue";
import AdminRecipes from "./views/admin/tables/Recipes.vue";
import AdminRecipeBuilder from "./views/admin/editors/Recipe.vue";

import Blank from "./views/Blank.vue";
import Home from "./views/app/Home.vue";
import Household from "./views/app/Household.vue";
import MealPlans from "./views/app/MealPlans.vue";
import Recipes from "./views/app/Recipes.vue";

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "Login",
    component: Login,
    meta: { layout: "empty" },
  },
  {
    path: "/register",
    name: "Register",
    component: Register,
    meta: { layout: "empty" },
  },
  {
    path: "/",
    name: "Landing",
    component: Blank,
    meta: { layout: "app" },
  },
  {
    path: "/home",
    name: "Home",
    component: Home,
    meta: { layout: "app" },
  },
  {
    path: "/household",
    name: "Household",
    component: Household,
    meta: { layout: "app" },
  },
  {
    path: "/plans",
    name: "MealPlans",
    component: MealPlans,
    meta: { layout: "app" },
  },
  {
    path: "/recipes",
    name: "Recipes",
    component: Recipes,
    meta: { layout: "app" },
  },
  {
    path: "/admin",
    redirect: "/admin/dashboard",
  },
  {
    path: "/admin/dashboard",
    name: "Dashboard",
    component: Dashboard,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/users",
    name: "AdminUsers",
    component: AdminUsers,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/users/:userID",
    name: "AdminUser",
    component: AdminUsers,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/valid_ingredients",
    name: "AdminValidIngredients",
    component: AdminValidIngredients,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/valid_ingredients/:ingredientID",
    name: "AdminValidIngredient",
    component: AdminValidIngredient,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/valid_ingredients/new",
    name: "AdminValidIngredientCreator",
    component: AdminValidIngredient,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/valid_instruments",
    name: "AdminValidInstruments",
    component: AdminValidInstruments,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/valid_instruments/new",
    name: "AdminValidInstrumentCreator",
    component: AdminValidInstrument,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/valid_instruments/:instrumentID",
    name: "AdminValidInstrument",
    component: AdminValidInstrument,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/valid_preparations",
    name: "AdminValidPreparations",
    component: AdminValidPreparations,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/valid_preparations/new",
    name: "AdminValidPreparationCreator",
    component: AdminValidPreparation,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/valid_preparations/:preparationID",
    name: "AdminValidPreparation",
    component: AdminValidPreparation,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/recipes",
    name: "AdminRecipes",
    component: AdminRecipes,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/recipes/:recipeID",
    name: "AdminRecipe",
    component: AdminRecipes,
    meta: { layout: "admin" },
  },
  {
    path: "/admin/recipes/new",
    name: "AdminRecipeBuilder",
    component: AdminRecipeBuilder,
    meta: { layout: "admin" },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes: routes,
});

export default router;
