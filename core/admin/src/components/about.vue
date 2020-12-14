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
    <h1>About</h1>
    <div class="row about">
      <div class="col-sm-3">
        <img src="/images/z-logo-500.png"/>
      </div>
      <div class="col-sm-9">
        <p>Peer-Z is a peer-to-peer network layer first envisioned by software developer and musician PeerGum in 2015
          and originally called Peer-X. Project went inactive quickly after its start, by lack of time, but was
          restarted
          in 2020, with a bunch of new ideas, including the integration of basic services like messaging (mail), chat,
          an
          embedded light web server, a personal profile, and two games: checkers and chess.</p>
        <h3>Team</h3>
        <table class="team">
          <tr>
            <td>Lead Project Manager<br/>
              Principal Developer<br/>
              Lead Designer
            </td>
            <td>PeerGum</td>
          </tr>
          <tr>
            <td>Sponsors</td>
            <td>TBD</td>
          </tr>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "about"
}
</script>

<style scoped>

</style>