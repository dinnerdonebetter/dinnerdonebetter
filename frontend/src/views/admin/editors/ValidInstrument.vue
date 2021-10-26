<template>
  <div class="mt-4">
    <div class="p-6 bg-white rounded-md shadow-md">
      <h2 class="text-lg font-semibold text-gray-700 capitalize">Valid Instrument</h2>
      <div class="flex rounded-md">

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="name">name</label>
          <input id="name" v-model="instrument.name" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="email">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="variant">variant</label>
          <input id="variant" v-model="instrument.variant" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="email">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="description">description</label>
          <input id="description" v-model="instrument.description" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="email">
        </div>

        <div class="flex-initial p-2">
          <label class="text-gray-700" for="description">iconPath</label>
          <input id="iconPath" v-model="instrument.iconPath" class="w-full mt-2 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" type="email">
        </div>

        </div>
      <div class="flex mt-4 justify-between">
        <button class="px-4 py-2 text-gray-100 rounded-md bg-red-500 hover:bg-red-600" :disabled="canDelete" @click="deleteInstrument"> Delete </button>
        <button class="float-right px-4 py-2 text-gray-100 rounded-md" :class="dataChanged ? 'bg-gray-800' : 'bg-gray-200'" :disabled="!dataChanged" @click="saveInstrument"> Save </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { useRouter } from "vue-router";
import format from "string-format";

import {ValidInstrument} from "../../../models";
import axios, {AxiosError, AxiosResponse} from "axios";
import {backendRoutes} from "../../../constants";
import {settings} from "../../../settings/settings";

export default defineComponent({
  data() {
    return {
      creationMode: true,
      instrumentAPIPath: '',
      originalInstrument: new ValidInstrument(),
      instrument: new ValidInstrument(),
    }
  },
  computed: {
    canDelete(): boolean {
      return !(!this.creationMode && this.instrument.id !== 0);
    },
    dataChanged(): boolean {
      return this.instrument.name !== this.originalInstrument.name ||
          this.instrument.variant !== this.originalInstrument.variant ||
          this.instrument.description !== this.originalInstrument.description ||
          this.instrument.iconPath !== this.originalInstrument.iconPath;
    },
  },
  methods: {
    deleteInstrument(): void {
      axios.delete(this.instrumentAPIPath)
          .then(() => {
            this.$router.push(`/admin/valid_instruments`);
          })
          .catch((err: AxiosError) => {
            console.error(err);
          });
    },
    saveInstrument(): void {
      const path = this.creationMode ? `${settings.API_SERVER_URL}${backendRoutes.VALID_INSTRUMENTS}` : this.instrumentAPIPath;
      const requestPromise = this.creationMode ? axios.post(path, this.instrument) : axios.put(path, this.instrument)

      requestPromise
          .then((res: AxiosResponse<ValidInstrument>) => {
            this.originalInstrument = { ...res.data };
            this.instrument = { ...res.data };

            if (this.creationMode) {
              this.$router.push(`/admin/valid_instruments/${this.instrument.id}`)
            }
          })
          .catch((err: AxiosError) => {
            console.error(err);
          });
    }
  },
  beforeMount: function () {
    const instrumentID = this.$route.params.instrumentID;

    if (instrumentID) {
      this.creationMode = false;
      this.instrumentAPIPath = format(`${settings.API_SERVER_URL}${backendRoutes.VALID_INSTRUMENT}`, instrumentID.toString());

      axios.get(this.instrumentAPIPath)
          .then((res: AxiosResponse<ValidInstrument>) => {
            this.originalInstrument = { ...res.data };
            this.instrument = { ...res.data };
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