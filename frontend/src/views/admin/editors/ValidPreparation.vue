<template>
  <div class="mt-4">
    <div class="p-6 bg-white rounded-md shadow-md">
      <h2 class="text-lg font-semibold text-gray-700 capitalize">Valid Preparation</h2>
      <div class="flex rounded-md">

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="name">name</label>
          <input id="name" v-model="preparation.name" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="email">
        </div>


        <div class="flex-initial p-2">
          <label class="text-gray-700" for="description">description</label>
          <input id="description" v-model="preparation.description" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="email">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="iconPath">iconPath</label>
          <input id="iconPath" v-model="preparation.iconPath" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="email">
        </div>

      </div>

      <div class="flex mt-4 justify-between">
        <button class="px-4 py-2 text-gray-100 rounded-md bg-red-500 hover:bg-red-600" :disabled="canDelete" @click="deletePreparation"> Delete </button>
        <button class="float-right px-4 py-2 text-gray-100 rounded-md" :class="dataChanged ? 'bg-gray-800' : 'bg-gray-200'" :disabled="!dataChanged" @click="savePreparation"> Save </button>
      </div>

    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { useRouter } from "vue-router";
import format from "string-format";

import {ValidPreparation} from "../../../models";
import axios, {AxiosError, AxiosResponse} from "axios";
import {backendRoutes} from "../../../constants";

export default defineComponent({
  data() {
    return {
      creationMode: true,
      preparationAPIPath: '',
      originalPreparation: new ValidPreparation(),
      preparation: new ValidPreparation(),
    }
  },
  computed: {
    canDelete(): boolean {
      return !(!this.creationMode && this.preparation.id !== 0);
    },
    dataChanged(): boolean {
      return this.preparation.name !== this.originalPreparation.name ||
          this.preparation.description !== this.originalPreparation.description ||
          this.preparation.iconPath !== this.originalPreparation.iconPath;
    },
  },
  methods: {
    deletePreparation(): void {
      axios.delete(this.preparationAPIPath)
          .then(() => {
            this.$router.push(`/admin/valid_preparations`);
          })
          .catch((err: AxiosError) => {
            console.error(err);
          });
    },
    savePreparation(): void {
      const path = this.creationMode ? backendRoutes.VALID_PREPARATIONS : this.preparationAPIPath;
      const requestPromise = this.creationMode ? axios.post(path, this.preparation) : axios.put(path, this.preparation)

      requestPromise
          .then((res: AxiosResponse<ValidPreparation>) => {
            this.originalPreparation = { ...res.data };
            this.preparation = { ...res.data };

            if (this.creationMode) {
              this.$router.push(`/admin/valid_preparations/${this.preparation.id}`)
            }
          })
          .catch((err: AxiosError) => {
            console.error(err);
          });
    }
  },
  beforeMount: function () {
    const preparationID = this.$route.params.preparationID;

    if (preparationID) {
      this.creationMode = false;
      this.preparationAPIPath = format(backendRoutes.VALID_PREPARATION, preparationID.toString());

      axios.get(this.preparationAPIPath)
          .then((res: AxiosResponse<ValidPreparation>) => {
            this.originalPreparation = { ...res.data };
            this.preparation = { ...res.data };
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