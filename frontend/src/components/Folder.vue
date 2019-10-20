<template>
    <div>
        <v-list-group v-if="!parent" no-action>
            <template v-slot:activator>
                <v-list-item-content>
                    <v-list-item-title v-text="row.fv_name"></v-list-item-title>
                </v-list-item-content>
            </template>

            <c-row
                v-for="r in row.fct_children"
                :key="r.fk_id"
                :row="r"
                :parent="row"
            />
        </v-list-group>
        <v-list-group v-if="parent" sub-group no-action>
            <template v-slot:activator>
                <v-list-item-content>
                    <v-list-item-title v-text="row.fv_name"></v-list-item-title>
                </v-list-item-content>
            </template>

            <c-row
                v-for="r in row.fct_children"
                :key="r.fk_id"
                :row="r"
                :parent="row"
            />
        </v-list-group>
    </div>
</template>
<script lang="ts">
import Vue from "vue";
export default Vue.extend({
    props: ["row", "parent"],
});
</script>
