<template>
    <v-form ref="form" @submit.prevent="login">
        <v-text-field
            variant="outlined"
            density="compact"
            rounded
            v-model="formData.email"
            label="Email"
            :rules="emailRules"
        >
        </v-text-field>
        <v-text-field
            variant="outlined"
            density="compact"
            rounded
            v-model="formData.password"
            label="Password"
            type="password"
            :rules="passwordRules"
        ></v-text-field>
        <v-row v-if="errorMessage" class="ma-0 justify-center">
            <span style="color: red">{{ errorMessage }}</span>
        </v-row>
        <v-row justify="center" class="ma-0 mt-12">
            <v-btn type="submit" color="black">LOG IN</v-btn>
        </v-row>
    </v-form>
</template>

<script>
import axios from 'axios'
import { createWSConnection } from '../../conn.js'
import { getUserConversations } from '../../chat.js'
export default {
    data() {
        return {
            errorMessage: '',
            formData: {
                email: '',
                password: ''
            },
            emailRules: [
                (v) => !!v || 'Email is required',
                (v) => /.+@.+\..+/.test(v) || 'Invalid email format'
            ],
            passwordRules: [
                (v) => !!v || 'Password is required',
                (v) => (v && v.length >= 6) || 'Password must be at least 6 characters'
            ]
        }
    },
    methods: {
        login() {
            this.errorMessage = ''
            this.$refs.form.validate().then((data) => {
                if (data.valid)
                    axios
                        .post('/api/login', {
                            email: this.formData.email,
                            password: this.formData.password
                        })
                        .then(() => {
                            createWSConnection()
                            getUserConversations()
                            this.$router.push({ path: '/' })
                        })
                        .catch((err) => (this.errorMessage = err.response.data))
            })
        }
    }
}
</script>

<style>
.v-field {
    margin-top: 8px;
}
</style>
