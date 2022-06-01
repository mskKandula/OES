import Vue from "vue";
import VueRouter from "vue-router";
// import Home from '../views/Home.vue'
import AppLayout from "@/components/common/AppLayout";

Vue.use(VueRouter);
function loadView(view) {
  return () =>
    import(/* webpackChunkName: "view-[request]" */ `@/views/${view}.vue`);
}

function loadCharts(view) {
  return () =>
    import(/* webpackChunkName: "view-[request]" */ `@/views/Charts/${view}.vue`);
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
        component: loadView("Dashboard"),
        meta: {
          title: "Dashboard",
        },
      },
      {
        path: "/studentDashboard",
        name: "StudentDashboard",
        component: loadView("StudentDashboard"),
        meta: {
          title: "StudentDashboard",
        },
      },
      {
        path: "/multipleStudentsRegistration",
        name: "MultipleStudentsRegistration",
        component: loadView("MultipleStudentsRegistration"),
        meta: {
          title: "MultipleStudentsRegistration",
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
        path: "/onlineExam",
        name: "Exam",
        component: loadView("Exam"),
        meta: {
          title: "Exam",
        },
      },
      {
        path: "/wordCounter",
        name: "WordCounter",
        component: loadView("WordCounter"),
        meta: {
          title: "WordCounter",
        },
      },
      {
        path: "/bar",
        name: "Bar",
        component: loadCharts("Bar"),
        meta: {
          title: "Bar",
        },
      },
      {
        path: "/uploadVideo",
        name: "UploadVideo",
        component: loadView("UploadVideo"),
        meta: {
          title: "UploadVideo",
        },
      },
      {
        path: "/fetchVideos",
        name: "VideoContent",
        component: loadView("VideoContent"),
        meta: {
          title: "VideoContent",
        },
      },
      {
        path: "/playVideo",
        name: "VideoPlay",
        component: loadView("VideoPlay"),
        meta: {
          title: "VideoPlay",
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
