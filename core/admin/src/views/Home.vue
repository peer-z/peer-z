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
import {onMounted, ref} from "vue";
import axios from 'axios'

const peers = ref([]);
const services = ref([]);
const me = ref({Name: ""});

onMounted(() => {
  getMe();
  getServices();
  getPeers();
  console.log("OK");
})

function getPeers() {
  peers.value = [];
  axios.get('/api/peers')
      .then(function (response) {
        console.log(response.data);
        if (response.data != null) {
          peers.value = response.data;
        }
      })
      .catch(function (error) {
        console.log(error);
      })
      .then(function () {
        // always executed
      });
}

function getServices() {
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
  // services.value = [{
  //   info: {
  //     name: "Messaging",
  //     component: "Messaging",
  //   },
  // }];
  axios.get('/api/services')
      .then(function (response) {
        services.value = response.data;
      })
      .catch(function (error) {
        console.log(error);
      })
      .then(function () {
        // always executed
      });
}

function getMe() {
  me.value = {};
  axios.get('/api/me')
      .then(function (response) {
        me.value = response.data;
      })
      .catch(function (error) {
        console.log(error);
      })
      .then(function () {
        // always executed
      });
}
</script>

<template>
  <div class="w-full grid grid-cols-1 lg:grid-cols-2 gap-4">
    <div class="w-full p-2">
      <h2>Me</h2>
      <table class="peers w-full">
        <tr>
          <th class="center">Name</th>
          <th>Address</th>
          <th>IP</th>
          <th>Port</th>
          <th>Peers</th>
        </tr>
        <tr v-if="me.name">
          <td class="center">{{ me.name }}</td>
          <td class="">{{ me.address }}</td>
          <td class="">{{ me.ip }}</td>
          <td class="">{{ me.port }}</td>
          <td class="">{{ me.peers }}</td>
        </tr>
      </table>
    </div>

    <div class="lg:hidden p-2">
      <h2>Services</h2>
      <table v-if="services.length>0" class="services w-full">
        <tr>
          <th class="center">#</th>
          <th>Name</th>
          <th>Description</th>
          <th>Version</th>
        </tr>
        <tr v-for="service in services">
          <td class="counter center">{{ service.info.id }}</td>
          <td class="">
            <div class="btn-link" v-on:click="$emit('update:content',service.info.name.toLowerCase())">
              {{ service.info.name }}x
            </div>
          </td>
          <td class="">{{ service.info.description }}</td>
          <td class="">{{ service.info.version / 256 }}.{{ service.info.version % 256 }}</td>
        </tr>
      </table>
    </div>
    <div class="p-2 w-full">
      <h2>Peers</h2>
      <table v-if="peers.length>0" class="peers w-full">
        <tr>
          <th class="center">#</th>
          <th>Name</th>
          <th>Address</th>
          <th>IP</th>
          <th>Local IP</th>
          <th>Port</th>
          <th>Distance</th>
          <th>Peers</th>
        </tr>
        <tr v-for="(peer,index) in peers">
          <td class="counter center">{{ index + 1 }}</td>
          <td class="">{{ peer.name }}</td>
          <td class="">{{ peer.address }}</td>
          <td class="">{{ peer.ip }}</td>
          <td class="">{{ peer.localIp }}</td>
          <td class="">{{ peer.port }}</td>
          <td class="center">{{ peer.distance }}</td>
          <td class="center">{{ peer.peers }}</td>
        </tr>
      </table>
    </div>
    <div class="max-md:hidden services p-2 w-full">
      <h2>Services</h2>
      <table v-if="services.length>0" class="services w-full">
        <tr>
          <th class="center">#</th>
          <th>Name</th>
          <th>Description</th>
          <th>Version</th>
        </tr>
        <tr v-for="service in services">
          <td class="counter center">{{ service.info.id }}</td>
          <td class="">
            <RouterLink class="btn-link" :to="service.info.path">{{ service.info.name }}</RouterLink>
            <!--            <button v-on:click="$emit('update:content',service.info.name.toLowerCase())">-->
            <!--              {{ service.info.name }}-->
            <!--            </button>-->
          </td>
          <td class="">{{ service.info.description }}</td>
          <td class="">{{ service.info.version / 256 }}.{{ service.info.version % 256 }}</td>
        </tr>
      </table>
    </div>

  </div>
</template>

<style scoped>
/*
table {
  &.peers, &.services {
    @apply border border-gray-700;
    & td {
      @apply justify-self-center;
    }
  }
}
*/
.btn-link {
  @apply text-blue-500 hover:text-blue-900 hover:border-b-2 mx-2 border-black;
}

</style>
