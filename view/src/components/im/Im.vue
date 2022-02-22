<template>
  <div class="imui-center">
    <div class="chat-box">
      <lemon-imui
        v-loading="loading"
        ref="IMUI"
        :user="this.user"
        @pull-messages="handlePullMessages"
        @send="handleSend"
        @menu-avatar-click="handleMenuAvatarClick"
        @change-menu="handleChangeMenu"
        @change-contact="handleChangeContact"
        @message-click="handleClickMessage"
        :theme="theme"
      >
        <template #message-title="contact">
          <span class="lemon-container__displayname">{{
            contact.displayName
          }}</span>
          <i class="more el-icon-more" @click="changeDrawer(contact)"></i>
          <br />
        </template>
      </lemon-imui>
    </div>
    <el-drawer title="编辑用户资料" :visible.sync="userInfoBox" direction="rtl">
      <UserFormVue
        @userInfoClose="userInfoClose"
        @userInfoEdit="userInfoEdit"
      ></UserFormVue>
    </el-drawer>
    <el-image-viewer
      v-if="showViewer"
      :on-close="
        () => {
          showViewer = false;
        }
      "
      :url-list="imgList"
    />
    <div
      class="mask"
      v-if="showVideo"
      @click="
        () => {
          showVideo = false;
        }
      "
    ></div>
    <div class="videomasks" v-if="showVideo">
      <VideoPlayer :options="videoOptions" />
    </div>
  </div>
</template>
<script>
import socket from "../../utils/socket";
import Contact from "../../utils/contact";
import User from "../../utils/user";
import Message from "../../utils/message";
import API from "../../api/api";
//使用了element样式
import UserFormVue from "../userinfo/UserForm.vue";
import ContactInfoVue from "../contactInfo/ContactInfo.vue";
import VideoPlayer from "@/components/videoPlayer/VideoPlayer.vue";
import emoji from "../../assets/emoji";

