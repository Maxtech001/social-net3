<template>
    <div v-if="user" class="d-flex flex-column" style="height: 100%">
        <HeaderComponent :user="user"></HeaderComponent>
        <div class="d-flex flex-grow-1" style="height: 100%">
            <v-sheet class="pa-4" min-width="360">
                <v-form v-model="isValid" class="d-flex flex-column" style="height: 100%">
                    <div>
                        <h2>Create group</h2>
                        <v-row class="ma-0 my-4 align-center">
                            <v-icon v-if="user.avatar" class="mr-4" size="36">
                                <v-img
                                    :src="user.avatar"
                                    :width="36"
                                    :height="36"
                                    style="border-radius: 50%"
                                    cover
                                />
                            </v-icon>
                            <v-icon v-else size="36" class="mr-4">mdi-account-circle</v-icon>
                            <div>
                                <span class="d-block">{{
                                    `${user.firstName} ${user.lastName}`
                                }}</span>
                                <span>Administrator</span>
                            </div>
                        </v-row>
                        <v-text-field
                            v-model="groupName"
                            rounded
                            variant="outlined"
                            label="Group name"
                            :rules="[(v) => !!v || 'Enter a group name']"
                        ></v-text-field>
                        <v-text-field
                            v-model="description"
                            rounded
                            variant="outlined"
                            label="Group description"
                            :rules="[(v) => !!v || 'Enter a group description']"
                        ></v-text-field>
                        <v-autocomplete
                            v-model="selectedFollowers"
                            :items="
                                followers.map((v) => ({
                                    title: v.firstName + ' ' + v.lastName,
                                    value: v
                                }))
                            "
                            rounded
                            variant="outlined"
                            multiple
                            clearable
                            label="Invite followers (optional)"
                        ></v-autocomplete>
                    </div>
                    <v-spacer></v-spacer>
                    <v-divider></v-divider>
                    <v-btn class="mt-4" color="blue" :disabled="!isValid" @click="createGroup()"
                        >Create</v-btn
                    >
                </v-form>
            </v-sheet>
            <div
                class="d-flex align-center justify-center flex-grow-1 pa-4"
                style="min-width: 600px"
            >
                <v-card class="pa-4" max-width="972">
                    <h4 style="font-weight: 500">Preview</h4>
                    <div style="border: 0.25px rgb(214, 214, 214) solid; border-radius: 8px">
                        <v-img
                            width="938"
                            height="300"
                            cover
                            src="/group-stock-image.jpg"
                            style="border-radius: 8px"
                        ></v-img>
                        <h1 class="ml-4" style="font-weight: bold">
                            {{ groupName.length ? groupName : 'Group name' }}
                        </h1>
                        <span class="ml-4"
                            >{{ description.length ? description : 'Description' }} • 1 member</span
                        >
                        <v-divider class="mt-4 mx-4"></v-divider>
                        <v-tabs class="ml-4" v-model="tab" disabled hide-slider>
                            <v-tab value="posts">Posts</v-tab>
                            <v-tab value="members">Members</v-tab>
                            <v-tab value="events">Events</v-tab>
                        </v-tabs>
                        <v-divider></v-divider>
                        <v-window disabled v-model="tab">
                            <v-window-item value="posts">
                                <div class="pa-4 d-flex">
                                    <v-overlay
                                        contained
                                        persistent
                                        :model-value="true"
                                        style="
                                            opacity: 0.4;
                                            border-bottom-left-radius: 8px;
                                            border-bottom-right-radius: 8px;
                                        "
                                    ></v-overlay>
                                    <PostSection :user="user"></PostSection>
                                </div>
                            </v-window-item>
                        </v-window>
                    </div>
                </v-card>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import HeaderComponent from '../components/header/HeaderComponent.vue'
import PostSection from '../components/forum/PostSection.vue'
export default {
    components: {
        HeaderComponent,
        PostSection
    },
    created() {
        axios.get('/api/getCurrentUser').then((res) => {
            this.user = res.data
            if (!this.user) return this.$router.push({ path: '/login' })
            axios.get('/api/getFollowers').then((res) => (this.followers = res.data))
        })
    },
    data() {
        return {
            user: null,
            groupName: '',
            description: '',
            isValid: false,
            tab: 'posts',
            followers: [],
            selectedFollowers: null
        }
    },
    methods: {
        createGroup() {
            axios
                .post('/api/createGroup', {
                    groupName: this.groupName,
                    description: this.description,
                    invitedUsers: this.selectedFollowers
                })
                .then((res) =>
                    this.$router.push({ name: 'group', params: { groupId: res.data.groupId } })
                )
        }
    }
}
</script>
