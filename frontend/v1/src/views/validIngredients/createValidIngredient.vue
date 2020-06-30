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
      @click="saveIngredient"
    >
      save
    </button>
    <p>id: {{ currentIngredient.id }}</p>
    <p>
      name:
      <input
        v-model="currentIngredient.name"
        type="text"
        :disabled="!editing"
        @keyup="fieldChange"
      >
    </p>
    <p>
      variant:
      <input
        v-model="currentIngredient.variant"
        type="text"
        :disabled="!editing"
        @keyup="fieldChange"
      >
    </p>
    <p>
      description:
      <input
        v-model="currentIngredient.description"
        type="text"
        :disabled="!editing"
        @keyup="fieldChange"
      >
    </p>
    <p>
      warning:
      <input
        v-model="currentIngredient.warning"
        type="text"
        :disabled="!editing"
        @keyup="fieldChange"
      >
    </p>
    <p>
      containsEgg:
      <input
        v-model="currentIngredient.containsEgg"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      containsDairy:
      <input
        v-model="currentIngredient.containsDairy"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      containsPeanut:
      <input
        v-model="currentIngredient.containsPeanut"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      containsTreeNut:
      <input
        v-model="currentIngredient.containsTreeNut"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      containsSoy:
      <input
        v-model="currentIngredient.containsSoy"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      containsWheat:
      <input
        v-model="currentIngredient.containsWheat"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      containsShellfish:
      <input
        v-model="currentIngredient.containsShellfish"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      containsSesame:
      <input
        v-model="currentIngredient.containsSesame"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      containsFish:
      <input
        v-model="currentIngredient.containsFish"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      containsGluten:
      <input
        v-model="currentIngredient.containsGluten"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      animalFlesh:
      <input
        v-model="currentIngredient.animalFlesh"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      animalDerived:
      <input
        v-model="currentIngredient.animalDerived"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      consideredStaple:
      <input
        v-model="currentIngredient.measurableByVolume"
        type="checkbox"
        :disabled="!editing"
        @change="fieldChange"
      >
    </p>
    <p>
      icon:
      <input
        v-model="currentIngredient.icon"
        type="text"
        :disabled="!editing"
        @keyup="fieldChange"
      >
    </p>
    <p>createdOn: {{ renderUnixTime(currentIngredient.createdOn) }}</p>
    <p>updatedOn: {{ renderUnixTime(currentIngredient.updatedOn) }}</p>
    <p>archivedOn: {{ renderUnixTime(currentIngredient.archivedOn) }}</p>
  </div>
</template>

<script lang="ts">
import axios from 'axios';
import { Component, Vue } from 'vue-property-decorator';

import { backendRoutes, statusCodes } from '@/constants';
import { ValidIngredient } from '@/models';
import { renderUnixTime } from '@/utils/time';

@Component({
  name: 'ValidIngredientCreationComponent',
})
export default class ValidIngredientCreationComponent extends Vue {
  private hasChanged = false;
  private editing = true;

  private currentIngredient: ValidIngredient = new ValidIngredient();

  private saveIngredient(): void {
    this.editing = false;

    axios.post(backendRoutes.VALID_INGREDIENTS, this.currentIngredient)
      .then((response) => {
      if (response.status === statusCodes.CREATED) {
        const vi = response.data as ValidIngredient;

        this.$router.push({
          path: `/admin/enumerations/valid_ingredients/${vi.id}`,
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
    this.hasChanged = this.currentIngredient.name !== '';
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
