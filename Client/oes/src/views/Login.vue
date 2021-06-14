<template>
  <section id="login">
    <div class="container-fluid login-content">
      <div class="row justify-content-center no-gutters">
        <div class="col-lg-5 d-none d-lg-block">
          <div class="card custom-card login-image-wrapper">
            <div class="card-body">
              <swiper :options="loginSlider" ref="loginSlider">
                <!-- slides -->
                <swiper-slide>
                  <div class="login-image-inner">
                    <img
                      class="login-image"
                      src="../assets/images/login-img.jpg"
                      alt
                    />
                  </div>
                  <div class="slider-content">
                    Learning Management System
                  </div>
                </swiper-slide>
              </swiper>
              <div class="swiper-control">
                <!-- Optional controls -->
                <div class="swiper-prev" slot="button-prev">
                  <span class="mdi mdi-chevron-left" />
                </div>
                <div class="swiper-next" slot="button-prev">
                  <span class="mdi mdi-chevron-right" />
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="col-12 col-sm-8 col-lg-4">
          <div class="card custom-card login-box-wrapper">
            <div class="card-body">
              <div class="logo mb-4">
                <img
                  class="logo-img"
                  src="../assets/images/logo_mkcl.svg"
                  alt="logo here"
                />
              </div>
              <div id="loginFormId" class="form-group">
                <span class="mdi mdi-account-outline login-icon" />
                <input
                  v-model="log.email"
                  id="email"
                  type="email"
                  name="email"
                  class="form-control"
                  aria-describedby="emailHelp"
                  placeholder="Email ID"
                />
                <!-- <small
                  id="emailHelp"
                  class="form-text text-muted"
                >We'll never share your email with anyone else.</small>-->
              </div>
              <div class="form-group">
                <span class="mdi mdi-lock-outline login-icon" />
                <input
                  v-model="log.password"
                  type="password"
                  name="password"
                  class="form-control"
                  @focus="showPassInfo = true"
                  id="password"
                  @blur="showPassInfo = false"
                  placeholder="Password"
                />
              </div>

              <!-- <button
                id="loginBtn"
                type="submit"
                name="btnlogin"
                class="btn btn-info px-4 mr-2"
                @click="authenticate()"
              >Submit</button>
              <button
                type="button"
                name="btnlogin"
                class="btn btn-info px-4"
                @click="show()"
              >Show</button>-->
              <div class="row text-center mt-4">
                <!-- <button
                  id="loginBtn"
                  type="submit"
                  name="btnlogin"
                  class="btn btn-theme btn-rounded mx-auto"
                  @click="authenticate()"
                >Login</button>-->
                <div class="col-12">
                  <b-button
                    id="loginBtn"
                    type="submit"
                    name="btnlogin"
                    varient="link"
                    class="btn btn-theme btn-rounded mx-auto"
                    @click="login"
                  >
                    Login
                  </b-button>
                </div>
                <div class="col-12 mt-4">
                  <a
                    v-b-modal.forgotPassward
                    href="javascript:void(0)"
                    class="btn forgot-link text-sm text-lt"
                    >Forgot Password?</a
                  >
                </div>
                <div class="col-12">
                  <div>
                    <span class="text-sm text-lt">Are you a new User?</span>
                    <a
                      v-b-modal.register
                      href="javascript:void(0)"
                      class="forgot-link text-sm ml-2 font-weight-bold"
                      >Create an account</a
                    >
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="col-lg-4">
          <div class="aro-pswd_info">
            <div v-show="showPassInfo" id="pswd_info">
              <h4>Password must be requirements</h4>
              <ul>
                <li id="letter" class="invalid">
                  At least
                  <strong>one letter</strong>
                </li>
                <li id="capital" class="invalid">
                  At least
                  <strong>one capital letter</strong>
                </li>
                <li id="number" class="invalid">
                  At least
                  <strong>one number</strong>
                </li>
                <li id="length" class="invalid">
                  Be at least
                  <strong>8 characters</strong>
                </li>
                <li id="space" class="invalid">
                  be
                  <strong>use [~,!,@,#,$,%,^,&,*,-,=,.,;,']</strong>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
    <b-modal
      id="forgotPassward"
      hide-footer
      title
      size="md"
      modal-class="custom-modal"
      ok-variant="theme btn-rounded"
      cancel-variant="light btn-rounded"
      footer-class="border-0"
    >
      <form>
        <div class="form-group">
          <label for="exampleInputPassword1">Old Password</label>
          <input
            type="password"
            id="exampleInputPassword1"
            class="form-control"
          />
        </div>
        <div class="form-group">
          <label for="exampleInputPassword1">New Password</label>
          <input
            type="password"
            id="exampleInputPassword1"
            class="form-control"
          />
        </div>
        <button
          type="submit"
          class="btn btn-primary btn-theme btn-rounded mt-4"
        >
          Submit
        </button>
      </form>
    </b-modal>
    <!-- register -->
    <b-modal
      id="register"
      hide-footer
      title
      size="md"
      modal-class="custom-modal"
      ok-variant="theme btn-rounded"
      cancel-variant="light btn-rounded"
      footer-class="border-0"
    >
      <form>
        <div class="form-group">
          <label for="exampleInputName">UserName</label>
          <input
            type="text"
            id="exampleInputName"
            class="form-control"
            v-model="details.name"
          />
        </div>
        <div class="form-group">
          <label for="exampleInputAge">Age</label>
          <input
            type="number"
            id="exampleInputAge"
            class="form-control"
            v-model.number="details.age"
          />
        </div>
        <div class="form-group">
          <label for="exampleInputEmail">EMail</label>
          <input
            type="email"
            id="exampleInputEmail"
            class="form-control"
            v-model="details.email"
          />
        </div>
        <div class="form-group">
          <label for="exampleInputPhoneNo">PhoneNumber</label>
          <input
            type="text"
            id="exampleInputPhoneNo"
            class="form-control"
            v-model="details.phoneNo"
          />
        </div>
        <div class="form-group">
          <label for="exampleInputPassword1">Password</label>
          <input
            type="password"
            id="exampleInputPassword1"
            class="form-control"
            v-model="details.password"
          />
        </div>
        <div class="form-check">
          <input type="checkbox" class="form-check-input" id="exampleCheck1" />
          <label class="form-check-label" for="exampleCheck1"
            >I accept the terms and condition</label
          >
          <a
            v-b-modal.termsCondition
            href="javascript:void(0)"
            class="forgot-link text-sm ml-2 font-weight-bold"
            >Terms & Condition</a
          >
        </div>
        <button
          type="submit"
          class="btn btn-primary btn-theme btn-rounded mt-4"
          @click="register"
        >
          Submit
        </button>
      </form>
    </b-modal>
    <b-modal
      id="termsCondition"
      ok-only
      title="Terms and Condition for Industry Partner(IP)"
      size="lg"
      modal-class="custom-modal"
      ok-variant="theme btn-rounded"
      cancel-variant="light btn-rounded"
      footer-class="border-0"
    >
      <p>
        The interested Companies in collaborating with MKCL for this Work-based
        Learning initiative in skilling the youth of the Nation for making them
        employable are registered under IP agreement for Work-based Degree
        Program–BBA(SM) in collaboration with IGNOU & MKCL.
      </p>
      <h6>[A] RESPONSIBILITIES</h6>
      <ol class="list-style-no list-style-decimal">
        <li>
          <b>Selection and Admission</b>
          <ol class="list-style-decimal">
            <li>
              Convey to MKCL the requirement of number and suitable profile of
              Students/Interns who can perform at the workplace preferably in
              October/November or March/April every year along with their
              eligibility and objective selection criteria, terms and conditions
              of scholarship to be accepted by the selected candidate, a
              template of the letter including amount of monthly scholarship
              (offer letter), location of workplace, working hours, working
              days, and other facilities, etc.;
            </li>
            <li>
              Convey to MKCL if any in-service employees working in your company
              are applying for this degree program. Roles and responsibilities
              of IP will be same towards in-service employees.
            </li>
            <li>
              Participate, as per the schedule stipulated by MKCL, in the
              selection process;
            </li>
            <li>
              Make a formal offer of scholarship (for a period of three years)
              to the selected Students/Interns in form of a formal internship
              letter or any other form of the letter by specifically mentioning
              the conditions for the scholarship and seek their acceptance and
              complete their joining formalities at the specified Workplace.
            </li>
          </ol>
        </li>

        <li>
          <b>Work Environment:</b>
          <ol class="list-style-decimal">
            <li>
              Ensure availability of the real workplace as a
              working–cum-learning environment for Students/Interns equipped
              with requisite infrastructure, internet bandwidth and
              telecommunication and computing facilities and other necessary
              equipment, etc.,
            </li>

            <li>
              Ensure that the Students/Interns are being meaningfully engaged in
              the daily work and functions (real-life work) related to the
              service industry viz. facilitation, administration, operations,
              management for the number of hours of work as per the organization
              policy (schedules/hours of practical);
            </li>

            <li>
              Undertake the induction training and orientation of the
              Students/Interns;
            </li>

            <li>
              Stipulate and share with MKCL and the Students/Interns the
              performance appraisal system for assessing Students/Interns’
              performance at the workplace, for the award of Work Ratings to be
              forwarded to the University on monthly basis;
            </li>

            <li>
              Ensure that the Appraisers assess the performance of the
              Students/Interns as per industry norms and IP/WP’s performance
              appraisal system and orient them to award Work Ratings to the
              Students/Interns based on the performance.
            </li>

            <li>
              Establish periodic job rotation mechanism internally to ensure
              role progression at the workplace and advanced career growth for
              the Students/Interns, to the extent to which it is possible.
            </li>

            <li>
              In case of repeated incidences of
              non-performance/under-performance supported by objective
              evidences, a letter of improvement may be given to the concerned
              Student/Intern. Incase expected improvement in performance is not
              observed inspite of repeated feed for improvement, then the same
              can be reflected in the ratings along with specific observations.
            </li>

            <li>
              In case the organization has a written code of conduct and other
              applicable policies, which the Students/Interns have to comply
              with, a copy of the same be made available to them and an
              awareness and orientation session regarding the same be conducted.
            </li>

            <li>
              If the Student/Intern is found to be absconding and/or absent at
              the workplace in breach of applicable organization policies, rules
              and regulation for any reasons whatsoever, you shall intimate MKCL
              about the same within 2 working days to enable to take further
              appropriate action.
            </li>
          </ol>
        </li>
        <li>
          <b>Learning Environment</b>
          <ol class="list-style-decimal">
            <li>
              Nominate Mentors, if possible, in the ratio 1:16 (Mentor:
              Students/Interns) at the workplace and seek their commitment for
              mentorship. Mentors will perform following responsibilities:
              <ol class="list-style-no" type="a">
                <li>
                  Conduct meaningful interactions, discussion sessions with
                  Student/Intern around actions and reflections, and facilitate
                  and lead the Student to derive theory out of practice;
                </li>
                <li>
                  Develop and nurture Student/Intern to attain improved skill
                  efficiency for better performance at the workplace;
                </li>
                <li>
                  Ingrain in the Student/Intern corporate culture, corporate
                  values and core values to ensure that s/he not only leads an
                  ethical professional life but also follows core values in
                  her//his personal life;
                </li>
                <li>
                  Sensitize the Student/Intern about various environmental,
                  social and economic concerns, challenges and the positive role
                  s/he can play even as an individual for making the planet a
                  better place;
                </li>
                <li>
                  Assess performance, offer corrective and objective feedback to
                  Student/Intern and consolidate Work Ratings;
                </li>
                <li>
                  Use ePlatform effectively to collaborate with
                  Students/Interns, and facilitate them to build their
                  portfolios.
                </li>
                <li>
                  To attend meetings convened by MKCL based on mutual
                  convenience.
                </li>
                <li>
                  Share student performance (work ratings) with the project
                  coordinator on a monthly basis.
                </li>
              </ol>
            </li>
          </ol>
        </li>
        <li>
          <b>Monthly scholarship/stipend</b>
          <ol class="list-style-decimal">
            <li>
              The Industry Partner/Workplace Partner shall, in consultation with
              MKCL, stipulate per Student/Intern per month scholarship and other
              facilities to be provided to the Student/Intern for three years of
              the Work-based Degree program and pay the scholarship on a monthly
              basis.
            </li>

            <li>
              For in service employee the scholarship will be replace by their
              salary
            </li>
          </ol>
        </li>
        <li>
          <b>Project Coordinator at IP location</b>
          <ol class="list-style-decimal">
            <li>
              IP/WP shall nominate with adequate authority and accountability a
              coordinator who will act as a single point of contact (SPOC) with
              MKCL for successful implementation of Work-based Degree Program
              and perform the following responsibilities. It shall also nominate
              an alternative person to the said officer who shall be equally
              responsible in the absence of the first nominated officer
              <ol class="list-style-no" type="a">
                <li>
                  coordination with MKCL for implementing the program at
                  respective IP/WP and setting up eLearning environment and
                  ensuring access,
                </li>

                <li>
                  Ensuring consolidation of monthly Work Ratings and submission
                  to MKCL,
                </li>

                <li>
                  Addressing Student/Intern grievances if any and referring
                  student feedback / grievance if any related to the program to
                  MKCL,
                </li>

                <li>
                  Coordination for organizing quarterly meetings with
                  appraisers, mentors to be conducted by MKCL
                </li>
              </ol>
            </li>
          </ol>
        </li>
      </ol>
      <h6>[B] VALIDITY AND TERMINATION</h6>
      <ol class="list-style-no list-style-decimal">
        <li>
          The empanelment shall be effective from day of signing the empanelment
          letter and shall be valid unless terminated as below -
          <ol type="a">
            <li>
              by either party by giving the other party a notice in writing of
              three months of its intention to do so, but without dishonoring
              any commitment entered into prior to the date of termination
              notice and no party shall leave its commitment unfinished which
              may result in tangible losses to the Students/Interns, and/or the
              other party; or
            </li>

            <li>
              for any reasons such as legal processes, acts of the State or
              similar such exigencies beyond the normal control of the party
              concerned and which disable any of the parties hereto from
              functioning further; or
            </li>

            <li>by both parties by mutual consent.</li>
          </ol>
        </li>
        <li>
          Notwithstanding anything stated above, incase of notice of termination
          by the IP, the Students/Interns of ongoing Batch/es with the IP shall
          continue till the completion of their 3 year BBA(SM) and the IP shall,
          during such period discharge all the responsibilities diligently.
          Incase, the IP is not able to continue for any force majeure
          conditions, then both the IP and MKCL shall make utmost efforts to
          ensure the continuity of the education of the Students/Interns with
          alternative options.
        </li>
        <li>
          The clauses of this collaboration, which by nature are intended to
          survive termination shall remain in effect after such termination.
        </li>
      </ol>
      <h6>[C] INTELLECTUAL PROPERTY RIGHTS</h6>
      <ol class="list-style-no">
        <li>
          ntellectual Property in the context of this agreement shall refer to
          all such patents, trademarks, copyrights in respect of any hardware,
          software, product documentation, design document, or any other
          document, whether in printed or in electronic, digital or any other
          format which is an integral part of the hardware/software or is
          supplied along with such products which forms the subject matter of
          this agreement and shall also include study material, course material,
          educational and promotional content whether in printed or in
          electronic, digital or any other format and all business data
          generated during the period of validity of this Agreement.
        </li>
        <li>
          All the intellectual property rights, to and in the course name,
          content, methodology, assignments, question banks, etc. are the
          exclusive intellectual property of party that developed it and any
          third party components licensed by it shall remain the property of
          that third party. The data regarding the students applied for the
          course/s shall be the property of IGNOU and MKCL shall have right to
          access thereto only to the extent of and for performing its
          responsibilities thereunder.
        </li>
        <li>
          The software frameworks for the delivery of the program are developed
          by Maharashtra Knowledge Corporation Ltd. As such, the software code,
          whether compiled or un-compiled, in printed or electronic format, with
          software design logic, graphical user interfaces (GUI) and their
          design, look and feel, are explicit Intellectual Property of
          Maharashtra Knowledge Corporation Ltd. only. Any third-party
          components licensed by it, if any, shall remain the property of those
          respective third-parties.
        </li>
      </ol>
      <h6>[D] NON-DISCLOSURE</h6>
      <ol class="list-style-no">
        <li>
          Both the parties undertake to each other to keep confidential all
          information (written or oral) concerning the business and affairs of
          the other, which has been obtained or received during the course of
          performance hereunder, save that which is :
          <ol class="list-style-no" type="a">
            <li>Inconsequential or obvious;</li>
            <li>
              Already in its possession other than as a result of a breach of
              this clause; or
            </li>
            <li>
              In the hands of the public other than as a result of a breach of
              this clause.
            </li>
          </ol>
        </li>
        <li>
          In the event of any of the parties becoming legally compelled to
          disclose any confidential information, such party shall give
          sufficient notice to the other party so as to enable the other party
          to seek a timely protective order or any other appropriate relief. If
          such an order or other relief cannot be obtained, the party being
          required to make such a disclosure shall make the disclosure of the
          Confidential Information only to the extent that is legally required
          of it and no further.
        </li>
      </ol>
      <h6>[E] OTHERS</h6>
      <ol class="list-style-no">
        <li>
          All disputes and differences, whatsoever arising out of this
          collaboration shall be first attempted to be settled mutually
          amicably, otherwise referred to the courts at Pune which shall be the
          courts having jurisdiction to entertain and try the same.
        </li>
        <li>
          MKCL, is implementing Work-based Degree program in collaboration with
          Indira Gandhi National Open University. As such, Work-based Degree
          program is implemented in compliance of applicable University Grants
          Commission (UGC) Regulations and other applicable statutes.
        </li>
        <li>
          Each person signing this Letter of Empanelment represents and warrants
          that s/he is duly authorized and has legal capacity to enter into this
          collaboration. Each party represents and warrants to the other that
          the execution and delivery of this Letter of Empanelment and the
          performance of such party’s obligations hereunder have been duly
          authorized by all necessary corporate or other appropriate actions to
          execute this and that the terms and conditions of this empanelment are
          valid and binding on each party and enforceable in accordance with its
          terms.
        </li>
      </ol>
      <p>
        I, hereby declare that, I have read the terms and condition mentioned in
        the empanelment letter and rules related to admission mentioned in the
        prospectus of this program as well as on the University website. The
        information furnished by me related to company is true and complete to
        the best of my knowledge
      </p>
    </b-modal>
  </section>
</template>

<script>
import Vue from "vue";
import axios from "axios";
// import Response from "@/plugins/response.js";
// import "swiper/dist/css/swiper.css";
import { swiper, swiperSlide } from "vue-awesome-swiper";
export default {
  data() {
    return {
      details: {
        name: "",
        age: "",
        email: "",
        phoneNo: "",
        password: "",
      },
      log: {
        email: "",
        password: "",
      },
      showPassInfo: false,
      loginSlider: {
        navigation: {
          nextEl: ".swiper-next",
          prevEl: ".swiper-prev",
        },
        autoplay: {
          delay: 3000,
          disableOnInteraction: false,
        },
      },
    };
  },
  components: {
    swiper,
    swiperSlide,
  },
  methods: {
    register() {
      let self = this;
      axios
        .post("/signup", self.details)
        .then(function(res) {
          self.$bvToast.toast(`Registered successfully`, {
            title: "Success",
            variant: "success",
            autoHideDelay: 5000,
            solid: true,
            class: "toast",
          });
          console.log("93", res.data);
        })
        .catch(function() {
          console.log("FAILURE!!");
        });
    },
    login() {
      let self = this;
      axios
        .post("/login", self.log, {
          headers: {
            "Content-Type": "text/plain",
          },
        })
        .then(function(res) {
          console.log("58", res.config.data);
          if (res.config.data) {
            self.$router.push("/dashboard");
          }
        })
        .catch(function() {
          self.$bvToast.toast(`Please Enter Valid Credentials`, {
            title: "Not Valid",
             variant: 'danger',
            autoHideDelay: 5000,
            solid: true,
            class: "toast",
          });
          console.log("FAILURE!!");
        });
    },
    show() {
      let r = new Response({});
      r.showElement("loginFormId");
    },
    authenticate() {
      this.$store
        .dispatch("AUTH_REQUEST", {
          loginId: this.username,
          password: this.password,
        })
        .then((res) => {
          // Redirect to next page after suucessfull login
          alert("Login : " + res.isValid("MQLLogin"));
        })
        .catch((err) => {
          alert(err);
          Vue.$log.error(err);
        });

      // let req = {
      //   loginId: this.username,
      //   password: this.password
      // };
      // this.$MQLFetch('O.LoginService', req)
      //   .then(res => {
      //     // alert(JSON.stringify(res));
      //     this.$router.push("/");
      //   })
      //   .catch(error => {
      //     // Do in case of error
      //     Vue.error(error);
      //   });
    },
    validatePassword() {
      // validate password length
      if (this.password.length < 8) {
        let showLengthMsg = document.getElementById("length");
        showLengthMsg.classList.remove("valid");
        showLengthMsg.classList.add("invalid");
      } else {
        let showLengthMsg = document.getElementById("length");
        showLengthMsg.classList.remove("invalid");
        showLengthMsg.classList.add("valid");
      }

      // validate letter
      if (this.password.match(/[A-z]/)) {
        let showLengthMsg = document.getElementById("letter");
        showLengthMsg.classList.remove("invalid");
        showLengthMsg.classList.add("valid");
      } else {
        let showLengthMsg = document.getElementById("letter");
        showLengthMsg.classList.remove("valid");
        showLengthMsg.classList.add("invalid");
      }

      // validate capital letter
      if (this.password.match(/[A-Z]/)) {
        let showLengthMsg = document.getElementById("capital");
        showLengthMsg.classList.remove("invalid");
        showLengthMsg.classList.add("valid");
      } else {
        let showLengthMsg = document.getElementById("capital");
        showLengthMsg.classList.remove("valid");
        showLengthMsg.classList.add("invalid");
      }

      // validate number
      if (this.password.match(/\d/)) {
        let showLengthMsg = document.getElementById("number");
        showLengthMsg.classList.remove("invalid");
        showLengthMsg.classList.add("valid");
      } else {
        let showLengthMsg = document.getElementById("number");
        showLengthMsg.classList.remove("valid");
        showLengthMsg.classList.add("invalid");
      }

      // validate space
      if (this.password.match(/[^a-zA-Z0-9\-/]/)) {
        let showLengthMsg = document.getElementById("space");
        showLengthMsg.classList.remove("invalid");
        showLengthMsg.classList.add("valid");
      } else {
        let showLengthMsg = document.getElementById("space");
        showLengthMsg.classList.remove("valid");
        showLengthMsg.classList.add("invalid");
      }
    },
  },
  computed: {
    swiper() {
      return this.$refs.loginSlider.swiper;
    },
  },
};
</script>

<style lang="scss">
// @import "public/assets/css/variable.scss";
// @import "public/assets/css/mixin.scss";
#login {
  overflow: hidden;
  position: relative;
  &:before {
    content: "";
    width: 150%;
    height: 500px;
    background: $themeGradient;
    position: absolute;
    top: -97px;
    left: 50%;
    transform: translateX(-50%) rotate(0deg);
    z-index: -1;
    clip-path: polygon(0 0, 100% 0, 100% 75%, 0% 100%);
  }
  .login-content {
    margin-top: 80px;
    padding: 30px;
  }
  .login-image-wrapper {
    height: 100%;
    border-radius: 15px 0 0 15px;
    overflow: hidden;
    .swiper-container,
    .swiper-wrapper,
    .swiper-slide,
    .login-image-inner {
      height: 100%;
    }
    .card-body {
      padding: 0;
    }
    .login-image {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      height: 100%;
    }
    .login-image-inner {
      overflow: hidden;
      width: 100%;
      position: relative;
    }
    .slider-content {
      background: rgba($priColor, 0.9);
      padding: 15px;
      position: absolute;
      bottom: 0;
      left: 0;
      font-size: 14px;
      color: #fff;
      width: 100%;
      text-align: center;
    }
  }
  .login-box-wrapper {
    .card-body {
      padding: 40px 20px;
    }
    .logo {
      text-align: center;
    }
    .logo-img {
      width: 120px;
    }
    .form-control {
      padding-left: 40px;
      font-size: 14px;
      &::placeholder {
        // padding-left: 20px;
        font-size: 14px;
      }
    }
    .form-group {
      position: relative;
      .login-icon {
        position: absolute;
        top: 9px;
        left: 12px;
        font-size: 14px;
        color: #999;
      }
    }
    .other-controls {
      margin: 20px 0;
      font-size: 14px;
      color: #999;
      .forgot-link {
        color: #999;
        font-size: 14px;
        padding: 0;
        &:hover {
          color: $priColor;
          text-decoration: none;
        }
      }
    }
  }
  .swiper-slide-active {
    z-index: 9;
  }
  .swiper-control {
    .swiper-next,
    .swiper-prev {
      width: 35px;
      height: 40px;
      background: rgba($priColor, 0.7);
      position: absolute;
      z-index: 99;
      top: 50%;
      transform: translateY(-50%);
      border-radius: 50%;
      color: #fff;
      font-size: 30px;
      display: flex;
      justify-content: center;
      align-items: center;
      cursor: pointer;
      &:hover {
        background: rgba($priColor, 1);
      }
      &:focus {
        outline: none;
      }
    }
    .swiper-prev {
      left: 0px;
      padding-right: 5px;
      border-radius: 0px 50px 50px 0px;
    }
    .swiper-next {
      right: 0px;
      padding-left: 5px;
      border-radius: 50px 0px 0px 50px;
    }
  }
}
.list-style-no {
  padding-left: 30px;
  margin-bottom: 20px;
  ol {
    margin-bottom: 12px;
    padding-left: 15px;
    ol {
      margin-bottom: 8px;
      padding-left: 10px;
    }
  }
  li {
    line-height: 26px;
  }
}
.list-style-decimal {
  counter-reset: item;
  & > li {
    display: table;
    &:before {
      content: counters(item, ".") " ";
      counter-increment: item;
      margin-right: 8px;
      display: table-cell;
      padding-right: 10px;
    }
  }
}
@media (min-width: 992px) {
  #login {
    .login-box-wrapper {
      border-radius: 0 15px 15px 0;
      .card-body {
        padding: 50px;
      }
    }
    .login-content {
      margin-top: 100px;
      padding: 30px 50px;
    }
  }
}
// old css
#login {
  .login-card {
    box-shadow: 0 0 10px #ccc;
    padding: 30px;
  }
}
#pswd_info {
  background: #dfdfdf none repeat scroll 0 0;
  color: #fff;
  left: 20px;
  position: absolute;
  top: 115px;
}
#pswd_info h4 {
  background: black none repeat scroll 0 0;
  display: block;
  font-size: 14px;
  letter-spacing: 0;
  padding: 17px 0;
  text-align: center;
  text-transform: uppercase;
}
#pswd_info ul {
  list-style: outside none none;
}
#pswd_info ul li {
  padding: 10px 45px;
}

.valid {
  background: rgba(0, 0, 0, 0)
    url("https://s19.postimg.org/vq43s2wib/valid.png") no-repeat scroll 2px 6px;
  color: green;
  line-height: 21px;
  padding-left: 22px;
}

.invalid {
  background: rgba(0, 0, 0, 0)
    url("https://s19.postimg.org/olmaj1p8z/invalid.png") no-repeat scroll 2px
    6px;
  color: red;
  line-height: 21px;
  padding-left: 22px;
}

#pswd_info::before {
  background: #dfdfdf none repeat scroll 0 0;
  content: "";
  height: 25px;
  left: -13px;
  margin-top: -12.5px;
  position: absolute;
  top: 50%;
  transform: rotate(45deg);
  width: 25px;
}
</style>
