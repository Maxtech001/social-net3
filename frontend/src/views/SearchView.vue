<template>
    <div v-if="user" class="d-flex flex-column" style="height: 100%">
        <HeaderComponent :searchQuery="searchInput" :user="user"></HeaderComponent>
        <div class="d-flex flex-grow-1" style="height: 100%">
            <v-sheet class="pa-2" min-width="360">
                <v-row class="ma-0">
                    <h2 style="font-weight: bold">Search results</h2>
                </v-row>
                <v-list
                    class="mt-2"
                    style="
                        user-select: none;
                        -webkit-user-select: none;
                        -moz-user-select: none;
                        -khtml-user-select: none;
                        -ms-user-select: none;
                    "
                >
                    <v-list-item :active="view == 'people'" rounded @click="searchPeople()">
                        <v-row class="ma-0 align-center"
                            ><v-icon size="24" class="mr-2">mdi-account-multiple-outline</v-icon
                            ><span style="font-weight: 500; font-size: 18px">People</span></v-row
                        >
                    </v-list-item>
                    <v-list-item :active="view == 'groups'" rounded @click="searchGroups()">
                        <v-row class="ma-0 align-center"
                            ><v-icon size="24" class="mr-2">mdi-account-group-outline</v-icon
                            ><span style="font-weight: 500; font-size: 18px">Groups</span></v-row
                        >
                    </v-list-item>
                </v-list>
            </v-sheet>
            <div v-if="view == 'people'" class="pa-8" style="width: 100%">
                <div class="d-flex flex-column align-center">
                    <v-card v-for="user in resUsers" :key="user.userId" width="680" class="mb-4">
                        <v-row class="ma-0 pa-4 align-center">
                            <div
                                class="mr-4"
                                style="cursor: pointer"
                                @click="navigateProfile(user.userId)"
                            >
                                <v-icon v-if="user.avatar" size="60">
                                    <v-img
                                        :src="user.avatar"
                                        :width="60"
                                        :height="60"
                                        style="border-radius: 50%"
                                        cover
                                    />
                                </v-icon>
                                <v-icon v-else size="60">mdi-account-circle</v-icon>
                            </div>
                            <div>
                                <h3
                                    class="group-header"
                                    style="font-weight: 500; cursor: pointer"
                                    @click="navigateProfile(user.userId)"
                                >
                                    {{ `${user.firstName} ${user.lastName}` }}
                                </h3>
                            </div>
                        </v-row>
                    </v-card>
                </div>
            </div>
            <div v-else-if="view == 'groups'" class="pa-8" style="width: 100%">
                <div class="d-flex flex-column align-center">
                    <v-card
                        v-for="group in resGroups"
                        :key="group.groupId"
                        width="680"
                        class="mb-4"
                    >
                        <v-row class="ma-0 pa-4 align-center">
                            <div
                                class="mr-4"
                                style="cursor: pointer"
                                @click="navigateGroupView(group.groupId)"
                            >
                                <v-img
                                    src="/group-stock-image.jpg"
                                    height="60"
                                    width="60"
                                    cover
                                    style="border-radius: 8px"
                                ></v-img>
                            </div>
                            <div>
                                <h3
                                    class="group-header"
                                    style="font-weight: 500; cursor: pointer"
                                    @click="navigateGroupView(group.groupId)"
                                >
                                    {{ group.groupName }}
                                </h3>
                                <span>{{ group.memberCount }} members</span>
                            </div>
                        </v-row>
                    </v-card>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import HeaderComponent from '../components/header/HeaderComponent.vue'
import axios from 'axios'

export default {
    components: {
        HeaderComponent
    },
    data() {
        return {
            user: null,
            searchInput: this.$route.params.searchInput,
            view: 'people',
            resUsers: [],
            resGroups: []
        }
    },
    created() {
        axios.get('/api/getCurrentUser').then((res) => {
            this.user = res.data
            if (!this.user) return this.$router.push({ path: '/login' })
            axios
                .post('/api/getUsersByName', {
                    nameSearch: this.searchInput
                })
                .then((res) => (this.resUsers = res.data))
        })
    },
    methods: {
        searchPeople() {
            this.view = 'people'
            axios
                .post('/api/getUsersByName', {
                    nameSearch: this.searchInput
                })
                .then((res) => (this.resUsers = res.data))
        },
        searchGroups() {
            this.view = 'groups'
            axios
                .post('/api/getGroupsByName', {
                    groupName: this.searchInput
                })
                .then((res) => (this.resGroups = res.data))
        },
        navigateGroupView(groupId) {
            this.$router.push({ name: 'group', params: { groupId: groupId } })
        },
        navigateProfile(userId) {
            this.$router.push({ name: 'user', params: { userId: userId } })
        }
    }
}
</script>

<style scoped>
.group-header:hover {
    text-decoration: underline;
}
</style>
