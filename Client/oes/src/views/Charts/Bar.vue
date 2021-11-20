<template>
  <div class="main-content">
    <h3>Bar Chart</h3>
    <bar-chart :data12 = this.data12></bar-chart>
  </div>
</template>

<script>
import BarChart from '@/components/charts/BarChart'
import { mapGetters } from 'vuex'
export default {
  components: {
    BarChart
  },
   data() {
      return {
     data12 :[]
      }
   },
   methods:{

     getData12(){
        this.socketConn.send(JSON.stringify({
        type: 2,
        body: "Fetch Chart Data",
        id : "6666" }))

        this.socketConn.onmessage = (evt) =>{
          let x = JSON.parse(evt.data)
          this.data12 = x.Intarr
        }
      
     }
   },
    computed:{
    ...mapGetters({
      socketConn: 'getConn'
    }),
  },
    created() {
       this.getData12()
  },
}
</script>