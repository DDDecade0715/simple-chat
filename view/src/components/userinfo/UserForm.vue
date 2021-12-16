<template>
  <div class="demo-drawer__content">
    <el-form>
      <el-form-item label="用户账号" :label-width="formLabelWidth">
        <el-input v-model="name" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="用户头像" :label-width="formLabelWidth">
        <el-upload
          class="avatar-uploader"
          action="#"
          :show-file-list="false"
          :before-upload="beforeAvatarUpload"
        >
          <img v-if="imageUrl" :src="imageUrl" class="avatar" />
          <i v-else class="el-icon-plus avatar-uploader-icon"></i>
        </el-upload>
      </el-form-item>
    </el-form>
    <div class="demo-drawer__footer">
      <el-button @click="cancelForm">取 消</el-button>
      <el-button type="primary" @click="confirmSubmit" :loading="loading">{{
        loading ? "提交中 ..." : "确 定"
      }}</el-button>
    </div>
  </div>
</template>

<script>
import API from "../../api/api_user";
export default {
  data: () => ({
    name: "",
    formLabelWidth: "80px",
    timer: null,
    loading: false,
    imageUrl: "",
  }),
  mounted() {
    let info = JSON.parse(localStorage.getItem("access-user"));
    this.name = info.displayName ? info.displayName : "";
    this.imageUrl = info.avatar ? info.avatar : "";
  },
  methods: {
    confirmSubmit() {
      //把资料提交
      let params = {
        account: this.name,
        avatar: this.imageUrl,
      };
      let that = this;
      API.updateUserinfo(params).then((res) => {
        if (res.code === 0) {
          //替换
          let info = JSON.parse(localStorage.getItem("access-user"));
          info.displayName = res.data.account;
          info.avatar = res.data.avatar;
          localStorage.setItem("access-user", JSON.stringify(info));
          that.$emit("userInfoEdit", info);
          that.$emit("userInfoClose");
        }else{
          that.$message.error(res.msg)
        }
      });
    },
    cancelForm() {
      this.loading = false;
      this.$emit("userInfoClose");
      clearTimeout(this.timer);
    },
    beforeAvatarUpload(file) {
      let isJPG = file.type === "image/jpeg";
      let isLt2M = file.size / 1024 / 1024 < 2;

      if (!isJPG) {
        this.$message.error("上传头像图片只能是 JPG 格式!");
        return false;
      }
      if (!isLt2M) {
        this.$message.error("上传头像图片大小不能超过 2MB!");
        return false;
      }
      let config = {
        headers: { "Content-Type": "multipart/form-data" },
      };
      let params = new FormData();
      let that = this;
      params.append("file", file, file.name);
      API.uploadAvatar(params, config)
        .then((res) => {
          if (res.code === 0) {
            that.imageUrl = res.data.AccessUrl;
          }
        })
        .catch((_) => {
          that.$message.error("上传文件超时！");
        });
    },
  },
};
</script>

<style>
.el-drawer__body {
  padding: 20px;
}
.demo-drawer__content {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.demo-drawer__content form {
  flex: 1;
}
.demo-drawer__footer {
  display: flex;
}
.demo-drawer__footer button {
  flex: 1;
}
.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}
.avatar-uploader .el-upload:hover {
  border-color: #409eff;
}
.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  line-height: 178px;
  text-align: center;
}
.avatar {
  width: 178px;
  height: 178px;
  display: block;
}
</style>