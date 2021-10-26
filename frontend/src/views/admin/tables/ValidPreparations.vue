<template>
  <DataTable
      title="Valid Preparations"
      :headers=headers
      :row-data="tableData"
      :previous-page-disabled="currentPage <= 1"
      :next-page-disabled="lastPage"
      @rowClicked="navigateToPreparation"
      @newButtonClicked="navigateToCreationPage"
      @searchQueried="searchForValidPreparations"
      @previousPageRequested="fetchPreviousPage"
      @nextPageRequested="fetchNextPage"
  />
</template>


<script lang="ts">
import { defineComponent } from "vue";

import DataTable from "../../../components/admin/DataTable.vue";
import {backendRoutes} from "../../../constants";
import {QueryFilter, ValidPreparation, ValidPreparationList } from "../../../models";
import axios, {AxiosError, AxiosResponse} from "axios";
import {settings} from "../../../settings/settings";

function filterPreparationFields(input: ValidPreparation): string[] {
  return [
    input.id.toString(),
    input.name,
    input.description,
    input.iconPath,
    input.createdOn.toString(),
  ]
}

export default defineComponent({
  beforeMount() {
    this.fetchValidPreparationsFromAPI();
  },
  data() {
    return {
      validPreparations: new Array<ValidPreparation>(),
      headers: [
        "ID",
        "Name",
        "Description",
        "Icon Path",
        "Created At",
      ],
      currentPage: 0,
      totalCount: 0,
      lastPage: false,
      talkedToServer: false,
      loading: true,
      tableData: new Array<Array<string>>(),
    }
  },

  methods: {
    fetchPreviousPage(): void {
      if (this.currentPage === 0) {
        return;
      }
      this.currentPage -= 1;
      this.fetchValidPreparationsFromAPI();
    },
    fetchNextPage(): void {
      if (this.lastPage) {
        return;
      }
      this.currentPage += 1;
      this.fetchValidPreparationsFromAPI();
    },
    navigateToCreationPage(): void {
      this.$router.push(`/admin/valid_preparations/new`);
    },
    navigateToPreparation(preparationIndex: number): void {
      const preparation = this.validPreparations[preparationIndex];
      this.$router.push(`/admin/valid_preparations/${preparation.id}`);
    },
    fetchValidPreparationsFromAPI(): void {
      console.log("fetchValidPreparationsFromAPI called")

      this.loading = true;
      const u = new URL(
          `${settings.API_SERVER_URL}${backendRoutes.USERS_SEARCH}${backendRoutes.VALID_PREPARATIONS}${location.search}`,
      );
      const qf = new QueryFilter(u.searchParams);
      qf.page = this.currentPage;
      qf.modifyURL(u);

      axios.get(u.toString())
          .then((res: AxiosResponse<ValidPreparationList>) => {
            this.validPreparations = res.data?.validPreparations || [];
            this.totalCount = res.data?.totalCount;
            this.currentPage = res.data?.page;

            this.tableData = this.validPreparations.map(filterPreparationFields);

            this.lastPage = this.validPreparations.length < res.data?.limit;
          })
          .catch((err: AxiosError) => {
            console.error(err)
          })
          .finally(() => {
            this.loading = false;
          });
    },
    searchForValidPreparations(searchQuery: string) {
      if (searchQuery.length < 2) {
        return
      }

      const u = new URL(
          `${settings.API_SERVER_URL}${backendRoutes.USERS_SEARCH}${backendRoutes.VALID_PREPARATIONS_SEARCH}?q=${encodeURIComponent(searchQuery)}`,
      );

      axios.get(u.toString())
          .then((res: AxiosResponse<ValidPreparation[]>) => {
            this.validPreparations = res.data;
            this.tableData = this.validPreparations.map(filterPreparationFields);
          })
          .catch((err: AxiosError) => {
            console.error(err)
          })
          .finally(() => {
            this.loading = false;
          });
    },
  },
  components: {
    DataTable,
  },
});
</script>
