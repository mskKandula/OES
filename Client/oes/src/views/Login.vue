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
                  <div class="slider-content">Learning Management System</div>
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
            v-model="details.mobileNo"
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
      title="Terms and Condition for User Registration"
      size="lg"
      modal-class="custom-modal"
      ok-variant="theme btn-rounded"
      cancel-variant="light btn-rounded"
      footer-class="border-0"
    >
      <p>
        Lorem ipsum dolor sit, amet consectetur adipisicing elit. Eius nam
        provident totam adipisci! Temporibus mollitia odio, dolore reiciendis
        similique obcaecati commodi molestias eum dolor voluptas harum alias
        fugit velit animi.
      </p>
      <h6>[A] RESPONSIBILITIES</h6>
      <ol class="list-style-no list-style-decimal">
        <li>
          <b>Selection and Admission</b>
          <ol class="list-style-decimal">
            <li>
              Lorem ipsum, dolor sit amet consectetur adipisicing elit. Numquam
              aliquam ad pariatur quod quos quas omnis, error magni labore, nam
              et eaque cumque ea iusto porro eligendi recusandae vero impedit.
            </li>
            <li>
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Aliquam
              deleniti esse voluptatibus velit facere, a architecto tenetur
              magni dolor cum iure provident deserunt assumenda. Et pariatur non
              consequuntur dolore minus.
            </li>
            <li>
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Dolorem
              illum nostrum dolore distinctio, aspernatur eaque accusamus.
              Molestiae sapiente facere, eaque, soluta, nobis culpa sunt iure a
              architecto labore est quod.
            </li>
            <li>
              Lorem ipsum, dolor sit amet consectetur adipisicing elit. Odit et
              explicabo eveniet distinctio, amet numquam pariatur adipisci nisi
              atque repellat nam dolorem? Explicabo eos voluptates itaque magni
              eum nulla impedit.
            </li>
          </ol>
        </li>

        <li>
          <b>Work Environment:</b>
          <ol class="list-style-decimal">
            <li>
              Lorem, ipsum dolor sit amet consectetur adipisicing elit.
              Dignissimos eveniet recusandae nobis architecto quam vitae
              molestias voluptates voluptatum enim, perferendis pariatur fugit
              velit similique. Odio soluta corporis alias? Repellendus, ducimus!
            </li>

            <li>
              Lorem, ipsum dolor sit amet consectetur adipisicing elit. Itaque
              at odio quasi exercitationem ab dicta doloribus assumenda
              cupiditate voluptas nisi dolorem totam enim, animi rerum nam
              pariatur impedit beatae perferendis!
            </li>

            <li>
              Lorem ipsum dolor sit amet consectetur, adipisicing elit. At optio
              deleniti fugiat sint dolore, ullam unde quibusdam autem nemo ab
              voluptates accusamus commodi repellendus quidem dignissimos
              tempore, assumenda recusandae beatae.
            </li>

            <li>
              Lorem ipsum dolor sit amet consectetur adipisicing elit. At, qui
              sint modi similique vero ab laboriosam quibusdam nisi aliquam
              autem voluptatum soluta a esse necessitatibus quasi impedit maxime
              atque! Odit!
            </li>

            <li>
              Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ut
              ratione perspiciatis non saepe quas consectetur repellat maxime
              officia veniam! Ipsa ducimus saepe consequuntur nihil
              necessitatibus pariatur perferendis illum mollitia perspiciatis.
            </li>

            <li>
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Soluta
              quasi possimus ullam quia animi quisquam maiores recusandae in
              amet consequuntur eum expedita optio exercitationem temporibus
              eius minima, minus, quidem corporis.
            </li>

            <li>
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Voluptas
              doloremque nostrum est at sint veniam alias, quod necessitatibus
              rem et, accusamus quam, eaque dolor! Fuga aspernatur illum
              eligendi laudantium magnam.
            </li>
          </ol>
        </li>
        <li>
          <b>Learning Environment</b>
          <ol class="list-style-decimal">
            <li>
              Lorem ipsum dolor, sit amet consectetur adipisicing elit. Mollitia
              maiores impedit dignissimos, quo, adipisci velit voluptatum qui
              vero aliquam alias ea libero fugiat, nihil perferendis. Ut
              architecto aliquid dolores excepturi.
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
              </ol>
            </li>
          </ol>
        </li>
        <li>
          <b>Monthly scholarship/stipend</b>
          <ol class="list-style-decimal">
            <li>
              Lorem ipsum dolor, sit amet consectetur adipisicing elit. Nisi
              perspiciatis inventore aspernatur eos voluptatum ullam, saepe,
              adipisci ratione quibusdam rerum cum. Ipsum dolorem tempore
              aperiam impedit alias reiciendis in optio.
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
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Impedit
              beatae neque aperiam nihil consectetur ab unde id iste corporis,
              cupiditate similique reiciendis facere maiores dolorum at quos,
              sit maxime velit.
              <ol class="list-style-no" type="a">
                <li>
                  Lorem ipsum dolor sit amet consectetur adipisicing elit.
                  Repellat hic placeat corrupti quos beatae ab veniam aliquid
                  accusantium soluta illum, tempora officiis facilis corporis
                  itaque ad laudantium. Architecto, officiis voluptatem.
                </li>
                <li>
                  Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsum
                  asperiores quibusdam recusandae officiis libero maxime totam
                  eius similique? Sequi tempore porro mollitia, earum officiis
                  provident! Atque qui temporibus pariatur dolorum!
                </li>
              </ol>
            </li>
          </ol>
        </li>
      </ol>
      <h6>[B] VALIDITY AND TERMINATION</h6>
      <ol class="list-style-no list-style-decimal">
        <li>
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Consectetur
          officiis odit libero nihil eos eveniet nulla eaque earum neque
          doloremque dolores tempora laudantium architecto expedita labore
          deleniti consequatur, aliquid recusandae.
          <ol type="a">
            <li>
              Lorem ipsum dolor sit amet consectetur adipisicing elit.
              Consequuntur, enim quia, vitae adipisci optio impedit temporibus
              provident suscipit facere cupiditate at voluptas minima maiores
              magnam, explicabo sapiente dolorem repudiandae vero.
            </li>

            <li>
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Pariatur
              temporibus voluptatum eaque ipsam a facere quam eum vitae sunt?
              Ducimus corrupti id maxime nostrum eos soluta inventore debitis
              eveniet temporibus!
            </li>

            <li>by both parties by mutual consent.</li>
          </ol>
        </li>
        <li>
          Lorem ipsum, dolor sit amet consectetur adipisicing elit. Minus, velit
          illum, maiores voluptatibus modi beatae magni optio fugit impedit,
          officia nostrum esse eum ipsam sit? Voluptatem autem tempore deleniti
          omnis.
        </li>
        <li>
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Harum
          incidunt reiciendis eveniet ex molestiae hic aliquam quibusdam
          nesciunt suscipit, exercitationem assumenda reprehenderit rerum
          voluptatem, cumque, sed atque? Sapiente, vitae! Id.
        </li>
      </ol>
      <h6>[C] INTELLECTUAL PROPERTY RIGHTS</h6>
      <ol class="list-style-no">
        <li>
          Lorem ipsum dolor sit amet consectetur, adipisicing elit. Nulla
          voluptatibus ex, et nemo eum asperiores officia culpa assumenda earum
          ipsa? Neque quae natus sapiente atque dolore aut, nobis id deleniti.
        </li>
        <li>
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Totam magnam
          aliquam eveniet voluptates nulla quis modi ad adipisci fugiat rerum,
          laborum sit praesentium voluptas, impedit repellat vero, architecto
          veritatis quaerat!
        </li>
        <li>
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Consequatur,
          excepturi autem? Accusamus aut consequuntur excepturi debitis. Totam,
          iure. Ducimus eius dolor tenetur aliquam tempora culpa ratione
          incidunt, accusamus optio. Pariatur?
        </li>
      </ol>
      <h6>[D] NON-DISCLOSURE</h6>
      <ol class="list-style-no">
        <li>
          Lorem, ipsum dolor sit amet consectetur adipisicing elit. Impedit eius
          officiis qui dolorem et necessitatibus ea fugiat alias ex sed nulla,
          ipsum sapiente ab facilis voluptatum deleniti iure suscipit unde?
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
          Lorem ipsum dolor sit amet, consectetur adipisicing elit. Accusamus,
          obcaecati. Voluptatum, possimus nostrum aliquam corrupti error sequi
          fugit dolorum laborum nam vitae, ab quasi ipsam minus omnis nihil!
          Laboriosam, sed?
        </li>
      </ol>
      <h6>[E] OTHERS</h6>
      <ol class="list-style-no">
        <li>
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Nobis
          repellendus incidunt, nisi id maiores magni non libero consectetur
          velit delectus laboriosam aspernatur consequatur blanditiis quisquam
          dignissimos ipsa eaque veniam quod!
        </li>
        <li>
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Voluptatem
          earum quia quasi suscipit error facere corrupti quisquam veniam animi
          quo ullam ratione inventore aut ad, culpa consectetur cupiditate
          explicabo minus?
        </li>
        <li>
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Eum, voluptas
          est. Corporis quo earum ratione neque fugit exercitationem sequi
          possimus voluptatum explicabo asperiores vel rerum, ipsam doloremque
          suscipit aspernatur sed?
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
        mobileNo: "",
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
      this.$http
        .post("/api/o/signUp", self.details)
        .then(function (res) {
          console.log(res);
          self.$bvToast.toast(`Registered successfully`, {
            title: "Success",
            variant: "success",
            autoHideDelay: 5000,
            solid: true,
            class: "toast",
          });
        })
        .catch(function () {
          console.log("FAILURE!!");
        });
    },
    login() {
      let self = this;
      this.$http
        .post("/api/o/login", self.log, {
          headers: {
            "Content-Type": "text/plain",
          },
        })
        .then(function (res) {
          if (res.data) {
            if (res.data.userType === "User") {
              self.$router.push({ path: "/dashboard" });
            } else {
              self.$router.push({ path: "/studentDashboard" });
            }
          }
        })
        .catch(function () {
          self.$bvToast.toast(`Please Enter Valid Credentials`, {
            title: "Not Valid",
            variant: "danger",
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
@import "../assets/scss/variable.scss";
@import "../assets/scss/mixin.scss";
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
