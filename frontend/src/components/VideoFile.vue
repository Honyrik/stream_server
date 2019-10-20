<template>
    <v-list-item>
        <v-list-item-icon>
            <v-icon @click="onClick(arg, $event)">mdi-video</v-icon>
        </v-list-item-icon>

        <v-list-item-title>{{ row.fv_name }}</v-list-item-title>
    </v-list-item>
</template>

<script lang="ts">
import Vue from "vue";
export default Vue.extend({
    props: ["row"],
    methods: {
        onClick() {
            this.$http.post("/api/getVideoUuid", this.row).then((res) => {
                this.$store.dispatch("setUuid", res.data.fv_uuid);
                this.$router.push("/video");
            });
        },
    },
});
</script>
