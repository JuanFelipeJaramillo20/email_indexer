<template>
  <Navbar />
  <SearchBar @search-results="handleSearchResults" />
  <EmailVisualizer :emails="emailList" :fetchData="fetchMoreData" />
</template>

<script>
import { ref } from 'vue'
import EmailVisualizer from './components/EmailVisualizer.vue'
import SearchBar from './components/SearchBar.vue'
import Navbar from './components/Navbar.vue'

export default {
  components: {
    EmailVisualizer,
    SearchBar,
    Navbar
  },
  setup() {
    const emailList = ref([])
    let page = 1

    const handleSearchResults = (data) => {
      emailList.value = data.hits
      page = 1 // Reset page number after initial search
    }

    const fetchMoreData = async () => {
      // Increment page number for pagination
      page++
      // Perform API request for next page
      //const newData = await fetchEmailData(page)
      // Append new data to existing email list
      //emailList.value = emailList.value.concat(newData.hits)
      console.log('EL USUARIO LLEGO AL FINAL DE LA TABLA')
    }

    const fetchEmailData = async (page) => {
      // Perform API request to fetch email data for the specified page
      // Replace this with your actual API request logic
      const response = await fetch(`http://localhost:8080/search?page=${page}`)
      const data = await response.json()
      return data
    }

    return {
      emailList,
      handleSearchResults,
      fetchMoreData
    }
  }
}
</script>
