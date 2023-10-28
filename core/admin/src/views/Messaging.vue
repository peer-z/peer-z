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
import TrashIcon from "../components/icons/TrashIcon.vue";
import ForwardIcon from "../components/icons/ForwardIcon.vue";
import ReplyIcon from "../components/icons/ReplyIcon.vue";
import ComposeIcon from "../components/icons/ComposeIcon.vue";

const mailboxes = ref([])
const messages = ref([])
const currentBox = ref({})
const currentMessage = ref({})
const box = ref({})

onMounted(() => {
  getMailboxes();
})

const getMailboxes = () => {
  mailboxes.value = [{
    "ID": 1,
    "CreatedAt": "2020-12-17T21:01:14.924528-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924528-05:00",
    "DeletedAt": null,
    "name": "Inbox",
    "Mails": null
  }, {
    "ID": 2,
    "CreatedAt": "2020-12-17T21:01:14.924528-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924528-05:00",
    "DeletedAt": null,
    "name": "Outbox",
    "Mails": null
  }, {
    "ID": 3,
    "CreatedAt": "2020-12-17T21:01:14.924528-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924528-05:00",
    "DeletedAt": null,
    "name": "Sent",
    "Mails": null
  }, {
    "ID": 4,
    "CreatedAt": "2020-12-17T21:01:14.924528-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924528-05:00",
    "DeletedAt": null,
    "name": "Trash",
    "Mails": null
  }];
  currentBox.value = mailboxes.value[0];
  getMessages(currentBox.value)
  return;
  currentBox.value = {};
  currentMessage.value = {};
  axios.get('/mailboxes')
      .then(function (response) {
        mailboxes.value = response.data;
        if (mailboxes.value.length > 0) {
          currentBox.value = mailboxes.value[0]
          getMessages(currentBox.value)
        }
      })
      .catch(function (error) {
        console.log(error);
      })
      .then(function () {
        // always executed
      });
}

const getMessages = (mailbox) => {
  currentBox.value = mailbox;
  currentMessage.value = {};
  messages.value = [{
    "ID": 1,
    "CreatedAt": "2020-12-17T21:01:14.924737-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924737-05:00",
    "DeletedAt": null,
    "dateSent": "0001-01-01T00:00:00Z",
    "dateReceived": "0001-01-01T00:00:00Z",
    "dateRead": "0001-01-01T00:00:00Z",
    "subject": "Welcome to your new mailbox - msg 0",
    "body": "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
    "flags": 0,
    "Mailboxes": null
  }, {
    "ID": 2,
    "CreatedAt": "2020-12-17T21:01:14.924737-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924737-05:00",
    "DeletedAt": null,
    "dateSent": "0001-01-01T00:00:00Z",
    "dateReceived": "0001-01-01T00:00:00Z",
    "dateRead": "0001-01-01T00:00:00Z",
    "subject": "Welcome to your new mailbox - msg 1",
    "body": "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
    "flags": 0,
    "Mailboxes": null
  }, {
    "ID": 3,
    "CreatedAt": "2020-12-17T21:01:14.924737-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924737-05:00",
    "DeletedAt": null,
    "dateSent": "0001-01-01T00:00:00Z",
    "dateReceived": "0001-01-01T00:00:00Z",
    "dateRead": "0001-01-01T00:00:00Z",
    "subject": "Welcome to your new mailbox - msg 2",
    "body": "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
    "flags": 0,
    "Mailboxes": null
  }, {
    "ID": 4,
    "CreatedAt": "2020-12-17T21:01:14.924737-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924737-05:00",
    "DeletedAt": null,
    "dateSent": "0001-01-01T00:00:00Z",
    "dateReceived": "0001-01-01T00:00:00Z",
    "dateRead": "0001-01-01T00:00:00Z",
    "subject": "Welcome to your new mailbox - msg 3",
    "body": "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
    "flags": 0,
    "Mailboxes": null
  }, {
    "ID": 5,
    "CreatedAt": "2020-12-17T21:01:14.924737-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924737-05:00",
    "DeletedAt": null,
    "dateSent": "0001-01-01T00:00:00Z",
    "dateReceived": "0001-01-01T00:00:00Z",
    "dateRead": "0001-01-01T00:00:00Z",
    "subject": "Welcome to your new mailbox - msg 4",
    "body": "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
    "flags": 0,
    "Mailboxes": null
  }, {
    "ID": 6,
    "CreatedAt": "2020-12-17T21:01:14.924737-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924737-05:00",
    "DeletedAt": null,
    "dateSent": "0001-01-01T00:00:00Z",
    "dateReceived": "0001-01-01T00:00:00Z",
    "dateRead": "0001-01-01T00:00:00Z",
    "subject": "Welcome to your new mailbox - msg 5",
    "body": "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
    "flags": 0,
    "Mailboxes": null
  }, {
    "ID": 7,
    "CreatedAt": "2020-12-17T21:01:14.924737-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924737-05:00",
    "DeletedAt": null,
    "dateSent": "0001-01-01T00:00:00Z",
    "dateReceived": "0001-01-01T00:00:00Z",
    "dateRead": "0001-01-01T00:00:00Z",
    "subject": "Welcome to your new mailbox - msg 6",
    "body": "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
    "flags": 0,
    "Mailboxes": null
  }, {
    "ID": 8,
    "CreatedAt": "2020-12-17T21:01:14.924737-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924737-05:00",
    "DeletedAt": null,
    "dateSent": "0001-01-01T00:00:00Z",
    "dateReceived": "0001-01-01T00:00:00Z",
    "dateRead": "0001-01-01T00:00:00Z",
    "subject": "Welcome to your new mailbox - msg 7",
    "body": "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
    "flags": 0,
    "Mailboxes": null
  }, {
    "ID": 9,
    "CreatedAt": "2020-12-17T21:01:14.924737-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924737-05:00",
    "DeletedAt": null,
    "dateSent": "0001-01-01T00:00:00Z",
    "dateReceived": "0001-01-01T00:00:00Z",
    "dateRead": "0001-01-01T00:00:00Z",
    "subject": "Welcome to your new mailbox - msg 8",
    "body": "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
    "flags": 0,
    "Mailboxes": null
  }, {
    "ID": 10,
    "CreatedAt": "2020-12-17T21:01:14.924737-05:00",
    "UpdatedAt": "2020-12-17T21:01:14.924737-05:00",
    "DeletedAt": null,
    "dateSent": "0001-01-01T00:00:00Z",
    "dateReceived": "0001-01-01T00:00:00Z",
    "dateRead": "0001-01-01T00:00:00Z",
    "subject": "Welcome to your new mailbox - msg 9",
    "body": "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
    "flags": 0,
    "Mailboxes": null
  }]
  currentMessage.value = messages.value[0];
  return
  axios.get('/mailboxes/' + mailbox.ID)
      .then(function (response) {
        messages.value = response.data;
        if (messages.value.length > 0) {
          currentMessage.value = messages.value[0];
        }
        console.log(messages.value, currentMessage.value)
      })
      .catch(function (error) {
        console.log(error);
      })
      .then(function () {
        // always executed
      });
}

