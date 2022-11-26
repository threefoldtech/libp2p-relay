<script lang="ts">
// import { enable } from '@libp2p/logger'
import { defineComponent } from 'vue'
import { createLibp2p } from 'libp2p'
import { multiaddr } from '@multiformats/multiaddr'
import { webSockets } from '@libp2p/websockets'
import { all } from '@libp2p/websockets/filters'
import { noise } from '@chainsafe/libp2p-noise'
import { mplex } from '@libp2p/mplex'
import { yamux } from '@chainsafe/libp2p-yamux'
import { preSharedKey } from 'libp2p/pnet'
import { Buffer } from 'buffer'
// import { pipe } from 'it-pipe'
import { toString as uint8ArrayToString } from 'uint8arrays/to-string'
import { fromString as uint8ArrayFromString } from 'uint8arrays/from-string'

const Protocol = '/echo/1.0.0'

export default defineComponent({
  data () {
    return {
      relayAddress: '/ip4/127.0.0.1/tcp/61503/ws/p2p/12D3KooWPo5j8T2fxEGeUmVDtf2gi3mNtypMTfFqzAQYPz5ii7mw',
      psk: '1ab7e23edf1a951da91cab2d5d77b434936d85fda6bf0fd984e7aed557aab2a0',
      daemonID: '12D3KooWJLNzacc9wcN8xYbAWwAnpisk2ixhjNHz8WQP9iXt3Ru4',
      response: ''
    }
  },
  methods: {
    async sayHello () {
      // enable('libp2p:*')
      console.log('Should say hello to', this.daemonID, 'using', this.relayAddress)
      const swarmKey = Buffer.from('/key/swarm/psk/1.0.0/\n/base16/\n' + this.psk)
      const node = await createLibp2p({
        transports: [webSockets({
          filter: all
        })
        ],
        connectionEncryption: [noise()],
        streamMuxers: [mplex(), yamux()],
        connectionProtector: preSharedKey({
          psk: swarmKey
        }),
        relay: {
          enabled: true
        }
      }
      )
      node.connectionManager.addEventListener('peer:connect', (evt) => {
        console.log('Connected to %s', evt.detail.remotePeer.toString()) // Log connected peer
      })
      node.connectionManager.addEventListener('peer:disconnect', (evt) => {
        console.log('Disconnected from %s', evt.detail.remotePeer.toString()) // Log connected peer
      })
      // start libp2p
      await node.start()
      console.log('libp2p has started')
      const targetMA = multiaddr(this.relayAddress + '/p2p-circuit/p2p/' + this.daemonID)
      try {
        const stream = await node.dialProtocol(targetMA, Protocol)
        console.log('created stream')
        stream.sink([uint8ArrayFromString('Hello from a browser\n')])
        console.log('wrote to stream')
        // For each chunk of data
        for await (const data of stream.source) {
          const reply = uint8ArrayToString(data.subarray())
          console.log('received echo:', reply)
          this.response = reply
        }
      } catch (err) {
        if (err instanceof AggregateError) {
          console.log(err.errors)
        } else {
          console.log(err)
        }
      }
      // stop libp2p
      // console.log('going to stoplibp2p')
      // await node.stop()
      // console.log('libp2p has stopped')
    }
  }
})
</script>

<template>
  <main>
    <label>Relay</label>
    <input v-model="relayAddress"/>
    <br/>
    <label>PSK</label>
    <input v-model="psk"/>
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
