<template>
  <DataTable
      title="Recipes"
      :headers=headers
      :row-data="tableData"
      :previous-page-disabled="currentPage <= 1"
      :next-page-disabled="lastPage"
      @newButtonClicked="navigateToCreationPage"
      @searchQueried="searchForRecipes"
      @previousPageRequested="fetchPreviousPage"
      @nextPageRequested="fetchNextPage"
  />
</template>

<script lang="ts">
import { defineComponent } from "vue";

import DataTable from "../../../components/admin/DataTable.vue";
import {backendRoutes} from "../../../constants";
import {QueryFilter, Recipe, RecipeList} from "../../../models";
import axios, {AxiosError, AxiosResponse} from "axios";

function filterRecipeFields(input: Recipe): string[] {
  return [
    input.id.toString(),
    input.name,
    input.description,
    input.createdOn.toString(),
  ]
}

export default defineComponent({
  data() {
    return {
      recipes: new Array<Recipe>(),
      headers: [
        "ID",
        "Name",
        "Description",
        "Created At",
      ],
      currentPage: 0,
      totalCount: 0,
      searchQuery: "",
      lastPage: false,
      talkedToServer: false,
      loading: true,
      tableData: new Array<Array<string>>(),
    }
  },
  created() {
    this.fetchRecipesFromAPI();
  },
  methods: {
    fetchPreviousPage(): void {
      if (this.currentPage === 0) {
        return;
      }
      this.currentPage -= 1;
      this.fetchRecipesFromAPI();
    },
    fetchNextPage(): void {
      if (this.lastPage) {
        return;
      }
      this.currentPage += 1;
      this.fetchRecipesFromAPI();
    },
    navigateToIngredient(id: number): void {
      this.$router.push(`/admin/recipes/${id}`);
    },
    navigateToCreationPage(): void {
      this.$router.push(`/admin/recipes/new`);
    },
    fetchRecipesFromAPI(): void {
      this.loading = true;
      const u = new URL(
          `${location.protocol}//${location.host}${backendRoutes.RECIPES}${location.search}`,
      );
      const qf = new QueryFilter(u.searchParams);
      qf.page = this.currentPage;
      qf.modifyURL(u);

      axios.get(u.toString())
          .then((res: AxiosResponse<RecipeList>) => {
            this.recipes = res.data?.recipes || [];
            this.totalCount = res.data?.totalCount;
            this.currentPage = res.data?.page;

            this.tableData = this.recipes.map(filterRecipeFields);

            this.lastPage = this.recipes.length < res.data?.limit;
          })
          .catch((err: AxiosError) => {
            console.error(err)
          })
          .finally(() => {
            this.loading = false;
          });
    },
    searchForRecipes() {
      if (this.searchQuery.length < 2) {
        return
      }

      const u = new URL(
          `${location.protocol}//${location.host}${backendRoutes.RECIPES_SEARCH}?q=${encodeURIComponent(this.searchQuery)}`,
      );

      axios.get(u.toString())
          .then((res: AxiosResponse<Recipe[]>) => {
            this.recipes = res.data || [];
            this.tableData = this.recipes.map(filterRecipeFields);
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
