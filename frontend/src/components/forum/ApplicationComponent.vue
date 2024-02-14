<template>
    <main>
        <HeaderComponent
            :user="user"
            :isHomeView="true"
            @homeFeed="getHomeFeedPost"
            @groupFeed="getGroupFeedPost"
        ></HeaderComponent>
        <PostSection
            v-if="posts"
            class="mt-6 mb-6"
            :propOpenPostCeateDialog="propOpenPostCeateDialog"
            :propPosts="posts"
            :user="user"
            :isGroupView="isGroupView"
        ></PostSection>
    </main>
</template>

<script>
import PostSection from './PostSection.vue'
import HeaderComponent from '../header/HeaderComponent.vue'
import axios from 'axios'
export default {
    components: {
        PostSection,
        HeaderComponent
    },
    props: {
        user: Object,
        propOpenPostCeateDialog: {
            type: Boolean,
            default: false
        }
    },
    created() {
        this.getHomeFeedPost()
    },
    data() {
        return {
            value: true,
            posts: null,
            isGroupView: false
        }
    },
    methods: {
        getHomeFeedPost() {
            axios.get('/api/getAllPosts').then((res) => {
                this.posts = res.data
                this.isGroupView = false
            })
        },
        getGroupFeedPost() {
            axios.get('/api/getAllGroupPosts').then((res) => {
                this.posts = res.data
                this.isGroupView = true
            })
        }
    }
}
</script>
