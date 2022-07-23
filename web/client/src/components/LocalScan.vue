<template>
  <div v-if="loading || data.length > 0">
    <h3>{{ name }} <b-spinner v-if="loading" type="grow"></b-spinner></h3>

    <b-table outlined :items="data" v-show="data.length > 0"></b-table>
  </div>
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
export default class GoogleSearch extends Vue {
  id = "local";
  name = "Local scan";
  data: LocalScanResponse[] = [];
  loading = false;
  computed = {
    ...mapState(["number"]),
    ...mapMutations(["pushError"]),
  };

  @Prop() scan!: Vue;

  mounted(): void {
    this.scan.$on("scan", async () => {
      this.loading = true;

      try {
        await this.run();
      } catch (e) {
        this.$store.commit("pushError", { message: `${this.name}: ${e}` });
      }

      this.loading = false;
    });
    this.scan.$on("clear", this.clear);
  }

  private clear() {
    this.data = [];
  }

  private async run(): Promise<void> {
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
  }
}
</script>
