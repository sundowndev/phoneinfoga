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
          v-model="number"
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
        <b-icon-play-fill></b-icon-play-fill>Run scan
      </b-button>

      <b-button
        variant="light"
        size="sm"
        v-on:click="clearData"
        :disabled="loading"
        >Clear results</b-button
      >
    </b-form>

    <hr />

    <LocalScan :scan="scanEvent" />
    <NumverifyScan :scan="scanEvent" />
    <GoogleSearch :scan="scanEvent" />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapMutations } from "vuex";
import LocalScan from "../components/LocalScan.vue";
import NumverifyScan from "../components/NumverifyScan.vue";
import GoogleSearch from "../components/GoogleSearch.vue";

interface Scanner {
  id: string;
  name: string;
  data: object[];
  loading: boolean;
}

interface Data {
  loading: boolean;
  number: string;
  scanEvent: Vue;
}

export default Vue.extend({
  components: { LocalScan, GoogleSearch, NumverifyScan },
  computed: {
    ...mapMutations(["pushError"])
  },
  data(): Data {
    return {
      loading: false,
      number: "",
      scanEvent: new Vue()
    };
  },
  methods: {
    clearData() {
      this.scanEvent.$emit("clear");
    },
    async runScans(): Promise<void> {
      if (this.number.length < 2) {
        this.$store.commit("pushError", { message: "Number is not valid." });
        return;
      }

      this.loading = true;

      this.$store.commit("setNumber", this.number);

      this.scanEvent.$emit("scan");

      this.scanEvent.$on("finished", () => {
        this.loading = false;
      });
    },
    onSubmit(evt: any) {
      evt.preventDefault();
    }
  }
});
</script>
