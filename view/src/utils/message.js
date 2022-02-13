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
            //上传文
            await uploadVideo(file);
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


// 每个文件切片大小定为10M
var chunksize = 1024 * 1024 * 5;
// 定义上传总切片数
var chunktotal;
// 设置上传成功数量记录
var successTotal = 0

let uploadVideo = function (file) {
    // var file = document.getElementById("file").files[0];
    var start = 0;
    var end;
    var index = 0;
    var filesize = file.size;
    var filename = file.name;

    // 计算总的切片数
    chunktotal = Math.ceil(filesize / chunksize);
    while (start < filesize) {
        end = start + chunksize;
        if (end > filesize) {
            end = filesize;
        }

        var chunk = file.slice(start, end);//切割文件
        var chunkindex = index;
        var formData = new FormData();
        // 新增切片文件
        formData.append("file", chunk, filename);
        // 切片索引
        formData.append("chunkindex", chunkindex);
        // 切片总数
        formData.append("chunktotal", chunktotal);
        // 文件总大小
        formData.append("filesize", filesize)
        var config = {
            headers: {
                "Content-Type":
                    "multipart/form-data; boundary=----WebKitFormBoundaryVCFSAonTuDbVCoAN",
            },
        };
        API.uploadChatVideo(formData, config).then((res) => {
            if (res.code == 0) {
                successTotal = successTotal + 1
            }
        })
        start = end;
        index++;
    }
    console.log("上传数：", successTotal)
    // if (chunktotal == successTotal) {
    //     alert("上传成功")
    // } else {
    //     alert("上传失败")
    // }
}
export default { handleMessage, handleMessageVideo }