export default {
  components: {
    UserFormVue,
    ContactInfoVue,
    "el-image-viewer": () =>
      import("element-ui/packages/image/src/image-viewer"),
    VideoPlayer,
  },
  data() {
    return {
      theme: "default",
      user: {
        id: 0,
        displayName: "",
        avatar: "",
      },
      messages: [],
      contacts: [],
      //隐藏用户编辑框
      userInfoBox: false,
      //某个群的信息
      groupInfo: {},
      //预览图片
      showViewer: false,
      imgList: [],
      //播放视频
      showVideo: false,
      videoOptions: {
        autoplay: true,
        controls: true,
        sources: [
          {
            src: "",
            type: "video/mp4",
          },
        ],
      },
      //获取消息历史记录
      page: 0,
      //加载
      loading: true,
    };
  },
  mounted() {
    this.initIm();
  },
  methods: {
    //初始化
    async initIm() {
      const i = this.$refs.IMUI;
      //初始化socket
      socket.connetSocket();
      //获取登录用户信息
      User.getUserInfo(this);
      //获取登录用户的联系人列表
      this.getUserContacts(i);

      var t2 = setInterval(() => {
        // console.log(socket.getSocketStatus());
        this.loading = !socket.getSocketStatus();
        if (socket.getSocketStatus()) {
          //获取消息
          this.getImMessage();
        }
      }, 2000);

      //初始化表情包。
      i.initEmoji(emoji);
    },

    //获取消息
    getImMessage() {
      var i = this.$refs.IMUI;
      socket.getMsg((event) => {
        let data = JSON.parse(event.data);
        if (data.type == "add_member" || data.type == "create_group") {
          let arr = [];
          data.members.forEach((val) => {
            arr.push(val.uuid);
          });
          this.contacts.forEach((v, k) => {
            //已经存在
            if (v.id == data.uuid) {
              this.contacts.splice(k, 1);
            }
          });
          let contact = {
            id: data.uuid,
            displayName: data.name,
            avatar: data.avatar,
            index: "群聊",
            unread: 0,
            type: "group",
            members: data.members,
            onMember: arr,
            lastContent: i.lastContentRender({
              type: data.message.type,
              content: data.message.content,
            }),
            lastSendTime: data.message.send_time,
          };
          this.contacts.push(contact);
        } else {
          i.appendMessage(data, true);
        }
      });
    },

    //获取联系人列表
    getUserContacts(im) {
      if (this.user.id != 0) {
        API.getContacts().then((res) => {
          if (res.code === 0) {
            res.data.forEach((v) => {
              let contact = {
                id: v.uuid,
                displayName: v.username,
                avatar: v.avatar,
                index: "好友",
                unread: 0,
                type: "user",
                onMember: [v.uuid],
              };
              if (v.message && v.message.type != "") {
                contact.lastContent = im.lastContentRender({
                  type: v.message.type,
                  content: v.message.content,
                });
                contact.lastSendTime = v.message.send_time;
              }
              this.contacts.push(contact);
            });
          }
        });
        API.getGroupContacts().then((res) => {
          if (res.code === 0) {
            res.data.forEach((v) => {
              let info = {
                id: v.uuid,
                displayName: v.name,
                avatar: v.avatar,
                index: "群聊",
                unread: 0,
                type: "group",
                members: v.members,
                onMember: [],
              };
              v.members.forEach((val) => {
                info.onMember.push(val.uuid);
              });
              if (v.message && v.message.type != "") {
                info.lastContent = im.lastContentRender({
                  type: v.message.type,
                  content: v.message.content,
                });
                info.lastSendTime = v.message.send_time;
              }
              this.contacts.push(info);
            });
          }
        });
        im.initContacts(this.contacts);
      }
    },

    //切换聊天联系人
    handleChangeContact(c, i) {
      this.page = 0;
      Contact.handleChangeContact(c, i);
    },

    //点击消息
    handleClickMessage(event, key, message, instance) {
      let that = this;
      //预览图片
      if (message.type == "image") {
        this.imgList = [];
        this.showViewer = true;
        this.imgList.push(message.content);
      }

      if (message.type == "file") {
        Message.handleMessageVideo(message, that);
      }
    },

    //切换导航栏
    handleChangeMenu() {
      //切换导航栏
      this.$refs.IMUI.closeDrawer();
    },

    //点击自己头像
    handleMenuAvatarClick() {
      //自己头像点击事件
      this.userInfoBox = !this.userInfoBox;
    },

    //获取联系人发送的消息
    handlePullMessages(contact, next) {
      this.page += 1;
      Contact.getContactMessage(contact, next, this);
    },

    //发送消息
    handleSend(message, next, file) {
      //通过接口存储消息
      Message.handleMessage(message, file, next, socket);
    },

    userInfoClose() {
      this.userInfoBox = !this.userInfoBox;
    },

    userInfoEdit(params) {
      this.user.displayName = params.displayName;
      this.user.avatar = params.avatar;
    },

    changeDrawer(contact) {
      let IMUI = this.$refs.IMUI;
      if (IMUI.drawerVisible === true) {
        IMUI.closeDrawer();
      } else {
        this.$store.commit("changeContact", contact);
        let params = {
          render: (contact) => {
            return <ContactInfoVue contact={contact}></ContactInfoVue>;
          },
        };
        params.width = 250;
        params.offsetY = 1;
        IMUI.openDrawer(params);
      }
    },
  },
};
</script>
<style lang="stylus">
* {
  margin: 0;
}

.imui-center {
  width: 100vw;
  height: 100vh;
  background: linear-gradient(45deg, palegoldenrod, pink, plum);
  // animation: hueRotate 10s infinite alternate;
}

@keyframes hueRotate {
  100% {
    filter: hue-rotate(360deg);
  }
}

.chat-box {
  position: absolute;
  top: 15%;
  left: 25%;
}

.more {
  font-size: 20px;
  line-height: 24px;
  height: 24px;
  position: absolute;
  top: 14px;
  right: 14px;
  cursor: pointer;
  user-select: none;
  display: inline-block;
  border-radius: 4px;
  padding: 0 8px;
}

.lemon-message-image .lemon-message__content img {
  max-width: 100%;
  min-width: auto;
  display: block;
}

.lemon-editor__inner {
  text-align: left;
}

.mask {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  right: 0;
  z-index: 10;
  background-color: #000000;
  opacity: 0.6;
}

.videomasks {
  max-width: 1200px;
  position: fixed;
  left: 50%;
  top: 50%;
  z-index: 20;
  transform: translate(-50%, -50%);
}
</style>
