<template>
    <v-menu width="300" v-model="isOpen">
        <template v-slot:activator="{ props }">
            <v-btn
                icon
                height="40"
                width="40"
                class="mr-2"
                variant="outlined"
                v-bind="props"
                style="position: relative"
            >
                <v-icon size="20" color="black">mdi-chat</v-icon>
                <div
                    v-if="
                        Array.from(this.conversations.values()).filter((v) => !v.lastMessage.isRead)
                            .length
                    "
                    style="
                        padding: 2px;
                        position: absolute;
                        top: -6px;
                        right: -2px;
                        background-color: red;
                        color: white;
                        border-radius: 50%;
                        font-size: 14px;
                        line-height: 14px;
                    "
                >
                    {{
                        Array.from(this.conversations.values()).filter((v) => !v.lastMessage.isRead)
                            .length
                    }}
                </div>
            </v-btn>
        </template>
        <v-list
            @click.stop.prevent
            class="pa-2"
            style="
                user-select: none;
                -webkit-user-select: none;
                -moz-user-select: none;
                -khtml-user-select: none;
                -ms-user-select: none;
            "
        >
            <v-row class="ma-0 pa-2">
                <h2 style="font-weight: bold">Chats</h2>
            </v-row>
            <div
                v-if="!conversations.size"
                class="d-flex justify-center align-center flex-column pa-4"
            >
                <div class="mb-5">
                    <v-icon size="112" style="color: rgb(101, 103, 107)">mdi-chat-outline</v-icon>
                </div>
                <span
                    style="
                        color: rgb(101, 103, 107);
                        font-weight: 700;
                        line-height: 24px;
                        font-size: 20px;
                    "
                    >You have no chats</span
                >
            </div>
            <v-list-item
                v-for="(key, index) in conversations.keys()"
                :key="index"
                class="px-2"
                rounded
                @click="openChat(key, conversations.get(key))"
            >
                <v-row class="ma-0">
                    <v-icon v-if="conversations.get(key).receiver.groupId" class="mr-2" size="56">
                        <v-img
                            src="/group-stock-image.jpg"
                            :width="56"
                            :height="56"
                            style="border-radius: 50%"
                            cover
                        />
                    </v-icon>
                    <v-icon
                        v-else-if="conversations.get(key).receiver.avatar"
                        class="mr-2"
                        size="56"
                    >
                        <v-img
                            :src="conversations.get(key).receiver.avatar"
                            :width="56"
                            :height="56"
                            style="border-radius: 50%"
                            cover
                        />
                    </v-icon>
                    <v-icon v-else size="56" class="mr-2">mdi-account-circle</v-icon>
                    <div>
                        <span style="font-size: 15px; font-weight: 700">{{
                            conversations.get(key).receiver.groupName ??
                            conversations.get(key).receiver.firstName +
                                ' ' +
                                conversations.get(key).receiver.lastName
                        }}</span>
                        <br />
                        <span
                            style="font-size: 14px; font-weight: 400; color: rgb(101, 103, 107)"
                            >{{
                                `${
                                    user.userId == conversations.get(key).lastMessage.senderId
                                        ? 'You: '
                                        : ''
                                }${conversations.get(key).lastMessage.message}`
                            }}</span
                        >
                    </div>
                    <v-spacer />
                    <div style="width: 48px">
                        <v-icon
                            v-if="!conversations.get(key).lastMessage.isRead"
                            size="48"
                            color="blue"
                            >mdi-circle-small</v-icon
                        >
                    </div>
                </v-row>
            </v-list-item>
        </v-list>
    </v-menu>
</template>

<script>
import chat from '../../chat.js'

export default {
    props: {
        user: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            conversations: chat.getAllConversations(),
            isOpen: false
        }
    },
    methods: {
        openChat(key, obj) {
            obj.lastMessage.isRead = true
            this.isOpen = false
            chat.removeMinimizedChat(key)
            chat.openChat(obj.receiver)
        }
    }
}
</script>
