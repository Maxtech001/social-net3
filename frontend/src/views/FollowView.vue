<template>
    <div v-if="user" class="d-flex flex-column" style="height: 100%">
        <HeaderComponent :user="user"></HeaderComponent>
        <div class="d-flex flex-grow-1" style="height: 100%">
            <v-sheet class="pa-2" min-width="360">
                <v-row class="ma-0">
                    <h2 style="font-weight: bold">Follows</h2>
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
                    <v-list-item
                        :active="view == 'followRequests'"
                        rounded
                        @click="view = 'followRequests'"
                    >
                        <v-row class="ma-0 align-center"
                            ><v-icon size="24" class="mr-2">mdi-account-plus-outline</v-icon
                            ><span style="font-weight: 500; font-size: 18px"
                                >Follow requests</span
                            ></v-row
                        >
                    </v-list-item>
                    <v-list-item :active="view == 'followers'" rounded @click="view = 'followers'">
                        <v-row class="ma-0 align-center"
                            ><v-icon size="24" class="mr-2">mdi-account-arrow-left-outline</v-icon
                            ><span style="font-weight: 500; font-size: 18px">Followers</span></v-row
                        >
                    </v-list-item>
                    <v-list-item :active="view == 'following'" rounded @click="view = 'following'">
                        <v-row class="ma-0 align-center"
                            ><v-icon size="24" class="mr-2">mdi-account-arrow-right-outline</v-icon
                            ><span style="font-weight: 500; font-size: 18px">Following</span></v-row
                        >
                    </v-list-item>
                </v-list>
            </v-sheet>
            <div v-if="view == 'followRequests'" class="pa-5">
                <div class="px-4">
                    <h3 class="pt-5 pb-4" style="font-weight: bold; font-size: 20px">
                        Follow requests ({{ followRequests.length }})
                    </h3>
                    <div class="d-flex flex-row flex-wrap">
                        <div
                            class="d-flex flex-grow-1 flex-shrink-1 pa-1"
                            v-for="user in followRequests"
                            :key="user.userId"
                            style="max-width: 250px; min-width: 200px; flex-basis: 0"
                        >
                            <v-card style="width: 100%">
                                <div style="cursor: pointer" @click="navigateProfile(user.userId)">
                                    <img
                                        v-if="user.avatar"
                                        :src="user.avatar"
                                        style="max-height: 242px"
                                    />
                                    <img
                                        v-else
                                        src="/account.svg"
                                        style="background-color: #c9ccd1"
                                        cover
                                    />
                                </div>
                                <div class="pa-3 d-flex align-center flex-column">
                                    <h4
                                        style="cursor: pointer"
                                        class="profile-name"
                                        @click="navigateProfile(user.userId)"
                                    >
                                        {{ `${user.firstName} ${user.lastName}` }}
                                    </h4>
                                    <v-btn
                                        color="blue"
                                        block
                                        class="mt-2"
                                        @click="acceptRequest(user)"
                                        >Accept</v-btn
                                    >
                                    <v-btn block class="mt-2" @click="discardRequest(user)"
                                        >Decline</v-btn
                                    >
                                </div>
                            </v-card>
                        </div>
                    </div>
                </div>
            </div>
            <div v-if="view == 'followers'" class="pa-5" style="width: 100%">
                <div class="px-4">
                    <h3 class="pt-5 pb-4" style="font-weight: bold; font-size: 20px">
                        Followers ({{ followers.length }})
                    </h3>
                    <div class="d-flex flex-row flex-wrap">
                        <div
                            class="d-flex flex-grow-1 flex-shrink-1 pa-1"
                            v-for="user in followers"
                            :key="user.userId"
                            style="max-width: 250px; min-width: 200px; flex-basis: 0"
                        >
                            <v-card style="width: 100%">
                                <div style="cursor: pointer" @click="navigateProfile(user.userId)">
                                    <img
                                        v-if="user.avatar"
                                        :src="user.avatar"
                                        style="max-height: 242px"
                                    />
                                    <img
                                        v-else
                                        src="/account.svg"
                                        style="background-color: #c9ccd1"
                                    />
                                </div>
                                <div class="pa-3 d-flex justify-center">
                                    <h4
                                        style="cursor: pointer"
                                        class="profile-name"
                                        @click="navigateProfile(user.userId)"
                                    >
                                        {{ `${user.firstName} ${user.lastName}` }}
                                    </h4>
                                </div>
                            </v-card>
                        </div>
                    </div>
                </div>
            </div>
            <div v-if="view == 'following'" class="pa-5" style="width: 100%">
                <div class="px-4">
                    <h3 class="pt-5 pb-4" style="font-weight: bold; font-size: 20px">
                        Following ({{ following.length }})
                    </h3>
                    <div class="d-flex flex-row flex-wrap">
                        <div
                            class="d-flex flex-grow-1 flex-shrink-1 pa-1"
                            v-for="user in following"
                            :key="user.userId"
                            style="max-width: 250px; min-width: 200px; flex-basis: 0"
                        >
                            <v-card style="width: 100%">
                                <div style="cursor: pointer" @click="navigateProfile(user.userId)">
                                    <img
                                        v-if="user.avatar"
                                        :src="user.avatar"
                                        style="max-height: 242px"
                                    />
                                    <img
                                        v-else
                                        src="/account.svg"
                                        style="background-color: #c9ccd1"
                                    />
                                </div>
                                <div class="pa-3 d-flex justify-center">
                                    <h4
                                        style="cursor: pointer"
                                        class="profile-name"
                                        @click="navigateProfile(user.userId)"
                                    >
                                        {{ `${user.firstName} ${user.lastName}` }}
                                    </h4>
                                </div>
                            </v-card>
                        </div>
                    </div>
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
            view: 'followRequests',
            followers: [],
            following: [],
            followRequests: []
        }
    },
    created() {
        axios.get('/api/getCurrentUser').then((res) => {
            this.user = res.data
            if (!this.user) return this.$router.push({ path: '/login' })
            axios.get('/api/getFollowers').then((res) => (this.followers = res.data))
            axios.get('/api/getFollowed').then((res) => (this.following = res.data))
            axios.get('/api/getFollowRequests').then((res) => (this.followRequests = res.data))
        })
    },
    methods: {
        navigateProfile(userId) {
            this.$router.push({ name: 'user', params: { userId: userId } })
        },
        acceptRequest(user) {
            axios
                .post('/api/acceptRequest', {
                    followedId: this.user.userId,
                    followerId: user.userId
                })
                .then(() => {
                    this.followRequests = this.followRequests.filter((v) => v.userId != user.userId)
                    this.followers.push(user)
                })
        },
        discardRequest(user) {
            axios
                .post('/api/discardRequest', {
                    followedId: this.user.userId,
                    followerId: user.userId
                })
                .then(() => {
                    this.followRequests = this.followRequests.filter((v) => v.userId != user.userId)
                })
        }
    }
}
</script>

<style scoped>
.profile-name:hover {
    text-decoration: underline;
}
</style>
