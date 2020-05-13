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
        :stacked="data.length == 1"
        :items="data"
        v-show="data.length > 0"
      ></b-table>
    </b-collapse>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import axios, { AxiosResponse } from "axios";
import { mapMutations } from "vuex";
import config from "@/config";

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

    this.scan.$emit("finished");
    this.loading = false;
  }
}
</script>
