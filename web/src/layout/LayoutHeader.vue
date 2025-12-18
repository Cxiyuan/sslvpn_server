<template>
  <div class="layout-header">
    <div class="header-left">
      <i @click="toggleClick" :class="is_active ? 'el-icon-s-fold' : 'el-icon-s-unfold'" class="toggle-icon"></i>

      <el-breadcrumb separator="/" class="app-breadcrumb">
        <el-breadcrumb-item v-for="(item, index) in route_name" :key="index">{{ item }}</el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <div class="header-right">
      <el-dropdown trigger="click" @command="handleCommand" class="user-dropdown">
        <div class="user-info">
          <div class="avatar">
            <i class="el-icon-user-solid"></i>
          </div>
          <span class="username">{{ admin_user }}</span>
          <i class="el-icon-arrow-down"></i>
        </div>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item command="changePassword">
            <i class="el-icon-key"></i>
            修改密码
          </el-dropdown-item>
          <el-dropdown-item command="logout">
            <i class="el-icon-switch-button"></i>
            退出登录
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
import {getUser, removeToken} from "@/plugins/token";

export default {
  name: "Layoutheader",
  props: ['route_name'],
  data() {
    return {
      is_active: true
    }
  },
  computed: {
    admin_user() {
      return getUser();
    },
  },
  methods: {
    toggleClick() {
      this.is_active = !this.is_active
      this.$emit('update:is_active', this.is_active)
    },
    handleCommand(command) {
      if (command === 'logout') {
        removeToken()
        this.$router.push("/login");
      } else if (command === 'changePassword') {
        this.$router.push("/admin/set/other");
      }
    },
  }
}
</script>

<style scoped>
.layout-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.toggle-icon {
  font-size: 20px;
  color: #606266;
  cursor: pointer;
  transition: all 0.3s;
  padding: 8px;
  border-radius: 4px;
}

.toggle-icon:hover {
  background: #f5f7fa;
  color: #409eff;
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.brand-logo i {
  font-size: 24px;
  color: #667eea;
}

.app-breadcrumb {
  font-size: 14px;
}

.app-breadcrumb >>> .el-breadcrumb__inner {
  color: #606266;
  font-weight: normal;
}

.app-breadcrumb >>> .el-breadcrumb__item:last-child .el-breadcrumb__inner {
  color: #303133;
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-dropdown {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: 20px;
  transition: all 0.3s;
}

.user-info:hover {
  background: #f5f7fa;
}

.avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar i {
  font-size: 16px;
  color: #fff;
}

.username {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}

.user-info .el-icon-arrow-down {
  font-size: 12px;
  color: #909399;
  transition: transform 0.3s;
}

.user-dropdown:hover .el-icon-arrow-down {
  transform: rotate(180deg);
}

/* 下拉菜单样式 */
.user-dropdown >>> .el-dropdown-menu__item {
  padding: 10px 20px;
}

.user-dropdown >>> .el-dropdown-menu__item i {
  margin-right: 8px;
  color: #909399;
}

.user-dropdown >>> .el-dropdown-menu__item:hover i {
  color: #409eff;
}
</style>