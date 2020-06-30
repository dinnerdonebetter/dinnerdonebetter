<template>
  <div class="dashboard-container">
    <div class="dashboard-text">
      <img
        alt="chef logo"
        class="chef"
        src="@/assets/prixfixeheroiconschef.svg"
      >
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

@Component({
  name: 'Dashboard',
})
export default class extends Vue {
}
</script>

<style lang="scss" scoped>
.dashboard {
  &-container {
    margin: 30px;
  }

  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
