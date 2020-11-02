<template>
  <div v-if="loading || data.length > 0">
    <hr />
    <h3>{{ name }} <b-spinner v-if="loading" type="grow"></b-spinner></h3>

    <b-button
      size="sm"
      variant="dark"
      v-b-toggle.numverify-collapse
      v-show="data.length > 0 && !loading"
      >Toggle results</b-button
    >
    <b-collapse id="numverify-collapse" class="mt-2">
      <b-table
        outlined
        :stacked="data.length === 1"
        :items="data"
        v-show="data.length > 0"
      ></b-table>
    </b-collapse>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import axios from "axios";
import { mapMutations } from "vuex";
import config from "@/config";
import { ScanResponse } from "@/views/Scan.vue";

interface NumverifyScanResponse {
  valid: boolean;
  number: string;
  localFormat: string;
  internationalFormat: string;
  countryPrefix: string;
  countryCode: string;
  countryName: string;
  location: string;
  carrier: string;
  lineType: string;
}

@Component
export default class GoogleSearch extends Vue {
  id = "numverify";
  name = "Numverify scan";
  data: NumverifyScanResponse[] = [];
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
      this.scan.$emit("finished");
    });
    this.scan.$on("clear", this.clear);
  }

  private clear() {
    this.data = [];
  }

  private async run(): Promise<void> {
    const res: ScanResponse<NumverifyScanResponse> = await axios.get(
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
