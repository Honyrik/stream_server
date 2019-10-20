import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import "./registerServiceWorker";
import axios from "axios";
import VueAxios from "vue-axios";

Vue.use(VueAxios, axios);
import vuetify from "./plugins/vuetify";

Vue.config.productionTip = false;

new Vue({
    router,
    store,
    vuetify,
    render: (h: any) => h(App),
} as any).$mount("#app");
