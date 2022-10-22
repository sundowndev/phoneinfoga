<template>
  <b-container class="mb-3">
    <b-row align-h="between" align-v="center">
      <h3>{{ name }}</h3>
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
    <b-collapse id="scanner-collapse" class="mt-2">
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

// Vue.use(JsonViewer);

@Component({
  components: {
    JsonViewer,
  },
})
export default class Scanner extends Vue {
  data = {};
  loading = false;
  error = null;
  computed = {
    ...mapState(["number"]),
    ...mapMutations(["pushError"]),
  };

  @Prop() scanId!: string;
  @Prop() name!: string;

  private async runScan(): Promise<void> {
    this.loading = true;
    try {
      const res = await axios.get(
        `${config.apiUrl}/numbers/${this.$store.state.number}/scan/${this.scanId}`,
        {
          validateStatus: () => true,
        }
      );

      if (!res.data.success && res.data.error) {
        throw res.data.error;
      }
      this.data = res.data.result;
      this.$root.$emit("bv::toggle::collapse", "scanner-collapse");
    } catch (error) {
      this.error = error;
    }

    this.loading = false;
  }
}
</script>
