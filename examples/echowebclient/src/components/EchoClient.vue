<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent({
  data () {
    return {
      relayAddress: '',
      daemonID: '',
      response: ''
    }
  },
  methods: {
    sayHello () {
      console.log('Should say hello to', this.daemonID, 'using', this.relayAddress)
    }
  }
})
</script>

<template>
  <main>
    <label>Relay</label>
    <input v-model="relayAddress"/>
    <br/>
    <label>echo daemon ID</label>
    <input v-model="daemonID"/>
    <br/>
    <button @click="sayHello">Say 'Hello'</button>
    <pre>Response: {{ response}}</pre>
  </main>
</template>
<style>
main {
    display: block;
    align-content: flex-start;
}
</style>
