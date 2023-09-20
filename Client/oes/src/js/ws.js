import store from '@/store'

export function connectWebsocket() {

    const url = new URL("ws://localhost:9000/ws");
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
        data = data.split(/\r?\n/);
        mutateData(data[0]);
        data = "";
    };
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