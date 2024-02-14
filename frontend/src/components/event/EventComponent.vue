<template>
    <div v-if="event">
        <div
            class="d-flex justify-center mb-4"
            style="background-color: white; box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1)"
        >
            <div style="min-width: 1000px; max-width: 1200px">
                <v-img
                    height="500"
                    cover
                    src="/event-stock-image.jpg"
                    style="border-radius: 8px"
                ></v-img>
                <h3 class="ml-4" style="font-weight: bold; color: red">{{ event.startDate }}</h3>
                <h1 class="ml-4" style="font-weight: bold">
                    {{ event.title }}
                </h1>
                <v-divider class="mt-4"></v-divider>
                <div class="d-flex justify-end my-3">
                    <div v-if="event.userAnswer == 0">
                        <v-btn
                            @click="setGoing()"
                            variant="flat"
                            rounded
                            color="success"
                            class="mr-2"
                            ><v-icon class="mr-1">mdi-check-circle-outline</v-icon>Going</v-btn
                        >
                        <v-btn @click="setNotGoing()" variant="flat" rounded color="red"
                            ><v-icon class="mr-1">mdi-close-circle-outline</v-icon>Not going</v-btn
                        >
                    </div>
                    <div v-else>
                        <v-menu width="300">
                            <template v-slot:activator="{ props }">
                                <v-btn v-if="event.userAnswer == 2" v-bind="props" color="success"
                                    ><v-icon class="mr-1">mdi-check-circle-outline</v-icon
                                    >Going<v-icon>mdi-chevron-down</v-icon></v-btn
                                >
                                <v-btn v-else v-bind="props" color="red"
                                    ><v-icon class="mr-1">mdi-close-circle-outline</v-icon>Not
                                    going<v-icon>mdi-chevron-down</v-icon></v-btn
                                >
                            </template>
                            <v-list class="pa-2">
                                <v-list-item @click="setGoing()" rounded>
                                    <v-row class="ma-0">
                                        <span>Going</span>
                                        <v-spacer></v-spacer>
                                        <v-icon v-if="event.userAnswer == 2">mdi-check</v-icon>
                                    </v-row>
                                </v-list-item>
                                <v-list-item @click="setNotGoing()" rounded>
                                    <v-row class="ma-0">
                                        <span>Not going</span>
                                        <v-spacer></v-spacer>
                                        <v-icon v-if="event.userAnswer == 1">mdi-check</v-icon>
                                    </v-row>
                                </v-list-item>
                            </v-list>
                        </v-menu>
                    </div>
                </div>
            </div>
        </div>
        <div class="d-flex justify-center">
            <div class="d-flex" style="min-width: 1000px; max-width: 1200px">
                <v-col cols="7" class="pa-0 mr-4">
                    <v-card width="auto" class="pa-4">
                        <h2 class="mb-2" style="font-weight: bold; line-height: 1">Details</h2>
                        <span class="d-block"
                            ><v-icon class="mr-4">mdi-account-multiple</v-icon
                            >{{ event.goingMembers.length }} people going</span
                        >
                        <span class="d-block"
                            ><v-icon class="mr-4">mdi-account</v-icon>Event that was created by
                            {{ `${event.author.firstName} ${event.author.lastName}` }}</span
                        >
                        <span class="d-block"
                            ><v-icon class="mr-4">mdi-account-group</v-icon>Group •
                            {{ event.group.groupName }}</span
                        >
                        <span class="d-block mt-4">{{ event.description }}</span>
                    </v-card>
                </v-col>
                <v-col cols="5" class="pa-0">
                    <v-card width="auto" class="pa-4">
                        <v-row class="ma-0">
                            <h2 style="font-weight: bold; line-height: 1">Guest list</h2>
                            <v-spacer></v-spacer>
                            <h4
                                class="see-all"
                                style="color: rgb(0, 110, 255); cursor: pointer"
                                @click="openDialogWithTab('going')"
                            >
                                See all
                            </h4>
                        </v-row>
                        <v-row class="ma-0 pt-2">
                            <v-col cols="6" class="pa-0">
                                <div
                                    class="d-flex align-center flex-column py-4 px-2 list-option"
                                    style="border-radius: 5px; cursor: pointer"
                                    @click="openDialogWithTab('going')"
                                >
                                    <h3 style="font-weight: bold">
                                        {{ event.goingMembers.length }}
                                    </h3>
                                    <span>GOING</span>
                                </div>
                            </v-col>
                            <v-col cols="6" class="pa-0">
                                <div
                                    class="d-flex align-center flex-column py-4 px-2 list-option"
                                    style="border-radius: 5px; cursor: pointer"
                                    @click="openDialogWithTab('invited')"
                                >
                                    <h3 style="font-weight: bold">
                                        {{ event.invitedMembers.length }}
                                    </h3>
                                    <span>INVITED</span>
                                </div>
                            </v-col>
                        </v-row>
                    </v-card>
                </v-col>
                <v-dialog v-model="openDialog" width="548">
                    <v-card>
                        <v-row class="ma-0 mt-4" style="height: 60px">
                            <h2 class="d-flex flex-grow-1 justify-center" style="font-weight: bold">
                                Guests
                            </h2>
                            <v-btn icon variant="flat" @click="openDialog = false"
                                ><v-icon>mdi-close</v-icon></v-btn
                            >
                        </v-row>
                        <v-divider class="mx-n4"></v-divider>
                        <v-tabs class="mt-1" v-model="tab">
                            <v-tab value="going">Going ({{ event.goingMembers.length }})</v-tab>
                            <v-tab value="notGoing"
                                >Not going ({{ event.notGoingMembers.length }})</v-tab
                            >
                            <v-tab value="invited"
                                >Invited ({{ event.invitedMembers.length }})</v-tab
                            >
                        </v-tabs>
                        <v-window v-model="tab" class="pa-4">
                            <v-window-item value="going">
                                <v-row
                                    class="ma-0 align-center"
                                    v-for="member in event.goingMembers"
                                    :key="member.memberId"
                                >
                                    <v-icon v-if="member.member.avatar" class="mr-4" size="36">
                                        <v-img
                                            :src="member.member.avatar"
                                            :width="36"
                                            :height="36"
                                            style="border-radius: 50%"
                                            cover
                                        />
                                    </v-icon>
                                    <v-icon v-else class="mr-4" size="36"
                                        >mdi-account-circle</v-icon
                                    >
                                    <span>{{
                                        `${member.member.firstName} ${member.member.lastName}`
                                    }}</span>
                                </v-row>
                            </v-window-item>
                            <v-window-item value="notGoing">
                                <v-row
                                    class="ma-0 align-center"
                                    v-for="member in event.notGoingMembers"
                                    :key="member.memberId"
                                >
                                    <v-icon v-if="member.member.avatar" class="mr-4" size="36">
                                        <v-img
                                            :src="member.member.avatar"
                                            :width="36"
                                            :height="36"
                                            style="border-radius: 50%"
                                            cover
                                        />
                                    </v-icon>
                                    <v-icon v-else class="mr-4" size="36"
                                        >mdi-account-circle</v-icon
                                    >
                                    <span>{{
                                        `${member.member.firstName} ${member.member.lastName}`
                                    }}</span>
                                </v-row>
                            </v-window-item>
                            <v-window-item value="invited">
                                <v-row
                                    class="ma-0 align-center"
                                    v-for="member in event.invitedMembers"
                                    :key="member.memberId"
                                >
                                    <v-icon v-if="member.member.avatar" class="mr-4" size="36">
                                        <v-img
                                            :src="member.member.avatar"
                                            :width="36"
                                            :height="36"
                                            style="border-radius: 50%"
                                            cover
                                        />
                                    </v-icon>
                                    <v-icon v-else class="mr-4" size="36"
                                        >mdi-account-circle</v-icon
                                    >
                                    <span>{{
                                        `${member.member.firstName} ${member.member.lastName}`
                                    }}</span>
                                </v-row>
                            </v-window-item>
                        </v-window>
                    </v-card>
                </v-dialog>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
export default {
    props: {
        propEvent: {
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
            event: null,
            openDialog: false,
            tab: 'going'
        }
    },
    created() {
        this.event = this.propEvent
    },
    methods: {
        setGoing() {
            axios
                .post('/api/setGoing', {
                    eventId: this.event.eventId
                })
                .then(() => {
                    this.event.userAnswer = 2
                    this.event.goingMembers.push({ member: this.user })
                    this.event.notGoingMembers = this.event.notGoingMembers.filter(
                        (v) => v.member.userId != this.user.userId
                    )
                })
        },
        setNotGoing() {
            axios
                .post('/api/setNotGoing', {
                    eventId: this.event.eventId
                })
                .then(() => {
                    this.event.userAnswer = 1
                    this.event.goingMembers = this.event.goingMembers.filter(
                        (v) => v.member.userId != this.user.userId
                    )
                    this.event.notGoingMembers.push({ member: this.user })
                })
        },
        openDialogWithTab(tabName) {
            this.openDialog = true
            this.tab = tabName
        }
    }
}
</script>

<style scoped>
.see-all:hover {
    text-decoration: underline;
}
.list-option:hover {
    background-color: rgb(236, 236, 236);
}
</style>
