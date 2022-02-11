import API from "../api/api_user";
import Common from "../utils/common";

let handleMessage = async function (message, file, next, socket) {
    console.log(message);
    //优先判断消息类型
    switch (message.type) {
        case 'image':
            //上传文件
            var config = {
                headers: {
                    "Content-Type":
                        "multipart/form-data; boundary=----WebKitFormBoundaryVCFSAonTuDbVCoAN",
                },
            };
            var params = new FormData();
            params.append("file", file, file.name);
            params.append("chat_id", message.id);
            await API.uploadChatImage(params, config).then((res) => {
                if (res.code != 0) {
                    //执行到next消息会停止转圈，如果接口调用失败，可以修改消息的状态 next({status:'failed'});
                    next({ status: "failed" });
                    return false;
                }
            });
            break;
        case 'file':
            //上传文件
            var config = {
                headers: {
                    "Content-Type":
                        "multipart/form-data; boundary=----WebKitFormBoundaryVCFSAonTuDbVCoAN",
                },
            };
            var params = new FormData();
            params.append("file", file, file.name);
            params.append("chat_id", message.id);
            await API.uploadChatImage(params, config).then((res) => {
                if (res.code != 0) {
                    //执行到next消息会停止转圈，如果接口调用失败，可以修改消息的状态 next({status:'failed'});
                    next({ status: "failed" });
                    return false;
                }
            });
            break;
    }
    //通过接口存储消息
    await API.saveMessage(message).then((res) => {
        if (res.code != 0) {
            //执行到next消息会停止转圈，如果接口调用失败，可以修改消息的状态 next({status:'failed'});
            next({ status: "failed" });
            return false;
        }
    })

    //暂停一秒后发送socket
    setTimeout(() => {
        next();
        //发送
        if (!socket.sendMsg(message)) {
            //定时发送
            var t1 = setInterval(() => {
                this.socketSuccess = socket.sendMsg(message);
            }, 500);
            if (this.socketStatus) {
                clearInterval(t1);
            }
        }
    }, 1000);
}

//获取视频地址
let handleMessageVideo = function (message, that) {
    let url = message.content;
    if (!url) {
        //通过接口获取
        API.getVideoUrl({ message_id: message.id }).then((res) => {
            if (res.code === 0) {
                videoShow(message, res.data.url, that);
            }
        });
    } else {
        videoShow(message, url, that);
    }
}
let videoShow = function (message, url, that) {
    //获取文件类型
    var fileType = Common.fileType(message.fileName);
    if (fileType == "video") {
        that.videoOptions.sources = [];
        that.showVideo = true;
        let source = {
            src: url,
            type: "video/mp4",
        };
        that.videoOptions.sources.push(source);
    }
}
export default { handleMessage, handleMessageVideo }