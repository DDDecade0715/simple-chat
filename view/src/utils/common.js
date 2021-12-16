//根据后缀名判断文件类型
let fileType = function (fileName) {
    // 后缀获取
    var suffix = '';
    // 获取类型结果
    var result = '';
    try {
        var flieArr = fileName.split('.');
        suffix = flieArr[flieArr.length - 1];
    } catch (err) {
        suffix = '';
    }
    // fileName无后缀返回 false
    if (!suffix) {
        result = false;
        return result;
    }
    // 图片格式
    var imglist = ['png', 'jpg', 'jpeg', 'bmp', 'gif'];
    // 进行图片匹配
    result = imglist.some(function (item) {
        return item == suffix;
    });
    if (result) {
        result = 'image';
        return result;
    }
    ;
    // 匹配txt
    var txtlist = ['txt'];
    result = txtlist.some(function (item) {
        return item == suffix;
    });
    if (result) {
        result = 'txt';
        return result;
    }
    ;
    // 匹配 excel
    var excelist = ['xls', 'xlsx'];
    result = excelist.some(function (item) {
        return item == suffix;
    });
    if (result) {
        result = 'excel';
        return result;
    }
    ;
    // 匹配 word
    var wordlist = ['doc', 'docx'];
    result = wordlist.some(function (item) {
        return item == suffix;
    });
    if (result) {
        result = 'word';
        return result;
    }
    ;
    // 匹配 pdf
    var pdflist = ['pdf'];
    result = pdflist.some(function (item) {
        return item == suffix;
    });
    if (result) {
        result = 'pdf';
        return result;
    }
    ;
    // 匹配 ppt
    var pptlist = ['ppt'];
    result = pptlist.some(function (item) {
        return item == suffix;
    });
    if (result) {
        result = 'ppt';
        return result;
    }
    ;
    // 匹配 视频
    var videolist = ['mp4', 'm2v', 'mkv'];
    result = videolist.some(function (item) {
        return item == suffix;
    });
    if (result) {
        result = 'video';
        return result;
    }
    ;
    // 匹配 音频
    var radiolist = ['mp3', 'wav', 'wmv'];
    result = radiolist.some(function (item) {
        return item == suffix;
    });
    if (result) {
        result = 'radio';
        return result;
    }
    // 其他 文件类型
    result = 'other';
    return result;
}

export default { fileType }