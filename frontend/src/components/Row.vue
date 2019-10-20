<template>
    <div>
        <c-folder v-if="!row.fl_leaf" :row="row" :parent="parent" />
        <c-video-file v-if="row.fl_leaf" :row="row" />
    </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    props: ["row", "parent"],
});
</script>
