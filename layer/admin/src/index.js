/*
 * Copyright 2020 PeerGum
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// import Vue from 'vue';
// import VueI18n from "vue-i18n";
var axios = require('axios');
// import './logo.png';
import './style.scss';
// import VueLoader from "vue-loader";
import './images/logo.png';
import './images/black-pawn.png';
import './images/white-pawn.png';
import './images/black-king.png';
import './images/white-king.png';
import './images/b-pawn.png';
import './images/w-pawn.png';
import './images/b-rook.png';
import './images/w-rook.png';
import './images/b-knight.png';
import './images/w-knight.png';
import './images/b-bishop.png';
import './images/w-bishop.png';
import './images/b-queen.png';
import './images/w-queen.png';
import './images/b-king.png';
import './images/w-king.png';
import './images/z-logo-50.png';
import './images/z-logo-500.png';

// Install BootstrapVue
// Vue.use(BootstrapVue)
// Optionally install the BootstrapVue icon components plugin
// Vue.use(IconsPlugin)
// Vue.use(VueI18n)

// const app = Vue.createApp({
//     el: '#app',
//     data: {
//         message: 'Hello Vue!'
//     }
// })

try {
    // window.Popper = require('popper.js').default;
    window.$ = window.jQuery = require('jquery');
    require('./jquery-ui-1.12.1.custom/jquery-ui.js')

    require('bootstrap');
} catch (e) {
}

window.axios = require('axios');
window.axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';

// var BLACK = 0, WHITE = 1;

const BLANK = 0, BLACK = 1, WHITE = 2;

// Vue.component('nav', require('./components/nav.vue').default);
// Vue.component('messaging', require('./components/messaging.vue').default);

Vue.component('menubar', require('./components/menubar.vue').default);
Vue.component('home', require('./components/home.vue').default);
Vue.component('about', require('./components/about.vue').default);
Vue.component('messaging', require('./components/messaging.vue').default);
Vue.component('profile', require('./components/profile.vue').default);
Vue.component('checkers', require('./components/checkers.vue').default);
Vue.component('chess', require('./components/chess.vue').default);

var app = new Vue({
    el: "#app",
})


