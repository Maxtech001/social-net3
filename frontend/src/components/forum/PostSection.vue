<template>
    <div class="d-flex justify-center" style="width: 100%">
        <div class="d-flex flex-column flex-grow-1" style="max-width: 680px">
            <v-card v-if="isViewedUser && !isGroupView" class="px-4 py-3">
                <v-form validate-on="submit">
                    <v-row class="ma-0 align-center">
                        <div class="pr-2">
                            <v-img
                                v-if="user.avatar"
                                :src="user.avatar"
                                :width="40"
                                :height="40"
                                cover
                                style="border-radius: 50%"
                            ></v-img>
                            <v-icon v-else size="40">mdi-account-circle</v-icon>
                        </div>
                        <v-btn
                            text
                            class="d-flex flex-grow-1 post-button justify-start"
                            rounded
                            variant="text"
                            @click="openPostCreateDialog = true"
                            style="
                                background-color: #f0f2f5;
                                text-align: left;
                                text-transform: unset !important;
                                letter-spacing: unset;
                            "
                            >Post your thoughts..</v-btn
                        >
                    </v-row>
                    <v-divider class="mt-2"></v-divider>
                    <v-row class="ma-0 justify-center mt-2">
                        <v-btn variant="text" @click="openPostCreateDialogWithPhotoSection()"
                            ><v-icon class="mr-2">mdi-image-multiple-outline</v-icon>Photo</v-btn
                        >
                    </v-row>
                </v-form>
            </v-card>
            <PostComponent
                v-for="post in posts"
                :key="post.postId"
                :user="user"
                :post="post"
                @openPost="openPost"
            ></PostComponent>
            <v-dialog v-model="openPostViewDialog" width="680"
                ><PostComponent
                    :user="user"
                    :post="selectedPost"
                    :showCommentSection="true"
                    @storedComment="updateCommentCount"
                ></PostComponent
            ></v-dialog>
            <v-dialog v-model="openPostCreateDialog" width="500">
                <v-card>
                    <v-row class="ma-0 mt-1 px-4 align-center" style="height: 60px">
                        <h2 class="d-flex flex-grow-1 justify-center" style="font-weight: bold">
                            Create post
                        </h2>
                        <v-btn icon variant="flat" @click="openPostCreateDialog = false"
                            ><v-icon>mdi-close</v-icon></v-btn
                        >
                    </v-row>
                    <v-divider></v-divider>
                    <div class="pa-4">
                        <v-row class="ma-0 align-center mt-n2">
                            <div class="pr-2">
                                <v-img
                                    v-if="user.avatar"
                                    :src="user.avatar"
                                    :width="40"
                                    :height="40"
                                    cover
                                    style="border-radius: 50%"
                                ></v-img>
                                <v-icon v-else size="40">mdi-account-circle</v-icon>
                            </div>
                            <span style="font-weight: 600">{{
                                `${user.firstName} ${user.lastName}`
                            }}</span>
                        </v-row>
                        <v-textarea
                            v-model="postContent"
                            rows="1"
                            variant="plain"
                            rounded
                            placeholder="Post your thoughts..."
                            hide-details
                            no-resize
                            @keydown.enter="storePost"
                        ></v-textarea>
                        <div
                            v-if="openPhotoSection"
                            class="pa-2"
                            style="border: 1px solid #ccc; border-radius: 8px; position: relative"
                            @dragenter.prevent="isDragged = true"
                        >
                            <v-file-input
                                ref="file"
                                style="display: none"
                                multiple
                                accept="image/*"
                                v-model="fileInput"
                                @input="handleInput($event)"
                            />
                            <div
                                v-if="isDragged"
                                @dragenter.prevent="isDragged = true"
                                @dragleave.prevent="isDragged = false"
                                @dragexit.prevent="isDragged = false"
                                @dragover.prevent
                                @drop.stop.prevent="onDrop($event)"
                                class="d-flex justify-center align-center"
                                style="
                                    position: absolute;
                                    top: 0;
                                    left: 0;
                                    height: 100%;
                                    width: 100%;
                                    z-index: 1;
                                    background-color: rgba(0, 0, 0, 0.1);
                                    backdrop-filter: blur(3px);
                                    font-size: 20px;
                                    font-weight: 500;
                                    border-radius: 8px;
                                "
                            >
                                Drop files
                            </div>
                            <div
                                v-if="!addedImageArray.length"
                                class="d-flex justify-center align-center flex-column file-input"
                                style="position: relative"
                                @click="$refs.file.click()"
                            >
                                <v-icon
                                    class="pa-6"
                                    size="36"
                                    style="
                                        background-color: gray;
                                        border-radius: 50%;
                                        background-color: rgb(228, 230, 235);
                                    "
                                    >mdi-image-plus-outline</v-icon
                                >
                                <h3 style="font-weight: 600">Add photos</h3>
                                <span style="font-size: 12px">or drag and drop</span>
                                <v-btn
                                    class="ma-2"
                                    density="compact"
                                    icon
                                    style="position: absolute; top: 0; right: 0"
                                    @click.stop
                                    @click="closePhotoSection()"
                                    ><v-icon>mdi-close</v-icon></v-btn
                                >
                            </div>
                            <div v-else style="height: 221px; overflow: auto">
                                <div
                                    v-for="(image, index) in addedImageArray"
                                    :key="index"
                                    style="position: relative"
                                >
                                    <v-img cover :src="createURL(image)"></v-img>
                                    <v-btn
                                        class="ma-2"
                                        density="compact"
                                        icon
                                        style="position: absolute; top: 0; right: 0"
                                        @click.stop
                                        @click="addedImageArray.splice(index, 1)"
                                        ><v-icon>mdi-close</v-icon></v-btn
                                    >
                                </div>
                            </div>
                        </div>
                        <v-row v-if="addedImageArray.length" class="ma-0 mt-2">
                            <v-btn variant="tonal" class="mr-2" @click="$refs.file.click()"
                                ><v-icon class="mr-2">mdi-image-plus-outline</v-icon>Add
                                images</v-btn
                            >
                            <v-btn
                                variant="tonal"
                                style="background-color: rgb(255, 94, 94)"
                                @click="addedImageArray = []"
                                ><v-icon class="mr-2">mdi-cancel</v-icon>Remove all</v-btn
                            >
                        </v-row>
                        <v-row
                            class="ma-0 align-center pa-2 mt-4"
                            style="
                                border: 1px solid #ced0d4;
                                border-radius: 8px;
                                box-shadow:
                                    rgba(60, 64, 67, 0.3) 0px 1px 2px 0px,
                                    rgba(60, 64, 67, 0.15) 0px 1px 3px 1px;
                            "
                        >
                            <span>Add to post</span>
                            <v-spacer></v-spacer>
                            <v-btn
                                icon
                                variant="text"
                                @click="openPhotoSection = true"
                                :active="openPhotoSection"
                            >
                                <v-icon color="green">mdi-image-multiple-outline</v-icon>
                            </v-btn>
                        </v-row>
                        <v-row class="ma-0 mt-2 align-center">
                            <span>Post privacy: </span>
                            <v-menu width="200">
                                <template v-slot:activator="{ props }">
                                    <v-btn v-bind="props" variant="text" rounded :disabled="!!group"
                                        >{{ postPrivacy
                                        }}<v-icon v-if="!group">mdi-chevron-down</v-icon></v-btn
                                    >
                                </template>
                                <v-list class="pa-2">
                                    <v-list-item rounded @click="postPrivacy = 'Private'">
                                        <v-row class="ma-0">
                                            <span>Private</span>
                                            <v-spacer></v-spacer>
                                            <v-icon v-if="postPrivacy === 'Private'"
                                                >mdi-check</v-icon
                                            >
                                        </v-row>
                                    </v-list-item>
                                    <v-list-item rounded @click="postPrivacy = 'Public'">
                                        <v-row class="ma-0">
                                            <span>Public</span>
                                            <v-spacer></v-spacer>
                                            <v-icon v-if="postPrivacy === 'Public'"
                                                >mdi-check</v-icon
                                            >
                                        </v-row>
                                    </v-list-item>
                                    <v-list-item rounded @click="postPrivacy = 'Specified'">
                                        <v-row class="ma-0">
                                            <span>Specified</span>
                                            <v-spacer></v-spacer>
                                            <v-icon v-if="postPrivacy === 'Specified'"
                                                >mdi-check</v-icon
                                            >
                                        </v-row>
                                    </v-list-item>
                                </v-list>
                            </v-menu>
                        </v-row>
                        <div v-if="postPrivacy === 'Specified'">
                            <v-autocomplete
                                v-model="selectedFollowers"
                                clearable
                                placeholder="Search followers"
                                density="compact"
                                variant="outlined"
                                rounded
                                multiple
                                :items="
                                    followers.map((v) => ({
                                        title: v.firstName + ' ' + v.lastName,
                                        value: v
                                    }))
                                "
                            ></v-autocomplete>
                        </div>
                        <v-btn
                            block
                            color="blue"
                            class="mt-4"
                            :disabled="
                                !postContent.length ||
                                !(postPrivacy !== 'Specified' || selectedFollowers?.length)
                            "
                            @click="storePost()"
                            >Post</v-btn
                        >
                    </div>
                </v-card>
            </v-dialog>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import PostComponent from './PostComponent.vue'
