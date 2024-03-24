<template>
  <div class="observer" ref="observer"></div>
</template>

<script>
export default {
  mounted() {
    this.observer = new IntersectionObserver(([entry]) => {
      if (entry.isIntersecting) {
        this.$emit('intersect')
      }
    })

    this.observer.observe(this.$refs.observer)
  },
  beforeUnmount() {
    this.observer.disconnect()
  }
}
</script>

<style>
.observer {
  position: relative;
  width: 100%;
  height: 1px;
  background-color: transparent;
}
</style>
