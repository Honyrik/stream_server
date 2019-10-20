<template>
    <v-flex>
        <v-button @click="onBack">Назад</v-button>
        <vue-plyr v-if="isVideo">
            <video>
                <source
                    :src="urlVideo"
                    type='video/webm;codecs="vp8, vorbis"'
                />
            </video>
        </vue-plyr>
    </v-flex>
</template>

<script>
import { VuePlyr } from "vue-plyr";
import "vue-plyr/dist/vue-plyr.css";

export default {
    components: {
        "vue-plyr": VuePlyr,
    },
    mounted() {
        if (!this.$store.state.videoUuid) {
            this.$router.push("/");
        }
    },
    methods: {
        onBack() {
            this.$router.push("/");
        },
    },
    computed: {
        isVideo() {
            return typeof this.$store.state.videoUuid === "string";
        },
        urlVideo() {
            return `/api/getVideo?uuid=${this.$store.state.videoUuid}`;
        },
    },
};
</script>
