<template>
  <div class="main-content">

<h2>Videos List</h2>

  <div v-for="(video,index) in this.videosList"
  :key="index">{{video}}
  <video id="video" width="320" height="240" controls @click="playVideo(video.VideoUrl)">{{video.VideoUrl}}
</video>
  </div>
  </div>
</template>
<script>
import axios from "axios";
export default {
  data() {
    return {
      videosList: [{
          Name : "docker",
          VideoUrl:"/assets/Videos/docker/index.m3u8",
          ThumbnailPath:"",
          Description:"Docker"
      }]
    };
  },
  methods:{
     getVideos(){
       let self = this
       axios
        .get("/getVideos")
        .then(function(res) {
          if (res.data) {
            console.log("26",res)
            self.videosList = res.data.videos;
          }
        })
        .catch(function() {
          console.log("FAILURE!!");
        });
     },
     playVideo(url) {
          var video = document.getElementById('video');
  if(Hls.isSupported()) {
      console.log("44",url)
    var hls = new Hls();
    hls.loadSource(url);
    hls.attachMedia(video);
  //   hls.on(Hls.Events.MANIFEST_PARSED,function() {
  //     video.play();
  // });
 } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
    video.src = url;
    // video.addEventListener('loadedmetadata',function() {
    //   video.play();
    // });
  }
     }
     },
 mounted() {
    //  this.getVideos();
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
