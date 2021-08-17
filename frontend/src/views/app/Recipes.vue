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
        <div class="container mx-auto shadow-lg rounded-lg max-w-md hover:shadow-2xl transition duration-300" v-for="recipe in recipes">
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

import {Recipe} from "../../models";

const lasagnaRecipe: Recipe = new Recipe();
lasagnaRecipe.name = "Lasagna Strabiliante";
lasagnaRecipe.displayImageURL = "https://thestayathomechef.com/wp-content/uploads/2017/08/Most-Amazing-Lasagna-2-e1574792735811.jpg";

const babaGanoushRecipe: Recipe = new Recipe();
babaGanoushRecipe.name = "Baba Ganoush";
babaGanoushRecipe.displayImageURL = "https://healthyfitnessmeals.com/wp-content/uploads/2020/12/Homemade-baba-ganoush-5.jpg";

const hareALaRoyaleRecipe: Recipe = new Recipe();
hareALaRoyaleRecipe.name = "Hare a la Royale";
hareALaRoyaleRecipe.displayImageURL = "https://static01.nyt.com/images/2016/12/07/dining/07BURNER-DUCK/07BURNER-DUCK-superJumbo.jpg";

const rasElHanoutChickenRecipe: Recipe = new Recipe();
rasElHanoutChickenRecipe.name = "Ras el Hanout Chicken";
rasElHanoutChickenRecipe.displayImageURL = "https://img.delicious.com.au/BXooNiDP/del/2016/02/ras-el-hanout-yoghurt-and-lime-grilled-chicken-27350-1.jpg";

const tacosAlPastorRecipe: Recipe = new Recipe();
tacosAlPastorRecipe.name = "Tacos al Pastor";
tacosAlPastorRecipe.displayImageURL = "https://www.comedera.com/wp-content/uploads/2017/08/tacos-al-pastor-receta.jpg";

const chuCheeCurryChickenRecipe: Recipe = new Recipe();
chuCheeCurryChickenRecipe.name = "Chu Chee Curry Chicken";
chuCheeCurryChickenRecipe.displayImageURL = "https://www.wandercooks.com/wp-content/uploads/2020/06/choo-chee-chicken-curry-ft-1.jpg";

export default defineComponent({
  data() {
    return {
      recipes: [
        lasagnaRecipe,
        babaGanoushRecipe,
        hareALaRoyaleRecipe,
        rasElHanoutChickenRecipe,
        tacosAlPastorRecipe,
        chuCheeCurryChickenRecipe,
      ]
    }
  },
  methods: {
  },
});
</script>