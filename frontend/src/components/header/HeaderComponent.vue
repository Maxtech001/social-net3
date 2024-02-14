<template>
    <div style="height: 56px">
        <v-row class="header ma-0 align-center">
            <div
                v-if="isHomeView"
                class="d-flex justify-center align-center"
                style="position: fixed; width: 100%; height: 56px"
            >
                <v-tabs v-model="tab">
                    <v-tab
                        class="ml-4"
                        icon
                        height="40"
                        width="40"
                        value="homeFeed"
                        @click="$emit('homeFeed')"
                    >
                        <v-icon size="24" color="black">mdi-post</v-icon>
                    </v-tab>
                    <v-tab
                        class="ml-4"
                        icon
                        height="40"
                        width="40"
                        value="groupFeed"
                        @click="$emit('groupFeed')"
                    >
                        <v-icon size="24" color="black"
                            >mdi-account-supervisor-circle-outline</v-icon
                        >
                    </v-tab>
                </v-tabs>
            </div>
            <v-btn
                class="ml-4 mr-2"
                icon
                height="40"
                width="40"
                variant="outlined"
                @click="redirectHome()"
            >
                <v-icon size="20" color="black">mdi-home</v-icon>
            </v-btn>
            <div
                class="pl-3"
                style="
                    width: 280px;
                    z-index: 1;
                    background-color: rgb(240, 242, 245);
                    height: 40px;
                    border-radius: 50px;
                "
                @click="openDialog = true"
                v-click-outside="() => (openDialog = false)"
            >
                <v-icon class="mr-1">mdi-magnify</v-icon>
                <input
                    v-model="searchInput"
                    placeholder="Search"
                    style="height: inherit; width: 190px; outline: none"
                />
                <div
                    v-if="searchInput.length && openDialog"
                    tabindex="0"
                    class="px-1"
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
                        <v-list-item rounded @click="search()">
                            <v-icon class="mr-1">mdi-magnify</v-icon>
                            Search for {{ searchInput }}
                        </v-list-item>
                    </v-list>
                </div>
            </div>
            <v-spacer></v-spacer>
            <MenuComponent></MenuComponent>
            <MessageMenu :user="user"></MessageMenu>
            <NotificationMenu></NotificationMenu>
            <ProfileMenu :user="user"></ProfileMenu>
        </v-row>
        <div
            style="
                position: fixed;
                top: 56px;
                height: 2px;
                width: 100%;
                background-color: #ccc;
                filter: blur(3px);
                z-index: 1;
            "
        ></div>
    </div>
</template>

<script>
import ProfileMenu from './ProfileMenu.vue'
import MenuComponent from './MenuComponent.vue'
import MessageMenu from './MessageMenu.vue'
import NotificationMenu from './NotificationMenu.vue'

export default {
    components: {
        ProfileMenu,
        MenuComponent,
        MessageMenu,
        NotificationMenu
    },
    props: {
        user: Object,
        isHomeView: {
            type: Boolean,
            default: false
        },
        searchQuery: {
            type: String,
            default: ''
        }
    },
    data() {
        return {
            tab: 'homeFeed',
            searchInput: '',
            openDialog: false
        }
    },
    created() {
        this.searchInput = this.searchQuery
    },
    methods: {
        redirectHome() {
            this.$router.push({ path: '/' })
        },
        search() {
            this.$router.push({
                name: 'search',
                params: { searchInput: this.searchInput }
            })
        }
    }
}
</script>

<style scoped>
.header {
    width: 100%;
    height: 56px;
    background-color: white;
    border-bottom: 1px solid #ccc;
    position: fixed;
    top: 0;
    z-index: 2;
}
</style>
