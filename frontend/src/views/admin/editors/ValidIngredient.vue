<template>
  <div class="mt-4">
    <div class="p-6 bg-white rounded-md shadow-md">
      <h2 class="text-lg font-semibold text-gray-700 capitalize">Valid Ingredient</h2>
      <div class="flex rounded-md">

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="name">name</label>
          <input id="name" v-model="ingredient.name" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="text">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="variant">variant</label>
          <input id="variant" v-model="ingredient.variant" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="text">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="description">description</label>
          <input id="description" v-model="ingredient.description" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="text">
        </div>

        </div>
        <div class="flex rounded-md">

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="warning">warning</label>
          <input id="warning" v-model="ingredient.warning" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="text">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="iconPath">iconPath</label>
          <input id="iconPath" v-model="ingredient.iconPath" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="text">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="containsSoy">containsSoy</label>
          <input id="containsSoy" v-model="ingredient.containsSoy" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

        </div>
      <div class="flex rounded-md">

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="containsTreeNut">containsTreeNut</label>
          <input id="containsTreeNut" v-model="ingredient.containsTreeNut" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="containsShellfish">containsShellfish</label>
          <input id="containsShellfish" v-model="ingredient.containsShellfish" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="containsSesame">containsSesame</label>
          <input id="containsSesame" v-model="ingredient.containsSesame" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

      </div>
      <div class="flex rounded-md">

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="containsFish">containsFish</label>
          <input id="containsFish" v-model="ingredient.containsFish" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="containsGluten">containsGluten</label>
          <input id="containsGluten" v-model="ingredient.containsGluten" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="animalFlesh">animalFlesh</label>
          <input id="animalFlesh" v-model="ingredient.animalFlesh" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

      </div>
      <div class="flex rounded-md">

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="animalDerived">animalDerived</label>
          <input id="animalDerived" v-model="ingredient.animalDerived" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="volumetric">volumetric</label>
          <input id="volumetric" v-model="ingredient.volumetric" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="containsPeanut">containsPeanut</label>
          <input id="containsPeanut" v-model="ingredient.containsPeanut" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="containsDairy">containsDairy</label>
          <input id="containsDairy" v-model="ingredient.containsDairy" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="containsEgg">containsEgg</label>
          <input id="containsEgg" v-model="ingredient.containsEgg" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="containsWheat">containsWheat</label>
          <input id="containsWheat" v-model="ingredient.containsWheat" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="checkbox">
        </div>

      </div>
      <div class="flex mt-4 justify-between">
        <button class="px-4 py-2 text-gray-100 rounded-md bg-red-500 hover:bg-red-600" :disabled="canDelete" @click="deleteIngredient"> Delete </button>
        <button class="float-right px-4 py-2 text-gray-100 rounded-md" :class="dataChanged ? 'bg-gray-800' : 'bg-gray-200'" :disabled="!dataChanged" @click="saveIngredient"> Save </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { useRouter } from "vue-router";
import format from "string-format";

import {ValidIngredient} from "../../../models";
import axios, {AxiosError, AxiosResponse} from "axios";
import {backendRoutes} from "../../../constants";

export default defineComponent({
    data() {
      return {
        creationMode: true,
        ingredientAPIPath: '',
        originalIngredient: new ValidIngredient(),
        ingredient: new ValidIngredient(),
      }
    },
    computed: {
      canDelete(): boolean {
        return !(!this.creationMode && this.ingredient.id !== 0);
      },
      dataChanged(): boolean {
        return this.ingredient.name !== this.originalIngredient.name ||
            this.ingredient.variant !== this.originalIngredient.variant ||
            this.ingredient.description !== this.originalIngredient.description ||
            this.ingredient.warning !== this.originalIngredient.warning ||
            this.ingredient.iconPath !== this.originalIngredient.iconPath ||
            this.ingredient.containsSoy !== this.originalIngredient.containsSoy ||
            this.ingredient.containsTreeNut !== this.originalIngredient.containsTreeNut ||
            this.ingredient.containsShellfish !== this.originalIngredient.containsShellfish ||
            this.ingredient.containsSesame !== this.originalIngredient.containsSesame ||
            this.ingredient.containsFish !== this.originalIngredient.containsFish ||
            this.ingredient.containsGluten !== this.originalIngredient.containsGluten ||
            this.ingredient.animalFlesh !== this.originalIngredient.animalFlesh ||
            this.ingredient.animalDerived !== this.originalIngredient.animalDerived ||
            this.ingredient.volumetric !== this.originalIngredient.volumetric ||
            this.ingredient.containsPeanut !== this.originalIngredient.containsPeanut ||
            this.ingredient.containsDairy !== this.originalIngredient.containsDairy ||
            this.ingredient.containsEgg !== this.originalIngredient.containsEgg ||
            this.ingredient.containsWheat !== this.originalIngredient.containsWheat;
      },
    },
    methods: {
      deleteIngredient(): void {
        axios.delete(this.ingredientAPIPath)
          .then(() => {
            this.$router.push(`/admin/valid_ingredients`);
          })
          .catch((err: AxiosError) => {
            console.error(err);
          });
      },
      saveIngredient(): void {
        const path = this.creationMode ? backendRoutes.VALID_INGREDIENTS : this.ingredientAPIPath;
        const requestPromise = this.creationMode ? axios.post(path, this.ingredient) : axios.put(path, this.ingredient)

        requestPromise
          .then((res: AxiosResponse<ValidIngredient>) => {
            this.originalIngredient = { ...res.data };
            this.ingredient = { ...res.data };

            if (this.creationMode) {
              this.$router.push(`/admin/valid_ingredients/${this.ingredient.id}`)
            }
          })
          .catch((err: AxiosError) => {
            console.error(err);
          });
      }
    },
    beforeMount: function () {
      const ingredientID = this.$route.params.ingredientID;

      if (ingredientID) {
        this.creationMode = false;
        this.ingredientAPIPath = format(backendRoutes.VALID_INGREDIENT, ingredientID.toString());

        axios.get(this.ingredientAPIPath)
            .then((res: AxiosResponse<ValidIngredient>) => {
              this.originalIngredient = { ...res.data };
              this.ingredient = { ...res.data };
            })
            .catch((err: AxiosError) => {
              console.error(err);
            });
      }
    },
    setup(props) {
      const router = useRouter();

      return {
          router,
      };
    },
});
</script>