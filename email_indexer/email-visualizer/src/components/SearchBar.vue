<template>
  <form class="flex items-center mx-auto p-10">
    <label for="simple-search" class="sr-only">Search</label>
    <div class="relative w-full">
      <div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
        <svg
          class="w-4 h-4 text-green-500 dark:text-gray-400"
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 18 20"
        >
          <path
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M3 5v10M3 5a2 2 0 1 0 0-4 2 2 0 0 0 0 4Zm0 10a2 2 0 1 0 0 4 2 2 0 0 0 0-4Zm12 0a2 2 0 1 0 0 4 2 2 0 0 0 0-4Zm0 0V6a3 3 0 0 0-3-3H9m1.5-2-2 2 2 2"
          />
        </svg>
      </div>
      <input
        v-model="searchTerm"
        type="text"
        id="simple-search"
        class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full ps-10 p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
        placeholder="Search emails..."
        required
      />
    </div>
    <button
      type="button"
      class="p-2.5 ms-2 text-sm font-medium text-white bg-green-500 rounded-lg border border-green-700 hover:bg-green-800 focus:ring-4 focus:outline-none focus:ring-green-300 dark:bg-green-600 dark:hover:bg-green-700 dark:focus:ring-green-800"
      @click="searchEmails"
    >
      <svg
        class="w-4 h-4"
        aria-hidden="true"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 20 20"
      >
        <path
          stroke="currentColor"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"
        />
      </svg>
      <span class="sr-only">Search</span>
    </button>
  </form>
</template>
<script>
import { ref } from 'vue'

export default {
  setup(_, ctx) {
    const searchTerm = ref('')
    const page = ref(1)
    const searchEmails = async () => {
      const url = new URL('http://localhost:8080/search')
      const params = { term: searchTerm.value, page: page.value }
      url.search = new URLSearchParams(params).toString()
      const options = {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' }
      }
      try {
        const response = await fetch(url, options)
        const data = await response.json()
        ctx.emit('searchResults', data.hits)
      } catch (error) {
        console.error('Error:', error)
      }
    }

    return {
      searchTerm,
      searchEmails
    }
  }
}
</script>
