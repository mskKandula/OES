<template>
  <header id="page-header">
    <div class="container">
      <!-- top header -->
      <div class="top-header mt-2 mb-3 mb-lg-0">
        <div class="row justify-content-between no-gutters">
          <div class="col-auto">
            <div class="logo">
              <img
                class="logo-img"
                src="../../assets/images/logo_mkcl_w.svg"
                alt="logo here"
              />
            </div>
          </div>

          <div
            class="
              col-6
              text-right
              d-flex
              justify-content-end
              align-items-center
            "
          >
            <div class="notification-dropdown">
              <b-button variant="link" to="/notification" class="btn">
                <span class="mdi mdi-bell-outline notification-icon" />
                <span class="notification-count" v-if="notificationCount > 0">{{
                  notificationCount
                }}</span>
              </b-button>
            </div>
            <!-- <b-dropdown class="notification-dropdown header-dropdown mr-3">
              <template v-slot:button-content>
                <span class="mdi mdi-bell-outline notification-icon"></span>
                <span class="notification-count">1</span>
              </template>
              <b-dropdown-item href="#" class="text-truncate">
                <span class="mdi mdi-bell-ring-outline list-icon"></span> This is dynamic geneator notification text
              </b-dropdown-item>
              <b-dropdown-item href="#" class="text-truncate">
                <span class="mdi mdi-bell-ring-outline list-icon"></span> This is dynamic geneator notification text
              </b-dropdown-item>
            </b-dropdown>-->
            <b-dropdown right class="user-info-dropdown header-dropdown">
              <template v-slot:button-content>
                <div class="user-image-wrapper">
                  <div class="user-image-inner">
                    <img
                      class="user-image"
                      src="../../assets/images/dummy-image.png"
                      alt
                    />
                  </div>
                </div>
              </template>
              <b-dropdown-item href="#">
                <span class="mdi mdi-account-outline list-icon" /> My profile
              </b-dropdown-item>
              <b-dropdown-item @click="logout">
                <span class="mdi mdi-logout list-icon" /> Log Out
              </b-dropdown-item>
            </b-dropdown>
          </div>
        </div>
      </div>

      <!-- dynamic routing -->
      <div class="navigation">
        <div class="row no-gutters">
          <div class="col-md-12">
            <b-navbar toggleable="lg" type="dark" variant="transparent">
              <b-navbar-toggle target="nav-collapse" />
              <b-collapse id="nav-collapse" is-nav>
                <b-navbar-nav>
                  <b-nav-item
                    class="nav-item active"
                    v-for="(route, index) in this.routes"
                    :key="index"
                    @click="goTo(route)"
                  >
                    <span :class="icons[index]" />
                    {{ route.name }}
                  </b-nav-item>
                </b-navbar-nav>
              </b-collapse>
            </b-navbar>
          </div>
        </div>
      </div>

      <div>
        <ul>
          <NestetedMainMenus
            v-for="(node, index) in menulist"
            :key="index"
            :index="index"
            :children="node"
            :depth="1"
            :disabled-link="forceToChange"
          />
        </ul>
      </div>
      <!-- top navigation -->
      <!-- <div class="navigation">
        <div class="row no-gutters">
          <div class="col-md-12">
            <b-navbar toggleable="lg" type="dark" variant="transparent">
              <b-navbar-toggle target="nav-collapse" />
              <b-collapse id="nav-collapse" is-nav>
                <b-navbar-nav>
                  <b-nav-item class="nav-item active" to="/dashboard">
                    <span class="mdi mdi-chart-bell-curve list-icon" />
                    Dashboard
                  </b-nav-item> -->

      <!-- <b-nav-item class="nav-item" to="/formElements">
                    <span class="mdi mdi-bookmark-multiple-outline list-icon" />
                    Form Element
                  </b-nav-item>-->

      <!-- <b-nav-item-dropdown>
                    <template #button-content>
                      <span
                        class="mdi mdi-bookmark-multiple-outline list-icon"
                      ></span>
                      Custom Elements
                    </template>
                    <b-dropdown-item to="/formElements"
                      >Form Elements</b-dropdown-item
                    >
                    <b-dropdown-item to="/plugins">Plugins</b-dropdown-item>
                    <b-dropdown-item to="/alerts">Alerts</b-dropdown-item>
                    <b-dropdown-item to="/buttons"
                      >Buttons & Badges</b-dropdown-item
                    >
                    <b-dropdown-item to="/toster">Toster</b-dropdown-item>
                  </b-nav-item-dropdown>
                </b-navbar-nav> -->

      <!-- Right aligned nav items -->
      <!-- <b-navbar-nav class="ml-lg-auto search-input-wrapper">
                  <b-nav-form class="search-input">
                    <b-form-input
                      size="sm"
                      class="mr-sm-2"
                      placeholder="Search"
                    />
                  </b-nav-form>
                </b-navbar-nav>
              </b-collapse>
            </b-navbar>
          </div>
        </div>
      </div> -->
    </div>
  </header>
