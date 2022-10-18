<template>
  <b-container>
    <b-row align-h="between" align-v="center">
      <h3>GoogleSearch</h3>
      <b-button
        v-if="!error && !loading"
        @click="runScan"
        variant="outline-primary"
        size="lg"
        >Run</b-button
      >
      <b-spinner v-if="loading && !error" type="grow"></b-spinner>
      <b-alert
        class="m-0"
        v-if="error && !loading"
        show
        variant="danger"
        fade
        >{{ error }}</b-alert
      >
    </b-row>
  </b-container>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import axios from "axios";
import { mapState, mapMutations } from "vuex";
import config from "@/config";
import { ScanResponse } from "@/views/Scan.vue";

interface LocalScanResponse {
  number: string;
  dork: string;
  URL: string;
}

@Component
export default class Scanner extends Vue {
  id = "local";
  name = "Google";
  data: LocalScanResponse[] = [];
  loading = false;
  error = null;
  computed = {
    ...mapState(["number"]),
    ...mapMutations(["pushError"]),
  };

  @Prop() scan!: Vue;

  mounted(): void {
    this.scan.$on("clear", this.clear);
  }

  private clear() {
    this.data = [];
  }

  private async runScan(): Promise<void> {
    this.loading = true;
    try {
      const res: ScanResponse<LocalScanResponse> = await axios.get(
        `${config.apiUrl}/numbers/${this.$store.state.number}/scan/${this.id}`,
        {
          validateStatus: () => true,
        }
      );

      if (!res.data.success && res.data.error) {
        throw res.data.error;
      }
      this.data.push(res.data.result);
    } catch (error) {
      this.error = error;
    }
    this.loading = false;
  }
}
</script>
