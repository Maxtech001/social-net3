import axios from 'axios'
import chat, { getUserConversations } from './chat.js'
import { ref } from 'vue'

export var user = ref()

export let conn = createWSConnection()

class Event {
    // Each Event needs a Type
    // The payload is not required
    constructor(type, payload) {
        this.type = type
        this.payload = payload
    }
}

class SendMessageEvent {
    constructor(message, from, to, groupId) {
        this.message = message
        this.from = from
        this.to = to
        this.groupId = groupId
    }
}
/**
 * NewMessageEvent is messages comming from clients
 * */
class NewMessageEvent {
    constructor(message, from, sent, groupId, fromObj, isRead, messageId) {
        this.message = message
        this.from = from
        this.fromObj = fromObj
        this.sent = sent
        this.groupId = groupId
        this.isRead = isRead
        this.messageId = messageId
    }
}

export async function createWSConnection() {
    await axios.get('/api/getCurrentUser').then((res) => (user.value = res.data))
    if (!user.value) return null
    conn = new WebSocket('ws://' + document.location.host + '/api/ws')

    conn.onmessage = function (evt) {
        // parse websocket message as JSON
        const eventData = JSON.parse(evt.data)
        // Assign JSON data to new Event Object
        const event = Object.assign(new Event(), eventData)
        // Let router manage message
        routeEvent(event)
    }
}

function routeEvent(event) {
    if (event.type === undefined) {
        alert("no 'type' field in event")
    }
    switch (event.type) {
        case 'new_message': {
            const messageEvent = Object.assign(new NewMessageEvent(), event.payload)
            appendChatMessage(messageEvent)
            break
        }
        default: {
            alert('unsupported message type')
            break
        }
    }
}
function appendChatMessage(messageEvent) {
    if (!messageEvent.groupId) {
        appendMessageByChatId(messageEvent, 'u-' + messageEvent.from)
        appendMessageByChatId(messageEvent, 'u-' + messageEvent.to)
        getUserConversations()
    } else {
        appendMessageByChatId(messageEvent, 'g-' + messageEvent.groupId)
        getUserConversations()
    }
}

function appendMessageByChatId(messageEvent, id) {
    let result = chat.getOpenedChatArray().value.get(id)
    if (!result) result = chat.getMinimizedChatArray().value.get(id)
    if (result)
        result.messages.unshift({
            message: messageEvent.message,
            senderId: messageEvent.from,
            senderObj: messageEvent.fromObj,
            isRead: messageEvent.from == user.value.userId || messageEvent.isRead,
            messageId: messageEvent.messageId
        })
}

function sendEvent(eventName, payload) {
    // Create a event Object with a event named send_message
    const event = new Event(eventName, payload)
    // Format as JSON and send
    conn.send(JSON.stringify(event))
}

export function sendMessageConn(message, receiverObj) {
    if (!message) return
    let outgoingEvent = receiverObj.userId
        ? new SendMessageEvent(message, user.value.userId, receiverObj.userId)
        : new SendMessageEvent(message, user.value.userId, null, receiverObj.groupId)
    sendEvent('send_message', outgoingEvent)
}
