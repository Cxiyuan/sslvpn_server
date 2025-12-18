<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script>

export default {
  name: 'app',
  components: {},
  created() {
    const token = sessionStorage.getItem('jwtToken');
    console.log("App created", token)
    if (token) {
      this.$root.isLogin = true
    }
  },
  mounted() {
    console.log("App mounted")
  },
  data() {
    return {}
  },
  computed: {
    currentComponent: function () {
      var isLogin = this.$root.isLogin
      console.log("App isLogin", isLogin)
      if (isLogin) {
        return "layout";
      }
      return "login";
    },
  },
}
</script>

<style>
html, body {
  height: 100%;
  margin: 0;
}

#app {
  font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  height: 100%;
}

.hide {
  display: none;
}

.sh-10 {
  height: 10px;
}

.sh-20 {
  height: 20px;
}

.sw-10 {
  height: 1px;
  width: 10px;
}

.sw-20 {
  height: 1px;
  width: 20px;
}

.m-left-10 {
  margin-left: 10px;
}

/* 全局按钮颜色覆盖 */
.el-button--primary {
  background-color: #1425E8 !important;
  border-color: #1425E8 !important;
}

.el-button--primary:hover,
.el-button--primary:focus {
  background-color: #0d1ab3 !important;
  border-color: #0d1ab3 !important;
}

.el-button--danger {
  background-color: #E00B0F !important;
  border-color: #E00B0F !important;
}

.el-button--danger:hover,
.el-button--danger:focus {
  background-color: #b8090c !important;
  border-color: #b8090c !important;
}

</style>