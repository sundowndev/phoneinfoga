<template>
  <div>
    <b-card
      v-if="isLookup || showInformations"
      header="Informations"
      class="mb-3 mt-3 text-center"
    >
      <b-list-group flush>
        <b-list-group-item
          v-for="(value, name) in localData"
          :key="name"
          class="text-left d-flex"
        >
          <h5 class="text-capitalize m-0 mr-4">{{ formatString(name) }}:</h5>
          <p class="m-0">{{ value }}</p>
        </b-list-group-item>
      </b-list-group>
    </b-card>

    <b-card v-if="isLookup" header="Scanners" class="text-center">
      <Scanner
        v-for="(scanner, index) in scanners"
        :key="index"
        :name="scanner.name.charAt(0).toUpperCase() + scanner.name.slice(1)"
        :scanId="scanner.name"
      />
    </b-card>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapMutations, mapState } from "vuex";
import { formatNumber, isValid, formatString, getScanners } from "../utils";
import Scanner from "../components/Scanner.vue";
import axios, { AxiosResponse } from "axios";
import config from "@/config";

interface ScannerObject {
  name: string;
  description: string;
}

interface Data {
  loading: boolean;
  isLookup: boolean;
  showInformations: boolean;
  scanners: Array<ScannerObject>;
  localData: {
    valid: boolean;
    raw_local: string;
    local: string;
    e164: string;
    international: string;
    countryCode: number;
    country: string;
    carrier: string;
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
      showInformations: false,
      scanners: [],
      localData: {
        valid: false,
        raw_local: "",
        local: "",
        e164: "",
        international: "",
        countryCode: 33,
        country: "",
        carrier: "",
      },
    };
  },
  mounted() {
    this.runScans();
  },
  methods: {
    formatString: formatString,
    async getScanners() {
      try {
        this.scanners = await getScanners();
      } catch (error) {
        this.$store.commit("pushError", { message: error });
      }
    },
    async runScans(): Promise<void> {
      if (!isValid(this.$route.params.number)) {
        this.$store.commit("pushError", { message: "Number is not valid." });
        return;
      }

      this.loading = true;

      this.$store.commit("setNumber", formatNumber(this.$route.params.number));

      try {
        const res = await axios.post(`${config.apiUrl}/v2/numbers`, {
          number: this.$store.state.number,
        });

        this.localData = res.data;

        if (this.localData.valid) {
          this.getScanners();
          this.isLookup = true;
        } else {
          this.showInformations = true;
        }
      } catch (error) {
        this.$store.commit("pushError", { message: error });
      }

      this.loading = false;
    },
  },
});
</script>
