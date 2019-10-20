import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        videoUuid: null,
    },
    mutations: {
        setUuid(state, obj) {
            state.videoUuid = obj;
        },
    },
    actions: {
        setUuid({ commit }, uuid) {
            commit("setUuid", uuid);
        },
    },
    modules: {},
});
