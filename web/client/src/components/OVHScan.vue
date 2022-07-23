<template>
  <div v-if="loading || data.length > 0" class="mt-2">
    <hr />
    <h3>
      {{ name }}
      <b-spinner v-if="loading" type="grow"></b-spinner>
    </h3>

    <b-button
      size="sm"
      variant="dark"
      v-b-toggle.ovh-collapse
      v-show="data.length > 0 && !loading"
      >Toggle results
    </b-button>
    <b-collapse id="ovh-collapse" class="mt-2">
      <b-table
        outlined
        :stacked="data.length === 1"
        :items="data"
        v-show="data.length > 0"
      ></b-table>
    </b-collapse>

    <hr />
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import axios from "axios";
import { mapMutations } from "vuex";
import config from "@/config";
import { ScanResponse } from "@/views/Scan.vue";

interface OVHScanResponse {
  found: boolean;
  numberRange: string;
  city: string;
  zipCode: string;
}

@Component
export default class GoogleSearch extends Vue {
  id = "ovh";
  name = "OVH Telecom scan";
  data: OVHScanResponse[] = [];
  loading = false;
  computed = {
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
    const res: ScanResponse<OVHScanResponse> = await axios.get(
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
