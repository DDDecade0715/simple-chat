
<template>
  <div class="login_box">
    <vue-particles
      color="#fff"
      :particleOpacity="0.7"
      :particlesNumber="60"
      shapeType="circle"
      :particleSize="4"
      linesColor="#fff"
      :linesWidth="1"
      :lineLinked="true"
      :lineOpacity="0.4"
      :linesDistance="150"
      :moveSpeed="2"
      :hoverEffect="true"
      hoverMode="grab"
      :clickEffect="true"
      clickMode="push"
    >
    </vue-particles>

    <div class="login-s">
      <el-form label-width="0px" class="login_form">
        <el-form-item label="">
          <el-input
            type="text"
            prefix-icon="el-icon-user"
            class="qxs-icon"
            placeholder="用户名"
            v-model="userName"
          ></el-input
        ></el-form-item>
        <el-form-item label="">
          <el-input
            type="text"
            prefix-icon="el-icon-unlock"
            class="qxs-icon"
            placeholder="密码"
            v-model="password"
            show-password
          ></el-input>
        </el-form-item>
        <el-form-item label="">
          <el-button
            class="login_btn"
            @click="login"
            type="primary"
            round
            :loading="isBtnLoading"
            >登录</el-button
          >
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>
 
 
 
<script>
import API from "../../api/api_user";
import md5 from "js-md5";
export default {
  data() {
    return {
      userName: "",
      password: "",
      isBtnLoading: false,
    };
  },
  created() {
    if (
      JSON.parse(localStorage.getItem("access-user")) &&
      JSON.parse(localStorage.getItem("access-user")).displayName
    ) {
      this.userName = JSON.parse(
        localStorage.getItem("access-user")
      ).displayName;
    }
  },
  computed: {
    btnText() {
      if (this.isBtnLoading) return "登录中...";
      return "登录";
    },
  },
  methods: {
    login() {
      if (!this.userName) {
        this.$message.error("请输入用户名");
        return;
      }
      if (!this.password) {
        this.$message.error("请输入密码");
        return;
      }
      this.isBtnLoading = true;
      let params = {
        username: this.userName,
        password: md5(this.password),
      };
      let that = this;
      API.login(params)
        .then(function (res) {
          if (res.code === 0) {
            that.isBtnLoading = false;
            let userInfo = {
              id: res.data.id,
              displayName: res.data.username,
              avatar: res.data.avatar,
              token: res.data.token,
              timestamp: new Date().getTime(),
            };
            localStorage.setItem("access-user", JSON.stringify(userInfo)); // 将用户信息存到localStorage中
            that.$router.replace({ path: "/im" }); // 登录成功跳转
          } else {
            that.isBtnLoading = false;
            that.$message.error(res.msg); // elementUI消息提示
          }
        })
        .catch(function (error) {
          that.isBtnLoading = false;
        });
    },
  },
};
</script>
<style lang="stylus">

.login_box {
  background-image: linear-gradient(-180deg, #1a1454 0%, #0e81a5 100%);
  background-repeat: no-repeat;
  background-size: cover;
  height: 100%;
  text-align: center;
}

.login-s {
  position: absolute;
  width: 350px;
  height: 400px;
  background-color: rgba(0, 0, 0, 0.3);
  top: 40%;
  left: 50%;
  transform: translate(-50%, -50%);
  border-radius: 8px;
  z-index: 2;
  box-shadow: 10px 10px 20px 0px rgba(0, 0, 0, 0.7);
  text-align: center;
}

.login_form {
  position: relative;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  padding: 0 30px;
}

.login_logo {
  height: 100%;
}

.login_btn {
  width: 50%;
  font-size: 16px;
  background: -webkit-linear-gradient(
    left,
    #000099,
    #2154fa
  );
  /* Safari 5.1 - 6.0 */
  background: -o-linear-gradient(
    right,
    #000099,
    #2154fa
  );
  /* Opera 11.1 - 12.0 */
  background: -moz-linear-gradient(
    right,
    #000099,
    #2154fa
  );
  /* Firefox 3.6 - 15 */
  background: linear-gradient(to right, #000099, #2154fa); /* 标准的语法 */
  filter: brightness(1.4);
  margin-top: 3%;
}
</style>