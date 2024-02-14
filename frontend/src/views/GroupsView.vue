<template>
    <div v-if="user" style="height: 100%">
        <HeaderComponent :user="user"></HeaderComponent>
        <div class="d-flex flex-grow-1" style="height: 100%">
            <v-sheet class="pa-2" min-width="360">
                <v-row class="ma-0">
                    <h2 style="font-weight: bold">Groups</h2>
                </v-row>
                <v-text-field
                    v-model="searchInput"
                    class="mb-2 mt-2 px-2"
                    prepend-inner-icon="mdi-magnify"
                    variant="solo"
                    rounded
                    placeholder="Find groups"
                    density="compact"
                    hide-details
                    style="position: relative"
                    @click="openDialog = true"
                    v-click-outside="() => (openDialog = false)"
                ></v-text-field>
                <div
                    v-if="searchInput.length && openDialog"
                    tabindex="0"
                    class="px-4"
                    style="
                        position: absolute;
                        z-index: 1;
                        background-color: white;
                        width: 360px;
                        left: 0;
                        box-shadow: rgba(0, 0, 0, 0.2) 0px 12px 12px 0px;
                        border-bottom-left-radius: 8px;
                        border-bottom-right-radius: 8px;
                    "
                >
                    <v-list>
                        <v-list-item rounded @click="searchGroups()">
                            <v-icon class="mr-1">mdi-magnify</v-icon>
                            Search groups for {{ searchInput }}
                        </v-list-item>
                    </v-list>
                </div>
                <v-list
                    style="
                        user-select: none;
                        -webkit-user-select: none;
                        -moz-user-select: none;
                        -khtml-user-select: none;
                        -ms-user-select: none;
                    "
                >
                    <v-list-item :active="view == 'feed'" rounded @click="view = 'feed'">
                        <v-row class="ma-0 align-center"
                            ><v-icon size="24" class="mr-2">mdi-post</v-icon
                            ><span style="font-weight: 500; font-size: 18px">Your feed</span></v-row
                        >
                    </v-list-item>
                    <v-list-item :active="view == 'groups'" rounded @click="view = 'groups'">
                        <v-row class="ma-0 align-center"
                            ><v-icon size="24" class="mr-2">mdi-account-group</v-icon
                            ><span style="font-weight: 500; font-size: 18px"
                                >Your groups</span
                            ></v-row
                        >
                    </v-list-item>
                    <v-list-item
                        :active="view == 'groupRequests'"
                        rounded
                        @click="view = 'groupRequests'"
                    >
                        <v-row class="ma-0 align-center"
                            ><v-icon size="24" class="mr-2">mdi-account-multiple-plus</v-icon
                            ><span style="font-weight: 500; font-size: 18px"
                                >Group requests</span
                            ></v-row
                        >
                    </v-list-item>
                    <v-list-item
                        :active="view == 'groupInvites'"
                        rounded
                        @click="view = 'groupInvites'"
                    >
                        <v-row class="ma-0 align-center"
                            ><v-icon size="24" class="mr-2">mdi-account-multiple-check</v-icon
                            ><span style="font-weight: 500; font-size: 18px"
                                >Group invites</span
                            ></v-row
                        >
                    </v-list-item>
                </v-list>
                <v-btn block color="blue" @click="navigateCreateGroup()"
                    ><v-icon>mdi-plus</v-icon>Create group</v-btn
                >
            </v-sheet>
            <div v-if="view == 'feed'" class="pa-5 d-flex justify-center" style="width: 100%">
                <div class="d-flex flex-grow-1" style="min-width: 400px; max-width: 680px">
                    <PostSection :user="user" :isGroupView="true" :propPosts="posts"></PostSection>
                </div>
            </div>
            <div v-if="view == 'groups'" class="px-16" style="width: 100%">
                <h3 class="pt-5 pb-4" style="font-weight: bold; font-size: 20px">
                    All groups you've joined ({{ groups.length }})
                </h3>
                <div class="d-flex flex-grow-1">
                    <div
                        class="d-flex flex-grow-1 flex-shrink-1 pa-1"
                        v-for="group in groups"
                        :key="group.groupId"
                        style="max-width: 600px; min-width: 320px; flex-basis: 0"
                    >
                        <v-card class="pa-4" style="width: 100%">
                            <v-row class="ma-0 align-center">
                                <div
                                    class="mr-4"
                                    style="cursor: pointer"
                                    @click="navigateGroupView(group.groupId)"
                                >
                                    <v-img
                                        src="/group-stock-image.jpg"
                                        height="80"
                                        width="80"
                                        cover
                                        style="border-radius: 8px"
                                    ></v-img>
                                </div>
                                <h3
                                    class="group-header"
                                    style="font-weight: 500; cursor: pointer"
                                    @click="navigateGroupView(group.groupId)"
                                >
                                    {{ group.groupName }}
                                </h3>
                            </v-row>
                            <v-btn
                                class="mt-4"
                                color="blue"
                                block
                                @click="navigateGroupView(group.groupId)"
                                >View group</v-btn
                            >
                        </v-card>
                    </div>
                </div>
            </div>
            <div v-if="view == 'searchGroups'" class="pa-8" style="width: 100%">
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
            <div v-if="view == 'groupRequests'" class="px-16" style="width: 100%">
                <h3 class="pt-5 pb-4" style="font-weight: bold; font-size: 20px">
                    Group join requests ({{ groupRequests.length }})
                </h3>
                <div class="d-flex flex-grow-1">
                    <div
                        class="d-flex flex-grow-1 flex-shrink-1 pa-1"
                        v-for="(groupRequest, index) in groupRequests"
                        :key="index"
                        style="max-width: 600px; min-width: 320px; flex-basis: 0"
                    >
                        <v-card class="pa-4" style="width: 100%">
                            <v-row class="ma-0 align-center">
                                <div
                                    class="mr-4"
                                    style="cursor: pointer"
                                    @click="navigateGroupView(groupRequest.group.groupId)"
                                >
                                    <v-img
                                        src="/group-stock-image.jpg"
                                        height="80"
                                        width="80"
                                        cover
                                        style="border-radius: 8px"
                                    ></v-img>
                                </div>
                                <div>
                                    <h3
                                        class="group-header"
                                        style="font-weight: 500; cursor: pointer"
                                        @click="navigateGroupView(groupRequest.group.groupId)"
                                    >
                                        {{ groupRequest.group.groupName }}
                                    </h3>
                                    <span
                                        >Requester:
                                        <span
                                            class="group-header"
                                            style="cursor: pointer"
                                            @click="navigateProfileView(groupRequest.user.userId)"
                                            >{{
                                                `${groupRequest.user.firstName} ${groupRequest.user.lastName}`
                                            }}</span
                                        ></span
                                    >
                                </div>
                            </v-row>
                            <v-btn
                                class="mt-4"
                                color="blue"
                                block
                                @click="acceptGroupRequest(groupRequest)"
                                >Accept</v-btn
                            >
                            <v-btn class="mt-2" block @click="discardGroupRequest(groupRequest)"
                                >Decline</v-btn
                            >
                        </v-card>
                    </div>
                </div>
            </div>
            <div v-if="view == 'groupInvites'" class="px-16" style="width: 100%">
                <h3 class="pt-5 pb-4" style="font-weight: bold; font-size: 20px">
                    Group join invites ({{ groupInvites.length }})
                </h3>
                <div class="d-flex flex-grow-1">
                    <div
                        class="d-flex flex-grow-1 flex-shrink-1 pa-1"
                        v-for="(groupInvite, index) in groupInvites"
                        :key="index"
                        style="max-width: 600px; min-width: 320px; flex-basis: 0"
                    >
                        <v-card class="pa-4" style="width: 100%">
                            <v-row class="ma-0 align-center">
                                <div
                                    class="mr-4"
                                    style="cursor: pointer"
                                    @click="navigateGroupView(groupInvite.group.groupId)"
                                >
                                    <v-img
                                        src="/group-stock-image.jpg"
                                        height="80"
                                        width="80"
                                        cover
                                        style="border-radius: 8px"
                                    ></v-img>
                                </div>
                                <div>
                                    <h3
                                        class="group-header"
                                        style="font-weight: 500; cursor: pointer"
                                        @click="navigateGroupView(groupInvite.group.groupId)"
                                    >
                                        {{ groupInvite.group.groupName }}
                                    </h3>
                                    <span
                                        >Inviter:
                                        <span
                                            class="group-header"
                                            style="cursor: pointer"
                                            @click="
                                                navigateProfileView(groupInvite.inviterUser.userId)
                                            "
                                            >{{
                                                `${groupInvite.inviterUser.firstName} ${groupInvite.inviterUser.lastName}`
                                            }}</span
                                        ></span
                                    >
                                </div>
                            </v-row>
                            <v-btn
                                class="mt-4"
                                color="blue"
                                block
                                @click="acceptGroupInvite(groupInvite)"
                                >Accept</v-btn
                            >
                            <v-btn class="mt-2" block @click="discardGroupInvite(groupInvite)"
                                >Decline</v-btn
                            >
                        </v-card>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import HeaderComponent from '../components/header/HeaderComponent.vue'
