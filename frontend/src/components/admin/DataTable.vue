<template>
  <h2 class="text-xl font-semibold leading-tight text-gray-700">{{ pageTitle }}</h2>

  <div class="flex flex-row justify-between mt-3">
    <div class="flex">
      <div class="relative">
        <select class="block w-full h-full px-4 py-2 pr-8 leading-tight text-gray-700 bg-white border border-gray-400 rounded-l appearance-none focus:outline-none focus:bg-white focus:border-gray-500">
          <option>5</option>
          <option>10</option>
          <option>20</option>
        </select>
      </div>
    </div>

    <div class="relative block mt-2 sm:mt-0 w-full">
        <span class="absolute inset-y-0 left-0 flex items-center pl-2">
          <svg viewBox="0 0 24 24" class="w-4 h-4 text-gray-500 fill-current">
            <path d="M10 4a6 6 0 100 12 6 6 0 000-12zm-8 6a8 8 0 1114.32 4.906l5.387 5.387a1 1 0 01-1.414 1.414l-5.387-5.387A8 8 0 012 10z"/>
          </svg>
        </span>

      <input
          @keyup="searchQueried"
          v-model="searchQuery"
          :disabled="searchDisabled"
          :placeholder="searchDisabled ? 'disabled' : 'Search'"
          :class="searchDisabled ? 'bg-gray' : 'bg-white'"
          class="block w-full py-2 pl-8 pr-6 text-sm text-gray-700 placeholder-gray-400 border border-b border-gray-400 rounded-l rounded-r appearance-none sm:rounded-l-none focus:bg-white focus:placeholder-gray-600 focus:text-gray-700 focus:outline-none"
      />
    </div>
    <div v-if="!newButtonDisabled">
      <button @click="newButtonClicked" class="px-4 py-2 text-sm font-semibold text-white bg-indigo-600 rounded-l hover:bg-indigo-400">
        New
      </button>
    </div>
  </div>

  <div class="px-4 py-4 -mx-4 overflow-x-auto sm:-mx-8 sm:px-8">
    <div class="inline-block min-w-full overflow-hidden rounded-lg shadow">
      <table class="min-w-full leading-normal">
        <thead>
        <tr>
          <th v-for="(header, headerIndex) in headers" @click=headerClicked(header) class="px-5 py-3 text-xs font-semibold tracking-wider text-left text-gray-600 uppercase bg-gray-100 border-b-2 border-gray-200">
            {{ header }}
          </th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(row, rowIndex) in rowData" @click="rowClicked(rowIndex)" :key="rowIndex">
          <td v-for="(cell, cellIndex) in row" :key="cellIndex" class="px-5 py-5 text-sm bg-white border-b border-gray-200">
            {{ cell }}
          </td>
        </tr>
        </tbody>
      </table>
      <div class="flex flex-col items-center px-5 py-5 bg-white border-t xs:flex-row xs:justify-between">
        <span class="text-xs text-gray-900 xs:text-sm">Showing X of XX Entries</span>
        <div class="inline-flex mt-2 xs:mt-0">
          <button :disabled="previousPageDisabled" @click="fetchPreviousPage" class="px-4 py-2 text-sm font-semibold text-gray-800 bg-gray-300 rounded-l hover:bg-gray-400">
            Prev
          </button>
          <button :disabled="nextPageDisabled" @click="fetchNextPage" class="px-4 py-2 text-sm font-semibold text-gray-800 bg-gray-300 rounded-r hover:bg-gray-400">
            Next
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue';

export default defineComponent({
  props: {
    title: {
      type: String,
      required: true,
    },
    headers: {
      type: Array as PropType<Array<string>>,
      required: true,
    },
    rowData: {
      type: Array as PropType<Array<Array<string | number | boolean>>>,
      required: true,
    },
    previousPageDisabled: {
      type: Boolean,
      required: false,
    },
    nextPageDisabled: {
      type: Boolean,
      required: false,
    },
    searchDisabled: {
      type: Boolean,
      required: false,
    },
    newButtonDisabled: {
      type: Boolean,
      required: false,
    }
  },
  setup(props) {
    props.rowData.forEach((row) => {
      if (row.length !== props.headers.length) {
        throw new Error("header/cell length mismatch");
      }
    })

    return {
      newButtonDisabled: props.newButtonDisabled,
      searchDisabled: props.searchDisabled,
      pageTitle: props.title,
    }
  },
  methods: {
    searchQueried(): void {
      this.$emit('searchQueried', this.searchQuery);
    },
    newButtonClicked(): void {
      this.$emit('newButtonClicked');
    },
    headerClicked(headerValue: string): void {
      this.$emit('headerClicked', headerValue);
    },
    rowClicked(rowIndex: number): void {
      this.$emit('rowClicked', rowIndex);
    },
    fetchPreviousPage(): void {
      this.$emit('previousPageRequested');
    },
    fetchNextPage(): void {
      this.$emit('nextPageRequested');
    },
  },
  data() {
    return {
      searchQuery: "",
    }
  },
});
</script>
