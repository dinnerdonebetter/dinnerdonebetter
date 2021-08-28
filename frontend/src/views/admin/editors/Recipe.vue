<template>
  <div class="mt-4">
    <div class="p-6 bg-white rounded-md shadow-md">
      <h2 class="text-lg font-semibold text-gray-700 capitalize"> New Recipe </h2>
        <div class="flex">
          <div class="flex-initial p-2">
            <label class="text-gray-700" for="recipeName">Name</label>
            <input id="recipeName" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="text">
          </div>
          <div class="flex-initial p-2">
            <label class="text-gray-700" for="recipeDescription">Description</label>
            <input id="recipeDescription" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="email">
          </div>
        </div>

        <hr class="mb-2 mt-2" />

        <div v-for="(step, stepIndex) in recipe.steps" v-bind:key="stepIndex">
          <div class="flex rounded-md">
            <div class="flex-initial p-2">
              <label class="text-gray-700" :for="`preparation${stepIndex}`">Preparation</label>
              <input
                  :id="`preparation${stepIndex}`"
                  type="text"
                  :list="validPreparationSuggestionListID(stepIndex)"
                  @keyup="queryForValidPreparation(stepIndex)"
                  v-model="preparationQueries[stepIndex]"
                  class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500"
              >
              <ul :id="validPreparationSuggestionListID(stepIndex)">
                <div v-for="suggestion in preparationSuggestions[stepIndex]" v-bind:key="suggestion.id" @click="selectValidPreparationSuggestion(stepIndex, suggestion)">{{ suggestion.name }}</div>
              </ul>
            </div>
            <div class="flex-initial p-2">
              <label class="text-gray-700" :for="`why${stepIndex}`">Why</label>
              <input :id="`why${stepIndex}`" v-model="step.why" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="text">
            </div>
            <div class="flex-initial p-2">
              <label class="text-gray-700" :for="`description${stepIndex}`">Description</label>
              <input :id="`description${stepIndex}`" v-model="step.notes" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="email">
            </div>
          </div>

          <div v-for="(ingredient, ingredientIndex) in step.ingredients" v-bind:key="ingredient.ingredientID">
            <div class="flex rounded-md">
              <div class="flex-initial p-2">
                <label class="text-gray-700" :for="`ingredientName${ingredientIndex}`">Name</label>
                <input
                    :id="`ingredientName${ingredientIndex}`"
                    class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500"
                    type="text"
                    :list="`recipeStep_${stepIndex}_ValidIngredient_${ingredientIndex}_Suggestions`"
                    @keyup="queryForValidIngredient(stepIndex, ingredientIndex)"
                    v-model="ingredientNameQueries[stepIndex][ingredientIndex]"
                >
                <ul :id="`recipeStep_${stepIndex}_ValidIngredient_${ingredientIndex}_Suggestions`">
                  <div v-for="suggestion in validIngredientSuggestions[stepIndex][ingredientIndex]" v-bind:key="suggestion.id" @click="selectValidIngredientSuggestion(stepIndex, ingredientIndex, suggestion)">{{ suggestion.name }}, {{ suggestion.variant }}</div>
                </ul>
              </div>

              <div class="flex-initial p-2">
                <label class="text-gray-700" :for="`quantityValue${ingredientIndex}`">Quantity</label>
                <input
                    :id="`quantityValue${ingredientIndex}`"
                    class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500"
                    type="number"
                    min="0"
                    v-model.number="ingredient.quantityValue"
                >
              </div>

              <div class="flex-initial p-2">
                <span class="text-gray-700" :for="`quantityType${ingredientIndex}`">Type</span>
                <select :id="`quantityType${ingredientIndex}`" v-model="ingredient.quantityType" class="form-select mt-2 block w-full rounded-md">
                  <option selected>grams</option>
                  <option>fl. oz.</option>
                </select>
              </div>

              <div class="flex-initial p-2">
                <span class="" :for="`deleteButton${ingredientIndex}`">&nbsp;</span>
                <button @click="removeIngredient(stepIndex, ingredientIndex)" :id="`deleteButton${ingredientIndex}`" class="mt-8 px-4 py-2 text-gray-200 bg-white-500 rounded-md hover:bg-white-700 focus:outline-none focus:bg-white-700" >🗑️</button>
              </div>
            </div>
          </div>
          <button @click="removeStep(stepIndex)" class="px-4 py-2 m-2 text-sm text-center text-white bg-red-600 rounded-md focus:outline-none hover:bg-red-500">Remove Step</button>
          <button @click="addIngredient(stepIndex)" class="px-4 py-2 mt-2 mb-2 text-sm text-center text-white bg-indigo-600 rounded-md focus:outline-none hover:bg-indigo-500">Add Ingredient</button>

          <hr class="mb-2 mt-2" v-if="stepIndex !== recipe.steps.length-1" />

        </div>

        <div class="flex justify-end mt-4">
          <button @click="addStep" class="px-4 py-2 m-2 text-sm text-center text-white bg-indigo-600 rounded-md focus:outline-none hover:bg-indigo-500">Add Step</button>
        </div>
        <div class="flex justify-end mt-4">
          <button class="px-4 py-2 text-gray-200 bg-gray-800 rounded-md hover:bg-gray-700 focus:outline-none focus:bg-gray-700" :disabled="!readyForSubmission" @click="submitRecipe"> Save </button>
        </div>
    </div>
  </div>
