import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    onceFlag: true,
    wsConn: ""
  },
  mutations: {
    setConn (state, data) {
      state.wsConn = data
    },
    setOnlyOnce(state, data){
      state.onceFlag = data
    }
  },
  actions: {
  },
  getters: {
    getConn: (state) => state.wsConn,
    getOnce: (state) => state.onceFlag	
  },
  modules: {
  }
})
