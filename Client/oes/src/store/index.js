import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    onceFlag: true,
    wsConn: "",
    notificationData:"",
    chatData:"",
    whiteBoardData:"",
    notificationCount:0
  },
  mutations: {
    setConn (state, data) {
      state.wsConn = data
    },
    setOnlyOnce(state, data){
      state.onceFlag = data
    },
    setNotification(state,data){
      state.notificationData=data
    },
    setChat(state,data){
      state.chatData=data
    },
    setWhiteBoard(state,data){
      state.whiteBoardData=data
    },
    setNotificationCount(state,data){
      if (data ===0){
        state.notificationCount = data
      }else{
      state.notificationCount += data
      }
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
  },
  modules: {
  }
})
