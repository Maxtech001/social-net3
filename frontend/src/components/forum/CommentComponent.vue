<template>
    <v-row class="ma-0" style="width: 100%">
        <v-icon v-if="comment.authorAvatar" class="mr-2" size="32">
            <v-img
                :src="comment.authorAvatar"
                :width="32"
                :height="32"
                style="border-radius: 50%"
                cover
            />
        </v-icon>
        <v-icon v-else size="32" class="mr-2">mdi-account-circle</v-icon>
        <div>
            <div
                class="px-3 py-2"
                style="background-color: rgb(240, 242, 245); border-radius: 18px"
            >
                <v-row class="ma-0">
                    <span style="font-weight: 600">{{ comment.authorName }}</span>
                </v-row>
                <v-row class="ma-0">
                    <span>{{ comment.content }}</span>
                </v-row>
                <img
                    v-if="comment.imageArray.length"
                    :src="comment.imageArray[0].imagePath"
                    cover
                    style="max-height: 100px"
                />
            </div>
            <span style="font-size: 12px">{{ comment.createdAt }}</span>
        </div>
    </v-row>
</template>

<script>
export default {
    props: {
        comment: {
            type: Object,
            required: true
        }
    }
}
</script>
