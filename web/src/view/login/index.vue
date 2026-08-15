<template>
  <div id="userLayout" class="relative h-full w-full">
    <div
      class="flex h-full w-full items-center justify-evenly rounded-lg bg-white md:h-screen md:w-screen md:bg-[#194bfb]"
    >
      <div class="flex h-full w-10/12 items-center justify-evenly md:w-3/5">
        <div
          class="oblique absolute -ml-52 h-[130%] w-3/5 -rotate-12 transform bg-white dark:bg-slate-900"
        />
        <div
          class="z-[999] box-border flex w-full flex-col justify-between rounded-lg pt-12 pb-10 md:w-96"
        >
          <div>
            <div class="flex items-center justify-center">
              <Logo :size="6" />
            </div>
            <div class="mb-9">
              <p class="text-center text-4xl font-bold">
                {{ $GIN_VUE_ADMIN.appName }}
              </p>
              <p class="mt-2.5 text-center text-sm font-normal text-gray-500">
                管理后台
              </p>
            </div>
            <el-form
              ref="loginForm"
              :model="loginFormData"
              :rules="rules"
              :validate-on-rule-change="false"
              @keyup.enter="submitForm"
            >
              <el-form-item prop="username" class="mb-6">
                <el-input
                  v-model="loginFormData.username"
                  size="large"
                  placeholder="请输入用户名"
                  suffix-icon="user"
                />
              </el-form-item>
              <el-form-item prop="password" class="mb-6">
                <el-input
                  v-model="loginFormData.password"
                  show-password
                  size="large"
                  type="password"
                  placeholder="请输入密码"
                />
              </el-form-item>
              <el-form-item
                v-if="!useMockLogin && loginFormData.openCaptcha"
                prop="captcha"
                class="mb-6"
              >
                <div class="flex w-full justify-between">
                  <el-input
                    v-model="loginFormData.captcha"
                    placeholder="请输入验证码"
                    size="large"
                    class="mr-5 flex-1"
                  />
                  <div class="h-11 w-1/3 rounded bg-[#c3d4f2]">
                    <img
                      v-if="picPath"
                      class="h-full w-full"
                      :src="picPath"
                      alt="验证码"
                      @click="loginVerify"
                    />
                  </div>
                </div>
              </el-form-item>
              <el-form-item class="mb-4">
                <el-button
                  class="h-11 w-full shadow shadow-active"
                  type="primary"
                  size="large"
                  :loading="loading"
                  @click="submitForm"
                >
                  登录
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </div>
      <div class="hidden h-full w-1/2 float-right bg-[#194bfb] md:block">
        <img
          class="h-full"
          src="@/assets/login_right_banner.jpg"
          alt="banner"
        />
      </div>
    </div>

    <BottomInfo class="absolute right-0 bottom-3 left-0 z-20 mx-auto w-full" />
  </div>
</template>

<script setup>
  import { reactive, ref, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { captcha } from '@/api/user'
  import BottomInfo from '@/components/bottomInfo/bottomInfo.vue'
  import { useUserStore } from '@/pinia/modules/user'
  import Logo from '@/components/logo/index.vue'
  import { config } from '@/core/config.js'

  defineOptions({ name: 'Login' })

  const useMockLogin = config.useMockLogin
  const captchaRequiredLength = ref(6)
  const loginForm = ref(null)
  const picPath = ref('')
  const loading = ref(false)

  const loginFormData = reactive({
    username: '',
    password: '',
    captcha: '',
    captchaId: '',
    openCaptcha: false
  })

  const checkUsername = (rule, value, callback) => {
    if (!value || value.length < 2) {
      return callback(new Error('请输入用户名'))
    }
    callback()
  }

  const checkPassword = (rule, value, callback) => {
    if (!useMockLogin && (!value || value.length < 6)) {
      return callback(new Error('密码至少 6 位'))
    }
    callback()
  }

  const checkCaptcha = (rule, value, callback) => {
    if (!loginFormData.openCaptcha) return callback()
    const sanitizedValue = (value || '').replace(/\s+/g, '')
    if (!sanitizedValue) return callback(new Error('请输入验证码'))
    if (sanitizedValue.length < captchaRequiredLength.value) {
      return callback(new Error(`验证码至少 ${captchaRequiredLength.value} 位`))
    }
    callback()
  }

  const rules = reactive({
    username: [{ validator: checkUsername, trigger: 'blur' }],
    password: [{ validator: checkPassword, trigger: 'blur' }],
    captcha: [{ validator: checkCaptcha, trigger: 'blur' }]
  })

  const userStore = useUserStore()

  const loginVerify = async () => {
    try {
      const result = await captcha()
      captchaRequiredLength.value = Number(result.data?.captchaLength) || 6
      picPath.value = result.data?.picPath || ''
      loginFormData.captchaId = result.data?.captchaId || ''
      loginFormData.openCaptcha = !!result.data?.openCaptcha
    } catch {
      loginFormData.openCaptcha = false
    }
  }

  const submitForm = () => {
    if (loading.value) return
    loginForm.value.validate(async (valid) => {
      if (!valid) {
        ElMessage.error('请正确填写登录信息')
        return
      }
      loading.value = true
      try {
        const ok = await userStore.LoginIn({
          username: loginFormData.username,
          password: loginFormData.password,
          captcha: loginFormData.captcha,
          captchaId: loginFormData.captchaId
        })
        if (!ok && !useMockLogin) {
          await loginVerify()
        }
      } finally {
        loading.value = false
      }
    })
  }

  onMounted(async () => {
    if (!useMockLogin) {
      await loginVerify()
    }
  })
</script>
