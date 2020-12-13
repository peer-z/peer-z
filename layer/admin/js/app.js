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

// // import 'bootstrap/dist/css/bootstrap.css'
// // import 'bootstrap-vue/dist/bootstrap-vue.css'
//
// import '../sass/app.scss'

// import 'bootstrap/dist/css/bootstrap.css'
// import 'bootstrap-vue/dist/bootstrap-vue.css'
//
// import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'

// Install BootstrapVue
// Vue.use(BootstrapVue)
// Optionally install the BootstrapVue icon components plugin
// Vue.use(IconsPlugin)

let en = require('./lang/en.json');
let fr = require('./lang/fr.json');

Vue.use(VueI18n)

const messages = {
    en,
    fr
}

Vue.component('messaging', require('../src/components/messaging.vue').default);

// Create VueI18n instance with options
const i18n = new VueI18n({
    locale: 'en', // set locale
    messages, // set locale messages
})

const app = new Vue({
    i18n,
    el: '#app',
    data: function () {
        return {
            menuTitle: "",
            menuSubtitle: ""
        };
    }
});
