<template>
    <div v-if="currentUser && viewedUser">
        <HeaderComponent :user="currentUser"></HeaderComponent>
        <div
            class="d-flex justify-center"
            style="background-color: white; box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1)"
        >
            <div
                style="
                    background-image: url('/user-stock-image.png');
                    background-position: center;
                    background-repeat: no-repeat;
                    background-size: cover;
                    height: 500px;
                    width: 100%;
                    position: absolute;
                    filter: blur(3px);
                    -webkit-filter: blur(4px);
                    z-index: 0;
                "
            ></div>
            <v-row class="ma-0" style="min-width: 1000px; max-width: 1200px">
                <v-col cols="12" class="pa-0">
                    <v-img height="500" cover src="/user-stock-image.png"></v-img>
                </v-col>
                <div class="d-flex pa-4" style="width: 100%">
                    <div class="d-flex align-end" style="max-height: 100px">
                        <v-icon
                            v-if="viewedUser.avatar"
                            class="mr-4"
                            size="168"
                            style="z-index: 1; overflow: hidden"
                        >
                            <v-img
                                :src="viewedUser.avatar"
                                :width="168"
                                :height="168"
                                style="border-radius: 50%; border: 5px solid white"
                                cover
                            />
                        </v-icon>
                        <v-icon
                            v-else
                            size="168"
                            class="mr-4"
                            style="z-index: 1; border-radius: 50%; background-color: white"
                            >mdi-account-circle</v-icon
                        >
                    </div>
                    <v-row class="ma-0">
                        <div
                            class="d-flex"
                            style="flex-direction: column; justify-content: flex-end"
                        >
                            <h1 style="font-weight: bold">
                                {{ `${viewedUser.firstName} ${viewedUser.lastName}` }}
                            </h1>
                            <span>{{ `${viewedUser.followers} followers` }}</span>
                        </div>
                        <v-spacer></v-spacer>
                        <div v-if="!isViewedUser" class="d-flex mt-auto">
                            <v-btn
                                v-if="viewedUser.isFollowed || viewedUser.isPublic"
                                variant="flat"
                                color="black"
                                rounded
                                class="mr-2"
                                @click="openChat()"
                            >
                                <v-icon class="mr-1">mdi-chat-outline</v-icon>Message
                            </v-btn>
                            <v-btn
                                v-if="viewedUser.isRequested"
                                variant="flat"
                                rounded
                                color="red"
                                @click="discardRequest()"
                                ><v-icon class="mr-1">mdi-close</v-icon>Cancel Request</v-btn
                            >
                            <v-btn
                                v-else
                                variant="flat"
                                rounded
                                :color="viewedUser.isFollowed ? 'red' : 'blue'"
                                @click="viewedUser.isFollowed ? unfollowUser() : followUser()"
                            >
                                {{ viewedUser.isFollowed ? 'Unfollow' : 'Follow' }}</v-btn
                            >
                        </div>
                        <div v-else class="d-flex mt-auto">
                            <v-menu width="200">
                                <template v-slot:activator="{ props }">
                                    <v-btn v-bind="props" variant="flat" rounded
                                        ><v-icon>{{
                                            viewedUser.isPublic
                                                ? 'mdi-eye-outline'
                                                : 'mdi-eye-off-outline'
                                        }}</v-icon
                                        ><v-icon>mdi-chevron-down</v-icon></v-btn
                                    >
                                </template>
                                <v-list class="pa-2">
                                    <v-list-item rounded @click="setPublic()">
                                        <v-row class="ma-0">
                                            <span>Public</span>
                                            <v-spacer></v-spacer>
                                            <v-icon v-if="viewedUser.isPublic">mdi-check</v-icon>
                                        </v-row>
                                    </v-list-item>
                                    <v-list-item rounded @click="setPrivate()">
                                        <v-row class="ma-0">
                                            <span>Private</span>
                                            <v-spacer></v-spacer>
                                            <v-icon v-if="!viewedUser.isPublic">mdi-check</v-icon>
                                        </v-row>
                                    </v-list-item>
                                </v-list>
                            </v-menu>
                        </div>
                    </v-row>
                </div>
                <v-tabs
                    v-if="isViewedUser || viewedUser.isFollowed || viewedUser.isPublic"
                    class="ml-4"
                    v-model="tab"
                >
                    <v-tab value="posts">Posts</v-tab>
                    <v-tab value="followers">Followers</v-tab>
                    <v-tab value="following">Following</v-tab>
                </v-tabs>
            </v-row>
        </div>
        <div
            v-if="isViewedUser || viewedUser.isFollowed || viewedUser.isPublic"
            class="d-flex justify-center"
        >
            <v-window
                v-model="tab"
                class="pb-1"
                style="min-width: 1000px; max-width: 1200px; flex-grow: 1"
            >
                <v-window-item value="posts">
                    <v-row class="ma-0">
                        <v-col cols="4" class="mt-6 pa-0 pr-6">
                            <v-card class="pa-4">
                                <div v-if="viewedUser.aboutMe">
                                    <h1>About me</h1>
                                    <div>{{ viewedUser.aboutMe }}</div>
                                    <v-divider class="mt-3"></v-divider>
                                </div>
                                <h1>Info</h1>
                                <div>Email: {{ viewedUser.email }}</div>
                                <div>
                                    Name: {{ `${viewedUser.firstName} ${viewedUser.lastName}` }}
                                </div>
                                <div v-if="viewedUser.nickname">
                                    Nickname: {{ viewedUser.nickname }}
                                </div>
                                <div>Date of birth: {{ viewedUser.birthDate }}</div>
                            </v-card>
                        </v-col>
                        <v-col cols="8" class="pa-0" v-bind:class="isViewedUser ? 'mt-6' : ''">
                            <PostSection
                                v-if="posts"
                                :propPosts="posts"
                                :isViewedUser="isViewedUser"
                                :viewedUser="viewedUser"
                                :user="currentUser"
                            ></PostSection>
                        </v-col>
                    </v-row>
                </v-window-item>
                <v-window-item value="followers">
                    <v-row class="ma-0 justify-center">
                        <v-card class="mt-6 pa-4 flex-grow-1" max-width="600">
                            <h2 style="font-weight: 600">Followers • {{ viewedUser.followers }}</h2>
                            <v-divider class="mb-2 mt-1" />
                            <v-row
                                class="ma-0 align-center"
                                v-for="follower in followerUsers"
                                :key="follower.userId"
                            >
                                <v-icon v-if="follower.avatar" class="mr-4" size="60">
                                    <v-img
                                        :src="follower.avatar"
                                        :width="60"
                                        :height="60"
                                        style="border-radius: 50%"
                                        cover
                                    />
                                </v-icon>
                                <v-icon v-else class="mr-4" size="60">mdi-account-circle</v-icon>
                                <h3 style="font-weight: 500; line-height: 1">
                                    {{ `${follower.firstName} ${follower.lastName}` }}
                                </h3>
                            </v-row>
                        </v-card>
                    </v-row>
                </v-window-item>
                <v-window-item value="following">
                    <v-row class="ma-0 justify-center">
                        <v-card class="mt-6 pa-4 flex-grow-1" max-width="600">
                            <h2 style="font-weight: 600">
                                Following • {{ followingUsers.length }}
                            </h2>
                            <v-divider class="mb-2 mt-1" />
                            <v-row
                                class="ma-0 align-center"
                                v-for="following in followingUsers"
                                :key="following.userId"
                            >
                                <v-icon v-if="following.avatar" class="mr-4" size="60">
                                    <v-img
                                        :src="following.avatar"
                                        :width="60"
                                        :height="60"
                                        style="border-radius: 50%"
                                        cover
                                    />
                                </v-icon>
                                <v-icon v-else class="mr-4" size="60">mdi-account-circle</v-icon>
                                <h3 style="font-weight: 500; line-height: 1">
                                    {{ `${following.firstName} ${following.lastName}` }}
                                </h3>
                            </v-row>
                        </v-card>
                    </v-row>
                </v-window-item>
            </v-window>
        </div>
    </div>