export default {
    components: {
        PostComponent
    },
    props: {
        group: {
            type: Object,
            default: null
        },
        user: Object,
        viewedUser: Object,
        isViewedUser: {
            type: Boolean,
            default: true
        },
        propPosts: {
            type: Array,
            default: () => []
        },
        isGroupView: {
            type: Boolean,
            default: false
        },
        propOpenPostCeateDialog: {
            type: Boolean,
            default: false
        },
        propOpenPhotoSection: {
            type: Boolean,
            default: false
        }
    },
    created() {
        this.posts = this.propPosts
        this.openPostCreateDialog = this.propOpenPostCeateDialog
        this.openPhotoSection = this.propOpenPhotoSection
    },
    data() {
        return {
            postContent: '',
            postPrivacy: 'Private',
            posts: [],
            openPostViewDialog: false,
            openPostCreateDialog: false,
            selectedPost: null,
            openPhotoSection: false,
            addedImageArray: [],
            fileInput: [],
            isDragged: false,
            followers: [],
            selectedFollowers: [],
            privacyType: {
                Private: 0,
                Public: 1,
                Specified: 2
            }
        }
    },
    methods: {
        storePost(event) {
            if (event) event.preventDefault()
            if (!this.postContent.length) return
            let formData = new FormData()
            formData.append('content', this.postContent)
            formData.append('groupId', this.group ? this.group.groupId : null)
            formData.append('privacyType', this.privacyType[this.postPrivacy])
            if (this.privacyType[this.postPrivacy] == 2)
                formData.append('specifiedUsers', JSON.stringify(this.selectedFollowers))
            this.addedImageArray.forEach((img) => {
                formData.append('images', img)
            })
            axios
                .post('/api/post', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                })
                .then((res) => {
                    console.log(res.data)
                    this.postContent = ''
                    this.posts = [res.data, ...this.posts]
                    this.openPostCreateDialog = false
                })
        },
        openPost(post) {
            this.selectedPost = post
            this.openPostViewDialog = true
        },
        updateCommentCount() {
            this.selectedPost.totalComments += 1
        },
        onDrop(event) {
            this.isDragged = false
            Array.from(event.dataTransfer.files).forEach(
                (element) => element.type.includes('image/') && this.addedImageArray.push(element)
            )
        },
        closePhotoSection() {
            this.openPhotoSection = false
            this.addedImageArray = []
        },
        createURL(image) {
            return URL.createObjectURL(image)
        },
        handleInput(event) {
            Array.from(event.target.files).forEach(
                (f) => f.type.includes('image/') && this.addedImageArray.push(f)
            )
            this.fileInput = []
        },
        openPostCreateDialogWithPhotoSection() {
            this.openPostCreateDialog = true
            this.openPhotoSection = true
        }
    },
    watch: {
        propPosts() {
            this.posts = this.propPosts
        },
        propOpenPostCeateDialog() {
            this.openPostCreateDialog = this.propOpenPostCeateDialog
        },
        openPostCreateDialog() {
            if (!this.openPostCreateDialog) return (this.selectedFollowers = [])
            axios.get('/api/getFollowers').then((res) => (this.followers = res.data))
        }
    }
}
</script>

<style>
.v-input__control .v-field--no-label {
    margin-top: 0px;
}
.post-button:hover {
    background-color: rgba(0, 0, 0, 0.05);
}
.file-input {
    height: 221px;
    background-color: rgb(247, 248, 250);
    border-radius: 8px;
    cursor: pointer;
}
.file-input:hover {
    background-color: rgba(0, 0, 0, 0.05);
}
</style>
