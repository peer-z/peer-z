<!--
  - Copyright 2023 PeerGum
  -
  - Licensed under the Apache License, Version 2.0 (the "License");
  - you may not use this file except in compliance with the License.
  - You may obtain a copy of the License at
  -
  -     http://www.apache.org/licenses/LICENSE-2.0
  -
  - Unless required by applicable law or agreed to in writing, software
  - distributed under the License is distributed on an "AS IS" BASIS,
  - WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  - See the License for the specific language governing permissions and
  - limitations under the License.
  -->

<script setup>
import {Disclosure, DisclosureButton, DisclosurePanel, Menu, MenuButton, MenuItem, MenuItems} from '@headlessui/vue'
import {Bars3Icon, BellIcon, XMarkIcon} from '@heroicons/vue/24/outline'
import {useRouter} from "vue-router";
import {onMounted, ref} from "vue";
import axios from "axios";

const router = useRouter()

const props = defineProps({
  tab: String
})

const user = {
  name: 'Tom Cook',
  email: 'tom@example.com',
  imageUrl:
      'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
}
// const navigation = [
//   {name: 'Home', href: '/', current: true},
//   // {name: 'Team', href: '#', current: false},
//   // {name: 'Projects', href: '#', current: false},
//   // {name: 'Calendar', href: '#', current: false},
//   {name: 'About', href: '/', current: false},
// ]
const userNavigation = [
  {name: 'Your Profile', href: '#'},
  {name: 'Settings', href: '#'},
  {name: 'Sign out', href: '#'},
]

const routes = router.options.routes
const route = router.currentRoute;

const title = ref('---')

const services = ref([])

const getServices = () => {
  services.value = [{
    "info": {
      "id": 1,
      "version": 256,
      "name": "Messaging",
      "description": "A basic e-mail like messaging service for people to simply communicate",
      "path": "/messaging",
      "port": 0,
      "component": "messaging"
    }
  }, {
    "info": {
      "id": 2,
      "version": 256,
      "name": "Profile",
      "description": "A basic profile service to introduce oneself",
      "path": "/profile",
      "port": 0,
      "component": ""
    }
  }, {
    "info": {
      "id": 3,
      "version": 256,
      "name": "WebServer",
      "description": "A basic www server like service to share your projects, businesses and more",
      "path": "/web",
      "port": 0,
      "component": ""
    }
  }, {
    "info": {
      "id": 4,
      "version": 256,
      "name": "Chat",
      "description": "A basic realtime messaging service for people to chat",
      "path": "/chat",
      "port": 0,
      "component": ""
    }
  }, {
    "info": {
      "id": 5,
      "version": 256,
      "name": "Storage",
      "description": "A basic storage service to distribute data among peers",
      "path": "/storage",
      "port": 0,
      "component": ""
    }
  }, {
    "info": {
      "id": 6,
      "version": 256,
      "name": "Checkers",
      "description": "A basic checker game",
      "path": "/checkers",
      "port": 0,
      "component": ""
    }
  }, {
    "info": {
      "id": 7,
      "version": 256,
      "name": "Chess",
      "description": "A basic chess game",
      "path": "/chess",
      "port": 0,
      "component": ""
    }
  }];
  axios.get('/api/services')
      .then(function (response) {
        services.value = response.data;
        console.log(services.value)
      })
      .catch(function (error) {
        console.log(error);
      })
      .then(function () {
        // always executed
      });
}

onMounted(() => {
  getServices()
})

</script>

