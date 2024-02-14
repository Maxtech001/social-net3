<template>
    <v-card class="mt-6 pt-2 pl-4 pr-4" :class="{ 'pb-2': !showCommentSection }">
        <v-row class="ma-0 mb-4 align-end">
            <v-icon v-if="post.authorAvatar" class="mr-2" size="40">
                <v-img
                    :src="post.authorAvatar"
                    :width="40"
                    :height="40"
                    style="border-radius: 50%"
                    cover
                />
            </v-icon>
            <v-icon v-else size="40" class="mr-2">mdi-account-circle</v-icon>

            <span v-if="post.group?.groupId"
                >Group: {{ post.group.groupName }} <br />
                {{ post.authorName }} <br />
                {{ post.createdAt }}</span
            >
            <span v-else
                >{{ post.authorName }} <br />
                {{ post.createdAt }}</span
            >
        </v-row>
        <span>{{ post.content }}</span>
        <div>
            <v-img
                v-for="(image, index) in post.imageArray"
                :key="index"
                class="mt-2"
                :src="image.imagePath"
                cover
                style="max-height: 400px"
            />
        </div>
        <div class="d-flex justify-end">
            <span class="comment-count" style="font-size: 14px" @click="$emit('openPost', post)"
                >{{ post.totalComments }} comments</span
            >
        </div>
        <v-divider class="mt-2"></v-divider>
        <v-row class="ma-0">
            <v-col class="pa-0 pt-2" cols="6" align="center">
                <v-btn
                    variant="text"
                    style="text-transform: unset"
                    block
                    @click="$emit('openPost', post)"
                    ><v-icon class="mr-2">mdi-comment-outline</v-icon>Comment</v-btn
                >
            </v-col>
            <v-col class="pa-0 pt-2" cols="6" align="center">
                <v-btn variant="text" style="text-transform: unset" block
                    ><v-icon class="mr-2">mdi-share-outline</v-icon>Share</v-btn
                >
            </v-col>
        </v-row>
        <v-row v-if="showCommentSection" class="ma-0 pt-2">
            <v-row v-if="comments.length > 0" class="ma-0" style="width: 100%">
                <v-divider></v-divider>
                <div class="py-2">
                    <CommentComponent
                        v-for="comment in comments"
                        :key="comment.Id"
                        :comment="comment"
                    ></CommentComponent>
                </div>
            </v-row>
            <v-form
                ref="commentTextfield"
                validate-on="submit"
                class="pb-4"
                style="position: sticky; width: 100%; bottom: 0; background-color: white"
            >
                <v-divider></v-divider>
                <v-row class="pt-4 ma-0 align-start">
                    <v-icon v-if="user.avatar" class="mr-2" size="32">
                        <v-img
                            :src="user.avatar"
                            :width="32"
                            :height="32"
                            style="border-radius: 50%"
                            cover
                        />
                    </v-icon>
                    <v-icon v-else size="32" class="mr-2">mdi-account-circle</v-icon>
                    <v-textarea
                        v-model="commentContent"
                        variant="outlined"
                        rounded
                        placeholder="Write a comment..."
                        no-resize
                        rows="1"
                        auto-grow
                        :rules="[(v) => !!v]"
                        hide-details
                        density="compact"
                        @keydown.enter="storeComment"
                    ></v-textarea>
                    <v-file-input
                        ref="file"
                        style="display: none"
                        accept="image/*"
                        :rules="[(v) => !v || v[0].type.includes('image/')]"
                        v-model="fileInput"
                    />
                </v-row>
                <v-row class="ma-0 ml-10 mt-1">
                    <v-icon
                        @click="this.$refs.file.click()"
                        :color="
                            fileInput?.length && !fileInput[0].type.includes('image/') ? 'red' : ''
                        "
                        >mdi-camera-plus-outline</v-icon
                    >
                </v-row>
                <div
                    v-if="fileInput?.length && fileInput[0].type.includes('image/')"
                    class="ml-10 mt-4"
                    style="position: relative; width: fit-content"
                >
                    <img :src="createURL(fileInput[0])" style="max-height: 80px" />
                    <v-btn
                        density="compact"
                        class="ma-1"
                        icon
                        style="position: absolute; top: 0; right: 0; height: 20px; width: 20px"
                        @click.stop
                        @click="fileInput = null"
                        ><v-icon size="14" style="margin-left: 1px">mdi-close</v-icon></v-btn
                    >
                </div>
            </v-form>
        </v-row>
    </v-card>
</template>

<script>
import axios from 'axios'
import CommentComponent from './CommentComponent.vue'
export default {
    components: {
        CommentComponent
    },
    props: {
        user: {
            type: Object,
            required: true
        },
        post: {
            type: Object,
            required: true
        },
        showCommentSection: {
            type: Boolean,
            default: false
        }
    },
    data() {
        return {
            commentContent: '',
            comments: [],
            fileInput: null
        }
    },
    created() {
        if (!this.showCommentSection) return
        axios
            .post('/api/getComments', {
                postId: parseInt(this.post.postId)
            })
            .then((res) => (this.comments = res.data))
    },
    methods: {
        storeComment(event) {
            if (event) event.preventDefault()
            this.$refs.commentTextfield.validate().then((data) => {
                if (data.valid) {
                    let formData = new FormData()
                    formData.append('content', this.commentContent)
                    formData.append('postId', parseInt(this.post.postId))
                    formData.append('images', this.fileInput.length && this.fileInput[0])
                    axios
                        .post('/api/comment', formData, {
                            headers: {
                                'Content-Type': 'multipart/form-data'
                            }
                        })
                        .then((res) => {
                            this.commentContent = ''
                            this.comments = [res.data, ...this.comments]
                            this.$emit('storedComment')
                            this.fileInput = null
                        })
                }
            })
        },
        createURL(image) {
            return URL.createObjectURL(image)
        }
    },
    watch: {
        fileInput() {
            if (this.fileInput && !this.fileInput[0].type.includes('image/')) this.fileInput = null
        }
    }
}
</script>

<style scoped>
.comment-count:hover {
    text-decoration: underline;
    cursor: pointer;
}
</style>
