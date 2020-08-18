<template>
  <div class="app-container">
    <div v-if="currentPreparation !== null">
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
      =
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
      <hr>
      <el-button
        type="danger"
        @click="deletePreparation"
      >
        Delete
      </el-button>
    </div>
    <div v-else>
      <p>invalid preparation ID!</p>
    </div>
  </div>
</template>

<script lang="ts">
import axios, { AxiosResponse } from 'axios';
import { Component, Vue } from 'vue-property-decorator';
import format from 'string-format';

import { backendRoutes, statusCodes } from '@/constants';
import {fakeValidPreparationFactory, ValidPreparation} from '@/models';
import { renderUnixTime } from '@/utils/time';
import {AppModule} from "@/store/modules/app";

@Component({
  name: 'ValidPreparationComponent',
})
export default class ValidPreparationComponent extends Vue {
  private talkedToServer = false;

  private hasChanged = false;
  private editing = false;

  private currentPreparation: ValidPreparation = new ValidPreparation();
  private originalPreparation: ValidPreparation | null = null;

  private mounted(): void {
    if (AppModule.frontendDevMode) {
      this.currentPreparation = fakeValidPreparationFactory.build();
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
        .then((data: ValidPreparation) => {
          this.currentPreparation = data;
          if (this.originalPreparation === null) {
            this.originalPreparation = {
              ...this.currentPreparation,
            } as ValidPreparation;

            document.title = `${this.currentPreparation.name}`;
          }
        });
    }
  }

  private buildURL(): string {
    return format(
      backendRoutes.VALID_PREPARATION,
      this.$route.params.validPreparationID,
    );
  }

  private savePreparation(): void {
    if (AppModule.frontendDevMode) {
      this.editing = false;
    } else {
      axios.put(this.buildURL(), this.currentPreparation)
        .then((response: AxiosResponse) => {
          this.talkedToServer = true;
          if (response.status === statusCodes.OK) {
            return response.data;
          } else {
            throw "no response from server";
          }
        })
        .then((data: ValidPreparation) => {
          this.currentPreparation = data;
          this.editing = false;
        });
    }
  }

  private deletePreparation() {
    const nameAndVariant = `${this.currentPreparation.name}`;
    const confirmationText = prompt(`please enter the full name to confirm deletion: ${nameAndVariant}`);

    if (nameAndVariant === confirmationText) {
      if (AppModule.frontendDevMode) {
        this.$router.push({path: "/admin/enumerations/valid_preparations/"});
      } else {
        axios.delete(this.buildURL())
          .then((response: AxiosResponse) => {
            if (response.status === statusCodes.NO_CONTENT) {
              this.$router.push({path: "/admin/enumerations/valid_preparations/"});
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
    if (this.originalPreparation !== null) {
      this.hasChanged = !ValidPreparation.areEqual(
        this.currentPreparation,
        this.originalPreparation,
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
