<template>
  <el-container style="box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1)">
    <el-header height="auto">
      <el-row :gutter="20">
        <div class="info-box">
          <div class="grid-content">
            <el-button class="button-icon" @click="open"
              ><i class="el-icon-plus"></i
            ></el-button>
            <div class="displayName">添加</div>
          </div>
        </div>
        <template v-if="contact.type == 'user'">
          <div class="info-box">
            <div class="grid-content">
              <img
                v-if="contact.avatar"
                :src="contact.avatar"
                class="contact_avatar"
              />
              <span v-if="!contact.avatar" class="lemon-avatar contact_avatar"
                ><i class="lemon-icon-people"></i
              ></span>
              <div class="displayName">{{ contact.displayName }}</div>
            </div>
          </div>
        </template>
        <template v-if="contact.type == 'group'">
          <template v-for="v in contact.members">
            <div class="info-box">
              <div class="grid-content">
                <img
                  v-if="v.user_avatar"
                  :src="v.user_avatar"
                  class="contact_avatar"
                />
                <span v-if="!v.user_avatar" class="lemon-avatar contact_avatar"
                  ><i class="lemon-icon-people"></i
                ></span>
                <div class="displayName">{{ v.username }}</div>
              </div>
            </div>
          </template>
        </template>
      </el-row>
    </el-header>
    <el-divider></el-divider>
    <el-main>
      <template v-if="contact.type == 'group'">
        <el-form ref="form" label-width="60px">
          <el-form-item label="群头像">
            <el-upload
              class="avatar-uploader"
              action="#"
              :show-file-list="false"
              :before-upload="beforeAvatarUpload"
            >
              <img v-if="contact.avatar" :src="contact.avatar" class="avatar" />
              <i v-else class="el-icon-plus avatar-uploader-icon"></i>
            </el-upload>
          </el-form-item>
          <el-form-item label="群昵称">
            <el-input v-model="contact.displayName"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="formSubmit">立即提交</el-button>
          </el-form-item>
        </el-form>
      </template>
    </el-main>
    <el-dialog
      :visible.sync="selectVisible"
      center
      :append-to-body="true"
      :lock-scroll="false"
      style="text-align: center"
      :destroy-on-close="true"
      width="500px"
    >
      <SelectContactVue @close="close" :contact="contact"></SelectContactVue>
    </el-dialog>
  </el-container>
</template>

<script>
import SelectContactVue from "../selectContact/SelectContact.vue";
import API from "../../api/api_user";
export default {
  components: { SelectContactVue },
  data() {
    return {
      selectVisible: false,
    };
  },
  props: ["contact"],
  mounted() {},
  beforeUpdate() {},
  methods: {
    open() {
      this.selectVisible = true;
    },
    close() {
      this.selectVisible = false;
    },
    formSubmit() {
      if (this.contact.avatar == "") {
        this.$message.error("请上传头像");
        return false;
      }
      if (this.contact.displayName == "") {
        this.$message.error("请填写群昵称");
        return false;
      }
      let params = {
        group_id: this.contact.id,
        avatar: this.contact.avatar,
        name: this.contact.displayName,
      };
      API.updateGroup(params)
        .then((res) => {
          if (res.code === 0) {
            this.$message.success("修改成功");
          }
        })
        .catch((_) => {
          this.$message.error("修改失败");
        });
      return false;
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
      params.append("file", file, file.name);
      API.uploadAvatar(params, config)
        .then((res) => {
          if (res.code === 0) {
            this.contact.avatar = res.data.AccessUrl;
          }
        })
        .catch((_) => {
          this.$message.error("上传文件超时！");
        });
    },
  },
};
</script>
<style lang="stylus">
.el-container {
  width: 100%;
  height: 100%;
}

.el-header {
  .el-row {
    margin-bottom: 20px;

    &:first-child {
      margin-top: 6%;
    }

    &:last-child {
      margin-bottom: 0;
    }
  }

  .info-box {
    width: 20%;
    float: left;
    min-height: 20%;
    margin-top: 4%;
    margin-bottom: 2.5%;
    margin-left: 2.5%;
    margin-right: 2.5%;
  }

  .grid-content {
    text-align: center;
  }

  .button-icon {
    padding: 10px;
    border-radius: 4px;
    color: #8c939d;
  }

  .displayName {
    padding-top: 4px;
    font-size: 12px;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }

  .contact_avatar {
    width: 36px;
    height: 36px;
    border-radius: 4px;
    line-height: 36px;
  }
}

.el-main {
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
    font-size: 14px;
    color: #8c939d;
    width: 40px;
    height: 40px;
    line-height: 40px;
    text-align: center;
  }

  .avatar {
    width: 40px;
    height: 40px;
    display: block;
  }
}
</style>