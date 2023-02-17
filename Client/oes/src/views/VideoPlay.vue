<template>
  <div>
    <video ref="videoRef" controls></video>
  </div>
</template>
<script>
import Hls from "hls.js";

export default {
  data() {
    return {
      hls: null,
      videoUrl: "",
    };
  },
  created() {
    this.videoUrl = this.$route.params.videoUrl ? this.$route.params.videoUrl : "/media/video/liveData/test.m3u8" ;
  },
  mounted() {
    this.hls = new Hls();
    let video = this.$refs["videoRef"];
    let url = "/cdn" + this.videoUrl;
    this.hls.loadSource(url);
    this.hls.attachMedia(video);
    this.hls.on(Hls.Events.MANIFEST_PARSED, function () {
      video.play();
    });
  },
};
</script>
<style lang="scss">
table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%;
}

td,
th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(even) {
  background-color: #dddddd;
}
</style>
