<template>
  <div class="main-content">
    <canvas
      ref="canvas"
      @mousemove="draw"
      @mousedown="beginDrawing"
      @mouseup="stopDrawing"
      @mouseleave="cancelDrawing"
      :width="canvasWidth"
      :height="canvasHeight"
    />
    <div id="bar">
      <div class="bar-item">
        <button class="btn btn-theme mx-auto" @click="clearCanvas">
          clear
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  name: "Room",
  data() {
    return {
      x: 0,
      y: 0,
      canvasWidth: 0,
      canvasHeight: 0,
      isDrawing: false,
    };
  },
  methods: {
    drawLine(x1, y1, x2, y2) {
      let ctx = this.$refs.canvas.getContext("2d");
      ctx.beginPath();
      ctx.strokeStyle = "black";
      ctx.lineWidth = 7;
      ctx.moveTo(x1, y1);
      ctx.lineTo(x2, y2);
      ctx.stroke();
      ctx.closePath();
    },
    draw(e) {
      if (this.isDrawing) {
        this.drawLine(this.x, this.y, e.offsetX, e.offsetY);
        this.x = e.offsetX;
        this.y = e.offsetY;
      }
    },
    beginDrawing(e) {
      this.x = e.offsetX;
      this.y = e.offsetY;
      this.isDrawing = true;
    },
    stopDrawing(e) {
      if (this.isDrawing) {
        this.drawLine(this.x, this.y, e.offsetX, e.offsetY);
        this.isDrawing = false;
        this.socketConn.send(
          JSON.stringify({
            type: 3,
            body: this.$refs.canvas.toDataURL("image/png"),
            id: sessionStorage.getItem("clientId"),
          })
        );
      }
    },
    cancelDrawing() {
      this.isDrawing = false;
    },
    drawUpdate(url) {
      let image = new Image();
      let ctx = this.$refs.canvas.getContext("2d");
      image.onload = () => {
        ctx.drawImage(image, 0, 0);
      };
      image.src = url;
    },
    handleResize() {
      let state = this.$refs.canvas.toDataURL("image/png");
      this.canvasWidth = window.innerWidth;
      this.canvasHeight = window.innerHeight - 25;
      this.drawUpdate(state);
    },

    clearCanvas() {
      let ctx = this.$refs.canvas.getContext("2d");
      ctx.clearRect(0, 0, this.canvasWidth, this.canvasHeight);

      this.socketConn.send(
        JSON.stringify({
          type: 3,
          body: this.$refs.canvas.toDataURL("image/png"),
          id: sessionStorage.getItem("clientId"),
        })
      );
    },
  },
  computed: {
    ...mapGetters({
      socketConn: "getConn",
    }),
  },
  watch: {
    "$store.getters.getWhiteBoard": function () {
      this.drawUpdate(this.$store.getters.getWhiteBoard);
    },
    immediate: true,
  },
  mounted() {
    this.handleResize();
    window.addEventListener("resize", this.handleResize);
  },
};
</script>

<style scoped>
* {
  padding: 0;
  margin: 0;
}
#bar {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  color: white;
  height: 40px;
}
.bar-item {
  align-items: center;
  padding-left: 20px;
  padding-right: 20px;
}
</style>