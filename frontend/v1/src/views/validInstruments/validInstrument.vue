<template>
  <div class="app-container">
    <div v-if="currentInstrument !== null">
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
      <p>id: {{ currentInstrument.id }}</p>
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
        warning:
        <input
          v-model="currentInstrument.warning"
          type="text"
          :disabled="!editing"
          @keyup="fieldChange"
        >
      </p>
      <p>
        containsEgg:
        <input
          v-model="currentInstrument.containsEgg"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
        >
      </p>
      <p>
        containsDairy:
        <input
          v-model="currentInstrument.containsDairy"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
        >
      </p>
      <p>
        containsPeanut:
        <input
          v-model="currentInstrument.containsPeanut"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
        >
      </p>
      <p>
        containsTreeNut:
        <input
          v-model="currentInstrument.containsTreeNut"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
        >
      </p>
      <p>
        containsSoy:
        <input
          v-model="currentInstrument.containsSoy"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
        >
      </p>
      <p>
        containsWheat:
        <input
          v-model="currentInstrument.containsWheat"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
        >
      </p>
      <p>
        containsShellfish:
        <input
          v-model="currentInstrument.containsShellfish"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
        >
      </p>
      <p>
        containsSesame:
        <input
          v-model="currentInstrument.containsSesame"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
        >
      </p>
      <p>
        containsFish:
        <input
          v-model="currentInstrument.containsFish"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
        >
      </p>
      <p>
        containsGluten:
        <input
          v-model="currentInstrument.containsGluten"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
        >
      </p>
      <p>
        animalFlesh:
        <input
          v-model="currentInstrument.animalFlesh"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
        >
      </p>
      <p>
        animalDerived:
        <input
          v-model="currentInstrument.animalDerived"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
        >
      </p>
      <p>
        consideredStaple:
        <input
          v-model="currentInstrument.consideredStaple"
          type="checkbox"
          :disabled="!editing"
          @change="fieldChange"
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
      <hr>
      <el-button
        type="danger"
        @click="deleteInstrument"
      >
        Delete
      </el-button>
    </div>
    <div v-else>
      <p>invalid instrument ID!</p>
    </div>
  </div>
</template>

<script lang="ts">
import axios, { AxiosResponse } from 'axios';
import { Component, Vue } from 'vue-property-decorator';
import format from 'string-format';

import { backendRoutes, statusCodes } from '@/constants';
import {fakeValidInstrumentFactory, ValidInstrument, validInstrumentsAreEqual} from '@/models';
import { renderUnixTime } from '@/utils/time';
import {AppModule} from "@/store/modules/app";

@Component({
  name: 'ValidInstrumentComponent',
})
export default class ValidInstrumentComponent extends Vue {
  private talkedToServer = false;

  private hasChanged = false;
  private editing = false;

  private currentInstrument: ValidInstrument = new ValidInstrument();
  private originalInstrument: ValidInstrument | null = null;

  private mounted(): void {
    if (AppModule.frontendDevMode) {
      this.currentInstrument = fakeValidInstrumentFactory.build();
    } else {
      axios.get(this.buildURL())
        .then((response: AxiosResponse) => {
          this.talkedToServer = true;
          if (response.status === statusCodes.OK) {
            return response.data;
          } else {
            throw "no response from server";
          }
        })
        .then((data: ValidInstrument) => {
          this.currentInstrument = data;
          if (this.originalInstrument === null) {
            this.originalInstrument = {
              ...this.currentInstrument,
            } as ValidInstrument;

            document.title = `${this.currentInstrument.variant} - ${this.currentInstrument.name}`;
          }
        });
    }
  }

  private buildURL(): string {
    return format(
      backendRoutes.VALID_INSTRUMENT,
      this.$route.params.validInstrumentID,
    );
  }

  private saveInstrument(): void {
    if (AppModule.frontendDevMode) {
      this.editing = true;
    } else {
      axios.put(this.buildURL(), this.currentInstrument)
        .then((response: AxiosResponse) => {
          this.talkedToServer = true;
          if (response.status === statusCodes.OK) {
            return response.data;
          } else {
            throw "no response from server";
          }
        })
        .then((data: ValidInstrument) => {
          this.currentInstrument = data;
          this.editing = false;
        });
    }
  }

  private deleteInstrument() {
    const nameAndVariant = `${this.currentInstrument.name} - ${this.currentInstrument.variant}`;
    const confirmationText = prompt(`please enter the full name to confirm deletion: ${nameAndVariant}`);

    if (nameAndVariant === confirmationText) {
      if (AppModule.frontendDevMode) {
        this.$router.push({path: "/admin/enumerations/valid_instruments/"});
      } else {
        axios.delete(this.buildURL())
          .then((response: AxiosResponse) => {
            if (response.status === statusCodes.NO_CONTENT) {
              this.$router.push({path: "/admin/enumerations/valid_instruments/"});
            } else {
              console.error("something has gone awry");
            }
          });
      }
    }
  }

  private toggleEditing(): void {
    this.editing = !this.editing;
  }

  private renderUnixTime = renderUnixTime;

  private fieldChange(): void {
    if (this.originalInstrument !== null) {
      this.hasChanged = !validInstrumentsAreEqual(
        this.currentInstrument,
        this.originalInstrument,
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
