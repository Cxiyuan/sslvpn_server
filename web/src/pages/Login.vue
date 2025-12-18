<template>
  <div class="login-container">
    <div class="login-background">
      <div class="background-overlay"></div>
    </div>
    
    <div class="login-box">
      <div class="login-header">
        <div class="logo-icon">
          <i class="el-icon-lock"></i>
        </div>
        <h1 class="login-title">安全访问管理系统</h1>
      </div>

      <el-form :model="ruleForm" status-icon :rules="rules" ref="ruleForm" class="login-form">
        <el-form-item prop="admin_user">
          <el-input 
            v-model="ruleForm.admin_user" 
            placeholder="请输入管理员账号"
            prefix-icon="el-icon-user"
            clearable>
          </el-input>
        </el-form-item>
        <el-form-item prop="admin_pass">
          <el-input 
            type="password" 
            v-model="ruleForm.admin_pass" 
            placeholder="请输入密码"
            prefix-icon="el-icon-lock"
            autocomplete="off"
            show-password>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button 
            type="primary" 
            :loading="isLoading" 
            @click="submitForm('ruleForm')"
            class="login-button">
            登 录
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import qs from "qs";
import {setToken, setUser} from "@/plugins/token";

export default {
  name: "Login",
  mounted() {
    // 进入login，删除登录信息
    console.log("login created")
    //绑定事件
    window.addEventListener('keydown', this.keyDown);
  },
  destroyed(){
    window.removeEventListener('keydown',this.keyDown,false);
  },
  data() {
    return {
      ruleForm: {},
      rules: {
        admin_user: [
          {required: true, message: '请输入用户名', trigger: 'blur'},
          {max: 50, message: '长度小于 50 个字符', trigger: 'blur'}
        ],
        admin_pass: [
          {required: true, message: '请输入密码', trigger: 'blur'},
          {min: 6, message: '长度大于 6 个字符', trigger: 'blur'}
        ],
      },
    }
  },
  methods: {
    keyDown(e) {
      //如果是回车则执行登录方法
      if (e.keyCode === 13) {
        this.submitForm('ruleForm');
      }
    },
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (!valid) {
          console.log('error submit!!');
          return false;
        }
        this.isLoading = true

        // alert('submit!');
        axios.post('/base/login', qs.stringify(this.ruleForm)).then(resp => {
          var rdata = resp.data
          if (rdata.code === 0) {
            this.$message.success(rdata.msg);
            setToken(rdata.data.token)
            setUser(rdata.data.admin_user)
            this.$router.push("/home");
          } else {
            this.$message.error(rdata.msg);
          }
          console.log(rdata);
        }).catch(error => {
          this.$message.error('哦，请求出错');
          console.log(error);
        }).finally(() => {
              this.isLoading = false
            }
        );

      });
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
    },
  },
}
</script>

<style scoped>
.login-container {
  position: relative;
  width: 100%;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: #183661;
  z-index: 1;
}

.background-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1440 320"><path fill="%23ffffff" fill-opacity="0.1" d="M0,96L48,112C96,128,192,160,288,186.7C384,213,480,235,576,213.3C672,192,768,128,864,128C960,128,1056,192,1152,197.3C1248,203,1344,149,1392,122.7L1440,96L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"></path></svg>') no-repeat bottom;
  background-size: cover;
  opacity: 0.3;
}

.login-box {
  position: relative;
  z-index: 10;
  width: 450px;
  padding: 50px 40px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: slideUp 0.6s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
}

.logo-icon i {
  font-size: 40px;
  color: #fff;
}

.login-title {
  font-size: 28px;
  font-weight: bold;
  color: #333;
  margin: 0 0 10px 0;
  font-family: "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
}

.login-subtitle {
  font-size: 14px;
  color: #999;
  margin: 0;
  letter-spacing: 1px;
}

.login-form {
  margin-top: 30px;
}

.login-form >>> .el-form-item {
  margin-bottom: 24px;
}

.login-form >>> .el-input__inner {
  height: 48px;
  line-height: 48px;
  border-radius: 8px;
  border: 1px solid #e0e0e0;
  font-size: 14px;
  transition: all 0.3s;
}

.login-form >>> .el-input__inner:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
}

.login-form >>> .el-input__prefix {
  left: 15px;
  font-size: 16px;
  color: #999;
}

.login-form >>> .el-input--prefix .el-input__inner {
  padding-left: 45px;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: bold;
  border-radius: 8px;
  background: #183661;
  border: none;
  transition: all 0.3s;
  margin-top: 10px;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(24, 54, 97, 0.4);
}

.login-button:active {
  transform: translateY(0);
}
</style>