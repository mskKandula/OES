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
  created() {
    const onlyOnce = this.$store.getters.getOnce;
    console.log("21", onlyOnce);
    if (onlyOnce) {
      // let roleName = this.$store.getters.getUserRole
      // let clientId = this.$store.getters.getClientId
      const url = new URL("ws://localhost:9000/ws");
      url.searchParams.append("role", "User");
      url.searchParams.append("id", "6666");

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

        alert(JSON.parse(data[0]).body);

        data = "";
      };
    }
  },
};
</script>
<style lang="scss"></style>
