import API from "../api/api_user";
import Common from "../utils/common";

let handleMessage = function (message, file, that) {
    if (file && (message.type == 'image' || message.type == 'file')) {
        let config = {
            headers: {
                "Content-Type":
                    "multipart/form-data; boundary=----WebKitFormBoundaryVCFSAonTuDbVCoAN",
            },
        };
        let params = new FormData();
        params.append("file", file, file.name);
        params.append("chat_id", message.id);
        API.uploadChatImage(params, config).then((res) => {
            if (res.code != 0) {
                that.imageSuccess = false;
            }
        });
    }
}
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