import PostSection from '../components/forum/PostSection.vue'
export default {
    components: { HeaderComponent, PostSection },
    data() {
        return {
            user: null,
            searchInput: '',
            view: 'feed',
            groups: [],
            resGroups: [],
            openDialog: false,
            posts: [],
            groupRequests: [],
            groupInvites: []
        }
    },
    created() {
        axios.get('/api/getCurrentUser').then((res) => {
            this.user = res.data
            if (!this.user) return this.$router.push({ path: '/login' })
            axios.get('/api/getUserGroups').then((res) => (this.groups = res.data))
            axios.get('/api/getAllGroupPosts').then((res) => {
                this.posts = res.data
            })
            axios.get('/api/getGroupRequests').then((res) => (this.groupRequests = res.data))
            axios.get('/api/getGroupInvites').then((res) => (this.groupInvites = res.data))
        })
    },
    methods: {
        navigateCreateGroup() {
            this.$router.push({ path: '/groups/create' })
        },
        navigateGroupView(groupId) {
            this.$router.push({ name: 'group', params: { groupId: groupId } })
        },
        navigateProfileView(userId) {
            this.$router.push({ name: 'user', params: { userId: userId } })
        },
        searchGroups() {
            this.view = 'searchGroups'
            axios
                .post('/api/getGroupsByName', {
                    groupName: this.searchInput
                })
                .then((res) => (this.resGroups = res.data))
        },
        acceptGroupRequest(groupRequest) {
            axios
                .post('/api/acceptGroupRequest', {
                    groupId: groupRequest.group.groupId,
                    requesterId: groupRequest.user.userId
                })
                .then(() => {
                    this.groupRequests = this.groupRequests.filter(
                        (v) =>
                            v.user.userId != groupRequest.user.userId ||
                            v.group.groupId != groupRequest.group.groupId
                    )
                })
        },
        discardGroupRequest(groupRequest) {
            axios
                .post('/api/discardGroupRequest', {
                    groupId: groupRequest.group.groupId,
                    requesterId: groupRequest.user.userId
                })
                .then(() => {
                    this.groupRequests = this.groupRequests.filter(
                        (v) =>
                            v.user.userId != groupRequest.user.userId ||
                            v.group.groupId != groupRequest.group.groupId
                    )
                })
        },
        acceptGroupInvite(groupInvite) {
            axios
                .post('/api/acceptGroupInvite', {
                    groupId: groupInvite.group.groupId,
                    invitedId: groupInvite.invitedUser.userId,
                    inviterId: groupInvite.inviterUser.userId
                })
                .then(() => {
                    this.groupInvites = this.groupInvites.filter(
                        (v) =>
                            v.invitedUser.userId != groupInvite.invitedUser.userId ||
                            v.group.groupId != groupInvite.group.groupId ||
                            v.inviterUser.userId != groupInvite.inviterUser.userId ||
                            (this.groups.push(v.group) && false)
                    )
                })
        },
        discardGroupInvite(groupInvite) {
            axios
                .post('/api/discardGroupInvite', {
                    groupId: groupInvite.group.groupId,
                    invitedId: groupInvite.invitedUser.userId,
                    inviterId: groupInvite.inviterUser.userId
                })
                .then(() => {
                    this.groupInvites = this.groupInvites.filter(
                        (v) =>
                            v.invitedUser.userId != groupInvite.invitedUser.userId ||
                            v.group.groupId != groupInvite.group.groupId ||
                            v.inviterUser.userId != groupInvite.inviterUser.userId
                    )
                })
        }
    }
}
</script>

<style scoped>
.group-header:hover {
    text-decoration: underline;
}
</style>
