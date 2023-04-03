<template>
  <div>
    <video id="video-player" controls></video>
  </div>
</template>
<script>
export default {
  data() {
    return {
      videoUrl: "",
    };
  },
  created() {
    this.videoUrl = this.$route.params.videoUrl
      ? this.$route.params.videoUrl
      : "/media/video/liveData/test.m3u8";
  },
  methods: {
    callToPlay() {
      alert("25");
      const videoElement = document.getElementById("video-player");
      const mediaSource = new MediaSource();
      videoElement.src = window.URL.createObjectURL(mediaSource);

      // Wait for the MediaSource to open
      mediaSource.addEventListener("sourceopen", () => {
        const mimeType =
          'video/mp4; codecs="avc1.64001f, avc1.4d002a,  mp4a.40.2, avc1.4D401F, avc1.42E01E"';
        // const mimeType = "video/mp4";
        // Create a new SourceBuffer for the video
        const sourceBuffer = mediaSource.addSourceBuffer(mimeType);

        const url = "/cdn" + this.videoUrl;
        // Load the main video
        fetch(url)
          .then((response) => {
            console.log("34", response);
            return response.text();
          })
          .then((m3u8Text) => {
            console.log("38", m3u8Text);

            const urls = m3u8Text
              .split("\n")
              .filter((line) => line.trim().length > 0 && line[0] !== "#");

            let index = 0;
            let isPlay = true;
            const fetchSegment = () => {
              let indexUrl =
                "/cdn/media/video/70b3c009c99f46419a95ea5a84107f38/golang/" +
                urls[index];
              console.log("47", indexUrl);
              fetch(indexUrl)
                .then((response) => {
                  console.log("50", response);
                  return response.arrayBuffer();
                })
                .then((segmentBuffer) => {
                  console.log("54", segmentBuffer);
                  sourceBuffer.appendBuffer(segmentBuffer);
                  if (isPlay) {
                    videoElement.play();
                    isPlay = false;
                  }
                  index++;

                  if (index < urls.length) {
                    fetchSegment();
                  }
                })
                .catch((error) => console.error(error));
            };

            fetchSegment();
          });
      });
    },
  },
  mounted() {
    this.callToPlay();
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
