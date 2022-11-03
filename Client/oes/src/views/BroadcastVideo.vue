<template>
  <div class="flex flex-col h-screen">
    <videomenu :publisher="this.isPub" @event="startPublish" />
    <VideoCall ref="videoElements" :publisher="this.isPub" />
  </div>
</template>

<script>
import videomenu from "../components/VideoMenu.vue";
import VideoCall from "../components/VideoCall.vue";

import { LocalStream, Client } from "ion-sdk-js";
import { IonSFUJSONRPCSignal } from "ion-sdk-js/lib/signal/json-rpc-impl";

const config = {
  iceServers: [
    {
      urls: "stun:stun.l.google.com:19302",
    },
  ],
};

export default {
  data() {
    return {
      isPub: false,
      client: null,
    };
  },
  mounted() {
    const URL = new URLSearchParams(window.location.search).get("publish");
    console.log("url", URL);
    if (URL) {
      this.isPub = true;
    }
  },
  created() {
    let self = this;
    const signal = new IonSFUJSONRPCSignal("ws://localhost:7000/ws");
    self.client = new Client(signal, config);
    signal.onopen = () => self.client.join("test room");
    if (!this.isPub) {
      self.client.ontrack = (track, stream) => {
        console.log("got track: ", track.id, "for stream: ", stream.id);
        track.onunmute = () => {
          console.log("unmute");
          self.$refs.videoElements.$refs.sub_video.srcObject = stream;
          self.$refs.videoElements.$refs.sub_video.autoplay = true;
          self.$refs.videoElements.$refs.sub_video.muted = false;

          // when the publisher leave
          stream.onremovetrack = () => {
            self.$refs.videoElements.$refs.sub_video.srcObject = null;
          };
        };
      };
    }
  },
  methods: {
    startPublish(type) {
      if (type) {
        LocalStream.getUserMedia({
          resolution: "vga",
          audio: true,
          codec: "vp8",
        })
          .then((stream) => {
            this.$refs.videoElements.$refs.pub_video.srcObject = stream;
            this.client.publish(stream);
          })
          .catch(console.error);
      } else {
        LocalStream.getDisplayMedia({
          resolution: "vga",
          video: true,
          audio: true,
          codec: "vp8",
        })
          .then((stream) => {
            this.$refs.videoElements.$refs.pub_video.srcObject = stream;
            this.client.publish(stream);
          })
          .catch(console.error);
      }
    },
  },
  name: "BroadcastVideo",
  components: {
    videomenu,
    VideoCall,
  },
};
</script>

<style>
@import "https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css";
</style>
