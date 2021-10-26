<template>
  <div class="bg-gray-100 pt-4">
    <div class="container mx-auto py-4 px-4 bg-white">
      <div class="flex justify-between">
        <div class="flex justify-center">
          <a class="flex items-center px-3 py-1 mr-3 border font-medium rounded text-gray-700">New</a>
        </div>

        <div class="relative block mt-2 sm:mt-0 w-full">
        <span class="absolute inset-y-0 left-0 flex items-center pl-2">
          <svg viewBox="0 0 24 24" class="w-4 h-4 text-gray-500 fill-current">
            <path d="M10 4a6 6 0 100 12 6 6 0 000-12zm-8 6a8 8 0 1114.32 4.906l5.387 5.387a1 1 0 01-1.414 1.414l-5.387-5.387A8 8 0 012 10z"/>
          </svg>
        </span>

          <!-- @keyup="searchQueried" -->
          <!-- v-model="searchQuery" -->
          <!-- :disabled="searchDisabled" -->
          <!-- :placeholder="searchDisabled ? 'disabled' : 'Search'" -->
          <!-- :class="searchDisabled ? 'bg-gray' : 'bg-white'" -->
          <input
              class="block w-full py-2 pl-8 pr-6 text-sm text-gray-700 placeholder-gray-400 border border-b border-gray-400 rounded appearance-none focus:bg-white focus:placeholder-gray-600 focus:text-gray-700 focus:outline-none"
          />
        </div>

        <div class="ml-2">
          <button type="button" class="flex items-center text-gray-700 px-3 py-1 border font-medium rounded">
            <svg viewBox="0 0 24 24" preserveAspectRatio="xMidYMid meet" class="w-5 h-5 mr-1">
              <g class="">
                <path d="M0 0h24v24H0z" fill="none" class=""></path>
                <path d="M3 17v2h6v-2H3zM3 5v2h10V5H3zm10 16v-2h8v-2h-8v-2h-2v6h2zM7 9v2H3v2h4v2h2V9H7zm14 4v-2H11v2h10zm-6-4h2V7h4V5h-4V3h-2v6z" class=""></path>
              </g>
            </svg>
            Filter
          </button>
        </div>
      </div>
    </div>
    <div class="min-h-screen py-6 px-10">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 md:gap-x-10 xl-grid-cols-4 gap-y-10 gap-x-6">
        <div class="container mx-auto shadow-lg rounded-lg max-w-md hover:shadow-2xl transition duration-300" v-for="recipe in recipes" @click="navigateToRecipe(recipe.id)">
          <img v-if="recipe.displayImageURL !== ''" :src="recipe.displayImageURL" alt="" class="rounded-t-lg w-full">
          <div class="p-6">
            <h1 class="md:text-1xl text-xl hover:text-indigo-600 transition duration-200 font-bold text-gray-900">{{ recipe.name }}</h1>
            <p class="text-gray-700 my-2 hover-text-900">Lorem ipsum dolor sit amet consectetur adipisicing elit. Praesentium quis.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent} from "vue";

import {backendRoutes} from "../../constants";
import axios, {AxiosError, AxiosResponse} from "axios";
import {Recipe, RecipeList} from "../../models";
import {settings} from "../../settings/settings";

export default defineComponent({
  data() {
    return {
      recipes: new Array<Recipe>(),
    }
  },
  beforeMount: function () {
    axios.get(`${settings.API_SERVER_URL}${backendRoutes.RECIPES}`)
      .then((res: AxiosResponse<RecipeList>) => {
        this.recipes = res.data.recipes;
        console.dir(this.recipes);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  },
  methods: {
    navigateToRecipe(id: number): void {
      this.$router.push(`/recipes/${id}`);
    },
  },
});
</script>