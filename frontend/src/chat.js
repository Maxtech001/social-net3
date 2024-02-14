import { ref } from 'vue'
import axios from 'axios'
let openedChatArray = ref(new Map())
let minimizedChatArray = ref(new Map())
let chatHistoryMap = ref(new Map())
var user

getUserConversations()

export async function getUserConversations() {
    await axios.get('/api/getCurrentUser').then((res) => (user = res.data))
    if (!user) return
    chatHistoryMap.value = new Map()
    axios.get('/api/getUserConversations').then((res) => {
        res.data?.forEach((v) =>
            chatHistoryMap.value.set(
                v.userObj?.userId ? 'u-' + v.userObj.userId : 'g-' + v.groupObj.groupId,
                {
                    receiver: v.userObj?.userId ? v.userObj : v.groupObj,
                    lastMessage: v.lastMessage
                }
            )
        )
    })
}

function openChat(receiverObj) {
    if (
        isOverFlown() &&
        !openedChatArray.value.get(
            receiverObj.userId ? 'u-' + receiverObj.userId : 'g-' + receiverObj.groupId
        )
    )
        minimizeLastChat()
    axios
        .post('/api/getConversation', {
            ...(receiverObj.userId && { userId: receiverObj.userId }),
            ...(receiverObj.groupId && { groupId: receiverObj.groupId })
        })
        .then((res) => {
            if (res.data.length && !res.data[0].isRead) {
                res.data[0].isRead = true
                axios.post('/api/setMessageIsRead', {
                    messageId: res.data[0].messageId,
                    isRead: true
                })
            }
            openedChatArray.value.set(
                receiverObj.userId ? 'u-' + receiverObj.userId : 'g-' + receiverObj.groupId,
                {
                    receiver: receiverObj,
                    messages: res.data,
                    page: 0,
                    throttlePause: false
                }
            )
        })
}

function isOverFlown() {
    return document.getElementById('chat-component').scrollWidth + 338 > document.body.scrollWidth
}

function minimizeLastChat() {
    minimizeChat(openedChatArray.value.get(Array.from(openedChatArray.value.keys()).pop()).receiver)
    removeOpenedChat(Array.from(openedChatArray.value.keys()).pop())
}

function minimizeChat(receiverObj) {
    minimizedChatArray.value.set(
        receiverObj.userId ? 'u-' + receiverObj.userId : 'g-' + receiverObj.groupId,
        {
            receiver: receiverObj
        }
    )
}

function getOpenedChatArray() {
    return openedChatArray
}

function getMinimizedChatArray() {
    return minimizedChatArray
}

function getAllConversations() {
    return chatHistoryMap
}

function removeOpenedChat(key) {
    openedChatArray.value.delete(key)
}

function removeMinimizedChat(key) {
    minimizedChatArray.value.delete(key)
}

function removeAllChats() {
    openedChatArray.value = new Map()
    minimizedChatArray.value = new Map()
}

function minimizeAllChats() {
    minimizedChatArray.value = new Map(openedChatArray.value)
    openedChatArray.value = new Map()
}

export default {
    openChat: openChat,
    minimizeChat: minimizeChat,
    getOpenedChatArray: getOpenedChatArray,
    getMinimizedChatArray: getMinimizedChatArray,
    getAllConversations: getAllConversations,
    removeOpenedChat: removeOpenedChat,
    removeMinimizedChat: removeMinimizedChat,
    removeAllChats: removeAllChats,
    minimizeAllChats: minimizeAllChats,
    isOverFlown: isOverFlown,
    minimizeLastChat: minimizeLastChat
}
