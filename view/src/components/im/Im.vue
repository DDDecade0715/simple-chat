<template>
  <div class="imui-center">
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
    <div class="chat-box">
      <lemon-imui
        ref="IMUI"
        :user="this.user"
        @pull-messages="handlePullMessages"
        @send="handleSend"
        @menu-avatar-click="handleMenuAvatarClick"
        @change-menu="handleChangeMenu"
        @change-contact="handleChangeContact"
        @message-click="handleClickMessage"
        :theme="theme"
        width="850px"
        height="600px"
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
import API from "../../api/api_user";
//使用了element样式
import UserFormVue from "../userinfo/UserForm.vue";
import ContactInfoVue from "../contactInfo/ContactInfo.vue";
import VideoPlayer from "@/components/videoPlayer/VideoPlayer.vue";
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
      imageSuccess: true,
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
    };
  },
  mounted() {
    const { IMUI } = this.$refs;
    socket.createSocket(socket.wsUrl);
    //获取登录用户信息
    User.getUserInfo(this);
    //获取登录用户的联系人列表
    let that = this;
    setTimeout(function () {
      that.getUserContacts(IMUI);
    }, 1000);
    //初始化表情包。
    // IMUI.initEmoji(...);
    setTimeout(() => {
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
            lastContent: IMUI.lastContentRender({
              type: data.message.type,
              content: data.message.content,
            }),
            lastSendTime: data.message.send_time,
          };
          this.contacts.push(contact);
        } else {
          IMUI.appendMessage(data, true);
        }
      });
    }, 1000);
  },
  methods: {
    //获取联系人
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
      Contact.getContactMessage(contact, next, this);
    },

    //发送消息
    handleSend(message, next, file) {
      Message.handleMessage(message, file, this);
      //如果判断失败
      //执行到next消息会停止转圈，如果接口调用失败，可以修改消息的状态 next({status:'failed'});
      if (this.imageSuccess) {
        setTimeout(() => {
          next();
          socket.sendMsg(message);
        }, 1000);
      } else {
        next({ status: "failed" });
      }
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
.imui-center {
  background-image: linear-gradient(-180deg, #1a1454 0%, #0e81a5 100%);
  background-repeat: no-repeat;
  background-size: cover;
  height: 100%;
  text-align: center;
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

  &:active {
    background: #999;
  }
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
