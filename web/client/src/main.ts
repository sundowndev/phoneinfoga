import Vue from "vue";
import { BootstrapVue, IconsPlugin } from "bootstrap-vue";

import App from "./App.vue";
import router from "./router";
import store from "./store";

Vue.config.productionTip = false;

// Install BootstrapVue
Vue.use(BootstrapVue);
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin);

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount("#app");
