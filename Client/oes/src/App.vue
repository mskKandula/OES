<template>
  <div id="apps" ref="apps">
    <transition name="fade">
      <router-view />
    </transition>
    <transition name="bounce">
      <button
        v-if="moveToDown"
        id="myBtn"
        @click="topFunction()"
        class="btn-primary btn btn-circle-l scroll-top-btn"
      >
        <span class="mdi mdi-chevron-up arrow-icon"></span>
      </button>
    </transition>
  </div>
</template>

<script>
export default {
  name: "App",
  components: {},
  data() {
    return {
      moveToDown: false
    };
  },
  created() {
    window.addEventListener("scroll", this.showScrollButton);
  },
  destroyed() {
    window.removeEventListener("scroll", this.showScrollButton);
  },
  methods: {
    test() {
      console.log("test2");
    },
    topFunction() {
      let to = this.moveToDown ? this.$refs.apps.offsetTop - 60 : 0;
      window.scroll({
        top: to,
        left: 0,
        behavior: "smooth"
      });
      //   this.moveToDown = !this.moveToDown;
    },
    showScrollButton() {
      var scroll = window.pageYOffset;
      if (scroll > 300) {
        this.moveToDown = true;
      } else {
        this.moveToDown = false;
      }
    }
  },
  mounted() {
    // this.showScrollButton();
  }
};
</script>

<style lang="scss">
#apps {
  .scroll-top-btn {
    position: fixed;
    right: 50px;
    bottom: 50px;
    transform: translateY(0px);
    // border: 4px solid #fff;
    // box-shadow: 0px 4px 15px rgba(0, 0, 0, 0.4);
    .arrow-icon {
      font-size: 24px;
      line-height: 20px;
    }
  }
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}
.bounce-enter-active {
  animation: slide-in 0.5s;
}
.bounce-leave-active {
  animation: slide-in 0.5s reverse;
}
@keyframes slide-in {
  0% {
    transform: translateY(0);
  }

  100% {
    transform: translateY(30px);
  }
}
</style>
