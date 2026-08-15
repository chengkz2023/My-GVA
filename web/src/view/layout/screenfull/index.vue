<template>
  <div @click="clickFull">
    <div v-if="isShow" class="gvaIcon gvaIcon-fullscreen-expand" />
    <div v-else class="gvaIcon gvaIcon-fullscreen-shrink" />
  </div>
</template>

<script setup>
  import screenfull from 'screenfull' // 引入screenfull
  import { onMounted, onUnmounted, ref } from 'vue'

  defineOptions({
    name: 'Screenfull'
  })

  onMounted(() => {
    if (screenfull.isEnabled) {
      screenfull.on('change', changeFullShow)
    }
  })

  onUnmounted(() => {
    if (screenfull.isEnabled) {
      screenfull.off('change', changeFullShow)
    }
  })

  const clickFull = () => {
    if (screenfull.isEnabled) {
      screenfull.toggle()
    }
  }

  const isShow = ref(true)
  const changeFullShow = () => {
    isShow.value = !screenfull.isFullscreen
  }
</script>

<style scoped lang="scss">
  .screenfull-svg {
    width: 16px;
    height: 16px;
    cursor: pointer;
    vertical-align: middle;
    margin-right: 32px;
    fill: rgba(0, 0, 0, 0.45);
  }
</style>
