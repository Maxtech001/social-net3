<template>
    <div class="d-flex justify-center align-center" style="height: 100%">
        <v-card width="400">
            <v-tabs v-model="tab" bg-color="black">
                <v-tab value="login" style="width: 50%">Login</v-tab>
                <v-tab value="register" style="width: 50%">Register</v-tab>
            </v-tabs>

            <v-card-text>
                <v-window v-model="tab">
                    <v-window-item value="login">
                        <LoginInput />
                    </v-window-item>

                    <v-window-item value="register">
                        <RegisterInput />
                    </v-window-item>
                </v-window>
            </v-card-text>
        </v-card>
    </div>
</template>

<script>
import LoginInput from './LoginInput.vue'
import RegisterInput from './RegisterInput.vue'
export default {
    components: {
        LoginInput,
        RegisterInput
    },
    props: {
        propTab: String
    },
    created() {
        if (this.propTab) this.tab = this.propTab
    },
    data: () => ({
        tab: 'login'
    })
}
</script>
