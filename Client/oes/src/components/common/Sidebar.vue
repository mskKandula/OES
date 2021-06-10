<template>
  <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block  sidebar collapse">
    <div class="sidebar-sticky pt-3">
      <ul class="nav flex-column">
        <li class="nav-item">
          <b-link class="nav-link" to="/dashboard">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="feather feather-home"
            >
              <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
              <polyline points="9 22 9 12 15 12 15 22" />
            </svg>
            Dashboard
            <span class="sr-only">(current)</span>
          </b-link>
        </li>
        <li class="nav-item">
          <b-link class="nav-link" to="/formElements">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="feather feather-file"
            >
              <path
                d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"
              />
              <polyline points="13 2 13 9 20 9" />
            </svg>
            Form Element
          </b-link>
        </li>
      </ul>
    </div>
  </nav>
</template>
<style lang="scss">
#sidebarMenu {
  /*
 * Sidebar
 */
  background: rgba(255, 255, 255, 0.2);
  .feather {
    width: 16px;
    height: 16px;
    vertical-align: text-bottom;
  }

  &.sidebar {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    z-index: 100; /* Behind the navbar */
    padding: 48px 0 0; /* Height of navbar */
    box-shadow: inset -1px 0 0 rgba(0, 0, 0, 0.1);
  }

  @media (max-width: 767.98px) {
    &.sidebar {
      top: 5rem;
    }
  }

  .sidebar-sticky {
    position: relative;
    top: 0;
    height: calc(100vh - 48px);
    padding-top: 0.5rem;
    overflow-x: hidden;
    overflow-y: auto; /* Scrollable contents if viewport is shorter than content. */
  }

  @supports ((position: -webkit-sticky) or (position: sticky)) {
    .sidebar-sticky {
      position: -webkit-sticky;
      position: sticky;
    }
  }

  &.sidebar .nav-link {
    font-weight: 500;
    color: #333;
  }

  &.sidebar .nav-link .feather {
    margin-right: 4px;
    color: #999;
  }

  &.sidebar .nav-link.active {
    color: #007bff;
  }

  &.sidebar .nav-link:hover .feather,
  &.sidebar .nav-link.active .feather {
    color: inherit;
  }

  .sidebar-heading {
    font-size: 0.75rem;
    text-transform: uppercase;
  }
}
</style>
