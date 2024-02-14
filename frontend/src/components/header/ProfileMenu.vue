<template>
    <v-menu width="300">
        <template v-slot:activator="{ props }">
            <v-icon
                v-if="user.avatar"
                v-bind="props"
                class="mr-4"
                size="40"
                v-ripple
                style="border-radius: 50%"
            >
                <v-img
                    :src="user.avatar"
                    :width="40"
                    :height="40"
                    cover
                    style="border-radius: 50%"
                />
            </v-icon>
            <v-icon
                v-else
                v-bind="props"
                class="mr-4"
                size="40"
                color="black"
                v-ripple
                style="border-radius: 50%"
                >mdi-account-circle</v-icon
            >
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
            <v-card class="pa-1 mb-3" style="border-radius: 8px !important">
                <v-list-item @click="navigateProfile()" rounded>
                    <v-icon v-if="user.avatar" class="mr-2" size="36">
                        <v-img
                            :src="user.avatar"
                            :width="36"
                            :height="36"
                            style="border-radius: 50%"
                            cover
                        />
                    </v-icon>
                    <v-icon v-else class="mr-2" size="36">mdi-account-circle</v-icon>
                    {{ `${user.firstName} ${user.lastName}` }}
                </v-list-item>
            </v-card>
            <v-list-item @click="logout()" rounded>
                <v-icon class="mr-2" size="36">mdi-logout</v-icon> {{ 'Log out' }}
            </v-list-item>
        </v-list>
    </v-menu>
</template>

<script>
import axios from 'axios'
export default {
    props: {
        user: {
            type: Object,
            required: true
        }
    },
    methods: {
        navigateProfile() {
            this.$router.push({ name: 'user', params: { userId: this.user.userId } })
        },
        logout() {
            axios
                .get('/api/logout')
                .then(() => this.$router.go())
                .catch((err) => console.error(err.message))
        }
    }
}
</script>
