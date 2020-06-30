<template>
  <div class="app-container">
    <div v-if="currentIngredient !== null">
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
      <hr>
      <el-button
        type="danger"
        @click="deleteIngredient"
      >
        Delete
      </el-button>
    </div>
    <div v-else>
      <p>invalid ingredient ID!</p>
    </div>
  </div>
</template>

<script lang="ts">
import axios, { AxiosResponse } from 'axios';
import { Component, Vue } from 'vue-property-decorator';
import format from 'string-format';

import { backendRoutes, statusCodes } from '@/constants';
import { ValidIngredient, validIngredientsAreEqual } from '@/models';
import { renderUnixTime } from '@/utils/time';

@Component({
  name: 'ValidIngredientComponent',
})
export default class ValidIngredientComponent extends Vue {
  private talkedToServer = false;

  private hasChanged = false;
  private editing = false;

  private currentIngredient: ValidIngredient = new ValidIngredient();
  private originalIngredient: ValidIngredient | null = null;

  private mounted(): void {
    axios.get(this.buildURL())
      .then((response: AxiosResponse) => {
        this.talkedToServer = true;
        if (response.status === statusCodes.OK) {
          return response.data;
        } else {
          throw "no response from server";
        }
      })
      .then((data: ValidIngredient) => {
        this.currentIngredient = data;
        if (this.originalIngredient === null) {
          this.originalIngredient = {
            ...this.currentIngredient,
          } as ValidIngredient;

          document.title = `${this.currentIngredient.variant} - ${this.currentIngredient.name}`;
        }
      });
  }

  private buildURL(): string {
    return format(
      backendRoutes.VALID_INGREDIENT,
      this.$route.params.validIngredientID,
    );
  }

  private saveIngredient(): void {
    axios.put(this.buildURL(), this.currentIngredient)
      .then((response: AxiosResponse) => {
        this.talkedToServer = true;
        if (response.status === statusCodes.OK) {
          return response.data;
        } else {
          throw "no response from server";
        }
      })
      .then((data: ValidIngredient) => {
        this.currentIngredient = data;
        this.editing = false;
    });
  }

  private deleteIngredient() {
    const nameAndVariant = `${this.currentIngredient.name} - ${this.currentIngredient.variant}`;
    const confirmationText = prompt(`please enter the full name to confirm deletion: ${nameAndVariant}`);

    if (nameAndVariant === confirmationText) {
      axios.delete(this.buildURL())
        .then((response: AxiosResponse) => {
          if (response.status === statusCodes.NO_CONTENT) {
            this.$router.push({path: "/admin/enumerations/valid_ingredients/"});
          } else {
            console.error("something has gone awry");
          }
        });
    }
  }

  private toggleEditing(): void {
    this.editing = !this.editing;
  }

  private renderUnixTime = renderUnixTime;

  private fieldChange(): void {
    if (this.currentIngredient !== null && this.originalIngredient !== null) {
      this.hasChanged = !validIngredientsAreEqual(
        this.currentIngredient,
        this.originalIngredient,
      );
    }
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
