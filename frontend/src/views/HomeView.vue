<template>
    <main>
        <ApplicationComponent
            v-if="user"
            :user="user"
            :propOpenPostCeateDialog="propOpenPostCeateDialog"
        ></ApplicationComponent>
    </main>
</template>

<script>
import axios from 'axios'
import ApplicationComponent from '../components/forum/ApplicationComponent.vue'
export default {
    components: { ApplicationComponent },
    props: {
        propOpenPostCeateDialog: {
            type: Boolean,
            default: false
        }
    },
    created() {
        axios.get('/api/getCurrentUser').then((res) => {
            this.user = res.data
            if (!this.user) return this.$router.push({ path: '/login' })
        })
    },
    data() {
        return {
            user: null
        }
    }
}
</script>
