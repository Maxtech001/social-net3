<template>
    <div v-if="group && user">
        <HeaderComponent :user="user"></HeaderComponent>
        <GroupComponent :propGroup="group" :user="user"></GroupComponent>
    </div>
</template>

<script>
import axios from 'axios'
import HeaderComponent from '../components/header/HeaderComponent.vue'
import GroupComponent from '../components/group/GroupComponent.vue'
export default {
    components: { HeaderComponent, GroupComponent },
    created() {
        axios.get('/api/getCurrentUser').then((res) => {
            this.user = res.data
            if (!this.user) return this.$router.push({ path: '/login' })
            axios
                .post('/api/getGroup', {
                    groupId: parseInt(this.$route.params.groupId)
                })
                .then((res) => (this.group = res.data))
        })
    },
    data() {
        return {
            user: null,
            group: null
        }
    }
}
</script>
