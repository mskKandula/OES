<template>
  <div class="main-content">
    <section class="section-wrapper">
      <!-- setion-title -->
      <div class="row section-title">
        <div class="col-12">
          <h4>Upload a Text File To Generate Questions</h4>
        </div>
      </div>
      <div class="row section-body">
        <div class="col-12">
          <!-- section content -->
          <div class="form-group">
            <label for>Drop file</label>
            <vue-dropzone
              ref="myVueDropzone"
              id="dropzone"
              :options="dropzoneOptions"
              @vdropzone-complete="submitFile"
            ></vue-dropzone>
          </div>
        </div>
      </div>
    </section>
    <!-- <button
    @click="test">submitttt</button> -->
  </div>
</template>

<script>
// import Vue from 'vue';
import vue2Dropzone from "vue2-dropzone";
import "vue2-dropzone/dist/vue2Dropzone.min.css";
export default {
  components: { vueDropzone: vue2Dropzone },
  data() {
    return {
      dropzoneOptions: {
        url: "https://httpbin.org/post",
        thumbnailWidth: 150,
        maxFilesize: 0.01,
        maxFiles: 1,
        acceptedFiles: ".txt",
        headers: { "My-Awesome-Header": "header value" },
      },
      studentsList: []
    };
  },
  methods: {
    async submitFile(file) {
      let self = this;
      let textData = await self.processFile(file)
      let data={
        paragraph:textData
        }
      self.$http
        .post("/api/r/questionGeneration",data)
        .then(function (res) {
          console.log(res);
        })
        .catch(function () {
          console.log("FAILURE!!");
        });
    },
     processFile(file){
       return new Promise((resolve) => {
        var fr=new FileReader();
            fr.onload=function(){
                resolve(fr.result)
            }
            fr.readAsText(file);
       })
    }
  },
};
</script>

<style lang="scss"></style>
