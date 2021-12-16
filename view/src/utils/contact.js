import API from "../api/api_user";
//切换联系人
let handleChangeContact = function (contact, instance) {
    instance.updateContact({
        id: contact.id,
        unread: 0,
    });
    instance.closeDrawer();
}

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