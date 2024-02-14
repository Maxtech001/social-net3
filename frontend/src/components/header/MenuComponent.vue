<template>
    <v-menu width="auto" style="border-radius: 8px !important">
        <template v-slot:activator="{ props }">
            <v-btn icon height="40" width="40" class="mr-2" variant="outlined" v-bind="props">
                <v-icon size="20" color="black">mdi-dots-grid</v-icon>
            </v-btn>
        </template>
        <v-list
            @click.stop.prevent
            class="pa-4"
            style="
                user-select: none;
                -webkit-user-select: none;
                -moz-user-select: none;
                -khtml-user-select: none;
                -ms-user-select: none;
            "
        >
            <h2 class="px-2 pb-4">Menu</h2>
            <v-row class="ma-0">
                <v-col class="pa-0 mr-4">
                    <v-card class="pa-1" width="400">
                        <v-list-item rounded @click="navigateFollowView()">
                            <v-icon class="mr-2" size="36" color="blue"
                                >mdi-account-multiple</v-icon
                            >
                            Followers
                        </v-list-item>
                        <v-list-item rounded @click="navigateGroupsView()">
                            <v-icon class="mr-2" size="36" color="blue"
                                >mdi-account-group-outline</v-icon
                            >
                            Groups
                        </v-list-item>
                    </v-card>
                </v-col>
                <v-col class="pa-0">
                    <v-card class="pa-2" width="200">
                        <h2 class="pl-2">Create</h2>
                        <v-list-item rounded @click="navigatePostCreate()">
                            <v-icon class="mr-2" size="36">mdi-pencil-box-outline</v-icon>
                            Post
                        </v-list-item>
                        <v-list-item rounded @click="navigateGroupCreate()">
                            <v-icon class="mr-2" size="36">mdi-account-group-outline</v-icon>
                            Group
                        </v-list-item>
                    </v-card>
                </v-col>
            </v-row>
        </v-list>
    </v-menu>
</template>

<script>
export default {
    methods: {
        navigateGroupCreate() {
            this.$router.push({ path: '/groups/create' })
        },
        navigateFollowView() {
            this.$router.push({ path: '/follows' })
        },
        navigateGroupsView() {
            this.$router.push({ path: '/groups' })
        },
        navigatePostCreate() {
            if (this.$route.name == 'createPost') return this.$router.go()
            this.$router.push({ name: 'createPost' })
        }
    }
}
</script>
