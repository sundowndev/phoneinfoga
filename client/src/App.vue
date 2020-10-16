<template>
  <div id="app" style="padding-bottom: 50px">
    <div>
      <b-navbar toggleable="lg" type="dark" variant="dark">
        <b-container>
          <b-navbar-brand to="/">
            <img
              src="@/assets/logo.svg"
              class="d-inline-block align-top"
              width="30"
              height="30"
              alt="logo"
            />
            {{ config.appName }}
          </b-navbar-brand>

          <b-collapse id="nav-text-collapse" is-nav>
            <b-navbar-nav>
              <b-nav-text>{{ config.appDescription }}</b-nav-text>
            </b-navbar-nav>
          </b-collapse>

          <b-navbar-nav class="ml-auto">
            <b-collapse id="nav-collapse" is-nav>
              <b-navbar-nav>
                <b-nav-item
                  href="https://github.com/sundowndev/PhoneInfoga"
                  target="_blank"
                  >GitHub</b-nav-item
                >
                <b-nav-item
                  href="https://sundowndev.github.io/PhoneInfoga/resources/"
                  target="_blank"
                  >Resources</b-nav-item
                >
                <b-nav-item
                  href="https://sundowndev.github.io/PhoneInfoga/"
                  target="_blank"
                  >Documentation</b-nav-item
                >
              </b-navbar-nav>
            </b-collapse>
          </b-navbar-nav>
        </b-container>
      </b-navbar>
    </div>

    <b-container class="my-md-3">
      <b-row>
        <b-col cols="12">
          <b-alert
            v-for="(err, i) in errors"
            v-bind:key="i"
            show
            variant="danger"
            dismissible
            fade
            >{{ err.message }}</b-alert
          >

          <router-view />
        </b-col>
      </b-row>
    </b-container>

    <b-navbar
      toggleable="lg"
      type="light"
      variant="light"
      fixed="bottom"
      v-if="version !== ''"
    >
      <b-container>
        <b-navbar-nav class="ml-auto">
          <b-collapse id="nav-collapse" is-nav>
            <b-navbar-nav>
              <b-nav-item
                href="https://github.com/sundowndev/PhoneInfoga/releases"
                target="_blank"
                >{{ config.appName }} {{ version }}</b-nav-item
              >
            </b-navbar-nav>
          </b-collapse>
        </b-navbar-nav>
      </b-container>
    </b-navbar>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapState } from "vuex";
import config from "@/config";
import axios, { AxiosResponse } from "axios";

type HealthResponse = { success: boolean; version: string };

export default Vue.extend({
  data: () => ({ config, version: "" }),
  computed: {
    ...mapState(["number", "errors"]),
  },
  async created() {
    const res: AxiosResponse<HealthResponse> = await axios.get(config.apiUrl);

    this.version = res.data.version;
  },
});
</script>
