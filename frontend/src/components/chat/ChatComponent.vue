<template>
    <div
        v-if="user"
        class="position-fixed d-flex align-end"
        style="bottom: 0; right: 0"
        id="chat-component"
    >
        <div class="d-flex flex-row-reverse">
            <div v-if="newMessageClicked" class="chat-container">
                <v-row class="ma-0 pa-2">
                    <h3 style="font-weight: 600">New message</h3>
                    <v-spacer />
                    <v-btn icon density="compact" variant="flat" @click="newMessageClicked = false"
                        ><v-icon>mdi-close</v-icon></v-btn
                    >
                </v-row>
                <v-row class="ma-0 my-2 ml-4">
                    <span>To:</span>
                    <div class="d-flex px-4" style="flex: 1">
                        <input
                            id="contactSearch"
                            autocomplete="off"
                            @input="(v) => handleInput(v.target)"
                            style="outline: none; width: 100%"
                        />
                    </div>
                </v-row>
                <v-divider />
                <v-list class="px-1">
                    <v-list-item
                        v-for="(result, index) in visibleContacts"
                        :key="index"
                        rounded
                        @click="
                            openChat(
                                (result.Group && `g-${result.Group.groupId}`) ||
                                    (result.User && `u-${result.User.userId}`),
                                result.Group || result.User
                            )
                        "
                    >
                        <v-icon v-if="result.Group" class="mr-2" size="36">
                            <v-img
                                src="/group-stock-image.jpg"
                                :width="36"
                                :height="36"
                                style="border-radius: 50%"
                                cover
                            />
                        </v-icon>
                        <v-icon v-else-if="result.User.avatar" class="mr-2" size="36">
                            <v-img
                                :src="result.User.avatar"
                                :width="36"
                                :height="36"
                                style="border-radius: 50%"
                                cover
                            />
                        </v-icon>
                        <v-icon v-else size="36" class="mr-2">mdi-account-circle</v-icon>
                        <span>{{
                            (result.User && `${result.User.firstName} ${result.User.lastName}`) ||
                            (result.Group && result.Group.groupName)
                        }}</span>
                    </v-list-item>
                </v-list>
            </div>
            <div
                :key="index"
                v-for="(i, index) in openedChats.keys()"
                class="chat-container"
                @click="
                    readMessage(
                        openedChats.get(i).messages.length && openedChats.get(i).messages[0],
                        i
                    )
                "
            >
                <div>
                    <div
                        class="pa-2 d-flex align-center"
                        style="box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1)"
                    >
                        <v-icon v-if="openedChats.get(i).receiver.groupId" class="mr-2" size="32">
                            <v-img
                                src="/group-stock-image.jpg"
                                :width="32"
                                :height="32"
                                style="border-radius: 50%"
                                cover
                            />
                        </v-icon>
                        <v-icon
                            v-else-if="openedChats.get(i).receiver.avatar"
                            class="mr-2"
                            size="32"
                        >
                            <v-img
                                :src="openedChats.get(i).receiver.avatar"
                                :width="32"
                                :height="32"
                                style="border-radius: 50%"
                                cover
                            />
                        </v-icon>
                        <v-icon v-else size="32" class="mr-2">mdi-account-circle</v-icon>
                        <div
                            class="d-flex align-center"
                            style="font-weight: 500"
                            v-if="openedChats.get(i).receiver.userId"
                        >
                            {{
                                `${openedChats.get(i).receiver.firstName} ${
                                    openedChats.get(i).receiver.lastName
                                }`
                            }}
                        </div>
                        <div class="d-flex align-center" style="font-weight: 500" v-else>
                            {{ openedChats.get(i).receiver.groupName }}
                        </div>
                        <v-icon
                            v-if="
                                openedChats.get(i).messages.length &&
                                !openedChats.get(i).messages[0].isRead
                            "
                            color="red"
                            >mdi-circle-medium</v-icon
                        >
                        <v-spacer />
                        <v-btn
                            icon
                            density="compact"
                            variant="flat"
                            @click.stop
                            @click="minimizeChat(i, openedChats.get(i))"
                            ><v-icon>mdi-minus</v-icon></v-btn
                        >
                        <v-btn
                            icon
                            density="compact"
                            variant="flat"
                            @click.stop
                            @click="closeOpenedChat(i)"
                            ><v-icon>mdi-close</v-icon></v-btn
                        >
                    </div>
                    <v-divider />
                    <div>
                        <div
                            style="
                                display: flex;
                                height: 462px;
                                overflow: auto;
                                flex-direction: column-reverse;
                            "
                            class="px-2"
                        >
                            <div
                                v-for="(message, mindex) in openedChats.get(i).messages"
                                :key="mindex"
                                style="margin-top: 2px; max-width: 328px"
                                :style="user.userId == message.senderId ? { direction: 'rtl' } : {}"
                            >
                                <v-row class="ma-0 align-end">
                                    <div v-if="user.userId != message.senderId">
                                        <v-icon
                                            v-if="
                                                openedChats.get(i).receiver.avatar ||
                                                message.senderObj?.avatar
                                            "
                                            class="mr-2"
                                            size="28"
                                        >
                                            <v-img
                                                :src="
                                                    openedChats.get(i).receiver.avatar ??
                                                    message.senderObj.avatar
                                                "
                                                :width="28"
                                                :height="28"
                                                style="border-radius: 50%"
                                                cover
                                            />
                                        </v-icon>
                                        <v-icon v-else size="28" class="mr-2"
                                            >mdi-account-circle</v-icon
                                        >
                                    </div>
                                    <div
                                        class="py-2 px-3"
                                        style="
                                            direction: ltr;
                                            border-radius: 18px;
                                            width: fit-content;
                                            max-width: 230px;
                                            color: rgb(28, 30, 33);
                                            background-color: #f0f0f0;
                                            word-wrap: break-word;
                                        "
                                        :style="
                                            user.userId == message.senderId
                                                ? {
                                                      backgroundColor: 'rgb(44, 98, 239)',
                                                      color: 'white'
                                                  }
                                                : {}
                                        "
                                    >
                                        {{ message.message }}
                                    </div>
                                </v-row>
                            </div>
                        </div>
                        <div class="d-flex py-3" style="min-height: 60px">
                            <textarea
                                class="mx-2 pl-3 py-2"
                                :id="'textarea-' + index"
                                placeholder="Aa"
                                style="
                                    background-color: rgb(240, 242, 245);
                                    border-radius: 20px;
                                    resize: none;
                                    width: 100%;
                                    line-height: 20px;
                                    outline: none;
                                "
                            ></textarea>
                            <v-btn
                                icon
                                density="comfortable"
                                @click="sendMessage(openedChats.get(i), index)"
                                variant="flat"
                                ><v-icon>mdi-send</v-icon></v-btn
                            >
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="d-flex flex-column-reverse ml-3 mr-4 align-center">
            <v-btn icon class="mb-4" @click="newMessageClicked = true"
                ><v-icon size="28">mdi-pencil-box-outline</v-icon></v-btn
            >
            <div
                class="mb-2 d-flex"
                :key="index"
                v-for="(i, index) in minimizedChats.keys()"
                style="position: relative"
            >
                <v-hover>
                    <template v-slot:default="{ isHovering, props }">
                        <v-btn v-bind="props" icon @click="openChat(i, minimizedChats.get(i))">
                            <v-icon v-if="minimizedChats.get(i).receiver.groupId" size="48">
                                <v-img
                                    src="/group-stock-image.jpg"
                                    :width="48"
                                    :height="48"
                                    style="border-radius: 50%"
                                    cover
                                />
                            </v-icon>
                            <v-icon v-else-if="minimizedChats.get(i).receiver.avatar" size="48">
                                <v-img
                                    :src="minimizedChats.get(i).receiver.avatar"
                                    :width="48"
                                    :height="48"
                                    style="border-radius: 50%"
                                    cover
                                />
                            </v-icon>
                            <v-icon v-else size="48">mdi-account-circle</v-icon>
                        </v-btn>
                        <v-btn
                            v-bind="props"
                            icon
                            style="
                                position: absolute;
                                top: 0;
                                right: 0;
                                height: 18px;
                                width: 18px;
                                z-index: 1;
                            "
                            @click="closeMinimizedChat(i)"
                            v-if="isHovering"
                            ><v-icon size="18">mdi-close</v-icon></v-btn
                        >
                    </template>
                </v-hover>
            </div>
            <v-menu location="top" width="344">
                <template v-slot:activator="{ props }">
                    <v-btn
                        v-bind="props"
                        density="comfortable"
                        icon
                        class="mb-2"
                        v-if="minimizedChats.size || openedChats.size"
                        ><v-icon>mdi-dots-horizontal</v-icon></v-btn
                    >
                </template>
                <v-list class="pa-2" style="border-radius: 8px">
                    <v-list-item rounded @click="closeChats()">
                        <v-icon class="mr-2" size="24">mdi-close-circle-outline</v-icon>
                        Close all chats
                    </v-list-item>
                    <v-list-item rounded @click="minimizeChats()" v-if="openedChats.size">
                        <v-icon class="mr-2" size="24">mdi-minus-circle-outline</v-icon>
                        Minimize open chats
                    </v-list-item>
                </v-list>
            </v-menu>
        </div>
    </div>
