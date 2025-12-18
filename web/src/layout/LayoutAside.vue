<template>
  <div class="sidebar-container">
    <div class="sidebar-logo">
      <transition name="sidebar-logo-fade">
        <div v-if="is_active" class="sidebar-logo-content">
          <div class="logo-icon">
            <i class="el-icon-connection"></i>
          </div>
          <div class="logo-text">
            <div>安全访问</div>
            <div>管理系统</div>
          </div>
        </div>
        <div v-else class="sidebar-logo-mini">
          <i class="el-icon-connection"></i>
        </div>
      </transition>
    </div>

    <el-menu 
      :collapse="!is_active" 
      :default-active="route_path" 
      :style="is_active ? 'width:200px' : ''" 
      router
      class="layout-menu" 
      :collapse-transition="false" 
      background-color="transparent" 
      text-color="#fff"
      active-text-color="#fff">
      
      <el-menu-item index="/admin/home">
        <i class="el-icon-s-home"></i>
        <span slot="title">首页</span>
      </el-menu-item>

      <el-submenu index="1">
        <template slot="title">
          <i class="el-icon-setting"></i>
          <span slot="title">基础信息</span>
        </template>
        <el-menu-item index="/admin/set/system">系统信息</el-menu-item>
        <el-menu-item index="/admin/set/soft">软件配置</el-menu-item>
        <el-menu-item index="/admin/set/other">其他设置</el-menu-item>
        <el-menu-item index="/admin/set/audit">审计日志</el-menu-item>
      </el-submenu>

      <el-submenu index="2">
        <template slot="title">
          <i class="el-icon-user"></i>
          <span slot="title">用户信息</span>
        </template>
        <el-menu-item index="/admin/user/list">用户列表</el-menu-item>
        <el-menu-item index="/admin/user/policy">用户策略</el-menu-item>
        <el-menu-item index="/admin/user/online">在线用户</el-menu-item>
        <el-menu-item index="/admin/user/lockmanager">锁定管理</el-menu-item>
        <el-menu-item index="/admin/user/ip_map">IP映射</el-menu-item>
      </el-submenu>

      <el-submenu index="3">
        <template slot="title">
          <i class="el-icon-folder"></i>
          <span slot="title">用户组信息</span>
        </template>
        <el-menu-item index="/admin/group/list">用户组列表</el-menu-item>
      </el-submenu>

      <el-submenu index="4">
        <template slot="title">
          <i class="el-icon-data-analysis"></i>
          <span slot="title">调试信息</span>
        </template>
        <el-menu-item>
          <a href="/debug/pprof/" target="_blank">pprof</a>
        </el-menu-item>
        <el-menu-item>
          <a href="/debug/statsviz/" target="_blank">statsviz</a>
        </el-menu-item>
      </el-submenu>
    </el-menu>
  </div>
</template>

<script>
export default {
  name: "LayoutAside",
  data() {
    return {}
  },
  props: ['is_active', 'route_path'],
  mounted() {
  }
}
</script>

<style scoped>
.sidebar-container {
  height: 100%;
  background: linear-gradient(180deg, rgb(24, 54, 97) 0%, rgb(18, 43, 77) 100%);
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
}

.sidebar-logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.1);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.sidebar-logo-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-icon i {
  font-size: 20px;
  color: #fff;
}

.logo-text {
  font-size: 14px;
  font-weight: bold;
  color: #fff;
  font-family: "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
  line-height: 1.4;
  text-align: center;
}

.sidebar-logo-mini {
  display: flex;
  align-items: center;
  justify-content: center;
}

.sidebar-logo-mini i {
  font-size: 24px;
  color: #fff;
}

.sidebar-logo-fade-enter-active,
.sidebar-logo-fade-leave-active {
  transition: opacity 0.3s;
}

.sidebar-logo-fade-enter,
.sidebar-logo-fade-leave-to {
  opacity: 0;
}

.layout-menu {
  height: calc(100% - 64px);
  border-right: none;
  overflow-y: auto;
}

.layout-menu::-webkit-scrollbar {
  width: 6px;
}

.layout-menu::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
}

.layout-menu::-webkit-scrollbar-track {
  background: transparent;
}

/* 菜单项样式 */
.layout-menu >>> .el-menu-item,
.layout-menu >>> .el-submenu__title {
  height: 50px;
  line-height: 50px;
  color: rgba(255, 255, 255, 0.85);
  transition: all 0.3s;
}

.layout-menu >>> .el-menu-item:hover,
.layout-menu >>> .el-submenu__title:hover {
  background-color: rgba(255, 255, 255, 0.1) !important;
  color: #fff;
}

.layout-menu >>> .el-menu-item.is-active {
  background: linear-gradient(90deg, rgba(102, 126, 234, 0.3) 0%, transparent 100%) !important;
  border-left: 3px solid #667eea;
  color: #fff;
}

.layout-menu >>> .el-menu-item i,
.layout-menu >>> .el-submenu__title i {
  color: rgba(255, 255, 255, 0.85);
  font-size: 18px;
  margin-right: 8px;
}

.layout-menu >>> .el-submenu__icon-arrow {
  color: rgba(255, 255, 255, 0.6);
}

/* 子菜单样式 */
.layout-menu >>> .el-menu--inline {
  background-color: rgba(0, 0, 0, 0.1) !important;
}

.layout-menu >>> .el-menu--inline .el-menu-item {
  padding-left: 60px !important;
  background-color: transparent !important;
}

.layout-menu >>> .el-menu--inline .el-menu-item:hover {
  background-color: rgba(255, 255, 255, 0.05) !important;
}

.layout-menu >>> .el-menu-item a {
  display: block;
  color: #fff;
  text-decoration: none;
}

/* 折叠状态样式 */
.layout-menu.el-menu--collapse {
  width: 64px;
}

.layout-menu.el-menu--collapse >>> .el-menu-item i,
.layout-menu.el-menu--collapse >>> .el-submenu__title i {
  margin-right: 0;
}
</style>