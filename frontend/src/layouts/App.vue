<template>
  <Header />
  <main>
    <slot />
  </main>
</template>

<script lang="ts">
import { defineComponent } from "vue";

import Header from "../components/app/Header.vue";

export default defineComponent({
  components: {
    Header,
  },
});
</script>
