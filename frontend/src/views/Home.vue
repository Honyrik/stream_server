<template>
    <v-flex>
        <v-list>
            <c-row v-for="row in rows" :key="row.fk_id" :row="row" />
        </v-list>
    </v-flex>
</template>

<script lang="ts">
import Vue from "vue";
import VideoFile from "@/components/VideoFile.vue";
import Row from "@/components/Row.vue";
import Folder from "@/components/Folder.vue";
Vue.component("c-row", Row);
Vue.component("c-video-file", VideoFile);
Vue.component("c-folder", Folder);

export default Vue.extend({
    name: "home",
    data() {
        return {
            isLoaded: false,
            rows: [],
        };
    },
    mounted() {
        this.$http
            .get("/api/fileTree")
            .then((res) => {
                this.rows = res.data;
                this.isLoaded = true;
            })
            .catch(() => {
                this.isLoaded = true;
            });
    },
});
</script>
