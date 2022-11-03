<template>
  <video
    ref="pub_video"
    class="bg-black"
    autoplay
    controls
    v-if="this.publisher"
  ></video>
  <video id="sub" ref="sub_video" class="bg-black" controls v-else></video>
</template>

<script>
export default {
  props: ["publisher"],
  name: "VideoCall",
};
</script>
