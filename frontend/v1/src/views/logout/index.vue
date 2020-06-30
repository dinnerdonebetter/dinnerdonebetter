<template>
  <div />
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

@Component({name: 'Login'})
export default class extends Vue {
  mounted() {
    this.$store.dispatch("user/Logout").then(() => {
      this.$router.push({ path: '/login' });
    });
  }
}
</script>

<style lang="scss" scoped>
</style>
