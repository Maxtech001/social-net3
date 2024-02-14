<template>
    <v-menu width="360">
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
                <v-icon size="20" color="black">mdi-bell</v-icon>
                <div
                    v-if="notifications.filter((v) => !v.isRead).length"
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
                    {{ notifications.filter((v) => !v.isRead).length }}
                </div>
            </v-btn>
        </template>
        <v-list
            @click.stop.prevent
            style="
                user-select: none;
                -webkit-user-select: none;
                -moz-user-select: none;
                -khtml-user-select: none;
                -ms-user-select: none;
            "
        >
            <v-row class="ma-0 pa-2">
                <h2 style="font-weight: bold">Notifications</h2>
            </v-row>
            <div
                v-if="!notifications.length"
                class="d-flex justify-center align-center flex-column pa-4"
            >
                <div class="mb-5">
                    <v-icon size="112" style="color: rgb(101, 103, 107)"
                        >mdi-bell-badge-outline</v-icon
                    >
                </div>
                <span
                    style="
                        color: rgb(101, 103, 107);
                        font-weight: 700;
                        line-height: 24px;
                        font-size: 20px;
                    "
                    >You have no notifications</span
                >
            </div>
            <v-list-item
                v-for="notification in notifications"
                :key="notification.notificationId"
                rounded
                @click="handleClick(notification)"
            >
                <v-row class="ma-0 align-center">
                    <div class="d-flex mr-2" style="height: 56px; width: 56px">
                        <v-img
                            v-if="notification.notificationType == 0"
                            :src="
                                notification.follower.avatar.length
                                    ? notification.follower.avatar
                                    : '/account.svg'
                            "
                            style="background-color: #c9ccd1; border-radius: 50%"
                            cover
                        >
                        </v-img>
                        <v-img
                            v-if="
                                notification.notificationType == 1 ||
                                notification.notificationType == 2
                            "
                            :src="'/group-stock-image.jpg'"
                            style="background-color: #c9ccd1; border-radius: 50%"
                            cover
                        >
                        </v-img>
                        <v-img
                            v-if="notification.notificationType == 3"
                            :src="'/event-stock-image.jpg'"
                            style="background-color: #c9ccd1; border-radius: 50%"
                            cover
                        >
                        </v-img>
                    </div>
                    <div class="flex-grow-1" style="flex-basis: 0">
                        <span>{{ notification.content }}</span>
                    </div>
                    <div style="width: 48px">
                        <v-icon v-if="!notification.isRead" size="48" color="blue"
                            >mdi-circle-small</v-icon
                        >
                    </div>
                </v-row>
            </v-list-item>
        </v-list>
    </v-menu>
</template>

<script>
import axios from 'axios'
export default {
    data() {
        return {
            notifications: []
        }
    },
    created() {
        axios.get('/api/getAllUserNotifications').then((res) => (this.notifications = res.data))
    },
    methods: {
        handleClick(notification) {
            axios
                .post('/api/setNotificationIsRead', {
                    notificationId: notification.notificationId,
                    isRead: true
                })
                .then(() => {
                    notification.isRead = true
                    switch (notification.notificationType) {
                        case 0: {
                            return this.$router.push({ path: '/follows' })
                        }
                        case 1: {
                            return this.$router.push({ path: '/groups' })
                        }
                        case 2: {
                            return this.$router.push({ path: '/groups' })
                        }
                        case 3: {
                            return this.$router.push({
                                name: 'event',
                                params: { eventId: notification.eventId }
                            })
                        }
                    }
                })
        }
    }
}
</script>
