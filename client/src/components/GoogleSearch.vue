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
import axios from "axios";
import { mapMutations } from "vuex";
import config from "@/config";
import { ScanResponse } from "@/views/Scan.vue";

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
    this.scan.$on("scan", async () => {
      this.loading = true;

      try {
        await this.run();
      } catch (e) {
        this.$store.commit("pushError", { message: `${this.name}: ${e}` });
      }

      this.loading = false;
      this.scan.$emit("finished");
    });
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
    const res: ScanResponse<GoogleSearchScanResponse> = await axios.get(
      `${config.apiUrl}/numbers/${this.$store.state.number}/scan/${this.id}`,
      {
        validateStatus: () => true,
      }
    );

    if (!res.data.success && res.data.error) {
      throw res.data.error;
    }

    this.data = res.data.result;
  }

  openLinks(dork: GoogleSearchDork[]): void {
    for (const result of dork) {
      window.open(result.URL, "_blank");
    }
  }
}
</script>
