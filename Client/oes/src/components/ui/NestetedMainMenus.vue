<template>
  <li class="nav-item">
    <b-link
      v-if="!Array.isArray(children.nodes)"
      class="nav-link"
      :to="children.url"
      :disabled="disabledLink"
    >
      {{ children.label }}
    </b-link>
    <template v-else-if="isMobile">
      <b-button
        v-b-toggle="'collapse-'+ depth + index"
        variant="link"
        class="nav-link dropdown-btn btn-block text-left"
        @click.prevent.stop="preventToggleCollapse(index)"
      >
        {{ children.label }}
      </b-button>
      <b-collapse
        :id="'collapse-'+ depth + index"
        class="dropdown-collapse"
        :accordion="'my-accordion-'+ depth"
        role="tabpanel"
      >
        <ul
          class="navigation navbar-nav main-menu"
          :class="{'dropdown-menu':Array.isArray(children.nodes) }"
        >
          <template>
            <NestetedMainMenus
              v-for="(node, index) in children.nodes"
              :index="index"
              :children="node"
              :depth="depth + 1"
              :key="index + 1"
            />
          </template>
        </ul>
      </b-collapse>
    </template>
    <template v-else>
      <b-button
        v-b-toggle="'collapse-'+ depth + index"
        variant="link"
        class="nav-link dropdown-btn btn-block text-left"
        @click.prevent.stop="preventToggleCollapse(index)"
        :disabled="disabledLink"
      >
        {{ children.label }}
      </b-button>
      <b-collapse
        :id="'collapse-'+ depth + index"
        class="dropdown-collapse "
        :accordion="'my-accordion-'+ depth"
        role="tabpanel"
        style="position:absolute;"
      >
        <ul
          class="navigation navbar-nav main-menu container-fluid"
          :class="{'dropdown-menu':Array.isArray(children.nodes) }"
        >
          <template>
            <NestetedMainMenus
              v-for="(node, index) in children.nodes"
              :index="index"
              :children="node"
              :depth="depth + 1"
              :key="index + 1"
            />
          </template>
        </ul>
      </b-collapse>
    </template>

    <!-- Desktop -->
    <!-- <template v-else>
      <b-nav-item-dropdown
        :id="'dropdown-'+ depth + index"
        class="mx-0"
        variant="link"
        :dropright="depth === 0 ? false: true"
        v-click-outside="closeAllDropdowns"
      >
        <template #button-content>
          {{ children.label }}
        </template>
        <template>
          <b-dropdown-item-button>
            <NestetedMainMenus
              v-for="(node, index) in children.nodes"
              :index="index + 1"
              :children="node"
              :depth="depth + 1"
              :key="index + 1"
            />
          </b-dropdown-item-button>
        </template>
      </b-nav-item-dropdown>
    </template> -->
  </li>
</template>
<script>
// import NestetedMainMenus from './NestetedMainMenus'
export default {
  name: 'NestetedMainMenus',
  components: {
    // NestetedMainMenus
  },
  props:
  {
    children: Object,
    label: String,
    depth: Number,
    index: Number,
    disabledLink: Boolean

  },
  inject: ['isMobile'],
  data () {
    return {
      showChildren: this.showChild,
      checkValue1: [],
      isDropdown2Visible: false,
      listDynamicStyle: ''
    }
  },
  methods: {

    toggleSideBar () {
      this.isActive = !this.isActive
    },
    preventToggleCollapse () {
      var sidebarContainer = document.querySelector('.app-container ')
      if (sidebarContainer.classList.contains('app-container')) {
        // e.preventDefault()
        // this.$root.$emit('bv::toggle::collapse', (collapseId)) => {
        //   collapseId = hide
        // })
        this.$root.$on('bv::collapse::state', (collapseId, isJustShown) => {

          if (isJustShown) {
            this.$root.$emit('bv::hide::collapse', collapseId)
            // this.stopPropagation()


          }
        })
      }
    },
    // calcListHeight () {

    //   this.listDynamicStyle = this.$el.querySelector('').clientHeight + 'px'
    // }
    closeAllDropdowns () {
      this.isDropdown2Visible = false
    }
  },
  mounted: function () {
    this.$root.$on('bv::dropdown::show', bvEvent => {
      if (bvEvent.componentId === 'dropdown-' + this.depth + this.index) {

        this.isDropdown2Visible = true
      }
    })
    this.$root.$on('bv::dropdown::hide', bvEvent => {
      if (bvEvent.componentId === 'dropdown-' + this.depth + this.index) {
        this.isDropdown2Visible = false
      }
      if (this.isDropdown2Visible) {
        if (this.depth > 0) {
          bvEvent.preventDefault()
        }
      }
    })
  }
}

</script>
