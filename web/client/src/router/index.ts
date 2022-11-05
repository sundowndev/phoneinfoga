import Vue from "vue";
import VueRouter from "vue-router";
import Scan from "../views/Scan.vue";
import Number from "../views/Number.vue";
import NotFound from "../views/NotFound.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Scan",
    component: Scan,
  },
  {
    path: "/numbers/:number",
    name: "Number",
    component: Number,
  },
  {
    path: "*",
    name: "NotFound",
    component: NotFound,
  },
];

const router = new VueRouter({
  mode: "hash",
  base: process.env.BASE_URL,
  routes,
});

export default router;
