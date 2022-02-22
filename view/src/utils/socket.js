let wsUrl = process.env.WS_URL;
let socket = null;
let lockReconnet = false; //避免重复连接
let isReconnet = false;

//连接
let connetSocket = () => {
    //判断用户是否登陆
    var user = JSON.parse(localStorage.getItem('access-user'))
    if (!user) {
        return false
    }
    //实例websocket
    try {
        if ('WebSocket' in window) {
            socket = new WebSocket(wsUrl + "?token=" + user.token)
        } else if ('MozWebSocket' in window) {
            socket = new MozWebSocket(wsUrl + "?token=" + user.token)
        }
        initSocket()
    } catch (e) {
        reconnetSocket()
    }
}

//重新连接
let reconnetSocket = () => {
    if (lockReconnet) {
        return false
    }
    isReconnet = true
    lockReconnet = true

    setTimeout(() => {
        connetSocket()
        lockReconnet = false
    }, 1000)
}

//初始化websocket
let initSocket = function () {
    socket.onopen = () => {
        //heartCheck.reset().start() //后端说暂时不需要做心跳检测
        //执行全局回调函数
        if (isReconnet) {
            console.log('websocket 重新连接了')
            isReconnet = false
        } else {
            console.log('websocket 连接成功')
        }
    }

    socket.onerror = () => {
        console.log('websocket服务出错了 ---onerror');
        reconnetSocket()
    }

    socket.onclose = () => {
        console.log('websocket服务关闭了 ---onclose');
        reconnetSocket()
    }
}


//发送数据
let sendMsg = (data) => {
    if (socket.readyState === 1) {
        data = JSON.stringify(data);
        socket.send(data);
        console.log('socket发送成功')
        return true;
    } else {
        setTimeout(() => {
            console.log('等待socket链接成功')
            sendMsg(data)
        }, 1500)
        return false
    }
}

//接收数据
let getMsg = (callback) => {
    socket.onmessage = ev => {
        callback && callback(ev)
    }
}

//心跳检测
let heartCheck = {
    timeout: 60 * 1000,
    timeoutObj: null,
    serverTimeoutObj: null,
    reset() {
        clearTimeout(this.timeoutObj)
        clearTimeout(this.serverTimeoutObj)
        return this;
    },
    start() {
        let that = this;
        this.timeoutObj = setTimeout(() => {
            //发送数据，如果onmessage能接收到数据，表示连接正常,然后在onmessage里面执行reset方法清除定时器
            socket.send('heart check')
            this.serverTimeoutObj = setTimeout(() => {
                socket.close()
            }, that.timeout)
        }, this.timeout)
    }
}

let getSocketStatus = function () {
    if (socket.readyState != 1) {
        return false;
    }
    return true;
}

export default { sendMsg, getMsg, connetSocket, getSocketStatus }