import Vue from "vue";
import VueRouter from "vue-router";
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
      title: "Login",
    },
  },
  {
    // Layout
    path: "/layout",
    name: "layout",
    component: AppLayout,
    children: [
      {
        path: "/dashboard",
        name: "Dashboard",
        component: loadView("dashboard"),
        meta: {
          title: "Dashboard",
        },
      },
      {
        path: "/fileUpload",
        name: "File",
        component: loadView("File"),
        meta: {
          title: "File",
        },
      },

      {
        path: "/uploadQuestions",
        name: "UploadQuestions",
        component: loadView("UploadQuestions"),
        meta: {
          title: "UploadQuestions",
        },
      },
      {
        path: "/studentsList",
        name: "StudentsList",
        component: loadView("StudentsList"),
        meta: {
          title: "StudentsList",
        },
      },
      {
        path: "/questions",
        name: "Questions",
        component: loadView("Questions"),
        meta: {
          title: "Questions",
        },
      },
    ],
  },
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

export default router;
