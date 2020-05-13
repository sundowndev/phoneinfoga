<template>
  <div v-if="loading || data.length > 0">
    <h3>{{ name }} <b-spinner v-if="loading" type="grow"></b-spinner></h3>

    <b-table outlined :items="data" v-show="data.length > 0"></b-table>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import axios, { AxiosResponse } from "axios";
import { mapState, mapMutations } from "vuex";
import config from "@/config";

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