const showMessage = (cMsg) => {
  console.log(cMsg);
  currentMessage.value = cMsg;
}

</script>

<template>
  <div class="flex flex-row mail-content">
    <div class="w-2/12">

    </div>
    <div class="mailbox-name">
      {{ currentBox.name }}
    </div>
    <div class="w-6/12">
      <ul class="grid grid-cols-6 w-full list-group list-group-horizontal message-toolbox items-center gap-2">
        <li class="flex list-group-item flex-col button">
          <ComposeIcon/>
          New
        </li>
        <li class="list-group-item flex-1 spacer"></li>
        <li class="list-group-item flex-col button">
          <ReplyIcon/>
          Reply
        </li>
        <li class="list-group-item flex-col button">
          <ForwardIcon/>
          Forward
        </li>
        <li class="list-group-item flex-1 spacer"></li>
        <li class="list-group-item flex-col button">
          <TrashIcon/>
          Trash
        </li>
      </ul>
    </div>
  </div>
  <div class="flex flex-row gap-4">
    <div class="mailboxes">
      <div v-for="box in mailboxes" :class="{mailbox: true, selected: box.ID == currentBox.ID}"
           v-on:click="getMessages(box)">
        {{ box.name }}
      </div>
    </div>
    <div class="message-list">
      <div v-for="msg in messages" v-if="currentBox" :class="{selected: msg.ID == currentMessage.ID}"
           class="message-summary" v-on:click="showMessage(msg)">
        <div class="mail-from">{{ msg.from ? msg.from : "Unknown" }}</div>
        <div class="mail-date">{{ msg.dateSent.substring(0, 10) }}<br/>{{ msg.dateSent.substring(11, 19) }}</div>
        <div class="mail-subject">{{ msg.subject }}</div>
        <div class="mail-extract">{{ msg.body.trim().split('.')[0] }} (...)</div>
      </div>
    </div>
    <div v-if="currentMessage.ID" class="message">
      <div class="mail-header">
        <div class="mail-from">From: {{ currentMessage.from }}</div>
        <div class="mail-date">{{
            currentMessage.dateSent.substring(0, 10)
          }} {{ currentMessage.dateSent.substring(11, 19) }}
        </div>
        <div class="mail-to">To: TO</div>
        <div class="mail-subject">{{ currentMessage.subject }}</div>
      </div>
      <div class="mail-body">{{ currentMessage.body }}</div>
    </div>
  </div>

</template>

<style scoped>
.button {
  @apply text-sm flex bg-gray-200 border border-black shadow-md shadow-gray-400 w-full mx-auto h-full place-content-center items-center ;
}

