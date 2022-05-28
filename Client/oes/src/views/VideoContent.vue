<template>
  <div class="main-content">

<h2>Videos List</h2>

  <!-- <div v-for="(vid,index) in this.videosList"
  :key="index"> -->
  <video width="320" height="240" ref="videoRef" controls>
    <!-- <source :src="vid.VideoUrl" type="application/x-mpegURL"> -->
</video>
  <!-- </div> -->
  </div>
</template>
<script>
import axios from "axios";
import Hls from 'hls.js';

export default {
  data() {
    return {
      videosList: [{
          Name : "docker",
          VideoUrl:"https://bitdash-a.akamaihd.net/content/sintel/hls/playlist.m3u8",
          ThumbnailPath:"",
          Description:"Docker"
      }],
      hls: null
    };
  },
  methods:{
     getVideos(){
       let self = this
       axios
        .get("/getVideos")
        .then(function(res) {
          if (res.data) {
            self.videosList = res.data.videos;
          }
        })
        .catch(function() {
          console.log("FAILURE!!");
        });
     },
     playVideo(ref) {
  //      console.log(ref)
  //  let video = this.$refs[ref];
  //   this.hls.loadSource("http://127.0.0.1:8887/index.m3u8");
  //     console.log("37-----------",this.hls)
  //   this.hls.attachMedia(video);
  //   this.hls.on(Hls.Events.MANIFEST_PARSED, function () {
  //     console.log("41-----------")
  //     video.play();
  //   });
//           var video = document.getElementById('video');
//   if(Hls.isSupported()) {
//       console.log("44",url)
//     var hls = new Hls();
//     hls.loadSource(url);
//     hls.attachMedia(video);
//     hls.on(Hls.Events.MANIFEST_PARSED,function() {
//       video.play();
//   });
//  } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
//     video.src = url;
//     video.addEventListener('loadedmetadata',function() {
//       video.play();
//     });
//   }
     }
     },
 mounted() {
    //  this.getVideos();
    this.hls = new Hls();
    let video = this.$refs["videoRef"];
    this.hls.loadSource("http://127.0.0.1:8887/docker/index.m3u8");
    this.hls.attachMedia(video);
    this.hls.on(Hls.Events.MANIFEST_PARSED, function () {
      video.play();
    });
}
}
</script>
<style lang="scss">

table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(even) {
  background-color: #dddddd;
}

</style>
