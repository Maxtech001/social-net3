import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AuthenticationForm from '../components/authentication/AuthenticationForm.vue'
import UserView from '../views/UserView.vue'
import CreateGroupView from '../views/CreateGroupView.vue'
import GroupView from '../views/GroupView.vue'
import EventView from '../views/EventView.vue'
import FollowView from '../views/FollowView.vue'
import GroupsView from '../views/GroupsView.vue'
import SearchView from '../views/SearchView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: HomeView,
            props: true
        },
        {
            path: '/login',
            name: 'login',
            component: AuthenticationForm,
            props: { propTab: 'login' }
        },
        {
            path: '/register',
            name: 'register',
            component: AuthenticationForm,
            props: { propTab: 'register' }
        },
        {
            path: '/user/:userId',
            name: 'user',
            component: UserView
        },
        {
            path: '/groups/create',
            name: 'createGroup',
            component: CreateGroupView
        },
        {
            path: '/groups/:groupId',
            name: 'group',
            component: GroupView
        },
        {
            path: '/event/:eventId',
            name: 'event',
            component: EventView
        },
        {
            path: '/follows',
            name: 'follows',
            component: FollowView
        },
        {
            path: '/groups',
            name: 'groups',
            component: GroupsView
        },
        {
            path: '/search/:searchInput',
            name: 'search',
            component: SearchView
        },
        {
            path: '/createPost',
            name: 'createPost',
            component: HomeView,
            props: { propOpenPostCeateDialog: true }
        }
    ]
})

export default router
