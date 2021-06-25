<template>
  <div class="main-content">

<h2>Candidate List</h2>

<table>
  <tr>
    <th>Id</th>
    <th>Name</th>
    <th>Email</th>
    <th>Mobile</th>
  </tr>
  <tr v-for="(candidate,index) in this.candidateList"
  :key="index">
    <td>{{index +1 }}</td>
    <td>{{candidate.name}}</td>
    <td>{{candidate.email}}</td>
    <td>{{candidate.mobile}}</td>
  </tr>
</table>
  </div>
</template>
<script>
import axios from "axios";
export default {
  data() {
    return {
      candidateList: [
    //    {
    //        id : 1,
    //        name: "Mohana Kandula",
    //        email:"mskKandula@gmail.com",
    //        mobile:"9553563699"
    //    }
      ]
    };
  },
  methods:{
     getstudents(){
       axios
        .get("/getStudents")
        .then(function(res) {
          console.log("155", res.data);
          if (res.data) {
            self.candidateList = res.data;
            console.log("158", self.candidateList);
          }
        })
        .catch(function() {
          console.log("FAILURE!!");
        });
     }
     },
     created(){
        this.candidateList = this.$route.params.CandidateList
     }
//  mounted() {
//      this.getstudents();
// }
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