</template>

<script>
import { sendMessageConn, user } from '../../conn.js'
import chat from '../../chat.js'
import axios from 'axios'
export default {
    data() {
        return {
            openedChats: chat.getOpenedChatArray(),
            minimizedChats: chat.getMinimizedChatArray(),
            user: user,
            newMessageClicked: false,
            contacts: [],
            visibleContacts: []
        }
    },
    methods: {
        sendMessage(chatObj, index) {
            let message = document.getElementById('textarea-' + index).value
            document.getElementById('textarea-' + index).value = ''
            sendMessageConn(message, chatObj.receiver)
        },
        closeOpenedChat(key) {
            chat.removeOpenedChat(key)
        },
        closeMinimizedChat(key) {
            chat.removeMinimizedChat(key)
        },
        minimizeChat(key, obj) {
            chat.removeOpenedChat(key)
            chat.minimizeChat(obj.receiver)
        },
        async openChat(key, obj) {
            this.newMessageClicked = false
            await this.$nextTick()
            chat.removeMinimizedChat(key)
            chat.openChat(obj.receiver || obj)
        },
        closeChats() {
            chat.removeAllChats()
        },
        minimizeChats() {
            chat.minimizeAllChats()
        },
        handleInput(target) {
            if (!target.value) return (this.visibleContacts = [])
            this.visibleContacts = this.contacts.filter(
                (v) =>
                    (v.User &&
                        `${v.User.firstName} ${v.User.lastName}`
                            .toLowerCase()
                            .includes(target.value.toLowerCase())) ||
                    (v.Group &&
                        v.Group.groupName.toLowerCase().includes(target.value.toLowerCase()))
            )
        },
        readMessage(message, key) {
            if (!message || message.isRead) return
            chat.getAllConversations().value.get(key).lastMessage.isRead = true
            message.isRead = true
            axios.post('/api/setMessageIsRead', {
                messageId: message.messageId,
                isRead: true
            })
        }
    },
    watch: {
        async newMessageClicked() {
            if (!this.newMessageClicked) return (this.visibleContacts = [])
            if (chat.isOverFlown()) chat.minimizeLastChat()
            axios.get('/api/getContacts').then((res) => (this.contacts = res.data))
            await this.$nextTick()
            document.getElementById('contactSearch').focus()
        }
    }
}
</script>

<style scoped>
.chat-container {
    background-color: white;
    width: 328px;
    height: 591px;
    border-top-right-radius: 8px;
    border-top-left-radius: 8px;
    margin-left: 10px;
    box-shadow:
        0 12px 28px 0 rgba(0, 0, 0, 0.1),
        0 2px 4px 0 rgba(0, 0, 0, 0.1);
}
</style>
