<template>
  <div
    v-if="!hidden"
    :class="['menu-wrapper', isCollapsed ? 'simple-mode' : 'full-mode']"
  >
    <el-submenu
      :index="resolvePath(path)"
    >
      <template slot="title">
        <svg-icon
          v-if="icon"
          :name="icon"
        />
        <span
          slot="title"
        >{{ title }}</span>
      </template>
    </el-submenu>
  </div>
</template>

<script lang="ts">
import path from 'path';
import { Component, Prop, Vue } from 'vue-property-decorator';
import { isExternal } from '@/utils/validate';
import SidebarItemLink from './SidebarItemLink.vue';

@Component({
  // Set 'name' here to prevent uglifyjs from causing recursive component not work
  // See https://medium.com/haiiro-io/element-component-name-with-vue-class-component-f3b435656561 for detail
  name: 'NewSidebarItem',
  components: {
    SidebarItemLink,
  },
})
export default class extends Vue {
  @Prop({required: true}) private path!: string;
  @Prop({required: true}) private title!: string;
  @Prop({required: false}) private icon!: string;

  @Prop({required: true}) private hidden!: boolean;
  @Prop({required: true}) private isExternal!: boolean;
  @Prop({ default: false }) private isCollapsed!: boolean;

  private resolvePath(routePath: string) {
    if (isExternal(routePath)) {
      return routePath;
    }
    if (isExternal(this.path)) {
      return this.path;
    }
    return path.resolve(this.path, routePath);
  }
}
</script>

<style lang="scss">
.el-submenu.is-active > .el-submenu__title {
  color: $subMenuActiveText !important;
}

.full-mode {
  .nest-menu .el-submenu>.el-submenu__title,
  .el-submenu .el-menu-item {
    min-width: $sideBarWidth !important;
    background-color: $subMenuBg !important;

    &:hover {
      background-color: $subMenuHover !important;
    }
  }
}

.simple-mode {
  &.first-level {
    .submenu-title-noDropdown {
      padding: 0 !important;
      position: relative;

      .el-tooltip {
        padding: 0 !important;
      }
    }

    .el-submenu {
      overflow: hidden;

      &>.el-submenu__title {
        padding: 0px !important;

        .el-submenu__icon-arrow {
          display: none;
        }

        &>span {
          visibility: hidden;
        }
      }
    }
  }
}
</style>

<style lang="scss" scoped>
.svg-icon {
  margin-right: 16px;
}

.simple-mode {
  .svg-icon {
    margin-left: 20px;
  }
}
</style>
