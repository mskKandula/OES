<template>
  <div class="main-content">

<h2>Student List</h2>

<table>
  <tr>
    <th>Id</th>
    <th>Name</th>
    <th>Email</th>
    <th>Mobile</th>
  </tr>
  <tr v-for="(student,index) in this.studentList"
  :key="index">
    <td>{{index +1 }}</td>
    <td>{{student.name}}</td>
    <td>{{student.email}}</td>
    <td>{{student.mobile}}</td>
  </tr>
</table>
  </div>
</template>
<script>
import axios from "axios";
export default {
  data() {
    return {
      studentList: []
    };
  },
  methods:{
     getstudents(){
       axios
        .get("/getStudents")
        .then(function(res) {
          if (res.data) {
            self.studentList = res.data;
          }
        })
        .catch(function() {
          console.log("FAILURE!!");
        });
     }
     },
     created(){
        this.studentList = this.$route.params.studentList
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
