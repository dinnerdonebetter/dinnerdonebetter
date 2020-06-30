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
      @click="saveInstrument"
    >
      save
    </button>
    <p>id: {{ currentInstrument.id || 0 }}</p>
    <p>
      name:
      <input
        v-model="currentInstrument.name"
        type="text"
        :disabled="!editing"
        @keyup="fieldChange"
      >
    </p>
    <p>
      variant:
      <input
        v-model="currentInstrument.variant"
        type="text"
        :disabled="!editing"
        @keyup="fieldChange"
      >
    </p>
    <p>
      description:
      <input
        v-model="currentInstrument.description"
        type="text"
        :disabled="!editing"
        @keyup="fieldChange"
      >
    </p>
    <p>
      icon:
      <input
        v-model="currentInstrument.icon"
        type="text"
        :disabled="!editing"
        @keyup="fieldChange"
      >
    </p>
    <p>createdOn: {{ renderUnixTime(currentInstrument.createdOn) }}</p>
    <p>updatedOn: {{ renderUnixTime(currentInstrument.updatedOn) }}</p>
    <p>archivedOn: {{ renderUnixTime(currentInstrument.archivedOn) }}</p>
  </div>
</template>

<script lang="ts">
import axios from 'axios';
import { Component, Vue } from 'vue-property-decorator';

import { backendRoutes, statusCodes } from '@/constants';
import { ValidInstrument } from '@/models';
import { renderUnixTime } from '@/utils/time';

@Component({
  name: 'ValidInstrumentCreationComponent',
})
export default class ValidInstrumentCreationComponent extends Vue {
  private hasChanged = false;
  private editing = true;

  private currentInstrument: ValidInstrument = new ValidInstrument();

  private saveInstrument(): void {
    this.editing = false;

    axios.post(backendRoutes.VALID_INSTRUMENTS, this.currentInstrument)
      .then((response) => {
      if (response.status === statusCodes.CREATED) {
        const vi = response.data as ValidInstrument;

        this.$router.push({
          path: `/admin/enumerations/valid_instruments/${vi.id}`,
        });
      } else if (response.status === statusCodes.UNAUTHORIZED) {
        console.dir(response);
      }
    });
  }

  private toggleEditing(): void {
    this.editing = !this.editing;
  }

  private renderUnixTime = renderUnixTime;

  private fieldChange(): void {
    this.hasChanged = this.currentInstrument.name !== '';
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
