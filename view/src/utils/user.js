//获取个人信息
let getUserInfo = function (that) {
    let info = JSON.parse(localStorage.getItem("access-user"));
    if (info) {
        that.user.id = info.id;
        that.user.displayName = info.displayName;
        that.user.avatar = info.avatar;
    }
}
export default { getUserInfo }