<template>
  <div class="main-content">
    <section class="section-wrapper">
      <!-- setion-title -->
      <div class="row section-title">
        <div class="col-12">
          <h4>Upload an text file to import Questions</h4>
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
import { mapGetters } from "vuex";
export default {
  components: { vueDropzone: vue2Dropzone },
  data() {
    return {
      dropzoneOptions: {
        url: "https://httpbin.org/post",
        thumbnailWidth: 150,
        maxFilesize: 0.01,
        maxFiles: 1,
        //  acceptedFiles: ".xls,.xlsx",
        headers: { "My-Awesome-Header": "header value" },
      },
    };
  },
  methods: {
    submitFile(file) {
      const self = this;
      let formData = new FormData();

      /*
                Add the form data we need to submit
            */
      formData.append("questionFile", file);
      formData.append("examName", "ComputerBasics");
      formData.append("examType", "Descriptive");
      /*
          Make the request to the POST /single-file URL
        */
      this.$http
        .post("/api/r/uploadQuestionFile", formData, {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        })
        .then(function(res) {
          self.$refs.myVueDropzone.removeFile(file);
          if (res.data) {
            self.notifyStudents();
            self.$bvToast.toast(`Imported Successfully`, {
              title: "Success",
              variant: "success",
              autoHideDelay: 5000,
              solid: true,
              class: "toast",
            });
          } else {
            self.$refs.myVueDropzone.removeFile(file);
            self.$bvToast.toast(`File is too big to parse`, {
              title: "Failed",
              variant: "danger",
              autoHideDelay: 5000,
              solid: true,
              class: "toast",
            });
          }
        })
        .catch(function() {
          self.$refs.myVueDropzone.removeFile(file);
          self.$bvToast.toast(`Please Upload Proper File`, {
            title: "Failed",
            variant: "danger",
            autoHideDelay: 5000,
            solid: true,
            class: "toast",
          });
          console.log("FAILURE!!");
        });
    },
    notifyStudents() {
      this.socketConn.send(
        JSON.stringify({
          type: 1,
          body: { data: "Exam has started, Please check" },
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
};
</script>

<style lang="scss"></style>
