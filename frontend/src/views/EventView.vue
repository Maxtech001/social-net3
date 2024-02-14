<template>
    <div v-if="user && event">
        <HeaderComponent :user="user"></HeaderComponent>
        <EventComponent :propEvent="event" :user="user"></EventComponent>
    </div>
</template>

<script>
import HeaderComponent from '../components/header/HeaderComponent.vue'
import EventComponent from '../components/event/EventComponent.vue'
import axios from 'axios'
export default {
    components: {
        HeaderComponent,
        EventComponent
    },
    created() {
        axios.get('/api/getCurrentUser').then((res) => {
            this.user = res.data
            if (!this.user) return this.$router.push({ path: '/login' })
            axios
                .post('/api/getEvent', {
                    eventId: parseInt(this.$route.params.eventId)
                })
                .then((res) => (this.event = res.data))
        })
    },
    data() {
        return {
            user: null,
            event: null
        }
    }
}
</script>
