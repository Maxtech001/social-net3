<template>
    <div v-if="group">
        <div
            class="d-flex justify-center"
            style="background-color: white; box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1)"
        >
            <div style="min-width: 1000px; max-width: 1200px">
                <v-img
                    height="500"
                    cover
                    src="/group-stock-image.jpg"
                    style="border-radius: 8px"
                ></v-img>
                <v-row class="ma-0">
                    <div>
                        <h1 class="ml-4" style="font-weight: bold">
                            {{ group.groupName.length ? group.groupName : 'Group name' }}
                        </h1>
                        <span class="ml-4"
                            >{{ group.description.length ? group.description : 'Description' }} •
                            {{ group.memberCount }} members</span
                        >
                    </div>
                    <v-spacer></v-spacer>
                    <div v-if="group.adminId != user.userId" class="d-flex align-end">
                        <v-btn
                            v-if="group.isRequested"
                            variant="flat"
                            rounded
                            color="red"
                            @click="cancelRequest()"
                            ><v-icon class="mr-2">mdi-close</v-icon>Cancel request</v-btn
                        >
                        <v-btn
                            v-else-if="group.isMember"
                            variant="flat"
                            rounded
                            color="red"
                            @click="leaveGroup()"
                            >Leave group</v-btn
                        >
                        <v-btn v-else color="blue" @click="joinGroup()"
                            ><v-icon class="mr-2">mdi-account-multiple-plus</v-icon>Join
                            group</v-btn
                        >
                    </div>
                    <div class="d-flex align-end ml-2">
                        <v-btn
                            class="mr-2"
                            v-if="group.isMember"
                            variant="flat"
                            color="black"
                            rounded
                            @click="openChat()"
                            ><v-icon class="mr-2">mdi-chat-outline</v-icon>Message</v-btn
                        >
                        <v-btn
                            v-if="group.isMember"
                            variant="flat"
                            color="blue"
                            rounded
                            @click="OpenInviteDialog()"
                            ><v-icon class="mr-2">mdi-plus</v-icon>Invite</v-btn
                        >
                    </div>
                    <v-dialog v-model="openInviteDialog" width="500">
                        <v-card>
                            <v-row class="ma-0 mt-1 px-4 align-center" style="height: 60px">
                                <h2
                                    class="d-flex flex-grow-1 justify-center"
                                    style="font-weight: bold"
                                >
                                    Invite people to this group
                                </h2>
                                <v-btn icon variant="flat" @click="openInviteDialog = false"
                                    ><v-icon>mdi-close</v-icon></v-btn
                                >
                            </v-row>
                            <v-divider></v-divider>
                            <div class="pa-4">
                                <v-autocomplete
                                    v-model="selectedFollowers"
                                    :items="
                                        mutualFollowers.map((v) => ({
                                            title: v.firstName + ' ' + v.lastName,
                                            value: v
                                        }))
                                    "
                                    :loading="loading"
                                    rounded
                                    variant="solo"
                                    density="compact"
                                    prepend-inner-icon="mdi-magnify"
                                    multiple
                                    clearable
                                    placeholder="Search followers"
                                ></v-autocomplete>
                                <div class="d-flex justify-end">
                                    <v-btn
                                        :disabled="!selectedFollowers.length"
                                        color="blue"
                                        @click="sendInvites()"
                                        >Send invites</v-btn
                                    >
                                </div>
                            </div>
                        </v-card>
                    </v-dialog>
                </v-row>
                <v-divider class="mt-4 mx-4"></v-divider>
                <v-tabs class="ml-4" v-model="tab">
                    <v-tab value="posts">Posts</v-tab>
                    <v-tab value="members">Members</v-tab>
                    <v-tab value="events">Events</v-tab>
                </v-tabs>
            </div>
        </div>
        <div class="d-flex justify-center">
            <v-window v-model="tab" style="min-width: 1000px; max-width: 1200px">
                <v-window-item value="posts">
                    <div class="pa-4 d-flex justify-center">
                        <div style="width: 680px">
                            <PostSection
                                :user="user"
                                :group="group"
                                :propPosts="posts"
                            ></PostSection>
                        </div>
                    </div>
                </v-window-item>
                <v-window-item value="members">
                    <div class="pa-4 d-flex justify-center">
                        <v-card width="680" class="pa-4">
                            <span style="font-weight: 500; font-size: 18px">Members</span>
                            <span> • {{ group.memberCount }}</span>
                            <v-text-field
                                v-model="searchInput"
                                class="mb-4 mt-2"
                                prepend-inner-icon="mdi-magnify"
                                variant="solo"
                                rounded
                                placeholder="Find member"
                                density="compact"
                                hide-details
                            ></v-text-field>
                            <v-divider class="mb-2"></v-divider>
                            <div v-if="!searchInput.length">
                                <h6 class="mb-4" style="font-weight: 500; font-size: 16px">
                                    Admin
                                </h6>
                                <v-row class="ma-0 align-center">
                                    <v-icon v-if="group.admin.avatar" class="mr-4" size="60">
                                        <v-img
                                            :src="group.admin.avatar"
                                            :width="60"
                                            :height="60"
                                            style="border-radius: 50%"
                                            cover
                                        />
                                    </v-icon>
                                    <v-icon v-else class="mr-4" size="60"
                                        >mdi-account-circle</v-icon
                                    >
                                    <h3 style="font-weight: 500; line-height: 1">
                                        {{ `${group.admin.firstName} ${group.admin.lastName}` }}
                                    </h3>
                                </v-row>
                                <v-divider class="mt-6 mb-2"></v-divider>
                                <h6 class="mb-4" style="font-weight: 500; font-size: 16px">
                                    All members
                                </h6>
                                <v-row
                                    class="ma-0 align-center"
                                    v-for="member in group.members"
                                    :key="member.userId"
                                >
                                    <v-icon v-if="member.avatar" class="mr-4" size="60">
                                        <v-img
                                            :src="member.avatar"
                                            :width="60"
                                            :height="60"
                                            style="border-radius: 50%"
                                            cover
                                        />
                                    </v-icon>
                                    <v-icon v-else class="mr-4" size="60"
                                        >mdi-account-circle</v-icon
                                    >
                                    <h3 style="font-weight: 500; line-height: 1">
                                        {{ `${member.firstName} ${member.lastName}` }}
                                    </h3>
                                </v-row>
                            </div>
                            <div v-else>
                                <h6 class="mb-4" style="font-weight: 500; font-size: 16px">
                                    Search results
                                </h6>
                                <v-row
                                    class="ma-0 align-center"
                                    v-for="member in group.members.filter((v) =>
                                        `${v.firstName} ${v.lastName}`
                                            .toLowerCase()
                                            .includes(searchInput.toLowerCase())
                                    )"
                                    :key="member.userId"
                                >
                                    <v-icon v-if="member.avatar" class="mr-4" size="60">
                                        <v-img
                                            :src="member.avatar"
                                            :width="60"
                                            :height="60"
                                            style="border-radius: 50%"
                                            cover
                                        />
                                    </v-icon>
                                    <v-icon v-else class="mr-4" size="60"
                                        >mdi-account-circle</v-icon
                                    >
                                    <h3 style="font-weight: 500; line-height: 1">
                                        {{ `${member.firstName} ${member.lastName}` }}
                                    </h3>
                                </v-row>
                            </div>
                        </v-card>
                    </div>
                </v-window-item>
                <v-window-item value="events">
                    <div class="pa-4 d-flex align-center flex-column">
                        <v-card width="680" class="pa-4 mb-4">
                            <v-row class="ma-0 mt-n2">
                                <h2 style="font-weight: bold; font-size: 20px">Upcoming events</h2>
                                <v-spacer></v-spacer>
                                <v-btn @click="openDialog = true">Create event</v-btn>
                                <v-dialog v-model="openDialog" width="548">
                                    <v-card class="pa-4">
                                        <v-row class="ma-0" style="height: 60px">
                                            <h2
                                                class="d-flex flex-grow-1 justify-center"
                                                style="font-weight: bold"
                                            >
                                                Create event
                                            </h2>
                                            <v-btn icon variant="flat" @click="openDialog = false"
                                                ><v-icon>mdi-close</v-icon></v-btn
                                            >
                                        </v-row>
                                        <v-row class="ma-0">
                                            <v-icon v-if="user.avatar" class="mr-4" size="36">
                                                <v-img
                                                    :src="user.avatar"
                                                    :width="36"
                                                    :height="36"
                                                    style="border-radius: 50%"
                                                    cover
                                                />
                                            </v-icon>
                                            <v-icon v-else class="mr-4" size="36"
                                                >mdi-account-circle</v-icon
                                            >
                                            <div>
                                                <h3 style="font-weight: bold; line-height: 1">
                                                    {{ `${user.firstName} ${user.lastName}` }}
                                                </h3>
                                                <span>Host</span>
                                            </div>
                                        </v-row>
                                        <v-form v-model="isValid">
                                            <v-text-field
                                                class="pb-2"
                                                v-model="title"
                                                variant="outlined"
                                                label="Title"
                                                hide-details
                                                :rules="[(v) => !!v || 'Enter a title']"
                                            ></v-text-field>
                                            <v-row class="ma-0 pb-2">
                                                <v-text-field
                                                    v-model="startDate"
                                                    variant="outlined"
                                                    label="Start date"
                                                    type="date"
                                                    :rules="[(v) => !!v]"
                                                    @input="
                                                        (v) => {
                                                            if (startDate > endDate)
                                                                endDate = startDate
                                                        }
                                                    "
                                                    hide-details
                                                ></v-text-field>
                                                <v-text-field
                                                    v-model="endDate"
                                                    variant="outlined"
                                                    label="End date"
                                                    type="date"
                                                    :rules="[(v) => !!v]"
                                                    :min="startDate"
                                                    hide-details
                                                ></v-text-field>
                                            </v-row>
                                            <v-text-field
                                                v-model="description"
                                                class="pb-4"
                                                variant="outlined"
                                                label="Description"
                                                hide-details
                                                :rules="[(v) => !!v || 'Enter a description']"
                                            ></v-text-field>
                                            <v-btn
                                                block
                                                color="blue"
                                                :disabled="!isValid"
                                                @click="createEvent()"
                                                >Create event</v-btn
                                            >
                                        </v-form>
                                    </v-card>
                                </v-dialog>
                            </v-row>
                            <div
                                v-if="!upcomingEvents.length"
                                class="d-flex justify-center align-center flex-column pa-4"
                            >
                                <div class="mb-5">
                                    <v-icon size="112" style="color: rgb(101, 103, 107)"
                                        >mdi-calendar-month-outline</v-icon
                                    >
                                </div>
                                <span style="font-weight: 400; line-height: 20px; font-size: 15px"
                                    >There is no upcoming events.</span
                                >
                            </div>
                            <v-row
                                class="ma-0"
                                v-for="(event, index) in upcomingEvents"
                                :key="event.eventId"
                            >
                                <v-divider class="mb-3 mt-3" v-if="index > 0"></v-divider>
                                <div class="mr-4" style="max-width: 144px">
                                    <v-img
                                        src="/event-stock-image.jpg"
                                        height="144"
                                        width="144"
                                        cover
                                    ></v-img>
                                </div>
                                <div>
                                    <h4 style="color: red; font-weight: 500">
                                        {{ event.startDate }}
                                    </h4>
                                    <a
                                        :href="'/event/' + event.eventId"
                                        class="event-header"
                                        style="
                                            font-weight: bold;
                                            text-decoration: unset;
                                            display: block;
                                            color: black;
                                            font-size: 20px;
                                        "
                                        >{{ event.title }}</a
                                    >
                                    <v-row class="ma-0">
                                        <v-icon v-if="event.author.avatar" class="mr-1" size="24">
                                            <v-img
                                                :src="event.author.avatar"
                                                :width="24"
                                                :height="24"
                                                style="border-radius: 50%"
                                                cover
                                            />
                                        </v-icon>
                                        <v-icon v-else class="mr-1" size="24"
                                            >mdi-account-circle</v-icon
                                        >
                                        Author:
                                        {{
                                            `${event.author.firstName} ${event.author.lastName}`
                                        }}</v-row
                                    >
                                </div>
                            </v-row>
                        </v-card>
                        <v-card width="680" class="pa-4">
                            <v-row class="ma-0 mt-n3">
                                <h2 style="font-weight: bold; font-size: 20px">Past events</h2>
                            </v-row>
                            <div
                                v-if="!pastEvents.length"
                                class="d-flex justify-center align-center flex-column pa-4"
                            >
                                <div class="mb-5">
                                    <v-icon size="112" style="color: rgb(101, 103, 107)"
                                        >mdi-calendar-month-outline</v-icon
                                    >
                                </div>
                                <span style="font-weight: 400; line-height: 20px; font-size: 15px"
                                    >There is no past events.</span
                                >
                            </div>
                            <v-row
                                class="ma-0"
                                v-for="(event, index) in pastEvents"
                                :key="event.eventId"
                            >
                                <v-divider class="mb-3 mt-3" v-if="index > 0"></v-divider>
                                <div class="mr-4" style="max-width: 144px">
                                    <v-img
                                        src="/event-stock-image.jpg"
                                        height="144"
                                        width="144"
                                        cover
                                    ></v-img>
                                </div>
                                <div>
                                    <h4 style="color: red; font-weight: 500">
                                        {{ event.startDate }}
                                    </h4>
                                    <a
                                        :href="'/event/' + event.eventId"
                                        class="event-header"
                                        style="
                                            font-weight: bold;
                                            text-decoration: unset;
                                            display: block;
                                            color: black;
                                            font-size: 20px;
                                        "
                                        >{{ event.title }}</a
                                    >
                                    <v-row class="ma-0">
                                        <v-icon v-if="event.author.avatar" class="mr-1" size="24">
                                            <v-img
                                                :src="event.author.avatar"
                                                :width="24"
                                                :height="24"
                                                style="border-radius: 50%"
                                                cover
                                            />
                                        </v-icon>
                                        <v-icon v-else class="mr-1" size="24"
                                            >mdi-account-circle</v-icon
                                        >
                                        Author:
                                        {{
                                            `${event.author.firstName} ${event.author.lastName}`
                                        }}</v-row
                                    >
                                </div>
                            </v-row>
                        </v-card>
                    </div>
                </v-window-item>
            </v-window>
        </div>
    </div>
