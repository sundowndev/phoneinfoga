<template>
  <b-container class="mb-3">
    <b-row align-h="between" align-v="center">
      <h3>{{ name }}</h3>
      <b-button
        v-if="!error && !loading && !data"
        @click="runScan"
        variant="dark"
        size="lg"
        >Run</b-button
      >
      <b-spinner v-if="loading && !error" type="grow"></b-spinner>
      <b-row v-if="error && !loading">
        <b-alert class="m-0" show variant="danger" fade>{{ error }}</b-alert>
        <b-button
          v-if="!dryrunError"
          @click="runScan"
          variant="danger"
          size="lg"
          >Retry</b-button
        >
      </b-row>
    </b-row>
    <b-collapse :id="collapseId" class="mt-2 text-left">
      <JsonViewer :value="data"></JsonViewer>
    </b-collapse>
  </b-container>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import axios from "axios";
import { mapState, mapMutations } from "vuex";
import JsonViewer from "vue-json-viewer";
import config from "@/config";

@Component({
  components: {
    JsonViewer,
  },
})
export default class Scanner extends Vue {
  data = null;
  loading = false;
  dryrunError = false;
  error: unknown = null;
  computed = {
    ...mapState(["number"]),
    ...mapMutations(["pushError"]),
  };

  @Prop() scanId!: string;
  @Prop() name!: string;

  collapseId = "scanner-collapse" + this.scanId;

  mounted(): void {
    this.dryRun();
  }

  private async dryRun(): Promise<void> {
    try {
      const res = await axios.post(
        `${config.apiUrl}/v2/scanners/${this.scanId}/dryrun`,
        {
          number: this.$store.state.number,
        },
        {
          validateStatus: () => true,
        }
      );

      if (!res.data.success && res.data.error) {
        throw res.data.error;
      }
    } catch (error: unknown) {
      this.dryrunError = true;
      this.error = error;
    }
  }

  private async runScan(): Promise<void> {
    this.error = null;
    this.loading = true;
    try {
      const res = await axios.post(
        `${config.apiUrl}/v2/scanners/${this.scanId}/run`,
        {
          number: this.$store.state.number,
        },
        {
          validateStatus: () => true,
        }
      );

      if (!res.data.success && res.data.error) {
        throw res.data.error;
      }
      this.data = res.data.result;
      this.$root.$emit("bv::toggle::collapse", this.collapseId);
    } catch (error) {
      this.error = error;
    }

    this.loading = false;
  }
}
</script>
