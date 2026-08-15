<template>
  <div>
    <div class="flex h-screen w-full items-center justify-center bg-gray-50">
      <div class="flex flex-col items-center gap-4 text-2xl">
        <img class="w-1/3" src="../../assets/404.png" />
        <p class="text-lg">页面不存在或当前账号无权访问</p>
        <p class="text-lg">如果确认需要访问此页面，请在角色权限中检查菜单与接口授权配置。</p>
        <el-button @click="toDashboard">返回首页</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: 'Error'
})

const userStore = useUserStore()
const router = useRouter()

// 导航失败属于路由配置问题，不应误判为登录失效强制登出
const toDashboard = () => {
  const defaultRouter = userStore.userInfo?.authority?.defaultRouter
  if (defaultRouter && router.hasRoute(defaultRouter)) {
    router.push({ name: defaultRouter }).catch(() => {})
    return
  }
  if (router.hasRoute('authority')) {
    router.push({ name: 'authority' }).catch(() => {})
    return
  }
  router.push({ path: '/layout' }).catch(() => {})
}
</script>
