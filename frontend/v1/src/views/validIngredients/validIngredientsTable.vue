<template>
  <div class="app-container">
    <el-table
      v-loading="loading"
      :data="validIngredients"
      element-loading-text="Loading"
      empty-text="no data"
      border
      fit
      highlight-current-row
    >
      <el-table-column
        align="center"
        label="ID"
      >
        <template slot-scope="scope">
          <router-link
            :to="{
              name: 'validIngredient',
              params: {validIngredientID: scope.row.id}
            }"
          >
            {{ scope.row.id }}
          </router-link>
        </template>
      </el-table-column>

      <el-table-column
        label="name"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.name }}
        </template>
      </el-table-column>

      <el-table-column
        label="variant"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.variant }}
        </template>
      </el-table-column>

      <el-table-column
        label="description"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.description }}
        </template>
      </el-table-column>

      <el-table-column
        label="warning"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.warning }}
        </template>
      </el-table-column>

      <el-table-column
        label="containsEgg"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.containsEgg }}
        </template>
      </el-table-column>

      <el-table-column
        label="containsDairy"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.containsDairy }}
        </template>
      </el-table-column>

      <el-table-column
        label="containsPeanut"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.containsPeanut }}
        </template>
      </el-table-column>

      <el-table-column
        label="containsTreeNut"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.containsTreeNut }}
        </template>
      </el-table-column>

      <el-table-column
        label="containsSoy"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.containsSoy }}
        </template>
      </el-table-column>

      <el-table-column
        label="containsWheat"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.containsWheat }}
        </template>
      </el-table-column>

      <el-table-column
        label="containsShellfish"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.containsShellfish }}
        </template>
      </el-table-column>

      <el-table-column
        label="containsSesame"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.containsSesame }}
        </template>
      </el-table-column>

      <el-table-column
        label="containsFish"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.containsFish }}
        </template>
      </el-table-column>

      <el-table-column
        label="containsGluten"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.containsGluten }}
        </template>
      </el-table-column>

      <el-table-column
        label="animalFlesh"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.animalFlesh }}
        </template>
      </el-table-column>

      <el-table-column
        label="animalDerived"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.animalDerived }}
        </template>
      </el-table-column>

      <el-table-column
        label="consideredStaple"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.consideredStaple }}
        </template>
      </el-table-column>

      <el-table-column
        label="icon"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.icon }}
        </template>
      </el-table-column>

      <el-table-column
        label="createdOn"
        align="center"
      >
        <template slot-scope="scope">
          {{ renderUnixTime(scope.row.createdOn) }}
        </template>
      </el-table-column>

      <el-table-column
        label="updatedOn"
        align="center"
      >
        <template slot-scope="scope">
          {{ renderUnixTime(scope.row.updatedOn) }}
        </template>
      </el-table-column>
    </el-table>

    <button @click="goBackOnePage">
      &lt;
    </button>
    {{ currentPage }}
    <button @click="goForwardOnePage">
      &gt;
    </button>
    <button
      style="float: right;"
      @click="createNewIngredient"
    >
      +
    </button>
  </div>
</template>

<script lang="ts">
import axios, { AxiosResponse } from 'axios';
import { Component, Vue } from 'vue-property-decorator';

import { backendRoutes, ContentType, statusCodes } from "@/constants";
import { QueryFilter, ValidIngredient } from "@/models";
import { renderUnixTime } from '@/utils/time';

@Component({
  name: 'ValidIngredientsTable',
  filters: {
    parseTime: (timestamp?: string): string => {
      if (timestamp) {
        return new Date(timestamp).toISOString();
      }
      return "never";
    },
  },
})

export default class extends Vue {
  private talkedToServer = false;
  private validIngredients: ValidIngredient[] = [];

  // pagination vars
  private currentPage = 1;
  private totalCount = 20;
  private loading = false;

  private created(): void {
    this.fetchData();
  }

  private fetchData(): void {
    const u = new URL(
      `${location.protocol}//${location.host}${backendRoutes.VALID_INGREDIENTS}${location.search}`,
    );
    const qf = new QueryFilter(u.searchParams);
    qf.page = this.currentPage;
    qf.modifyURL(u);

    axios.get(u.toString(), {
      headers: {
        "Content-Type": ContentType,
      },
    })
      .then((response: AxiosResponse) => {
        this.talkedToServer = true;
        if (response.status === statusCodes.OK) {
          return response.data;
        } else {
          throw "no response from server";
        }
      })
      .then((data: {
        'validIngredients': ValidIngredient[];
        'totalCount': number;
        'page': number;
      }) => {
        this.validIngredients = data["validIngredients"];
        this.totalCount = data["totalCount"];
        this.currentPage = data["page"];
      });
    this.loading = false;
  }

  private renderUnixTime = renderUnixTime;

  private createNewIngredient(): void {
    this.$router.push({
      path: '/admin/enumerations/valid_ingredients/new',
    });
  }

  private goBackOnePage(): void {
    this.currentPage -= 1;
    this.fetchData();
  }

  private goForwardOnePage(): void {
    this.currentPage += 1;
    this.fetchData();
  }
}
</script>

<style lang="scss" scoped>
  a {
    color: blue;
  }
</style>
