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
  </div>
</template>

<script>
// import questionsData from "./questions.json";
import SimpleKeyboard from "./SimpleKeyboard";
import "./App.css";
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
        .get("/api/getQuestions")
        .then(function (res) {
          if (res.data) {
            self.questions = res.data.questions;
          }
        })
        .catch(function () {
          console.log("FAILURE!!");
        });
    },
  },
  mounted() {
    this.getQues();
  },
};
</script>

<style>
</style>