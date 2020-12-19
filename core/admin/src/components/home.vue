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

<template>
  <div>
    <div class="row">
      <div class="col">

        <h2>Me</h2>
        <table class="peers col-sm-12">
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

        <h2>Services</h2>
        <table class="services col-sm-12" v-if="services.length>0">
          <tr>
            <th class="center">#</th>
            <th>Name</th>
            <th>Description</th>
            <th>Version</th>
          </tr>
          <tr v-for="service in services">
            <td class="counter center">{{ service.info.id }}</td>
            <td class="">
              <div class="btn-link" v-on:click="$emit('update:content',service.info.name.toLowerCase())">{{ service.info.name }}</div>
            </td>
            <td class="">{{ service.info.description }}</td>
            <td class="">{{ service.info.version / 256 }}.{{ service.info.version % 256 }}</td>
          </tr>
        </table>
      </div>
      <div class="col">
        <h2>Peers</h2>
        <table class="peers col-sm-12" v-if="peers.length>0">
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
            <td class="counter center">{{ index+1 }}</td>
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
    </div>
  </div>
</template>

<script>
export default {
  props: [],
  data() {
    return {
      peers: [],
      services: [],
      me: { Name: "" },
    };
  },
  mounted: function() {
    this.getMe();
    this.getServices();
    this.getPeers();
    console.log("OK");
  },
  methods: {
    getPeers: function () {
      var vm = this;
      this.peers = [];
      axios.get('/api/peers')
          .then(function (response) {
            console.log(response.data);
            if (response.data != null) {
              vm.peers = response.data;
            }
          })
          .catch(function (error) {
            console.log(error);
          })
          .then(function () {
            // always executed
          });
    },
    getServices: function () {
      var vm = this;
      this.services = [];
      axios.get('/api/services')
          .then(function (response) {
            vm.services = response.data;
          })
          .catch(function (error) {
            console.log(error);
          })
          .then(function () {
            // always executed
          });
    },
    getMe: function () {
      var vm = this;
      this.me = {};
      axios.get('/api/me')
          .then(function (response) {
            vm.me = response.data;
          })
          .catch(function (error) {
            console.log(error);
          })
          .then(function () {
            // always executed
          });
    },
  }
};
</script>
