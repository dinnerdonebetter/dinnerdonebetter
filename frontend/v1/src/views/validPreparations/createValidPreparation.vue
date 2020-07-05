<template>
  <div class="app-container">
    <button
      v-if="!editing"
      @click="toggleEditing"
    >
      edit
    </button>
    <button
      v-else
      @click="toggleEditing"
    >
      stop editing
    </button>
    <button
      :disabled="!hasChanged"
      @click="savePreparation"
    >
      save
    </button>
    <p>id: {{ currentPreparation.id }}</p>
    <p>
      name:
      <input
        v-model="currentPreparation.name"
        type="text"
        :disabled="!editing"
        @keyup="fieldChange"
      >
    </p>
    <p>
      description:
      <input
        v-model="currentPreparation.description"
        type="text"
        :disabled="!editing"
        @keyup="fieldChange"
      >
    </p>
    <p>
      icon:
      <input
        v-model="currentPreparation.icon"
        type="text"
        :disabled="!editing"
        @keyup="fieldChange"
      >
    </p>
    <p>createdOn: {{ renderUnixTime(currentPreparation.createdOn) }}</p>
    <p>updatedOn: {{ renderUnixTime(currentPreparation.updatedOn) }}</p>
    <p>archivedOn: {{ renderUnixTime(currentPreparation.archivedOn) }}</p>
  </div>
</template>

<script lang="ts">
import faker from 'faker';
import axios from 'axios';
import { Component, Vue } from 'vue-property-decorator';

import { backendRoutes, statusCodes } from '@/constants';
import { ValidPreparation } from '@/models';
import { renderUnixTime } from '@/utils/time';
import {AppModule} from "@/store/modules/app";

@Component({
  name: 'ValidPreparationCreationComponent',
})
export default class ValidPreparationCreationComponent extends Vue {
  private hasChanged = false;
  private editing = true;

  private currentPreparation: ValidPreparation = new ValidPreparation();

  private savePreparation(): void {
    this.editing = false;

    if (AppModule.frontendDevMode) {
      this.$router.push({
        path: `/admin/enumerations/valid_preparations/${faker.random.number()}`,
      });
    } else {
      axios.post(backendRoutes.VALID_PREPARATIONS, this.currentPreparation)
        .then((response) => {
          if (response.status === statusCodes.CREATED) {
            const vp = response.data as ValidPreparation;

            this.$router.push({
              path: `/admin/enumerations/valid_preparations/${vp.id}`,
            });
          } else if (response.status === statusCodes.UNAUTHORIZED) {
            console.dir(response);
          }
        });
    }
  }

  private toggleEditing(): void {
    this.editing = !this.editing;
  }

  private renderUnixTime = renderUnixTime;

  private fieldChange(): void {
    this.hasChanged = this.currentPreparation.name !== '';
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  h3 {
    margin: 40px 0 0;
  }
  ul {
    list-style-type: none;
    padding: 0;
  }
  li {
    display: inline-block;
    margin: 0 10px;
  }
  a {
    color: #42b983;
  }
</style>
