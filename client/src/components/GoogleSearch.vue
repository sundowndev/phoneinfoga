<template>
  <div v-if="loading || data.socialMedia.length > 0">
    <hr />
    <h3>
      {{ name }}
      <b-spinner v-if="loading" type="grow"></b-spinner>
    </h3>

    <b-button
      size="sm"
      variant="dark"
      v-b-toggle.googlesearch-collapse
      v-show="data.socialMedia.length > 0 && !loading"
      >Toggle results</b-button
    >

    <b-collapse id="googlesearch-collapse" class="mt-2">
      <div class="my-3">
        <h4>
          General footprints
          <small>
            <b-button
              variant="outline-primary"
              size="sm"
              v-on:click="openLinks(data.general)"
              >Open all links</b-button
            >
          </small>
        </h4>

        <b-list-group>
          <b-list-group-item
            :href="value.URL"
            target="blank"
            v-for="(value, i) in data.general"
            v-bind:key="i"
            >{{ value.dork }}</b-list-group-item
          >
        </b-list-group>
      </div>

      <div class="my-3">
        <h4>
          Social networks footprints
          <small>
            <b-button
              variant="outline-primary"
              size="sm"
              v-on:click="openLinks(data.socialMedia)"
              >Open all links</b-button
            >
          </small>
        </h4>

        <b-list-group>
          <b-list-group-item
            :href="value.URL"
            target="blank"
            v-for="(value, i) in data.socialMedia"
            v-bind:key="i"
            >{{ value.dork }}</b-list-group-item
          >
        </b-list-group>
      </div>

      <div class="my-3">
        <h4>
          Individual footprints
          <small>
            <b-button
              variant="outline-primary"
              size="sm"
              v-on:click="openLinks(data.individuals)"
              >Open all links</b-button
            >
          </small>
        </h4>

        <b-list-group>
          <b-list-group-item
            :href="value.URL"
            target="blank"
            v-for="(value, i) in data.individuals"
            v-bind:key="i"
            >{{ value.dork }}</b-list-group-item
          >
        </b-list-group>
      </div>

      <div class="my-3">
        <h4>
          Reputation footprints
          <small>
            <b-button
              variant="outline-primary"
              size="sm"
              v-on:click="openLinks(data.reputation)"
              >Open all links</b-button
            >
          </small>
        </h4>

        <b-list-group>
          <b-list-group-item
            :href="value.URL"
            target="blank"
            v-for="(value, i) in data.reputation"
            v-bind:key="i"
            >{{ value.dork }}</b-list-group-item
          >
        </b-list-group>
      </div>

      <div class="my-3">
        <h4>
          Temporary number providers footprints
          <small>
            <b-button
              variant="outline-primary"
              size="sm"
              v-on:click="openLinks(data.disposableProviders)"
              >Open all links</b-button
            >
          </small>
        </h4>

        <b-list-group>
          <b-list-group-item
            :href="value.URL"
            target="blank"
            v-for="(value, i) in data.disposableProviders"
            v-bind:key="i"
            >{{ value.dork }}</b-list-group-item
          >
        </b-list-group>
      </div>
    </b-collapse>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import axios, { AxiosResponse } from "axios";
import { mapMutations } from "vuex";
import config from "@/config";

interface GoogleSearchScanResponse {
  socialMedia: GoogleSearchDork[];
  disposableProviders: GoogleSearchDork[];
  reputation: GoogleSearchDork[];
  individuals: GoogleSearchDork[];
  general: GoogleSearchDork[];
}

interface GoogleSearchDork {
  number: string;
  dork: string;
  URL: string;
}

@Component
export default class GoogleSearch extends Vue {
  id = "googlesearch";
  name = "Google search";
  data: GoogleSearchScanResponse = {
    socialMedia: [],
    disposableProviders: [],
    reputation: [],
    individuals: [],
    general: [],
  };
  loading = false;
  computed = {
    ...mapMutations(["pushError"]),
  };

  @Prop() scan!: Vue;

  mounted() {
    this.scan.$on("scan", this.run);
    this.scan.$on("clear", this.clear);
  }

  private clear() {
    this.data = {
      socialMedia: [],
      disposableProviders: [],
      reputation: [],
      individuals: [],
      general: [],
    };
  }

  private async run(): Promise<void> {
    this.loading = true;

    try {
      const res: AxiosResponse = await axios.get(
        `${config.apiUrl}/numbers/${this.$store.state.number}/scan/${this.id}`
      );

      this.data = res.data.result;

      console.log("google", this.data.socialMedia);
    } catch (e) {
      this.$store.commit("pushError", { message: e });
    }

    this.loading = false;
  }

  openLinks(dork: GoogleSearchDork[]): void {
    for (const result of dork) {
      window.open(result.URL, "_blank");
    }
  }
}
</script>
