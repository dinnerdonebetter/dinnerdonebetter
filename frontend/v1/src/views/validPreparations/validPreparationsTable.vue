<template>
  <div class="app-container">
    <el-table
      v-loading="loading"
      :data="validPreparations"
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
              name: 'validPreparation',
              params: {validPreparationID: scope.row.id}
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
        label="description"
        align="center"
      >
        <template slot-scope="scope">
          {{ scope.row.description }}
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
      @click="createNewPreparation"
    >
      +
    </button>
  </div>
</template>

<script lang="ts">
import axios, { AxiosResponse } from 'axios';
import { Component, Vue } from 'vue-property-decorator';

import { backendRoutes, ContentType, statusCodes } from "@/constants";
import {fakeValidPreparationFactory, QueryFilter, ValidPreparation} from "@/models";
import { renderUnixTime } from '@/utils/time';
import {AppModule} from "@/store/modules/app";

@Component({
  name: 'ValidPreparationsTable',
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
  private validPreparations: ValidPreparation[] = [];

  // pagination vars
  private currentPage = 1;
  private totalCount = 0;
  private perPageCount = 20;
  private loading = false;

  private created(): void {
    this.fetchData();
  }

  private fetchData(): void {
    if (AppModule.frontendDevMode) {
      this.validPreparations = fakeValidPreparationFactory.buildList(this.perPageCount);
    } else {
      const u = new URL(
        `${location.protocol}//${location.host}${backendRoutes.VALID_PREPARATIONS}${location.search}`,
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
          'validPreparations': ValidPreparation[];
          'totalCount': number;
          'page': number;
        }) => {
          this.validPreparations = data["validPreparations"];
          this.totalCount = data["totalCount"];
          this.currentPage = data["page"];
        });
      this.loading = false;
    }
  }

  private renderUnixTime = renderUnixTime;

  private createNewPreparation(): void {
    this.$router.push({
      path: '/admin/enumerations/valid_preparations/new',
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
