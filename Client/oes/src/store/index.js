import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    onceFlag: true,
    wsConn: "",
    notificationData: "",
    chatData: "",
    whiteBoardData: "",
    notificationCount: 0,
    onlineUsers: [],
    broadcastData: null
  },
  mutations: {
    setConn(state, data) {
      state.wsConn = data
    },
    setOnlyOnce(state, data) {
      state.onceFlag = data
    },
    setNotification(state, data) {
      state.notificationData = data.body
    },
    setChat(state, data) {
      state.chatData = data.user
    },
    setWhiteBoard(state, data) {
      state.whiteBoardData = data
    },
    setNotificationCount(state, data) {
      if (data === 0) {
        state.notificationCount = data
      } else {
        state.notificationCount += data
      }
    },
    setOnlineUsers(state, data) {
      state.onlineUsers.push(data.user)
    },
    setBroadcast(state, data) {
      state.broadcastData = data
    }
  },
  actions: {
  },
  getters: {
    getConn: (state) => state.wsConn,
    getOnce: (state) => state.onceFlag,
    getNotification: (state) => state.notificationData,
    getChat: (state) => state.chatData,
    getWhiteBoard: (state) => state.whiteBoardData,
    getNotificationCount: (state) => state.notificationCount,
    getOnlineUsers: (state) => state.onlineUsers,
    getBroadcast: (state) => state.broadcastData
  },
  modules: {
  }
})
