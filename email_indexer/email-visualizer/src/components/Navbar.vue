<template>
  <div class="bg-gray-900 text-gray-100 py-3.5 px-6 shadow md:flex justify-between items-center">
    <div class="flex items-center cursor-pointer">
      <span class="text-green-500 text-x1 mr-1">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" height="1em">
          <path
            fill="#ffffff"
            d="M502.3 190.8c3.9-3.1 9.7-.2 9.7 4.7V400c0 26.5-21.5 48-48 48H48c-26.5 0-48-21.5-48-48V195.6c0-5 5.7-7.8 9.7-4.7 22.4 17.4 52.1 39.5 154.1 113.6 21.1 15.4 56.7 47.8 92.2 47.6 35.7 .3 72-32.8 92.3-47.6 102-74.1 131.6-96.3 154-113.7zM256 320c23.2 .4 56.6-29.2 73.4-41.4 132.7-96.3 142.8-104.7 173.4-128.7 5.8-4.5 9.2-11.5 9.2-18.9v-19c0-26.5-21.5-48-48-48H48C21.5 64 0 85.5 0 112v19c0 7.4 3.4 14.3 9.2 18.9 30.6 23.9 40.7 32.4 173.4 128.7 16.8 12.2 50.2 41.8 73.4 41.4z"
          />
        </svg>
      </span>

      <h1 class="text-xl">Gojimail</h1>
    </div>

    <span @click="openMenu()" class="absolute md:hidden right-6 top-1 cursor-pointer text-4xl">
      <svg v-if="open" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 352 512" height="1em">
        <path
          fill="#ffffff"
          d="M242.7 256l100.1-100.1c12.3-12.3 12.3-32.2 0-44.5l-22.2-22.2c-12.3-12.3-32.2-12.3-44.5 0L176 189.3 75.9 89.2c-12.3-12.3-32.2-12.3-44.5 0L9.2 111.5c-12.3 12.3-12.3 32.2 0 44.5L109.3 256 9.2 356.1c-12.3 12.3-12.3 32.2 0 44.5l22.2 22.2c12.3 12.3 32.2 12.3 44.5 0L176 322.7l100.1 100.1c12.3 12.3 32.2 12.3 44.5 0l22.2-22.2c12.3-12.3 12.3-32.2 0-44.5L242.7 256z"
        />
      </svg>

      <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512" height="1em">
        <path
          fill="#ffffff"
          d="M16 132h416c8.8 0 16-7.2 16-16V76c0-8.8-7.2-16-16-16H16C7.2 60 0 67.2 0 76v40c0 8.8 7.2 16 16 16zm0 160h416c8.8 0 16-7.2 16-16v-40c0-8.8-7.2-16-16-16H16c-8.8 0-16 7.2-16 16v40c0 8.8 7.2 16 16 16zm0 160h416c8.8 0 16-7.2 16-16v-40c0-8.8-7.2-16-16-16H16c-8.8 0-16 7.2-16 16v40c0 8.8 7.2 16 16 16z"
        />
      </svg>
    </span>

    <ul
      class="md:flex md:items-center md:px-0 px-10 md:pb-0 pb-10 md:static absolute bg-gray-900 md:w-auto w-full top-14 duration-700 ease-in"
      :class="[open ? 'left-o' : 'left-[-100%]']"
    >
      <li class="md:mx-4 md:my-0 my-6" v-for="link in links">
        <a :href="link.link" class="text-xl hover:text-green-500">
          {{ link.name }}
        </a>
      </li>

      <Button>Get started</Button>
    </ul>
  </div>
</template>

<script>
import { ref } from 'vue'
import Button from './Button.vue'
export default {
  components: {
    Button
  },
  setup() {
    let open = ref(false)
    const links = [
      { name: 'Inicio', link: '#' },
      { name: 'Sobre el proyecto', link: 'about' }
    ]
    function openMenu() {
      open.value = !open.value
    }
    return { links, open, openMenu }
  }
}
</script>
