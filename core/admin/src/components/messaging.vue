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
      <div class="mailboxes col-2">
        <div class="column-top"></div>
        <div>
          <div v-for="box in mailboxes" :class="{mailbox: true, selected: box.ID == currentBox.ID}"
               v-on:click="getMessages(box)">
            {{ box.name }}
          </div>
        </div>
      </div>
      <div class="mails col-4">
        <div class="box-name align-items-center">
          {{ currentBox.name }}
        </div>
        <div class="mail-list">
          <div v-for="msg in messages" v-if="currentBox" :class="{selected: msg.ID == currentMessage.ID}"
               class="mail-summary" v-on:click="showMessage(msg)">
            <div class="mail-from">From</div>
            <div class="mail-date">{{ msg.dateSent.substring(0, 10) }}<br/>{{ msg.dateSent.substring(11, 19) }}</div>
            <div class="mail-subject">{{ msg.subject }}</div>
            <div class="mail-body">{{ msg.body }}</div>
          </div>
        </div>
      </div>
      <div class="mail-content col-6">
        <div>
          <ul class="list-group list-group-horizontal message-toolbox align-items-center">
            <li class="list-group-item col center button"><i class="fas fa-pen-square"></i> New</li>
            <li class="list-group-item col-1 center spacer"></li>
            <li class="list-group-item col center button"><i class="fas fa-reply"></i> Reply</li>
            <li class="list-group-item col center button"><i class="fas fa-share"></i> Forward</li>
            <li class="list-group-item col-1 center spacer"></li>
            <li class="list-group-item col center button"><i class="fas fa-trash"></i> Trash</li>
          </ul>
        </div>
        <div v-if="currentMessage.ID" class="message">
          <div class="mail-header">
            <div class="mail-from">From</div>
            <div class="mail-date">{{
                currentMessage.dateSent.substring(0, 10)
              }}<br/>{{ currentMessage.dateSent.substring(11, 19) }}
            </div>
            <div class="mail-subject">{{ currentMessage.subject }}</div>
            <div class="mail-to">To: TO</div>
          </div>
          <div class="mail-body">{{ currentMessage.body }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  // name: "messaging",
  props: [],
  data() {
    return {
      mailboxes: [],
      messages: [],
      currentBox: {},
      currentMessage: {},
      box: {}
    }
  },
  watch: {},
  mounted() {
    this.getMailboxes();
  },
  methods: {
    getMailboxes: function () {
      var vm = this;
      this.messages = [];
      this.currentBox = {};
      this.currentMessage = {};
      axios.get('/mailboxes')
          .then(function (response) {
            vm.mailboxes = response.data;
            if (vm.mailboxes.length > 0) {
              vm.currentBox = vm.mailboxes[0]
              vm.getMessages(vm.currentBox)
            }
          })
          .catch(function (error) {
            console.log(error);
          })
          .then(function () {
            // always executed
          });
    },
    getMessages: function (mailbox) {
      var vm = this;
      this.currentBox = mailbox;
      this.currentMessage = {};
      axios.get('/mailboxes/' + mailbox.ID)
          .then(function (response) {
            vm.messages = response.data;
            if (vm.messages.length > 0) {
              vm.currentMessage = vm.messages[0];
            }
          })
          .catch(function (error) {
            console.log(error);
          })
          .then(function () {
            // always executed
          });
    },
    showMessage: function (cMsg) {
      console.log(cMsg);
      this.currentMessage = cMsg;
    },
  }
};
</script>
