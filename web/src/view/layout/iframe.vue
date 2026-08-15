<template>
  <div class="h-screen w-screen bg-gray-50 text-slate-700 dark:bg-slate-800 dark:text-slate-500">
    <iframe
      id="gva-base-load-dom"
      class="gva-body-h w-full border-t border-gray-200 bg-gray-50 dark:border-slate-700 dark:bg-slate-800"
      :src="url"
    ></iframe>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import useResponsive from '@/hooks/responsive'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: 'GvaLayoutIframe'
})

useResponsive(true)

const userStore = useUserStore()
const route = useRoute()

// 只允许 http(s) 绝对地址，防止菜单数据注入 javascript:/data: 协议在 iframe 内执行脚本
function safeIframeUrl(raw) {
  const val = String(raw || '')
  if (/^https?:\/\//i.test(val)) {
    return val
  }
  return 'about:blank'
}

const url = safeIframeUrl(route.query.url)

onMounted(() => {
  if (userStore.loadingInstance) {
    userStore.loadingInstance.close()
  }
})
</script>