<template>
  <!--
    This example requires updating your template:

    ```
    <html class="h-full bg-gray-100">
    <body class="h-full">
    ```
  -->
  <!--  <div>-->
  <!--    {{services}}-->
  <!--  </div>-->
  <div class="min-h-full">
    <Disclosure as="nav" class="bg-gray-800" v-slot="{ open }">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex h-16 items-center justify-between">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <img class="h-12 w-12" src="@/assets/images/z-logo-50.png"/>
            </div>
            <div class="hidden md:block">
              <div class="ml-10 flex items-baseline space-x-4">
                <div v-for="item in routes">
                  <RouterLink v-if="Array('NotFound').indexOf(item.name)<0" :to="item.path" :href="item.path"
                              :class="[item.name == router.currentRoute.value.name ? 'bg-gray-900 border-gray-600 border text-white' : 'text-gray-300 border border-transparent hover:border-white hover:text-white', 'rounded-md px-3 py-2 text-sm font-medium']"
                              :aria-current="item.name == router.currentRoute.value.name ? 'page' : undefined">
                    {{ item.name }}
                  </RouterLink>
                </div>
                <!--                <div v-for="item in services">-->
                <!--                  <RouterLink :to="{name:'Services',params:{name:item.info.name}}"-->
                <!--                            :class="[item.info.name == route.params.name ? 'bg-gray-900 border-gray-600 border text-white' : 'text-gray-300 border border-transparent hover:border-white hover:text-white', 'rounded-md px-3 py-2 text-sm font-medium']"-->
                <!--                            :aria-current="(item.info.name == route.params.name) ? 'page' : undefined">-->
                <!--                    {{ item.info.name }}-->
                <!--                  </RouterLink>-->
                <!--                </div>-->
              </div>
            </div>
          </div>
          <div class="hidden md:block">
            <div class="ml-4 flex items-center md:ml-6">
              <button type="button"
                      class="relative rounded-full bg-gray-800 p-1 text-gray-400 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
                <span class="absolute -inset-1.5"/>
                <span class="sr-only">View notifications</span>
                <BellIcon class="h-6 w-6" aria-hidden="true"/>
              </button>

              <!-- Profile dropdown -->
              <Menu as="div" class="relative ml-3">
                <div>
                  <MenuButton
                      class="relative flex max-w-xs items-center rounded-full bg-gray-800 text-sm focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
                    <span class="absolute -inset-1.5"/>
                    <span class="sr-only">Open user menu</span>
                    <img class="h-8 w-8 rounded-full" :src="user.imageUrl" alt=""/>
                  </MenuButton>
                </div>
                <transition enter-active-class="transition ease-out duration-100"
                            enter-from-class="transform opacity-0 scale-95"
                            enter-to-class="transform opacity-100 scale-100"
                            leave-active-class="transition ease-in duration-75"
                            leave-from-class="transform opacity-100 scale-100"
                            leave-to-class="transform opacity-0 scale-95">
                  <MenuItems
                      class="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                    <MenuItem v-for="item in userNavigation" :key="item.name" v-slot="{ active }">
                      <a :href="item.href"
                         :class="[active ? 'bg-gray-100' : '', 'block px-4 py-2 text-sm text-gray-700']">{{
                          item.name
                        }}</a>
                    </MenuItem>
                  </MenuItems>
                </transition>
              </Menu>
            </div>
          </div>
          <div class="-mr-2 flex md:hidden">
            <!-- Mobile menu button -->
            <DisclosureButton
                class="relative inline-flex items-center justify-center rounded-md bg-gray-800 p-2 text-gray-400 hover:bg-gray-700 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
              <span class="absolute -inset-0.5"/>
              <span class="sr-only">Open main menu</span>
              <Bars3Icon v-if="!open" class="block h-6 w-6" aria-hidden="true"/>
              <XMarkIcon v-else class="block h-6 w-6" aria-hidden="true"/>
            </DisclosureButton>
          </div>
        </div>
      </div>

      <DisclosurePanel class="md:hidden">
        <div class="space-y-1 px-2 pb-3 pt-2 sm:px-3">
          <DisclosureButton v-for="item in navigation" :key="item.name" as="a" :href="item.href"
                            :class="[item.current ? 'bg-gray-900 text-white' : 'text-gray-300 hover:bg-gray-700 hover:text-white', 'block rounded-md px-3 py-2 text-base font-medium']"
                            :aria-current="item.current ? 'page' : undefined">{{ item.name }}
          </DisclosureButton>
        </div>
        <div class="border-t border-gray-700 pb-3 pt-4">
          <div class="flex items-center px-5">
            <div class="flex-shrink-0">
              <img class="h-10 w-10 rounded-full" :src="user.imageUrl" alt=""/>
            </div>
            <div class="ml-3">
              <div class="text-base font-medium leading-none text-white">{{ user.name }}</div>
              <div class="text-sm font-medium leading-none text-gray-400">{{ user.email }}</div>
            </div>
            <button type="button"
                    class="relative ml-auto flex-shrink-0 rounded-full bg-gray-800 p-1 text-gray-400 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
              <span class="absolute -inset-1.5"/>
              <span class="sr-only">View notifications</span>
              <BellIcon class="h-6 w-6" aria-hidden="true"/>
            </button>
          </div>
          <div class="mt-3 space-y-1 px-2">
            <DisclosureButton v-for="item in userNavigation" :key="item.name" as="a" :href="item.href"
                              class="block rounded-md px-3 py-2 text-base font-medium text-gray-400 hover:bg-gray-700 hover:text-white">
              {{ item.name }}
            </DisclosureButton>
          </div>
        </div>
      </DisclosurePanel>
    </Disclosure>

    <header class="bg-white shadow">
      <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <h1 class="text-3xl font-bold tracking-tight text-gray-900">
          {{ router.currentRoute.value.name }}
        </h1>
      </div>
    </header>
    <main>
      <div class="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
        <RouterView/>
      </div>
    </main>
  </div>
</template>

<style scoped>

</style>