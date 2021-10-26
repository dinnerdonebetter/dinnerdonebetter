<template>
  <DataTable
      title="Users"
      :headers=headers
      :row-data="tableData"
      :previous-page-disabled="currentPage <= 1"
      :next-page-disabled="lastPage"
      :new-button-disabled="true"
      @searchQueried="searchForUsers"
      @previousPageRequested="fetchPreviousPage"
      @nextPageRequested="fetchNextPage"
  />
</template>

<script lang="ts">
import axios, {AxiosError, AxiosResponse} from "axios";
import { defineComponent } from "vue";

import DataTable from "../../../components/admin/DataTable.vue";
import {backendRoutes} from "../../../constants";
import {QueryFilter, User, UserList} from "../../../models";
import {settings} from "../../../settings/settings";

function filterUserFields(input: User): string[] {
  return [
      input.id.toString(),
      input.username,
      input.serviceRoles.join(", "),
      input.reputation,
      input.createdOn.toString(),
  ]
}

export default defineComponent({
  beforeMount() {
    this.fetchUsersFromAPI();
  },
  data() {
    return {
      users: new Array<User>(),
      headers: [
        "ID",
        "Username",
        "Service Role",
        "Reputation",
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
    headerClicked(name: string): void {
      console.log('headerClicked', name);
    },
    rowClicked(rowIndex: number): void {
      console.log('rowClicked', rowIndex);
    },
    fetchPreviousPage(): void {
      if (this.currentPage === 0) {
        return;
      }
      this.currentPage -= 1;
      this.fetchUsersFromAPI();
    },
    fetchNextPage(): void {
      if (this.lastPage) {
        return;
      }
      this.currentPage += 1;
      this.fetchUsersFromAPI();
    },
    fetchUsersFromAPI(): void {
      this.loading = true;
      const u = new URL(
          `${settings.API_SERVER_URL}${backendRoutes.USERS}${location.search}`,
      );
      const qf = new QueryFilter(u.searchParams);
      qf.page = this.currentPage;
      qf.modifyURL(u);

      axios.get(u.toString())
          .then((res: AxiosResponse<UserList>) => {
            this.users = res.data?.users || [];
            this.totalCount = res.data?.totalCount;
            this.currentPage = res.data?.page;

            this.tableData = this.users.map(filterUserFields);

            this.lastPage = this.users.length < res.data?.limit;
          })
          .catch((err: AxiosError) => {
            console.error(err)
          })
          .finally(() => {
            this.loading = false;
          });
    },
    searchForUsers(searchQuery: string) {
      if (searchQuery.length < 2) {
        return
      }

      const u = new URL(
          `${settings.API_SERVER_URL}${backendRoutes.USERS_SEARCH}?q=${encodeURIComponent(searchQuery)}`,
      );

      axios.get(u.toString())
          .then((res: AxiosResponse<User[]>) => {
            this.users = res.data || [];
            this.tableData = this.users.map(filterUserFields);
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

