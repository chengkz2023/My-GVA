<template>
  <div class="mx-4 flex items-center gap-4">
    <el-tooltip effect="dark" content="搜索" placement="bottom">
      <span class="flex h-8 w-8 cursor-pointer items-center justify-center rounded-full border border-solid border-gray-200 p-2 shadow dark:border-gray-600">
        <el-icon @click="handleCommand">
          <Search />
        </el-icon>
      </span>
    </el-tooltip>

    <el-tooltip effect="dark" content="系统设置" placement="bottom">
      <span class="flex h-8 w-8 cursor-pointer items-center justify-center rounded-full border border-solid border-gray-200 p-2 shadow dark:border-gray-600">
        <el-icon @click="toggleSetting">
          <Setting />
        </el-icon>
      </span>
    </el-tooltip>

    <el-tooltip effect="dark" content="刷新" placement="bottom">
      <span class="flex h-8 w-8 cursor-pointer items-center justify-center rounded-full border border-solid border-gray-200 p-2 shadow dark:border-gray-600">
        <el-icon :class="showRefreshAnmite ? 'animate-spin' : ''" @click="toggleRefresh">
          <Refresh />
        </el-icon>
      </span>
    </el-tooltip>

    <el-tooltip effect="dark" content="切换主题" placement="bottom">
      <span class="flex h-8 w-8 cursor-pointer items-center justify-center rounded-full border border-solid border-gray-200 p-2 shadow dark:border-gray-600">
        <el-icon v-if="appStore.isDark" @click="appStore.toggleTheme(false)">
          <Sunny />
        </el-icon>
        <el-icon v-else @click="appStore.toggleTheme(true)">
          <Moon />
        </el-icon>
      </span>
    </el-tooltip>

    <gva-setting v-model:drawer="showSettingDrawer" />
    <command-menu ref="command" />
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import { emitter } from '@/utils/bus.js'
import { useAppStore } from '@/pinia'
import GvaSetting from '@/view/layout/setting/index.vue'
import CommandMenu from '@/components/commandMenu/index.vue'

const appStore = useAppStore()
const showSettingDrawer = ref(false)
const showRefreshAnmite = ref(false)
const command = ref()

const toggleRefresh = () => {
  showRefreshAnmite.value = true
  emitter.emit('reload')
  setTimeout(() => {
    showRefreshAnmite.value = false
  }, 1000)
}

const toggleSetting = () => {
  showSettingDrawer.value = true
}

const handleCommand = () => {
  command.value.open()
}

const handleKeyDown = (e) => {
  // 输入框聚焦时忽略 Ctrl+K，避免与打字冲突
  const target = e.target
  if (target && (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable)) {
    return
  }
  if (e.ctrlKey && e.key === 'k') {
    e.preventDefault()
    handleCommand()
  }
}

window.addEventListener('keydown', handleKeyDown)
onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>
