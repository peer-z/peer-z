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
import {defineAsyncComponent, onMounted, onUpdated, ref, watch, watchEffect} from "vue";
import axios from "axios";

const Messaging = defineAsyncComponent(() => import("./Messaging.vue"))
const Chess = defineAsyncComponent(() => import("./Chess.vue"))
const props = defineProps({
  name: String,
})

const emit = defineEmits([
  'update:service',
])

const Service = ref(null)
const service = ref(props.name)

const serviceDetails = () => {
  axios.get('/api/services/' + props.name)
      .then(function (response) {
        service.value = response.data;
        console.log(service.value)
      })
      .catch(function (error) {
        console.log(error);
      })
      .then(function () {
        // always executed
      });
}

const loadComponent = (name) => {
  switch(name) {
    case 'Messaging':
      return Messaging
    case 'Chess':
      return Chess
    default:
      break;
  }
}

onUpdated(() => {
  console.log("Mounted");

  Service.value = loadComponent(service.value)

  return {Service}
})

const updateComponent = () => {
  loadComponent(service.value)
}

emit('update:service', props.name)

</script>

<template>
  <div class="w-full items-center justify-center">
    <!--    {{ props.name }}-->
    {{ service }}
    <component :is="Service"/>
  </div>
</template>

<style scoped>
table {
  &.peers, &.services {
    @apply border border-gray-700;
  }
}

</style>
