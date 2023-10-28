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
    <nav class="container navbar navbar-expand-lg navbar-dark">
      <a class="navbar-brand" href="#">
        <img height="30" src="/images/z-logo-50.png" width="30"/>
      </a>
      <button aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation" class="navbar-toggler"
              data-target="#navbarSupportedContent" data-toggle="collapse" type="button">
        <span class="navbar-toggler-icon"></span>
      </button>

      <div id="navbarSupportedContent" class="collapse navbar-collapse">
        <ul class="navbar-nav">
          <li :class="{'nav-item':true, active:(service=='home')}">
            <div class="nav-link" href="#" v-on:click="setService('home')">Home</div>
          </li>
          <li :class="{'nav-item':true, active:(service=='messaging')}">
            <div class="nav-link" href="#" v-on:click="setService('messaging')">Messaging</div>
          </li>
          <li :class="{'nav-item':true, active:(service=='profile')}">
            <div class="nav-link" href="#" v-on:click="setService('profile')">Profile</div>
          </li>
          <li :class="{'nav-item':true, active:(service=='web')}">
            <div class="nav-link" href="#" v-on:click="setService('web')">Web</div>
          </li>
          <li :class="{'nav-item':true, active:(service=='checkers')}">
            <div class="nav-link" href="#" v-on:click="setService('checkers')">Checkers</div>
          </li>
          <li :class="{'nav-item':true, active:(service=='chess')}">
            <div class="nav-link" href="#" v-on:click="setService('chess')">Chess</div>
          </li>
          <li :class="{'nav-item':true, active:(service=='about')}">
            <div class="nav-link" href="#" v-on:click="setService('about')">About</div>
          </li>
          <!--          <li class="nav-item dropdown">-->
          <!--            <a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-toggle="dropdown"-->
          <!--               aria-haspopup="true" aria-expanded="false">-->
          <!--              Dropdown-->
          <!--            </a>-->
          <!--            <div class="dropdown-menu" aria-labelledby="navbarDropdown">-->
          <!--              <a class="dropdown-item" href="#">Action</a>-->
          <!--              <a class="dropdown-item" href="#">Another action</a>-->
          <!--              <div class="dropdown-divider"></div>-->
          <!--              <a class="dropdown-item" href="#">Something else here</a>-->
          <!--            </div>-->
          <!--          </li>-->
          <!--          <li class="nav-item">-->
          <!--            <a class="nav-link disabled" href="#" tabindex="-1" aria-disabled="true">Disabled</a>-->
          <!--          </li>-->
        </ul>
        <!--        <form class="form-inline my-2 my-lg-0">-->
        <!--          <input class="form-control mr-sm-2" type="search" placeholder="Search" aria-label="Search">-->
        <!--          <button class="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>-->
        <!--        </form>-->
      </div>
    </nav>
    <div class="container">
      <div v-if="service == 'home'">
        <home :content.sync="service"></home>
      </div>
      <div v-if="service == 'messaging'">
        <messaging></messaging>
      </div>
      <div v-if="service == 'profile'">
        <profile></profile>
      </div>
      <div v-if="service == 'webserver'">
        <web></web>
      </div>
      <div v-if="service == 'checkers'">
        <checkers></checkers>
      </div>
      <div v-if="service == 'chess'">
        <chess></chess>
      </div>
      <div v-if="service == 'about'">
        <about></about>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: [
    'content'
  ],
  data: function () {
    return {
      service: 'home',
    };
  },
  mounted() {
    // this.setService('main');
  },
  watch: {
    content: function (value) {
      this.service = value;
    }
  },
  methods: {
    setService: function (content) {
      console.log(content);
      this.service = content;
    }
  }
}
</script>

<style scoped>

</style>