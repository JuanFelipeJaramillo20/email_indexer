<template>
  <Navbar />
  <SearchBar
    @search-results="handleSearchResults"
    @search-results-next-page="handleNextPageFetch"
    :searchNextPage="searchNextPage"
    :resetObserver="resetObserver"
  />
  <EmailVisualizer @intersected="fetchMoreData" :emails="emailList" :fetchData="fetchMoreData" />
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
    let searchNextPage = ref(false)
    const handleSearchResults = (data) => {
      emailList.value = data.hits
    }

    const resetObserver = () => {
      searchNextPage.value = false
    }

    const handleNextPageFetch = (data) => {
      emailList.value = [...emailList.value, ...data.hits]
    }

    const fetchMoreData = async () => {
      searchNextPage.value = true
    }

    return {
      emailList,
      handleSearchResults,
      fetchMoreData,
      handleNextPageFetch,
      searchNextPage,
      resetObserver
    }
  }
}
</script>
