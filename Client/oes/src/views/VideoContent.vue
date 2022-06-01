<template>
  <div class="main-content">

<h2>Videos List</h2>

<div class="row" v-for="(vid,index) in this.videosList"
  :key="index"> 
  <div class="col-4">
  <img :src="'../../../../'+ vid.thumbnailPath" :alt="vid.description" width="300" height="300" @click="playVideo(vid.videoUrl)">
    </div> 
  </div>
  </div>
</template>
<script>
import axios from "axios";

export default {
  data() {
    return {
      videosList: [],
  
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
     playVideo(url) {
      this.$router.push({  name: 'VideoPlay', params: { videoUrl : url } })
     }
     },
 mounted() {
     this.getVideos();
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
