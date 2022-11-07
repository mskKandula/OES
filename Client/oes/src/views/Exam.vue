<template>
  <div class="main-content">
    <div id="examId">
      <h1>Online Exam</h1>
      <hr />
      <button class="btn btn-primary" v-on:click="key = !key">
        Virtual keyboard</button
      >&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
      <button class="btn btn-primary" @click="startTxtToSpeech">Speak</button>
      <!-- <div v-for="(word, index) in transcription_" :key="index">
        {{ word }}
        
      </div> -->

      <div v-for="(question, index) in questions" :key="index">
        <div v-show="index === questionIndex">
          <h4>{{ questionIndex + 1 }}:{{ question }}</h4>
          <textarea
            class="form-control"
            id="exampleFormControlTextarea1"
            rows="3"
            :value="input"
            @input="onInputChange"
          ></textarea>

          <br />

          <button
            class="btn btn-theme btn-rounded mx-auto"
            v-if="questionIndex > 0"
            v-on:click="prev"
          >
            prev
          </button>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
          <button
            class="btn btn-theme btn-rounded mx-auto"
            v-if="questionIndex < questions.length - 1"
            v-on:click="next"
          >
            next
          </button>
          <button
            class="btn btn-theme btn-rounded mx-auto"
            v-if="index == questions.length - 1"
            @click="submitAnswer"
          >
            Submit
          </button>
        </div>
      </div>
    </div>
    <SimpleKeyboard
      v-if="this.key"
      @onChange="onChange"
      @onKeyPress="onKeyPress"
      :input="input"
    />
    <canvas id="canvas" hidden></canvas>
    <video
      class="center"
      height="500px"
      controls
      autoplay
      id="video"
      hidden
    ></video>
  </div>
</template>

<script>
// import questionsData from "./questions.json";
import SimpleKeyboard from "./SimpleKeyboard";
import "./App.css";
import JSZip from "jszip";
export default {
  components: {
    SimpleKeyboard,
  },
  data() {
    return {
      runtimeTranscription_: "",
      transcription_: [],
      lang_: "en_US",
      input: "",
      key: false,
      questions: [],
      questionIndex: 0,
      answer: [],
      test: "",
      current: "",
      res: "",
      displayOptions: {
        video: {
          cursor: "always",
        },
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          sampleRate: 44100,
        },
      },
      videoTrack: null,
      blobsArray: [],
      new_zip: null,
    };
  },
  methods: {
    next: function () {
      this.questionIndex++;
      if (this.input.length > 0) {
        this.answer.push(this.input);
        this.input = "";
      } else {
        this.input = "";
        this.answer.push(this.input);
      }
    },
    // Go to previous question
    prev: function () {
      this.questionIndex--;
    },
    submitAnswer() {
      this.answer.push(this.input);
      this.res = this.$go.countser(JSON.stringify({ ans: this.answer }));
      if (this.res) {
        this.$router.push({ name: "WordCounter", params: { count: this.res } });
      }
    },
    onChange(input) {
      this.input = input;
    },
    onInputChange(input) {
      this.input = input.target.value;
    },
    startTxtToSpeech() {
      // initialisation of voicereco
      window.SpeechRecognition =
        window.SpeechRecognition || window.webkitSpeechRecognition;
      const recognition = new window.SpeechRecognition();
      recognition.lang = this.lang_;
      recognition.interimResults = true;
      // event current voice reco word
      recognition.addEventListener("result", (event) => {
        //  console.log("event",event)
        this.current = event.results[0][0].transcript;
        // console.log("this.current",this.current)
        // const text = Array.from(event.results)
        //   .map(result => result[0])
        //   .map(result => result.transcript)
        //   .join("");
        // this.runtimeTranscription_ = text;
        // console.log("this.runtimeTranscription_",this.runtimeTranscription_)
        // console.log("text",text)
      });
      // end of transcription
      recognition.addEventListener("end", () => {
        // this.transcription_.push(this.runtimeTranscription_);
        // console.log("this.transcription_",this.transcription_)
        this.input += this.current;
        this.runtimeTranscription_ = "";
        recognition.stop();
      });
      recognition.start();
    },
    startSpeechToTxt() {
      // start speech to txt
      var utterance = new SpeechSynthesisUtterance("Message Envoyé");
      window.speechSynthesis.speak(utterance);
    },
    getQues() {
      let self = this;
      this.$http
        .get("/api/r/getQuestions")
        .then(function (res) {
          if (res.data) {
            self.questions = res.data.questions;
          }
        })
        .catch(function () {
          console.log("FAILURE!!");
        });
    },
    async getStream() {
      const self = this;
      navigator.mediaDevices
        // Uncomment below line to capture screen
        // .getUserMedia(self.displayOptions)
        .getUserMedia(self.displayOptions)
        .then(async (stream) => {
          let vid = document.getElementById("video");
          vid.srcObject = stream;

          // Grab frame from stream
          self.videoTrack = stream.getVideoTracks()[0];

          // self.getscShot();
          for (let i = 0; i < 10; i++) {
            self.captureImage();
            await new Promise((r) => setTimeout(r, 5000));
          }

          self.uploadZip();
        })

        .catch((e) => console.log(e));
    },
    captureImage() {
      const self = this;
      self.imageCapture = new ImageCapture(self.videoTrack);
      self.imageCapture.grabFrame().then((bitmap) => {
        // Stop sharing
        // track.stop();
        let canvas = document.getElementById("canvas");
        // Draw the bitmap to canvas
        canvas.width = bitmap.width;
        canvas.height = bitmap.height;
        canvas.getContext("2d").drawImage(bitmap, 0, 0);

        // Grab blob from canvas
        canvas.toBlob((blob) => {
          // Do things with blob here
          blob.name = `studentId-${new Date().getTime()}.png`;

          self.blobsArray.push(blob);

          // To Display on Screen
          // let image = document.createElement("img");
          // image.setAttribute("style", "width: 150px; height: 150px;");
          //image.height="15px"
          // let url = window.URL.createObjectURL(blob);
          // image.src = url;
          // document.body.appendChild(image);

          // To Download the images
          //  let a = document.createElement("a");
          // let url = window.URL.createObjectURL(blob);
          // a.href = url;
          // a.download = blob.name;
          // a.click();
          //  window.URL.revokeObjectURL(url);
        });
      });
    },
    uploadZip() {
      const self = this;
      for (let i = 0; i < self.blobsArray.length; i++) {
        self.new_zip.file(self.blobsArray[i].name, self.blobsArray[i], {
          binary: true,
        });
      }
      self.new_zip
        .generateAsync({
          type: "blob",
        })
        .then(function (content) {
          // To download a Zip File
          // var a = document.createElement("a");
          // let url = window.URL.createObjectURL(content);
          // a.href = url;
          // a.download = "img_archives.zip";
          // a.click();
          // window.URL.revokeObjectURL(url);

          //generated zip content to file type
          var files = new File([content], "studentImgCapture.zip");

          var formData = new FormData();
          formData.append("zipFile", files);

          self.$http
            .post("/api/r/uploadExamProof", formData, {
              headers: {
                "Content-Type": "multipart/form-data",
              },
            })
            .then(function (res) {
              console.log(res);
              self.blobsArray = [];
            })
            .catch(function (e) {
              console.log("Failed to upload", e);
            });
        });
    },
  },
  mounted() {
    this.new_zip = new JSZip();
    this.getStream();
    this.getQues();
  },
};
</script>

<style>
</style>