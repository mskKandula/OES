<template>
  <div class="main-content">
    <h1 class="h2">Dashboard</h1>
    <h3 class="text-watermark-lg p-5 m-5 text-center">Welcome User</h3>
  </div>
</template>
<script>
export default {
  components: {},
  data() {
    return {
      ws: null,
    };
  },
  methods: {
    mutateData(data) {
      let parsedData = JSON.parse(data);
      let dataBody = parsedData.body;

      if (parsedData.type == 1) {
        this.$store.commit("setNotification", dataBody);
      } else if (parsedData.type == 2) {
        this.$store.commit("setChat", dataBody);
      } else if (parsedData.type == 3) {
        this.$store.commit("setWhiteBoard", dataBody);
      } else if (parsedData.type == 4) {
        this.$store.commit("setBroadcast", dataBody);
      } else if (parsedData.type == 5) {
        this.$store.commit("setOnlineUsers", dataBody);
      } else {
        this.$store.commit("setNotificationCount", 1);
      }
    },
  },
  created() {
    const onlyOnce = this.$store.getters.getOnce;
    if (onlyOnce) {
      // let roleName = this.$store.getters.getUserRole
      // let clientId = this.$store.getters.getClientId
      const url = new URL("ws://localhost:9000/ws");
      url.searchParams.append("role", "User");
      url.searchParams.append("id", sessionStorage.getItem("clientId"));
      url.searchParams.append("userId", sessionStorage.getItem("userId"));

      this.ws = new WebSocket(url.href);

      this.$store.commit("setConn", this.ws);
      // this.$store.commit('setFalse', false)
      this.$store.commit("setOnlyOnce", false);

      this.ws.onconnect = (evt) => {
        console.log("ws connected", evt);
      };

      this.ws.onclose = (evt) => {
        console.log("ws closed", evt);
      };

      this.ws.onmessage = (evt) => {
        let data = evt.data;

        data = data.split(/\r?\n/);

        this.mutateData(data[0]);

        data = "";
      };
    }
  },
};
</script>
<style lang="scss"></style>
