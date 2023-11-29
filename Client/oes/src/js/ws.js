import store from '@/store'

const HEARTBEAT_TIMEOUT = (1000 * 12);
const HEARTBEAT_VALUE = 1;

export function connectWebsocket() {

    const url = new URL("ws://" + window.location.host + "/r/ws");
    url.searchParams.append("role", sessionStorage.getItem("userType"));
    url.searchParams.append("id", sessionStorage.getItem("clientId"));
    url.searchParams.append("userId", sessionStorage.getItem("userId"));

    let ws = new WebSocket(url.href);

    store.commit("setConn", ws);

    ws.onconnect = (evt) => {
        console.log("ws connected", evt);
    };

    ws.onclose = (evt) => {
        console.log("ws closed", evt);
    };

    ws.onmessage = (evt) => {
        let data = evt.data;
        if (isBinary(data)) {
            heartBeat(ws)
        } else {
            data = data.split(/\r?\n/);
            mutateData(data[0]);
            data = "";
        }
    };
}

function isBinary(obj) {
    return typeof obj === 'object' && Object.prototype.toString.call(obj) === '[object Blob]';
}

function heartBeat(ws) {
    if (!ws) {
        return;
    } else if (ws.pingTimeout) {
        clearTimeout(ws.pingTimeout);
    }

    ws.pingTimeout = setTimeout(() => {
        ws.close();

        // business logic for deciding whether or not to reconnect
    }, HEARTBEAT_TIMEOUT);

    const data = new Uint8Array(1);
    data[0] = HEARTBEAT_VALUE;
    ws.send(data);
}

function mutateData(data) {
    let parsedData = JSON.parse(data);
    let dataBody = parsedData.body;

    if (parsedData.type == 1) {
        store.commit("setNotification", dataBody);
        store.commit("setNotificationCount", 1);
    } else if (parsedData.type == 2) {
        store.commit("setChat", dataBody);
    } else if (parsedData.type == 3) {
        store.commit("setWhiteBoard", dataBody);
    } else if (parsedData.type == 4) {
        store.commit("setBroadcast", dataBody);
    } else {
        store.commit("setOnlineUsers", dataBody);
    }

}