</template>
<script>
// import Vue from "vue";
import NestetedMainMenus from "../ui/NestetedMainMenus";
import { mapGetters } from "vuex";

// import lodashOrderBy from 'lodash/orderBy'

export default {
  name: "Header1",
  components: {
    NestetedMainMenus,
  },
  data() {
    return {
      routes: [],
      icons: [
        "mdi mdi-chart-bell-curve list-icon",
        "mdi mdi-account-plus list-icon",
        "mdi mdi-format-list-bulleted list-icon",
        "mdi mdi-file-upload list-icon",
        "mdi mdi-video-plus list-icon",
        "mdi mdi-monitor-dashboard list-icon",
      ],
    };
  },
  methods: {
    getRoutes() {
      let self = this;
      this.$http
        .get("/api/r/getRoutes")
        .then(function (res) {
          if (res.data) {
            self.routes = res.data.routes;
          }
        })
        .catch(function () {
          console.log("FAILURE!!");
        });
    },
    goTo(route) {
      let go = route.url;
      this.$router.push(go);
    },
    logout() {
      let self = this;
      this.$http
        .get("/api/r/logOut")
        .then(function () {
          self.$router.push("/");
        })
        .catch(function () {
          console.log("FAILURE!!");
        });
    },
  },
  // mounted() {
  //   let searchBox = document.querySelector(".search-input").parentElement;
  //   if (window.matchMedia("(min-width: 992px)").matches) {
  //     /* The viewport is less than, or equal to, 700 pixels wide */
  //     document
  //       .querySelector(".navigation .navbar-collapse")
  //       .appendChild(searchBox);
  //   } else {
  //     /* The viewport is greater than 700 pixels wide */
  //     document.querySelector(".navigation .navbar ").prepend(searchBox);
  //   }
  // },
  created() {
    this.getRoutes();
  },
  computed: {
    menuList() {
      let fields = [
        {
          label: "Charts",
          url: "",
          openInNewTab: false,
          nodes: [
            {
              label: "Bar",
              url: "/bar",
              openInNewTab: false,
            },
            {
              label: "Line",
              url: "/line",
              openInNewTab: false,
            },
          ],
        },
      ];
      return fields;
    },
    ...mapGetters({
      notificationCount: "getNotificationCount",
    }),
  },
};
</script>