</template>

<script>
import HeaderComponent from '../components/header/HeaderComponent.vue'
import PostSection from '../components/forum/PostSection.vue'
import axios from 'axios'
import chat from '../chat.js'
export default {
    components: {
        HeaderComponent,
        PostSection
    },
    created() {
        // axios interceptors
        axios.interceptors.request.use(
            (config) => {
                this.loading = true
                return config
            },
            (error) => {
                this.loading = false
                return Promise.reject(error)
            }
        )

        axios.interceptors.response.use(
            (response) => {
                this.loading = false
                return response
            },
            (error) => {
                this.loading = false
                return Promise.reject(error)
            }
        )

        axios.get('/api/getCurrentUser').then((res) => {
            this.currentUser = res.data
            if (!this.currentUser) return this.$router.push({ path: '/login' })
            axios
                .post('/api/getUser', {
                    userId: parseInt(this.$route.params.userId)
                })
                .then((res) => {
                    this.viewedUser = res.data
                    this.isViewedUser = this.viewedUser.userId == this.currentUser.userId
                })
            this.getPosts()
            this.getFollowers()
            this.getFollowing()
        })
    },
    data() {
        return {
            currentUser: null,
            viewedUser: null,
            isViewedUser: false,
            posts: null,
            loading: false,
            tab: 'posts',
            followerUsers: [],
            followingUsers: []
        }
    },
    methods: {
        followUser() {
            axios
                .post('/api/follow', {
                    followedId: this.viewedUser.userId,
                    followerId: this.currentUser.userId
                })
                .then((res) => {
                    this.viewedUser = res.data
                    this.getFollowers()
                    this.getPosts()
                })
        },
        unfollowUser() {
            axios
                .post('/api/unfollow', {
                    followedId: this.viewedUser.userId,
                    followerId: this.currentUser.userId
                })
                .then((res) => {
                    this.viewedUser = res.data
                    this.posts = null
                })
        },
        discardRequest() {
            axios
                .post('/api/discardRequest', {
                    followedId: this.viewedUser.userId,
                    followerId: this.currentUser.userId
                })
                .then(() => (this.viewedUser.isRequested = false))
        },
        getPosts() {
            axios
                .post('/api/getUserPosts', {
                    userId: parseInt(this.$route.params.userId)
                })
                .then((res) => (this.posts = res.data))
        },
        setPublic() {
            axios.get('/api/setPublic').then(() => (this.viewedUser.isPublic = true))
        },
        setPrivate() {
            axios.get('/api/setPrivate').then(() => (this.viewedUser.isPublic = false))
        },
        openChat() {
            chat.openChat(this.viewedUser)
        },
        getFollowers() {
            axios
                .post('/api/getUserFollowers', {
                    userId: parseInt(this.$route.params.userId)
                })
                .then((res) => (this.followerUsers = res.data))
        },
        getFollowing() {
            axios
                .post('/api/getUserFollowing', {
                    userId: parseInt(this.$route.params.userId)
                })
                .then((res) => (this.followingUsers = res.data))
        }
    }
}
</script>
