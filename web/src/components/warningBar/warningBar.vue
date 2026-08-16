<template>
  <div
    class="px-1.5 py-2 flex items-center rounded-sm mt-2 bg-amber-50 gap-2 mb-3 text-amber-500 dark:bg-amber-700 dark:text-gray-200"
    :class="href && 'cursor-pointer'"
    @click="open"
  >
    <el-icon class="text-xl">
      <warning-filled />
    </el-icon>
    <span>
      {{ title }}
    </span>
  </div>
</template>
<script setup>
  import { WarningFilled } from '@element-plus/icons-vue'
  const prop = defineProps({
    title: {
      type: String,
      default: ''
    },
    href: {
      type: String,
      default: ''
    }
  })

  const safeHref = () => {
    const value = prop.href
    if (!value) {
      return ''
    }
    try {
      const url = new URL(value, window.location.origin)
      if (url.protocol === 'http:' || url.protocol === 'https:') {
        return url.href
      }
    } catch {
      return ''
    }
    return ''
  }

  const open = () => {
    const target = safeHref()
    if (target) {
      window.open(target, '_blank', 'noopener,noreferrer')
    }
  }
</script>
