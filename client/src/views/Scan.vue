<template>
  <div>
    <b-form @submit="onSubmit">
      <b-form-group
        id="input-group-1"
        label="Phone number :"
        label-for="input-1"
        description="Only accepts E164 and International formats as input."
      >
        <b-form-input
          id="input-number"
          v-model="inputNumber"
          type="text"
          required
          placeholder="e.g: 33678132393"
          :disabled="loading"
        ></b-form-input>
      </b-form-group>

      <b-button
        size="sm"
        variant="dark"
        v-on:click="runScans"
        :disabled="loading"
        class="m-1"
      >
        <b-icon-play-fill></b-icon-play-fill>
        Run scan
      </b-button>

      <b-button
        variant="light"
        size="sm"
        v-on:click="clearData"
        v-show="number"
        :disabled="loading"
        >Reset
      </b-button>
    </b-form>

    <hr />

    <LocalScan :scan="scanEvent" />
    <NumverifyScan :scan="scanEvent" />
    <GoogleSearch :scan="scanEvent" />
    <OVHScan :scan="scanEvent" />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapMutations, mapState } from "vuex";
import { formatNumber, isValid } from "../utils";
import LocalScan from "../components/LocalScan.vue";
import NumverifyScan from "../components/NumverifyScan.vue";
import GoogleSearch from "../components/GoogleSearch.vue";
import OVHScan from "../components/OVHScan.vue";
import { AxiosResponse } from "axios";

interface Data {
  loading: boolean;
  inputNumber: string;
  scanEvent: Vue;
}

export type ScanResponse<T> = AxiosResponse<{
  success: boolean;
  result: T;
  error: string;
}>;

export default Vue.extend({
  components: { LocalScan, GoogleSearch, NumverifyScan, OVHScan },
  computed: {
    ...mapState(["number"]),
    ...mapMutations(["pushError"]),
  },
  data(): Data {
    return {
      loading: false,
      inputNumber: "",
      scanEvent: new Vue(),
    };
  },
  methods: {
    clearData() {
      this.scanEvent.$emit("clear");
      this.$store.commit("resetState");
    },
    async runScans(): Promise<void> {
      if (!isValid(this.inputNumber)) {
        this.$store.commit("pushError", { message: "Number is not valid." });
        return;
      }

      this.loading = true;

      this.$store.commit("setNumber", formatNumber(this.inputNumber));

      this.scanEvent.$emit("scan");

      this.scanEvent.$on("finished", () => {
        this.loading = false;
      });
    },
    onSubmit(evt: Event) {
      evt.preventDefault();
    },
  },
});
</script>
