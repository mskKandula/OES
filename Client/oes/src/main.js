import Go from './wasm_exec'
import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import BootstrapVue from "bootstrap-vue";
import VueAwesomeSwiper from "vue-awesome-swiper";
import VueCookies from 'vue-cookies';
// Link to css
// import "bootstrap/dist/css/bootstrap.css"
import "bootstrap-vue/dist/bootstrap-vue.css";
import "./assets/scss/bootstrap-custom.scss";
import "../node_modules/vue-multiselect/dist/vue-multiselect.min.css";
import "./assets/plugins/materialdesignicons/css/materialdesignicons.min.css";
// import "swiper/dist/css/swiper.css";
import "./assets/scss/template.scss";
import Axios from "axios";


const go = new Go()
WebAssembly.instantiateStreaming(
  fetch("main.wasm"), go.importObject
).then((result) => {
  go.run(result.instance)

  Vue.prototype.$go = {
    countser : this.countWords,
  }
})
Vue.prototype.$http = Axios;

Vue.config.productionTip = false;
// Install BootstrapVue
Vue.use(BootstrapVue);
Vue.use(VueAwesomeSwiper);
Vue.use(VueCookies);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
