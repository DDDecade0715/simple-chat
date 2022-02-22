import API from "../api/api";

//获取用户联系人
let getUserContact = function () {
    API.getContacts()
}

//获取群组联系人
let getGroupContact = function () {
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
    let params = { contact_id: contact.id, type: contact.type, page: that.page };
    API.getMessages(params).then(function (res) {
        if (res.code === 0) {
            if (Array.isArray(res.data.data) && res.data.data.length > 0) {
                that.messages = res.data.data;
                if (res.data.page == res.data.total_page) {
                    next(res.data.data, true);
                } else {
                    next(res.data.data, false);
                }
            } else {
                next([], true);
            }
        }
    });
}

export default { handleChangeContact, getContactMessage }