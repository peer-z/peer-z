import {createRouter, createWebHistory} from 'vue-router'
import Home from "../views/Home.vue";
import About from "../views/About.vue";
import NotFound from "../views/NotFound.vue";
// import Service from "../views/Service.vue";
import Checkers from "../views/Checkers.vue";
import Messaging from "../views/Messaging.vue";
import Profile from "../views/Profile.vue";
import Chess from "../views/Chess.vue";

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'Home',
            component: Home,
        },
        {
            path: '/messaging',
            name: 'Messaging',
            component: Messaging,
            props: false,
        },
        {
            path: '/profile',
            name: 'Profile',
            component: Profile,
            props: false,
        },
        {
            path: '/web',
            name: 'Web',
            // component: Web,
            component: NotFound,
            props: false,
        },
        {
            path: '/checkers',
            name: 'Checkers',
            component: Checkers,
            props: false,
        },
        {
            path: '/chess',
            name: 'Chess',
            component: Chess,
            props: false,
        },
        {
            path: '/chat',
            name: 'Chat',
            // component: Chat,
            component: NotFound,
            props: false,
        },
        {
            path: '/storage',
            name: 'Storage',
            // component: Storage,
            component: NotFound,
            props: false,
        },
        {
            path: '/about',
            name: 'About',
            component: About,
            // route level code-splitting
            // this generates a separate chunk (About.[hash].js) for this route
            // which is lazy-loaded when the route is visited.
            // component: () => import('../views/About.vue')
        },
        {
            path: '/:pathMatch(.*)*',
            name: 'NotFound',
            component: NotFound,
            props: true,
        },

    ]
})

export default router
