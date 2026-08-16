<template>
  <div />
</template>

<script setup>
  import { useRouter } from 'vue-router'

  defineOptions({
    name: 'Reload'
  })

  const router = useRouter()
  // 无历史记录时（如新标签页直接打开），后退是 no-op，兜底回首页
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.replace('/')
  }
</script>
