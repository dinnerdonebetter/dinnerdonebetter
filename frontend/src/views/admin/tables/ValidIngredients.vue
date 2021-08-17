<template>
  <DataTable
      title="Valid Ingredients"
      :headers=headers
      :row-data="tableData"
      :previous-page-disabled="currentPage <= 1"
      :next-page-disabled="lastPage"
      @rowClicked="navigateToIngredient"
      @newButtonClicked="navigateToCreationPage"
      @searchQueried="searchForValidIngredients"
      @previousPageRequested="fetchPreviousPage"
      @nextPageRequested="fetchNextPage"
  />
</template>

<script lang="ts">
import { defineComponent } from "vue";

import DataTable from "../../../components/admin/DataTable.vue";
import {backendRoutes} from "../../../constants";
import {QueryFilter, ValidIngredient, ValidIngredientList} from "../../../models";
import axios, {AxiosError, AxiosResponse} from "axios";

function filterIngredientFields(input: ValidIngredient): string[] {
  return [
      input.id.toString(),
      input.name,
      input.variant,
      input.description,
      input.warning,
      input.createdOn.toString(),
  ]
}

export default defineComponent({
  beforeMount() {
    this.fetchValidIngredientsFromAPI();
  },
  data() {
    return {
      validIngredients: new Array<ValidIngredient>(),
      headers: [
        "ID",
        "Name",
        "Variant",
        "Description",
        "Warning",
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
      this.fetchValidIngredientsFromAPI();
    },
    fetchNextPage(): void {
      if (this.lastPage) {
        return;
      }
      this.currentPage += 1;
      this.fetchValidIngredientsFromAPI();
    },
    navigateToCreationPage(): void {
      this.$router.push(`/admin/valid_ingredients/new`);
    },
    navigateToIngredient(ingredientIndex: number): void {
      const ingredient = this.validIngredients[ingredientIndex];
      this.$router.push(`/admin/valid_ingredients/${ingredient.id}`);
    },
    fetchValidIngredientsFromAPI(): void {
      this.loading = true;
      const u = new URL(
          `${location.protocol}//${location.host}${backendRoutes.VALID_INGREDIENTS}${location.search}`,
      );
      const qf = new QueryFilter(u.searchParams);
      qf.page = this.currentPage;
      qf.modifyURL(u);

      axios.get(u.toString())
          .then((res: AxiosResponse<ValidIngredientList>) => {
            this.validIngredients = res.data?.validIngredients || [];
            this.totalCount = res.data?.totalCount;
            this.currentPage = res.data?.page;

            this.tableData = this.validIngredients.map(filterIngredientFields);

            this.lastPage = this.validIngredients.length < res.data?.limit;
          })
          .catch((err: AxiosError) => {
            console.error(err)
          })
          .finally(() => {
            this.loading = false;
          });
    },
    searchForValidIngredients(searchQuery: string) {
      if (searchQuery.length < 2) {
        return
      }

      const u = new URL(
          `${location.protocol}//${location.host}${backendRoutes.VALID_INGREDIENTS_SEARCH}?q=${encodeURIComponent(searchQuery)}`,
      );

      axios.get(u.toString())
          .then((res: AxiosResponse<ValidIngredient[]>) => {
            this.validIngredients = res.data || [];
            this.tableData = this.validIngredients.map(filterIngredientFields);
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

