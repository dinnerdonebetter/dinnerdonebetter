<template>
  <DataTable
      title="Valid Instruments"
      :headers=headers
      :row-data="tableData"
      :previous-page-disabled="currentPage <= 1"
      :next-page-disabled="lastPage"
      @rowClicked="navigateToInstrument"
      @newButtonClicked="navigateToCreationPage"
      @searchQueried="searchForValidInstruments"
      @previousPageRequested="fetchPreviousPage"
      @nextPageRequested="fetchNextPage"
  />
</template>

<script lang="ts">
import { defineComponent } from "vue";

import DataTable from "../../../components/admin/DataTable.vue";
import {backendRoutes} from "../../../constants";
import {QueryFilter, ValidInstrument, ValidInstrumentList } from "../../../models";
import axios, {AxiosError, AxiosResponse} from "axios";

function filterInstrumentFields(input: ValidInstrument): string[] {
  return [
    input.id.toString(),
    input.name,
    input.variant,
    input.description,
    input.createdOn.toString(),
  ]
}

export default defineComponent({
  beforeMount() {
    this.fetchValidInstrumentsFromAPI();
  },
  data() {
    return {
      validInstruments: new Array<ValidInstrument>(),
      headers: [
        "ID",
        "Name",
        "Variant",
        "Description",
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
      this.fetchValidInstrumentsFromAPI();
    },
    fetchNextPage(): void {
      if (this.lastPage) {
        return;
      }
      this.currentPage += 1;
      this.fetchValidInstrumentsFromAPI();
    },
    navigateToCreationPage(): void {
      this.$router.push(`/admin/valid_instruments/new`);
    },
    navigateToInstrument(instrumentIndex: number): void {
      const instrument = this.validInstruments[instrumentIndex];
      this.$router.push(`/admin/valid_instruments/${instrument.id}`);
    },
    fetchValidInstrumentsFromAPI(): void {
      console.log("fetchValidInstrumentsFromAPI called")

      this.loading = true;
      const u = new URL(
          `${location.protocol}//${location.host}${backendRoutes.VALID_INSTRUMENTS}${location.search}`,
      );
      const qf = new QueryFilter(u.searchParams);
      qf.page = this.currentPage;
      qf.modifyURL(u);

      axios.get(u.toString())
          .then((res: AxiosResponse<ValidInstrumentList>) => {
            this.validInstruments = res.data?.validInstruments || [];
            this.totalCount = res.data?.totalCount;
            this.currentPage = res.data?.page;

            this.tableData = this.validInstruments.map(filterInstrumentFields);

            this.lastPage = this.validInstruments.length < res.data?.limit;
          })
          .catch((err: AxiosError) => {
            console.error(err)
          })
          .finally(() => {
            this.loading = false;
          });
    },
    searchForValidInstruments(searchQuery: string) {
      if (searchQuery.length < 2) {
        return
      }

      const u = new URL(
          `${location.protocol}//${location.host}${backendRoutes.VALID_INSTRUMENTS_SEARCH}?q=${encodeURIComponent(searchQuery)}`,
      );

      axios.get(u.toString())
          .then((res: AxiosResponse<ValidInstrument[]>) => {
            this.validInstruments = res.data;
            this.tableData = this.validInstruments.map(filterInstrumentFields);
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
