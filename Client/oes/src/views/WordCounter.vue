<template>
  <div class="main-content">
    <h1 class="h2">Counter</h1>
    <h5 class="p-5 m-5 text-center">
      Count of Words that exist in the Dictionary :
    </h5>
    <h2>{{ this.count }}</h2>
  </div>
</template>
<script>
export default {
  data() {
    return {
      count: 0,
    };
  },
  created() {
    this.count = this.$route.params.count;
  },
};
</script>
<style lang="scss"></style>