</template>

<script>
import PostSection from '../forum/PostSection.vue'
import axios from 'axios'
import chat from '../../chat.js'
export default {
    components: { PostSection },
    props: {
        propGroup: {
            type: Object,
            required: true
        },
        user: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            tab: 'posts',
            openDialog: false,
            searchInput: '',
            title: '',
            startDate: null,
            endDate: null,
            description: '',
            isValid: false,
            upcomingEvents: [],
            pastEvents: [],
            posts: [],
            group: null,
            openInviteDialog: false,
            selectedFollowers: [],
            mutualFollowers: [],
            loading: false
        }
    },
    created() {
        this.group = this.propGroup
        if (!this.group) return
        axios
            .post('/api/getUpcomingEvents', {
                groupId: this.group.groupId
            })
            .then((res) => (this.upcomingEvents = res.data))
        axios
            .post('/api/getPastEvents', {
                groupId: this.group.groupId
            })
            .then((res) => (this.pastEvents = res.data))
        axios
            .post('/api/getGroupPosts', {
                groupId: this.group.groupId
            })
            .then((res) => (this.posts = res.data))
    },
    methods: {
        createEvent() {
            axios
                .post('/api/createEvent', {
                    title: this.title,
                    description: this.description,
                    startDate: this.startDate,
                    endDate: this.endDate,
                    groupId: this.group.groupId
                })
                .then((res) => {
                    if (new Date().toISOString().split('T')[0] <= res.data.startDate) {
                        this.upcomingEvents.push(res.data)
                        this.upcomingEvents.sort((a, b) =>
                            a.startDate < b.startDate ? 1 : b.startDate < a.startDate ? -1 : 0
                        )
                    } else {
                        this.pastEvents.push(res.data)
                        this.pastEvents.sort((a, b) =>
                            a.startDate < b.startDate ? 1 : b.startDate < a.startDate ? -1 : 0
                        )
                    }
                    this.openDialog = false
                })
        },
        joinGroup() {
            axios
                .post('/api/joinGroup', {
                    groupId: this.group.groupId
                })
                .then(() => (this.group.isRequested = true))
        },
        leaveGroup() {
            axios
                .post('/api/leaveGroup', {
                    groupId: this.group.groupId
                })
                .then(() => (this.group.isMember = false))
        },
        cancelRequest() {
            axios
                .post('/api/cancelGroupRequest', {
                    groupId: this.group.groupId
                })
                .then(() => (this.group.isRequested = false))
        },
        OpenInviteDialog() {
            this.loading = true
            axios
                .post('/api/getFollowersNotInGroup', {
                    groupId: this.group.groupId
                })
                .then((res) => {
                    this.mutualFollowers = res.data
                })
                .finally(() => (this.loading = false))
            this.openInviteDialog = true
        },
        sendInvites() {
            axios
                .post('/api/storeInvites', {
                    users: this.selectedFollowers,
                    groupId: this.group.groupId
                })
                .then(() => {
                    this.selectedFollowers = []
                    this.openInviteDialog = false
                })
        },
        openChat() {
            chat.openChat(this.group)
        }
    },
    watch: {
        openDialog() {
            if (this.openDialog) return
            this.title = null
            this.description = null
            this.startDate = null
            this.endDate = null
        }
    }
}
</script>

<style scoped>
.event-header:hover {
    text-decoration: underline !important;
}
</style>
