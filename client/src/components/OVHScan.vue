<template>
  <div v-if="loading || data.length > 0" class="mt-2">
    <hr />
    <h3>{{ name }} <b-spinner v-if="loading" type="grow"></b-spinner></h3>

    <b-button
      size="sm"
      variant="dark"
      v-b-toggle.ovh-collapse
      v-show="data.length > 0 && !loading"
      >Toggle results</b-button
    >
    <b-collapse id="ovh-collapse" class="mt-2">
      <b-table
        outlined
        :stacked="data.length == 1"
        :items="data"
        v-show="data.length > 0"
      ></b-table>
    </b-collapse>

    <hr />
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import axios, { AxiosResponse } from "axios";
import { mapMutations } from "vuex";
import config from "@/config";

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

  mounted() {
    this.scan.$on("scan", this.run);
    this.scan.$on("clear", this.clear);
  }

  private clear() {
    this.data = [];
  }

  private async run(): Promise<void> {
    this.loading = true;

    try {
      const res: AxiosResponse = await axios.get(
        `${config.apiUrl}/numbers/${this.$store.state.number}/scan/${this.id}`
      );

      this.data.push(res.data.result);
    } catch (e) {
      this.$store.commit("pushError", { message: e });
    }

    this.loading = false;
  }
}
</script>
