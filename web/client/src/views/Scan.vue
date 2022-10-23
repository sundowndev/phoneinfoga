<template>
  <div>
    <b-form @submit="onSubmit">
      <b-form-group
        id="input-group-1"
        label="Phone number :"
        label-for="input-1"
        description="Only accepts E164 and International formats as input."
      >
        <b-input-group>
          <!-- <b-form-input
            id="input-number"
            v-model="inputNumber"
            type="text"
            required
            placeholder="e.g: 33678132393"
            :disabled="loading"
          ></b-form-input> -->
          <VuePhoneNumberInput
            v-model="inputNumberVal"
            :disabled="loading"
            @update="updateInputNumber"
          />
          <b-button
            size="sm"
            variant="dark"
            v-on:click="runScans"
            :disabled="loading"
          >
            <b-icon-play-fill></b-icon-play-fill>
            Lookup
          </b-button>

          <b-button
            variant="light"
            size="sm"
            v-on:click="clearData"
            v-show="number"
            :disabled="loading"
            >Reset
          </b-button>
        </b-input-group>
      </b-form-group>
    </b-form>

    <hr />

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

    <!-- <LocalScan :scan="scanEvent" />
    <NumverifyScan :scan="scanEvent" />
    <GoogleSearch :scan="scanEvent" />
    <OVHScan :scan="scanEvent" /> -->
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapMutations, mapState } from "vuex";
import { formatNumber, isValid } from "../utils";
import VuePhoneNumberInput from "vue-phone-number-input";
import Scanner from "../components/Scanner.vue";
// import LocalScan from "../components/LocalScan.vue";
// import NumverifyScan from "../components/NumverifyScan.vue";
// import GoogleSearch from "../components/GoogleSearch.vue";
// import OVHScan from "../components/OVHScan.vue";
import axios, { AxiosResponse } from "axios";
import config from "@/config";

interface Data {
  loading: boolean;
  isLookup: boolean;
  inputNumber: string;
  inputNumberVal: string;
  scanEvent: Vue;
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
  components: { Scanner, VuePhoneNumberInput },
  computed: {
    ...mapState(["number"]),
    ...mapMutations(["pushError"]),
  },
  data(): Data {
    return {
      loading: false,
      isLookup: false,
      inputNumber: "",
      inputNumberVal: "",
      scanEvent: new Vue(),
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
  methods: {
    clearData() {
      // this.scanEvent.$emit("clear");
      this.isLookup = false;
      this.$store.commit("resetState");
    },
    async runScans(): Promise<void> {
      if (!isValid(this.inputNumber)) {
        this.$store.commit("pushError", { message: "Number is not valid." });
        return;
      }

      this.loading = true;

      this.$store.commit("setNumber", formatNumber(this.inputNumber));

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
      // this.scanEvent.$emit("scan");

      // this.scanEvent.$on("finished", () => {
      //   this.loading = false;
      // });
    },
    onSubmit(evt: Event) {
      evt.preventDefault();
    },
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    updateInputNumber(val: any) {
      this.inputNumber = val.e164;
    },
  },
});
</script>

<style src="vue-phone-number-input/dist/vue-phone-number-input.css"></style>
