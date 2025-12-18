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

    <!-- 修改密码对话框 -->
    <el-dialog
      title="修改密码"
      :visible.sync="passwordDialogVisible"
      width="450px"
      :close-on-click-modal="false">
      <el-form :model="passwordForm" :rules="passwordRules" ref="passwordForm" label-width="100px">
        <el-form-item label="原密码" prop="old_pass">
          <el-input
            type="password"
            v-model="passwordForm.old_pass"
            placeholder="请输入原密码"
            autocomplete="off"
            show-password>
          </el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="new_pass">
          <el-input
            type="password"
            v-model="passwordForm.new_pass"
            placeholder="请输入新密码（至少6位）"
            autocomplete="off"
            show-password>
          </el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm_pass">
          <el-input
            type="password"
            v-model="passwordForm.confirm_pass"
            placeholder="请再次输入新密码"
            autocomplete="off"
            show-password>
          </el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="passwordDialogVisible = false">取 消</el-button>
        <el-button type="primary" :loading="passwordLoading" @click="submitPasswordChange">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {getUser, removeToken} from "@/plugins/token";
import axios from "axios";

export default {
  name: "Layoutheader",
  props: ['route_name'],
  data() {
    const validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入新密码'));
      } else if (value.length < 6) {
        callback(new Error('密码长度不能少于6位'));
      } else {
        if (this.passwordForm.confirm_pass !== '') {
          this.$refs.passwordForm.validateField('confirm_pass');
        }
        callback();
      }
    };
    const validateConfirmPass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入新密码'));
      } else if (value !== this.passwordForm.new_pass) {
        callback(new Error('两次输入密码不一致'));
      } else {
        callback();
      }
    };
    
    return {
      is_active: true,
      passwordDialogVisible: false,
      passwordLoading: false,
      passwordForm: {
        old_pass: '',
        new_pass: '',
        confirm_pass: ''
      },
      passwordRules: {
        old_pass: [
          { required: true, message: '请输入原密码', trigger: 'blur' }
        ],
        new_pass: [
          { required: true, validator: validatePass, trigger: 'blur' }
        ],
        confirm_pass: [
          { required: true, validator: validateConfirmPass, trigger: 'blur' }
        ]
      }
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
        this.passwordDialogVisible = true;
        this.$nextTick(() => {
          this.$refs.passwordForm.resetFields();
        });
      }
    },
    submitPasswordChange() {
      this.$refs.passwordForm.validate((valid) => {
        if (!valid) {
          return false;
        }
        
        this.passwordLoading = true;
        axios.post('/base/set_pwd', {
          old_pass: this.passwordForm.old_pass,
          new_pass: this.passwordForm.new_pass
        }).then(resp => {
          const rdata = resp.data;
          if (rdata.code === 0) {
            this.$message.success('密码修改成功，请重新登录');
            this.passwordDialogVisible = false;
            setTimeout(() => {
              removeToken();
              this.$router.push("/login");
            }, 1500);
          } else {
            this.$message.error(rdata.msg || '密码修改失败');
          }
        }).catch(error => {
          if (error.response && error.response.status === 401) {
            this.$message.error('原密码错误');
          } else {
            this.$message.error('密码修改失败');
          }
          console.log(error);
        }).finally(() => {
          this.passwordLoading = false;
        });
      });
    }
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