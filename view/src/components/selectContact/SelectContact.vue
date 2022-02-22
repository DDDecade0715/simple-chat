<template>
  <el-form style="height: 400px">
    <el-form-item
      style="
        float: left;
        height: 100%;
        width: 50%;
        overflow-y: auto;
        overflow-x: hidden;
      "
    >
      <el-checkbox-group v-model="selectOne">
        <el-checkbox
          v-for="(operate, index) in data"
          :label="operate.label"
          :key="index"
          @change="handleCheckedCitiesChange(operate)"
          style="width: 100%"
          :disabled="contact.onMember.includes(operate.key)"
          :checked="operate.checked"
        >
          <el-avatar
            :size="35"
            :src="operate.avatar"
            style="vertical-align: middle"
            shape="square"
          ></el-avatar>
          <span> {{ operate.label }}</span></el-checkbox
        >
      </el-checkbox-group>
    </el-form-item>
    <el-form-item
      style="
        float: right;
        height: 90%;
        width: 50%;
        overflow-y: auto;
        overflow-x: hidden;
      "
    >
      <div v-for="(item, key) in checkOne" style="width: 100%">
        <span class="el-checkbox__label"
          ><el-avatar
            :size="35"
            :src="item.avatar"
            style="vertical-align: middle"
            shape="square"
          ></el-avatar>
          <span> {{ item.label }}</span></span
        >
        <span
          style="
            cursor: pointer;
            outline: 0;
            line-height: 40px;
            vertical-align: middle;
            float: right;
          "
          @click="deleteMember(item)"
          ><i class="el-icon-circle-close"></i
        ></span>
      </div>
    </el-form-item>
    <el-form-item style="float: right; height: 10%; width: 50%">
      <el-button
        size="mini"
        @click="close"
        style="float: right; margin-left: 7px"
        >取消</el-button
      >
      <el-button
        type="success"
        @click="onCreate"
        v-show="contact.type === 'user'"
        size="mini"
        style="float: right"
        >立即创建</el-button
      >
      <el-button
        type="success"
        @click="onAdd"
        v-show="contact.type === 'group'"
        size="mini"
        style="float: right"
        >添加到群</el-button
      >
    </el-form-item>
    <!-- <el-form-item label="">
      <el-transfer
        v-model="onMember"
        :data="data"
        :titles="['联系人', '已选择联系人']"
      >
        <div slot-scope="{ option }" style="margin-bottom: 10px">
          <el-avatar
            :size="20"
            :src="option.avatar"
            style="vertical-align: middle"
            shape="square"
          ></el-avatar>
          <span> {{ option.label }}</span>
        </div></el-transfer
      >
    </el-form-item>
    <el-form-item>
      <el-button
        type="primary"
        @click="onCreate"
        v-show="contact.type === 'user'"
        >立即创建</el-button
      >
      <el-button type="primary" @click="onAdd" v-show="contact.type === 'group'"
        >添加到群</el-button
      >
      <el-button @click="close">取消</el-button>
    </el-form-item> -->
  </el-form>
</template>

<script>
import API from "../../api/api";
import socket from "../../utils/socket";
export default {
  data() {
    return {
      data: [],
      checkOne: [],
      selectOne: [],
    };
  },
  props: ["contact"],
  mounted() {
    this.load();
  },
  beforeUpdate() {},
  updated() {},
  watch: {},
  methods: {
    load() {
      API.getContacts().then((res) => {
        if (res.code === 0) {
          res.data.forEach((element) => {
            let i = {
              key: element.uuid,
              label: element.username,
              avatar: element.avatar,
              disabled: false,
              checked: false,
            };
            this.data.push(i);
          });
        }
      });
    },
    onCreate() {
      this.data.forEach((v,k)=>{
        if(v.key == this.contact.onMember[0]){
            this.checkOne.push(v)
        }
      })
      if (this.checkOne.length < 2) {
        this.$message.error("人数不够组群");
        return false;
      }

      if (this.contact.type === "group") {
        this.$message.error("不能创建群");
        return false;
      }
      let checkOneArray = [];
      this.checkOne.forEach((v, k) => {
        checkOneArray.push(v.key);
      });
      let paramas = { contact_ids: checkOneArray };
      API.createGroup(paramas).then((res) => {
        if (res.code === 0) {
          //发送给群里所有人
          let message = {
            type: "create_group",
            id: res.data.uuid,
            members: checkOneArray,
          };
          socket.sendMsg(message);
          this.$emit("close");
        }
      });
    },
    onAdd() {
      let checkOneArray = [];
      this.checkOne.forEach((v, k) => {
        checkOneArray.push(v.key);
      });
      let paramas = { group_id: this.contact.id, contact_ids: checkOneArray };
      API.addGroupMember(paramas).then((res) => {
        if (res.code === 0) {
          let message = {
            type: "add_member",
            id: res.data.uuid,
            members: checkOneArray,
          };
          socket.sendMsg(message);
          this.$emit("close");
        }
      });
    },
    close() {
      this.$emit("close");
    },
    handleCheckedCitiesChange(item) {
      if (this.checkOne.includes(item)) {
        let index = this.checkOne.indexOf(item);
        this.checkOne.splice(index, 1);
      } else {
        this.checkOne.push(item);
      }
    },
    deleteMember(item) {
      if (this.selectOne.includes(item.label)) {
        let index = this.selectOne.indexOf(item.label);
        this.selectOne.splice(index, 1);
      }
      if (this.checkOne.includes(item)) {
        let index = this.checkOne.indexOf(item);
        this.checkOne.splice(index, 1);
      }
    },
  },
};
</script>
<style lang="stylus"></style>