import Vue from 'vue'
import VueRouter from 'vue-router'
// import Home from '../views/Home.vue'
import AppLayout from "@/components/common/AppLayout";

Vue.use(VueRouter);
function loadView(view) {
  return () =>
    import(/* webpackChunkName: "view-[request]" */ `@/views/${view}.vue`);
}


const routes = [
  {
    path: "/",
    name: "Login",
    component: loadView("Login"),
    meta: {
      title: "Login"
    }
  },
  {
    // Layout
    path: "/layout",
    name: "layout",
    component: AppLayout,
    children: [
      {
        path: "/about",
        name: "About",
        component: loadView("About")
      }
    ]
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
