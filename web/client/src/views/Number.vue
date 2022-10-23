<template>
  <div>
    <b-container v-if="isLookup" class="border p-4 mb-3">
      <h3 class="text-center">Local</h3>
      <b-container>
        <b-row v-for="(value, name) in localData" :key="name" align-v="center">
          <h5 class="text-capitalize m-0 mr-4">{{ name }}:</h5>
          <p class="m-0">{{ value }}</p>
        </b-row>
      </b-container>
    </b-container>

    <b-container v-if="isLookup" class="border p-4">
      <h3 class="text-center">Scanners</h3>
      <Scanner name="GoogleSearch" scanId="googlesearch" />
      <Scanner name="Numverify Scan" scanId="numverify" />
      <Scanner name="OVH Telecom scan" scanId="ovh" />
    </b-container>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapMutations, mapState } from "vuex";
import { formatNumber, isValid } from "../utils";
import Scanner from "../components/Scanner.vue";
import axios, { AxiosResponse } from "axios";
import config from "@/config";

interface Data {
  loading: boolean;
  isLookup: boolean;
  localData: {
    raw_local: string;
    local: string;
    e164: string;
    international: string;
    country_code: string;
    country: string;
  };
}

export type ScanResponse<T> = AxiosResponse<{
  success: boolean;
  result: T;
  error: string;
}>;

export default Vue.extend({
  components: { Scanner },
  computed: {
    ...mapState(["number"]),
    ...mapMutations(["pushError"]),
  },
  data(): Data {
    return {
      loading: false,
      isLookup: false,
      localData: {
        raw_local: "",
        local: "",
        e164: "",
        international: "",
        country_code: "",
        country: "",
      },
    };
  },
  mounted() {
    this.runScans();
  },
  methods: {
    async runScans(): Promise<void> {
      if (!isValid(this.$route.params.number)) {
        this.$store.commit("pushError", { message: "Number is not valid." });
        return;
      }

      this.loading = true;

      this.$store.commit("setNumber", formatNumber(this.$route.params.number));

      try {
        const res = await axios.get(
          `${config.apiUrl}/numbers/${this.$store.state.number}/scan/local`,
          {
            validateStatus: () => true,
          }
        );

        this.localData = res.data.result;
      } catch (error) {
        this.$store.commit("pushError", { message: error });
      }

      this.isLookup = true;
      this.loading = false;
    },
  },
});
</script>
