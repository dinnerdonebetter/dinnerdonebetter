<template>
  <div class="grid gap-4 grid-cols-1 justify-items-center mt-3">
    <div class="w-1/2 p-6 rounded-lg shadow-lg bg-gray-50">
      <h1 class="text-xl text-black">{{ recipe.name }}</h1>
      <div class="w-full m-2 p-6 rounded-lg shadow-lg bg-white" v-for="step in recipe.steps">
        <h2 class="text-lg text-gray-900">{{ step.preparation.name }}</h2>
        <ul class="ml-4 list-disc">
          <li v-for="ingredient in step.ingredients">{{ ingredient.quantityValue }} {{ ingredient.quantityType }} {{ ingredient.ingredient.name }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent} from "vue";

import {backendRoutes} from "../../constants";
import axios, {AxiosError, AxiosResponse} from "axios";
import {FullRecipe} from "../../models";
import format from "string-format";

export default defineComponent({
  data() {
    return {
      recipe: new FullRecipe(),
      recipeAPIPath: "",
    }
  },
  beforeMount: function () {
    const recipeID = this.$route.params.recipeID;

    if (recipeID) {
      this.recipeAPIPath = format(backendRoutes.RECIPE, recipeID.toString());

      axios.get(this.recipeAPIPath)
        .then((res: AxiosResponse<FullRecipe>) => {
          this.recipe = { ...res.data };
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    }
  },
})
</script>