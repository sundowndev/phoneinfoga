<template>
  <div>
    <b-form @submit="onSubmit" class="d-flex justify-content-center mt-5">
      <b-form-group id="input-group-1" label-for="input-1">
        <b-input-group>
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

    <b-container v-if="isLookup || showInformations" class="border p-4 mb-3">
      <h3 class="text-center">Information</h3>
      <b-container>
        <b-row v-for="(value, name) in localData" :key="name" align-v="center">
          <h5 class="text-capitalize m-0 mr-4">{{ name }}:</h5>
          <p class="m-0">{{ value }}</p>
        </b-row>
      </b-container>
    </b-container>

    <b-container v-if="isLookup" class="border p-4">
      <h3 class="text-center">Scanners</h3>
      <Scanner
        v-for="(scanner, index) in scanners"
        :key="index"
        :name="scanner.name.charAt(0).toUpperCase() + scanner.name.slice(1)"
        :scanId="scanner.name"
      />
    </b-container>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapMutations, mapState } from "vuex";
import { formatNumber, isValid } from "../utils";
import VuePhoneNumberInput from "vue-phone-number-input";
import Scanner from "../components/Scanner.vue";
import axios, { AxiosResponse } from "axios";
import config from "@/config";

interface InputNumberObject {
  countryCallingCode: string;
  countryCode: string;
  e164: string;
  formatInternational: string;
  formatNational: string;
  formattedNumber: string;
  isValid: boolean;
  nationalNumber: string;
  phoneNumber: string;
  type: string;
  uri: string;
}

interface ScannerObject {
  name: string;
  description: string;
}

interface Data {
  loading: boolean;
  isLookup: boolean;
  showInformations: boolean;
  inputNumber: string;
  inputNumberVal: string;
  scanEvent: Vue;
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
  components: { Scanner, VuePhoneNumberInput },
  computed: {
    ...mapState(["number"]),
    ...mapMutations(["pushError"]),
  },
  data(): Data {
    return {
      loading: false,
      isLookup: false,
      showInformations: false,
      inputNumber: "",
      inputNumberVal: "",
      scanEvent: new Vue(),
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
  methods: {
    clearData() {
      this.isLookup = false;
      this.showInformations = false;
      this.$store.commit("resetState");
    },
    async runScans(): Promise<void> {
      this.clearData();
      if (!isValid(this.inputNumber)) {
        this.$store.commit("pushError", { message: "Number is not valid." });
        return;
      }

      this.loading = true;

      this.$store.commit("setNumber", formatNumber(this.inputNumber));

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
    onSubmit(evt: Event) {
      evt.preventDefault();
    },
    updateInputNumber(val: InputNumberObject) {
      this.inputNumber = val.e164;
    },
    async getScanners() {
      try {
        const res = await axios.get(`${config.apiUrl}/v2/scanners`);

        this.scanners = res.data.scanners.filter(
          (scanner: ScannerObject) => scanner.name !== "local"
        );
      } catch (error) {
        this.$store.commit("pushError", { message: error });
      }
    },
  },
});
</script>

<style src="vue-phone-number-input/dist/vue-phone-number-input.css"></style>
