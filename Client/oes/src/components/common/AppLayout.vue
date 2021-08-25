<template>
  <div class="app-container theme2">
    <div v-if="false" class="loader-wrapper">
      <div class="lds-roller">
        <div></div>
        <div></div>
        <div></div>
        <div></div>
      </div>
    </div>
    <header1></header1>
    <!-- main container -->
    <div class="main container">
      <transition name="fade">
        <router-view />
      </transition>
    </div>
    <footer class="footer">
      <p class="copyright">
        &copy; Copyright 2021
        <a href="" taget="_blank"
          >mskKandula.</a
        >
        All Rights Reserved.
      </p>
    </footer>
    <!-- <div class="menu-overlay" @click="toggleSidebar2()"></div> -->
  </div>
</template>
<script>
import header1 from "./Header1.vue";
export default {
  name: "AppLayout",
  components: {
    header1
  },
  data() {
    return {};
  },
  methods: {}
};
</script>
<style lang="scss">
@import "../../assets/scss/variable.scss";
@import "../../assets/scss/mixin.scss";
.footer {
  padding: 30px 40px;
  text-align: center;
  color: rgba(#fff, 0.7);
  font-size: 14px;
  margin-top: 100px;
  .copyright {
    margin-bottom: 0;
  }
  a {
    color: rgba(#fff, 0.9);
    &:hover {
      color: rgba(#fff, 1);
    }
  }

  &:before {
    content: "";
    width: 150%;
    height: 150px;
    background: linear-gradient(135deg, $darkGrey 0%, #34495e 100%);
    position: absolute;
    bottom: -50px;
    left: 50%;
    transform: translate(-50%) rotate(-3deg);
    z-index: -1;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}

.loader-wrapper {
  position: fixed;
  background: $priColor;
  opacity: 0.7;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  z-index: 9;
}
.lds-roller {
  display: inline-block;
  position: relative;
  width: 80px;
  height: 80px;
}
.lds-roller div {
  animation: lds-roller 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
  transform-origin: 40px 40px;
}
.lds-roller div:after {
  content: " ";
  display: block;
  position: absolute;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #fff;
  margin: -4px 0 0 -4px;
}
.lds-roller div:nth-child(1) {
  animation-delay: -0.036s;
}
.lds-roller div:nth-child(1):after {
  top: 63px;
  left: 63px;
}
.lds-roller div:nth-child(2) {
  animation-delay: -0.072s;
}
.lds-roller div:nth-child(2):after {
  top: 68px;
  left: 56px;
}
.lds-roller div:nth-child(3) {
  animation-delay: -0.108s;
}
.lds-roller div:nth-child(3):after {
  top: 71px;
  left: 48px;
}
.lds-roller div:nth-child(4) {
  animation-delay: -0.144s;
}
.lds-roller div:nth-child(4):after {
  top: 72px;
  left: 40px;
}
.lds-roller div:nth-child(5) {
  animation-delay: -0.18s;
}
.lds-roller div:nth-child(5):after {
  top: 71px;
  left: 32px;
}
.lds-roller div:nth-child(6) {
  animation-delay: -0.216s;
}
.lds-roller div:nth-child(6):after {
  top: 68px;
  left: 24px;
}
.lds-roller div:nth-child(7) {
  animation-delay: -0.252s;
}
.lds-roller div:nth-child(7):after {
  top: 63px;
  left: 17px;
}
.lds-roller div:nth-child(8) {
  animation-delay: -0.288s;
}
.lds-roller div:nth-child(8):after {
  top: 56px;
  left: 12px;
}
@keyframes lds-roller {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
