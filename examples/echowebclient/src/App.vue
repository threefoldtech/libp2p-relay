<template>
  <header>
    <img alt="Threefold logo"  class="logo" src="@/assets/tf_logo.svg" width="200">
    <HelloWorld msg="libp2p echo client example"/>
  </header>
  <EchoClient/>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import HelloWorld from './components/HelloWorld.vue'
import EchoClient from './components/EchoClient.vue'

export default defineComponent({
  name: 'App',
  components: {
    HelloWorld,
    EchoClient
  }
})
</script>

<style>

.logo {
  display: block;
  margin: 0 auto 2rem;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>
