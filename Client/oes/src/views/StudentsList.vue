<template>
  <div class="main-content">
    <h2>Students List</h2>

    <table>
      <tr>
        <th>Id</th>
        <th>Name</th>
        <th>Email</th>
        <th>Mobile</th>
        <th>Online</th>
        <th>Action</th>
      </tr>
      <tr v-for="(student, index) in this.studentsList" :key="index">
        <td>{{ student.id }}</td>
        <td>{{ student.name }}</td>
        <td>{{ student.email }}</td>
        <td>{{ student.mobile }}</td>
        <td>{{ findStudentId(student.id) }}</td>
        <td>
          <b-button
            id="requestBtn"
            type="submit"
            name="btnRequest"
            varient="link"
            class="btn btn-theme btn-rounded mx-auto"
            @click="requestVideo(student.id)"
          >
            Call
          </b-button>
        </td>
      </tr>
    </table>

    <video autoPlay id="partnerVideo"></video>
  </div>
</template>
<script>
import { mapGetters } from "vuex";
export default {
  data() {
    return {
      studentsList: [],
      peerRef: null,
      studentId: 0,
    };
  },
  methods: {
    handleTrackEvent(e) {
      console.log("Received Tracks", e);
      const remoteVideo = document.getElementById("partnerVideo");
      remoteVideo.srcObject = e.streams[0];
      // this.partnerVideo = e.streams[0];
    },
    createPeer() {
      return new Promise((resolve) => {
        console.log("Creating Peer Connection");
        const peer = new RTCPeerConnection({
          iceServers: [
            {
              urls: [
                "stun:stun.l.google.com:19302",
                "stun:stun2.l.google.com:19302",
              ],
            },
          ],
        });
        peer.onicecandidate = this.handleIceCandidateEvent;
        peer.ontrack = this.handleTrackEvent;
        resolve(peer);
      });
    },
    handleIceCandidateEvent(e) {
      console.log("Found Ice Candidate");
      if (e.candidate) {
        this.socketConn.send(
          JSON.stringify({
            type: 4,
            body: {
              iceCandidate: e.candidate,
            },
            id: sessionStorage.getItem("clientId"),
            to: this.studentId,
          })
        );
      }
    },
    async handleOffer(offer) {
      console.log("Received Offer, Creating Answer");
      this.peerRef = await this.createPeer();

      await this.peerRef.setRemoteDescription(new RTCSessionDescription(offer));

      const answer = await this.peerRef.createAnswer();
      await this.peerRef.setLocalDescription(answer);

      this.socketConn.send(
        JSON.stringify({
          type: 4,
          body: {
            answer: this.peerRef.localDescription,
          },
          id: sessionStorage.getItem("clientId"),
          to: this.studentId,
        })
      );
    },
    getstudents() {
      let self = this;
      this.$http
        .get("/api/r/getStudents")
        .then(function(res) {
          if (res.data) {
            self.studentsList = res.data.students;
          }
        })
        .catch(function() {
          console.log("FAILURE!!");
        });
    },
    findStudentId(id) {
      return this.onlineUsers.includes(id);
    },
    requestVideo(id) {
      this.studentId = id;
      this.socketConn.send(
        JSON.stringify({
          type: 4,
          body: { join: true },
          id: sessionStorage.getItem("clientId"),
          to: id,
        })
      );
    },
  },
  created() {
    this.studentsList = this.$route.params.studentsList;
  },
  mounted() {
    this.getstudents();
  },
  watch: {
    message: {
      deep: true,
      handler(newmessage) {
        if (newmessage.offer) {
          this.handleOffer(newmessage.offer);
        } else if (newmessage.iceCandidate) {
          console.log("Receiving and Adding ICE Candidate");
          try {
            this.peerRef.addIceCandidate(newmessage.iceCandidate);
          } catch (err) {
            console.log("Error Receiving ICE Candidate", err);
          }
        }
      },
    },
  },
  computed: {
    ...mapGetters({
      onlineUsers: "getOnlineUsers",
      socketConn: "getConn",
      message: "getBroadcast",
    }),
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