</template>

<script lang="ts">
import axios, {AxiosError, AxiosResponse} from "axios";
import { defineComponent } from "vue";
import { useRouter } from "vue-router";

import { Recipe, RecipeStep, RecipeStepIngredient } from "../../../models";
import {backendRoutes} from "../../../constants";

class SearchSuggestion {
    name: string;
    variant: string;
    id: number;

    constructor() {
        this.name = "";
        this.variant = "";
        this.id = 0;
    }
}

function initializeRecipe(): Recipe {
  const r = new Recipe();

  r.steps = [new RecipeStep()]
  r.steps[0].ingredients = [new RecipeStepIngredient()]
  r.steps[0].ingredients[0].quantityType = "grams"

  return r;
}

export default defineComponent({
    data() {
        return {
            recipe: initializeRecipe(),
            preparationQueries: new Array<string>(""),
            ingredientNameQueries: new Array<Array<string>>([]),
            preparationSuggestions: new Array<Array<SearchSuggestion>>(),
            validIngredientSuggestions: new Array<Array<Array<SearchSuggestion>>>([]),
        }
    },
    methods: {
      addStep() {
        this.recipe.steps.push(new RecipeStep());
        this.ingredientNameQueries.push([]);
        this.validIngredientSuggestions.push(new Array<Array<SearchSuggestion>>())
        console.log(`Step added!`);
      },
      removeStep(stepIndex: number) {
        this.recipe.steps.splice(stepIndex, 1);
      },
      readyForSubmission(): boolean {
        return this.recipe.steps.length > 0;
      },
      addIngredient(stepIndex: number) {
        let ingredient = new RecipeStepIngredient();
        this.recipe.steps[stepIndex].ingredients.push(ingredient);
        this.ingredientNameQueries.push([]);
        this.preparationQueries[stepIndex] = "";
        this.validIngredientSuggestions.push(new Array<Array<SearchSuggestion>>())
        console.log(`Ingredient added!`);
      },
      removeIngredient(stepIndex: number, ingredientIndex: number) {
        this.recipe.steps[stepIndex].ingredients.splice(ingredientIndex, 1);
      },
      validPreparationSuggestionListID(stepIndex: number) {
          return `recipeStep_${stepIndex}_ValidPreparationSuggestions`;
      },
      queryForValidPreparation(stepIndex: number) {
          let query = this.preparationQueries[stepIndex].trim();

          const searchURL = `${backendRoutes.VALID_PREPARATIONS_SEARCH}?q=${encodeURIComponent(query)}`;

          if (query.length > 1) {
            axios.get(searchURL)
                .then((res: AxiosResponse<SearchSuggestion[]>) => {
                  this.preparationSuggestions[stepIndex] = res.data || [];
                })
                .catch((err: AxiosError) => {
                  console.error(err);
                });
          } else {
              this.preparationSuggestions[stepIndex] = [];
          }
      },
      selectValidPreparationSuggestion(stepIndex: number, selection: SearchSuggestion) {
        let step = this.recipe.steps[stepIndex];
        step.preparationID = selection.id;

        this.preparationQueries[stepIndex] = `${selection.name}${selection.variant ? ' - ' : ''}${selection.variant || ''}`;
        this.preparationSuggestions[stepIndex] = [];
      },
      queryForValidIngredient(stepIndex: number, ingredientIndex: number) {
        let preparationID = this.recipe.steps[stepIndex].preparationID
        let query = this.ingredientNameQueries[stepIndex][ingredientIndex].trim();

        if (query.length > 1 && preparationID !== 0) {
          const searchURL = `${backendRoutes.VALID_INGREDIENTS_SEARCH}?q=${encodeURIComponent(query)}&pid=${preparationID}`;
          axios.get(searchURL)
              .then((res: AxiosResponse<SearchSuggestion[]>) => {
                this.validIngredientSuggestions[stepIndex][ingredientIndex] = res.data || [];
              })
              .catch((err: AxiosError) => {
                console.error(err);
              });
        }
      },
      selectValidIngredientSuggestion(stepIndex: number, ingredientIndex: number, selection: SearchSuggestion) {
          let ingredient = this.recipe.steps[stepIndex].ingredients[ingredientIndex];

          ingredient.ingredientID = selection.id;
          this.ingredientNameQueries[stepIndex][ingredientIndex] = `${selection.name}${selection.variant ? ' - ' : ''}${selection.variant || ''}`;
          this.validIngredientSuggestions[stepIndex][ingredientIndex] = [];
      },
      submitRecipe() {
        axios.post(backendRoutes.RECIPES, this.recipe)
        .then((res: AxiosResponse<Recipe>) => {
          console.dir(res.data);
        })
        .catch((err: AxiosError) =>{
          console.error(err);
        });
      }
    },
    setup() {
        const router = useRouter();

        return {
            router,
        };
    },
});
</script>