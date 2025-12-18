<template>
  <el-container style="height: 100%;">
    <!--侧边栏菜单-->
    <el-aside :width="is_active?'200':'64'">
      <LayoutAside :is_active="is_active" :route_path="route_path"/>
    </el-aside>

    <el-container>
      <!--正文头部内容-->
      <el-header>
        <!--监听子组件的变量事件-->
        <LayoutHeader :is_active.sync="is_active" :route_name="route_name"/>
      </el-header>
      <!--正文内容-->
      <el-main>
        <!-- 对应的组件内容渲染到router-view中 -->
        <!--子组件上报route信息-->
        <transition name="fade" mode="out-in">
          <router-view :route_path.sync="route_path" :route_name.sync="route_name"></router-view>
        </transition>
      </el-main>
      <el-footer>
        <div>
          Powered by DxP0rt
        </div>
      </el-footer>
    </el-container>
  </el-container>
</template>

<script>
import LayoutAside from "@/layout/LayoutAside";
import LayoutHeader from "@/layout/LayoutHeader";

export default {
  name: "Layout",
  components: {LayoutHeader, LayoutAside},
  data() {
    return {
      is_active: true,
      route_path: '/index',
      route_name: ['首页'],
    }
  },
  methods: {
  },
  watch: {
    route_path: function (val) {
      // var w = document.getElementById('layout-menu').clientWidth;
      window.console.log('is_active', val)
    },
  },
  created() {
    window.console.log('layout-route', this.$route)
  },
}
</script>

<style>
.el-header {
  background-color: #fff;
  color: #303133;
  line-height: 64px;
  height: 64px !important;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  z-index: 100;
}

.el-main {
  background-color: #f5f7fa;
  padding: 20px;
  min-height: calc(100vh - 124px);
}

.el-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  background-color: #fff;
  border-top: 1px solid #e4e7ed;
  height: 60px !important;
  font-size: 12px;
  color: #909399;
}

.el-aside {
  transition: width 0.3s;
  overflow: hidden;
}

/* 全局滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #c0c4cc;
}

::-webkit-scrollbar-track {
  background: #f5f7fa;
}

/* Element UI 组件优化 */
.el-card {
  border-radius: 8px;
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
}

.el-card:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
}

.el-table {
  border-radius: 8px;
  overflow: hidden;
}

/* 页面过渡动画 */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter, .fade-leave-to {
  opacity: 0;
}
</style>