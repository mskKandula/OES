import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
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
      if (data.add) {
        state.onlineUsers.push(data.user)
      } else {
        const index = state.onlineUsers.indexOf(data.user)
        if (index > -1) {
          state.onlineUsers.splice(index, 1)
        }
      }
    },
    setBroadcast(state, data) {
      state.broadcastData = data
    }
  },
  actions: {
  },
  getters: {
    getConn: (state) => state.wsConn,
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
