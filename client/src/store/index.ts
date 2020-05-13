import Vue from "vue";
import Vuex, { Store } from "vuex";

Vue.use(Vuex);

interface ErrorAlert {
  message: string;
}

interface StoreInterface {
  number: string;
  errors: ErrorAlert[];
}

const store: Store<StoreInterface> = new Vuex.Store({
  state: {
    number: "",
    errors: [] as ErrorAlert[],
  },
  mutations: {
    pushError(state, error: ErrorAlert): void {
      state.errors.push(error);
    },
    setNumber(state, number: string): void {
      state.number = number;
    },
    resetState(state) {
      state.number = "";
      state.errors = [];

      return state;
    },
  },
  getters: {},
  actions: {
    resetState(context): void {
      context.commit("resetState");
    },
  },
  modules: {},
});

export default store;