<style lang="scss">
@import "../../assets/scss/variable.scss";
@import "../../assets/scss/mixin.scss";
#page-header {
  .top-header {
    .logo-img {
      max-height: 65px;
    }
    .dropdown-menu {
      max-width: 200px;
      border: none;
      position: relative;
      &:before {
        content: "";
        width: 0;
        height: 0;
        border-style: solid;
        border-width: 0 7.5px 8px 7.5px;
        border-color: transparent transparent #ffffff transparent;
        position: absolute;
        top: -8px;
        left: 18px;
      }
    }
  }

  .header-dropdown {
    .btn {
      @extend .btn-transparent;
    }
    .dropdown-item {
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .list-icon {
      margin-right: 3px;
    }
  }
  .notification-dropdown {
    position: relative;
    .notification-icon {
      display: inline-block;
      vertical-align: middle;
      font-size: 26px;
      margin-right: 5px;
    }
    .btn {
      .notification-icon {
        opacity: 0.7;
        color: #fff;
      }
      &:hover .notification-icon {
        opacity: 1;
      }
    }
    .notification-count {
      width: 18px;
      height: 18px;
      background: #fff;
      border-radius: 50%;
      display: block;
      position: absolute;
      top: 8px;
      left: 25px;
      color: $priColor;
      font-size: 12px;
      font-weight: bold;
      text-align: center;
      padding-top: 1px;
    }
  }
  .user-info-dropdown {
    .user-image-wrapper {
      position: relative;
      display: inline-block;
      &:before {
        content: "";
        width: 50px;
        height: 50px;
        position: absolute;
        left: -5px;
        top: -5px;
        border: 1px solid #fff;
        border-radius: 50%;
        opacity: 0.7;
      }
    }
    .user-image-inner {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      position: relative;
      display: inline-block;
      vertical-align: middle;
      margin-right: 5px;
      overflow: hidden;
    }
    .user-image {
      position: absolute;
      width: 100%;
      top: 0;
      left: 0;
    }
    .dropdown-menu-right:before {
      left: auto;
      right: 50px;
    }
  }

  // Navigation
  .navigation {
    margin-top: 0px;
    margin-bottom: 30px;
    .navbar {
      padding: 0;
    }
    .navbar-dark .navbar-nav .nav-link {
      color: rgba(255, 255, 255, 0.75);
      font-size: 15px;
      margin-right: 10px;
      margin-right: 15px;
      .list-icon {
        font-size: 18px;
      }
    }
    .navbar-dark .navbar-nav .nav-link:focus,
    .navbar-dark .navbar-nav .nav-link:hover {
      color: rgba(255, 255, 255, 0.9);
    }
    .navbar-dark .navbar-nav .active > .nav-link {
      font-weight: bold;
      color: #fff;
    }
    .search-input {
      position: relative;
      .form-control {
        border-radius: 50px;
        background: rgba(255, 255, 255, 0.1);
        border: rgba(255, 255, 255, 0.2);
        padding: 0.7rem 1.2rem;
        color: rgba(255, 255, 255, 0.9);
        padding-left: 40px;
        &::placeholder {
          color: rgba(255, 255, 255, 0.7);
          font-size: 14px;
        }
      }
      &:before {
        font-family: "Material Design Icons";
        content: "\F349";
        position: absolute;
        left: 18px;
        top: 50%;
        transform: translateY(-50%);
        color: rgba(255, 255, 255, 0.7);
      }
    }
  }
}
@media (max-width: 991px) {
  #page-header {
    .navigation {
      .navbar-toggler {
        margin-left: auto;
      }
      .navbar-dark .navbar-collapse > .navbar-nav {
        background: #fff;
        box-shadow: 0 0 5px rgba(#000, 0.1);
        border-radius: 15px;
        padding: 20px;
        margin-top: 20px;
        .nav-link {
          color: $darkGrey;
          padding: 10px 12px;
          border-radius: 4px;
          margin-right: 0;
          &:hover {
            background: $priColor;
            color: #fff;
          }
        }
        .active .nav-link {
          background: $priColor;
          color: #fff;
        }
        .b-nav-dropdown {
          .nav-link {
            &:focus,
            &:hover {
              color: $priColor;
              background: transparent;
            }
          }
          .dropdown-menu {
            background: #eee;
            box-shadow: none;
          }
        }
      }
    }
  }
}
@media (min-width: 992px) {
  .navbar-expand-lg .navbar-nav .nav-link {
    padding-right: 0.75rem;
    padding-left: 0.75rem;
  }
  #page-header {
    .navigation {
      margin-top: 50px;
      margin-bottom: 30px;
    }
    .top-header .logo-img {
      max-height: 85px;
    }
    .user-info-dropdown .user-image-wrapper {
      .user-image-inner {
        width: 60px;
        height: 60px;
        border-radius: 50%;
        position: relative;
        display: inline-block;
        vertical-align: middle;
        margin-right: 5px;
        overflow: hidden;
      }
      &:before {
        width: 70px;
        height: 70px;
      }
    }
  }

  #notificationModal {
    .notification-icon-wrapper {
      height: 100%;
      background: #eee;
      border-radius: 4px;
    }
    .notification-icon {
      font-size: 120px;
      color: $priColor;
      opacity: 0.5;
    }
    .modal-body {
      padding: 30px;
    }
  }
}
</style>