.mail {
  padding: 0;
  margin: 0;
}

.mailbox-name {
  @apply flex w-4/12 items-center place-content-start px-4 text-lg font-bold;
}

.message-list {
  @apply w-4/12 grid grid-cols-1 gap-2 p-2;
}

.mailboxes {
  @apply w-2/12 flex flex-col text-lg gap-3 py-4 items-center;

  div {
    @apply flex shadow-md shadow-gray-500 border border-black w-full justify-center bg-gray-200 rounded-2xl hover:bg-blue-500 hover:text-white;
    margin-top: -2px;
    padding: 0;
    /* border: 1 px solid black;*/

    .mailbox {
      margin: 2px 5px;
      cursor: pointer;
      text-align: center;
      padding: 5px;
      border-radius: 5px;
      border: none;

      &.selected {
        font-weight: bold;
        /* border: 1 px inset grey;*/
        background-color: #e0e0e0;
      }
    }
  }
}

.mails {
  /* margin: 0 10 px;*/
  padding: 0;
  height: 80vh;
}

.box-name {
  height: 50px;
  border-bottom: 1px solid white;
  background-color: #fafafa;
  padding: 10px;
  font-weight: bold;
}

.mail-list {
  height: 600px;
  overflow-y: scroll;

  .mail-summary {
    @apply border-dotted border border-gray-200;

    &:hover {
      background-color: #f7f7f7;
    }

    &.selected {
      background-color: #0040a0;
      color: white;
    }

    cursor: pointer;
    /* margin: 1 px;*/
    padding: 10px;
    margin: 0;
    /* border-top: 1 px solid grey;*/
    border-bottom: 1px dotted lightgrey;
    border-radius: 5px;

    .mail-from {
      text-align: left;
      font-size: 13px;
    }

    .mail-from {
      float: left;
      font-weight: bold;
    }

    .mail-date {
      text-align: right;
      float: right;
      font-size: 13px;
      line-height: 16px;
      color: lightgrey;
    }

    .mail-subject {
      font-weight: bold;
      font-size: 13px;
      /* color: darkblue;*/
      clear: left;
    }

    &.selected {
      /* border: 1 px inset grey;*/
    }
  }
}


.mail-content {
  padding: 0;
  /* height: 100 vh;*/
}

.message-toolbox {
  @apply gap-2 pl-4 h-12;

  li {
    @apply flex flex-row gap-1 bg-gray-200 hover:bg-blue-500 rounded-lg hover:text-white;

    &.button {
      cursor: pointer;
    }

    &.spacer {
      @apply bg-transparent;
    }

    i {
      font-size: 18px;
    }
  }
}

.message-summary {
  @apply flex flex-wrap border border-gray-300 p-2 rounded-md
  bg-gray-50 shadow-sm mt-2 hover:border-black hover:bg-gray-200
  hover:shadow-sm hover:shadow-black;

  &.selected {
    @apply bg-blue-100 border border-black;
  }

  .mail-from {
    @apply w-3/4 font-bold;
  }

  .mail-date {
    @apply w-1/4 justify-end place-content-end border-l border-dotted border-gray-500 pl-2;
  }

  .mail-subject {
    @apply flex justify-center border-y border-gray-400 bg-gray-100 w-full my-2 py-1 font-bold;
  }

  .mail-extract {
    @apply flex justify-between
  }
}

.message {
  @apply flex flex-col w-6/12 mt-4 bg-gray-50 min-h-full border border-gray-300 rounded-md;
  /*
  @apply flex flex-wrap border border-gray-300 p-2 rounded-md
  bg-gray-50 shadow-sm mt-2 hover:border-black hover:bg-gray-200
  hover:shadow-sm hover:shadow-black;
  */

  &.selected {
    @apply bg-blue-100 border border-black;
  }

  .mail-header {
    @apply flex flex-wrap p-2 border border-dotted border-black bg-white m-2 rounded-sm;
  }

  .mail-from, .mail-to {
    @apply flex-grow w-2/3 font-bold;
  }

  .mail-date {
    @apply flex-shrink place-content-end border-l border-b border-dotted border-gray-500 pl-2;
  }

  .mail-subject {
    @apply flex justify-center border border-dotted border-black w-full my-2 py-1 font-bold;
  }

  .mail-body {
    @apply flex justify-between p-2;
  }
}

/*.message {

  .mail-header {

    border-bottom: 1px dotted lightgrey;
    padding-bottom: 10px;

    .mail-from {
      float: left;
      font-weight: bold;
      font-size: 16px;
    }

    .mail-date {
      float: right;
      text-align: right;
      line-height: 16px;
      color: grey;
    }

    .mail-subject {
      clear: left;
      font-size: 14px;
    }

    .mail-to {

    }
  }

  .mail-body {
    padding-top: 10px;
  }
}*/

</style>

