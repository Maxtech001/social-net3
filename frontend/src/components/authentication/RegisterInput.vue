<template>
    <v-form ref="form2" @submit.prevent="register">
        <v-text-field
            variant="outlined"
            density="compact"
            rounded
            v-model="formData.email"
            label="Email"
            required
            :rules="emailRules"
        ></v-text-field>
        <v-text-field
            variant="outlined"
            density="compact"
            rounded
            v-model="formData.password"
            label="Password"
            type="password"
            required
            :rules="passwordRules"
        ></v-text-field>
        <v-text-field
            variant="outlined"
            density="compact"
            rounded
            v-model="formData.firstName"
            label="First Name"
            required
            :rules="nameRules"
        ></v-text-field>
        <v-text-field
            variant="outlined"
            density="compact"
            rounded
            v-model="formData.lastName"
            label="Last Name"
            required
            :rules="nameRules"
        ></v-text-field>
        <v-text-field
            variant="outlined"
            density="compact"
            rounded
            v-model="formData.dob"
            label="Date of Birth"
            type="date"
            required
            :rules="dobRules"
        ></v-text-field>
        <v-file-input
            variant="outlined"
            density="compact"
            accept="image/*"
            v-model="formData.avatar"
            label="Avatar/Image (Optional)"
            :rules="avatarRules"
            hint="Only images under 5 MB allowed"
            persistent-hint
        ></v-file-input>
        <v-text-field
            variant="outlined"
            density="compact"
            rounded
            v-model="formData.nickname"
            label="Nickname (Optional)"
        ></v-text-field>
        <v-textarea
            variant="outlined"
            density="compact"
            rounded
            v-model="formData.aboutMe"
            label="About Me (Optional)"
        ></v-textarea>
        <v-row v-if="errorMessage" class="ma-0 justify-center">
            <span style="color: red">{{ errorMessage }}</span>
        </v-row>
        <v-row justify="center" class="ma-0 mt-12">
            <v-btn type="submit" color="black">Register</v-btn>
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
                password: '',
                firstName: '',
                lastName: '',
                dob: '',
                avatar: null,
                nickname: '',
                aboutMe: ''
            },
            emailRules: [
                (v) => !!v || 'Email is required',
                (v) => /.+@.+\..+/.test(v) || 'Invalid email format'
            ],
            passwordRules: [
                (v) => !!v || 'Password is required',
                (v) => (v && v.length >= 6) || 'Password must be at least 6 characters'
            ],
            nameRules: [(v) => !!v || 'Name is required'],
            dobRules: [(v) => !!v || 'Date of Birth is required'],
            avatarRules: [
                (v) =>
                    !v ||
                    v.length == 0 ||
                    v[0].type.includes('image/') ||
                    'Only image files allowed',
                (v) =>
                    !v ||
                    v.length == 0 ||
                    v[0].size < 5 * (1 << 20) ||
                    'Image must not exceed 5 MB limit'
            ]
        }
    },
    methods: {
        register() {
            this.errorMessage = ''
            this.$refs.form2.validate().then((data) => {
                if (data.valid) {
                    let formData = new FormData()
                    formData.append('email', this.formData.email)
                    formData.append('password', this.formData.password)
                    formData.append('firstName', this.formData.firstName)
                    formData.append('lastName', this.formData.lastName)
                    formData.append('birthDate', this.formData.dob)
                    if (this.formData.avatar && this.formData.avatar.length == 1)
                        formData.append('avatar', this.formData.avatar[0])
                    formData.append('nickname', this.formData.nickname)
                    formData.append('aboutMe', this.formData.aboutMe)
                    axios
                        .post('/api/register', formData, {
                            headers: {
                                'Content-Type': 'multipart/form-data'
                            }
                        })
                        .then(() => {
                            createWSConnection()
                            getUserConversations()
                            this.$router.push({ path: '/' })
                        })
                        .catch((err) => (this.errorMessage = err.response.data))
                }
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
