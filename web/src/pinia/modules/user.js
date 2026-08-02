import { login, getUserInfo } from '@/api/user'
import { jsonInBlacklist } from '@/api/jwt'
import router from '@/router/index'
import { ElLoading, ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRouterStore } from './router'
import { useCookies } from '@vueuse/integrations/useCookies'
import { useStorage } from '@vueuse/core'
import { useAppStore } from '@/pinia'
import { config } from '@/core/config.js'

export const useUserStore = defineStore('user', () => {
  const appStore = useAppStore()
  const loadingInstance = ref(null)

  const userInfo = ref({
    uuid: '',
    nickName: '',
    headerImg: '',
    authority: {}
  })
  const token = useStorage('token', '')
  const xToken = useCookies('x-token')
  const currentToken = computed(() => token.value || xToken.value || '')

  const setUserInfo = (val) => {
    userInfo.value = val
    if (val.originSetting) {
      Object.keys(appStore.config).forEach((key) => {
        if (val.originSetting[key] !== undefined) {
          appStore.config[key] = val.originSetting[key]
        }
      })
    }
  }

  const setToken = (val) => {
    token.value = val
    xToken.value = val
  }

  const NeedInit = async () => {
    await ClearStorage()
    await router.push({ name: 'Init', replace: true })
  }

  const ResetUserInfo = (value = {}) => {
    userInfo.value = {
      ...userInfo.value,
      ...value
    }
  }
  /* 获取用户信息 */
  const GetUserInfo = async () => {
    if (config.useMockLogin) {
      setUserInfo({
        ...userInfo.value,
        nickName: userInfo.value.nickName || 'Scaffold User',
        authority: { defaultRouter: 'authority' }
      })
      return { code: 0 }
    }
    const res = await getUserInfo()
    if (res.code === 0) {
      setUserInfo(res.data.userInfo)
    }
    return res
  }

  /* 登录 */
  const LoginIn = async (loginInfo) => {
    try {
      loadingInstance.value = ElLoading.service({
        fullscreen: true,
        text: '登录中，请稍候...'
      })

      if (config.useMockLogin) {
        setUserInfo({
          uuid: 'scaffold',
          nickName: loginInfo.username || 'Scaffold User',
          headerImg: '',
          authority: { defaultRouter: 'authority' }
        })
        setToken('mock-token-' + Date.now())
        const routerStore = useRouterStore()
        await routerStore.SetAsyncRouter()
        routerStore.asyncRouters.forEach((r) => router.addRoute(r))
        const redirect = router.currentRoute.value.query.redirect
        if (redirect) {
          await router.replace(redirect)
        } else {
          await router.replace({ name: 'authority' })
        }
        return true
      }

      const res = await login(loginInfo)
      if (res.code !== 0) {
        return false
      }
      setUserInfo(res.data.user)
      setToken(res.data.token)
      await GetUserInfo()

      const routerStore = useRouterStore()
      await routerStore.SetAsyncRouter()
      routerStore.asyncRouters.forEach((r) => router.addRoute(r))

      if (router.currentRoute.value.query.redirect) {
        await router.replace(router.currentRoute.value.query.redirect)
        return true
      }

      const defaultRouter = userInfo.value.authority?.defaultRouter
      if (defaultRouter && router.hasRoute(defaultRouter)) {
        await router.replace({ name: defaultRouter })
      } else {
        await router.replace({ name: 'authority' })
      }

      const isWindows = /windows/i.test(navigator.userAgent)
      window.localStorage.setItem('osType', isWindows ? 'WIN' : 'MAC')
      return true
    } catch (error) {
      console.error('LoginIn error:', error)
      return false
    } finally {
      loadingInstance.value?.close()
    }
  }
  /* 登出 */
  const LoginOut = async () => {
    if (!config.useMockLogin) {
      const res = await jsonInBlacklist()
      if (res.code !== 0) return
    }
    await ClearStorage()
    router.push({ name: 'Login', replace: true })
    window.location.reload()
  }
  /* 清理数据 */
  const ClearStorage = async () => {
    token.value = ''
    // 使用remove方法正确删除cookie
    xToken.remove()
    sessionStorage.clear()
    // 清理所有相关的localStorage项
    localStorage.removeItem('originSetting')
    localStorage.removeItem('token')
  }

  return {
    userInfo,
    token: currentToken,
    NeedInit,
    ResetUserInfo,
    GetUserInfo,
    LoginIn,
    LoginOut,
    setToken,
    loadingInstance,
    ClearStorage
  }
})
