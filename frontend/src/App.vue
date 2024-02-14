<template>
    <main style="height: 100%">
        <RouterView :key="$route.path" />
        <ChatComponent></ChatComponent>
    </main>
</template>

<script setup>
import { RouterView } from 'vue-router'
import ChatComponent from './components/chat/ChatComponent.vue'
</script>

<style scoped></style>
