import API from "../api/api_user";

//获取用户联系人
let getUserContact = function (){
    API.getContacts()
}

//获取群组联系人
let getGroupContact = function (){
    API.getGroupContacts()
}


//切换联系人
let handleChangeContact = function (contact, instance) {
    instance.updateContact({
        id: contact.id,
        unread: 0,
    });
    instance.closeDrawer();
}

//切换联系人后获取历史记录
let getContactMessage = function (contact, next, that) {
    //从后端请求消息数据，包装成下面的样子
    let params = { contact_id: contact.id, type: contact.type };
    API.getMessages(params).then(function (res) {
        if (res.code === 0) {
            if (Array.isArray(res.data) && res.data.length > 0) {
                that.messages = res.data;
            }
        }
        next(that.messages, true);
    });
}

export default { handleChangeContact, getContactMessage }