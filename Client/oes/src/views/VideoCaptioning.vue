<template>
  <div>
    
  <video ref="videoRef" controls src="../assets/docker.mp4" type="video/mp4">
</video>
 
  </div>
</template>
<script>

export default {
  components: {
  },
  data() {
    return {
      recognition:null,
      captionRecord:null,
      interim_transcript:null,
      recognizing:false,
      captionRecord:true

    }
  },
  methods:{
      
checkSpeech(){
       if('webkitSpeechRecognition' in window){
            console.log("speech recognition supported");
            this.recognition = new webkitSpeechRecognition();
        }else{
            console.log("speech this.recognition not supported");
            this.captionRecord = false;
        }
        console.log("captionRecord", this.captionRecord);
        this.recognition.interimResults = true;
        this.recognition.continuous = true;
         this.recognition.onstart = function() {
                this.recognizing = true;
            };

            this.recognition.onerror = function(event) {
                console.log ("there was a captioning error");
            };

            this.recognition.onend = function() {
                console.log ("captioning stopped");
                this.recognizing = false;
                
            };

            this.recognition.onresult = function(event) {
                //heres where I'd put where stuff goes in my app....
                for (var i = event.resultIndex; i < event.results.length; ++i) {
                if (event.results[i].isFinal) {
                    this.interim_transcript = "";
                } else {          
                        //append the words
                        this.interim_transcript = event.results[i][0].transcript;
                    console.log("74",this.interim_transcript);
                }
                }
            };
}
  },
  created() {
    this.checkSpeech()
  }
}
</script>