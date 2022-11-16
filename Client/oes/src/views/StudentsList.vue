<template>
  <div class="main-content">
    <h2>Students List</h2>

    <table>
      <tr>
        <th>Id</th>
        <th>Name</th>
        <th>Email</th>
        <th>Mobile</th>
      </tr>
      <tr v-for="(student, index) in this.studentsList" :key="index">
        <td>{{ index + 1 }}</td>
        <td>{{ student.name }}</td>
        <td>{{ student.email }}</td>
        <td>{{ student.mobile }}</td>
      </tr>
    </table>
  </div>
</template>
<script>
export default {
  data() {
    return {
      studentsList: [],
    };
  },
  methods: {
    getstudents() {
      let self = this;
      this.$http
        .get("/api/r/getStudents")
        .then(function (res) {
          if (res.data) {
            self.studentsList = res.data.students;
          }
        })
        .catch(function () {
          console.log("FAILURE!!");
        });
    },
  },
  created() {
    this.studentsList = this.$route.params.studentsList;
  },
  mounted() {
    this.getstudents();
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
