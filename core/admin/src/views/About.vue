<!--
  - Copyright 2020 PeerGum
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
import axios from "axios";

const info = ref({})

const getInfo = () => {
  info.value = {};
  axios.get('/api/info')
      .then(function (response) {
        info.value = response.data;
        console.log(info.value)
      })
      .catch(function (error) {
        console.log(error);
      })
      .then(function () {
        // always executed
      });
}

onMounted(() => {
  getInfo()
})
</script>


<template>
  <div>
    <div class="about">
      <div class="w-3/12">
        <img src="@/assets/images/z-logo-500.png"/>
      </div>
      <div class="w-9/12 flex flex-col">
        <p class="w-full">Peer-Z is a peer-to-peer network layer first envisioned by software developer and musician PeerGum in 2015
          and originally called Peer-X. Project went inactive quickly after its start, by lack of time, but was
          restarted
          in 2020, with a bunch of new ideas, including the integration of basic services like messaging (mail), chat,
          an
          embedded light web server, a personal profile, and two games: checkers and chess.</p>
        <h3>Team</h3>
        <table class="team">
          <tr>
            <th>Lead Project Manager<br/>
              Principal Developer<br/>
              Lead Designer
            </th>
            <td>PeerGum</td>
          </tr>
          <tr>
            <th>Sponsors</th>
            <td>(pending)</td>
          </tr>
        </table>

        <h3>Build Info</h3>
        <table class="info">
          <tr>
            <th>Name</th>
            <td>{{ info.name }}</td>
          </tr>
          <tr>
            <th>Version</th>
            <td>{{ info.version }}</td>
          </tr>
          <tr>
            <th>Revision</th>
            <td>{{ info.revision }}<br/>{{ info.time }}</td>
          </tr>
          <tr>
            <th>OS</th>
            <td>{{ info.os }}</td>
          </tr>
          <tr>
            <th>Arch</th>
            <td>{{ info.arch }}</td>
          </tr>
        </table>

      </div>
    </div>
  </div>
</template>

<style scoped>
.about {
  @apply flex flex-row items-center gap-4;

  h3 {
    @apply font-bold mt-4;
  }

  div:last-child {
    @apply border-l border-l-gray-200 p-8;
  }

  .team, .info {
    /*@apply border border-gray-200 w-full;*/

    td {
      @apply w-1/2 border border-gray-200 px-4;
    }
  }
}
</style>