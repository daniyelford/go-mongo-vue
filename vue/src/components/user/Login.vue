<template>
  <b-container class="p-4" style="max-width: 400px;">
    <h2 class="mb-4">Login</h2>
    <LoginMobileInput
      v-if="step === 1"
      :countries-code="countriesCode"
      :country="country"
      :mobile="mobile"
      :loading="loading"
      :disable-send="disableSend"
      :timer-end="sendTimerEnd"
      @update:country="country = $event"
      @update:mobile="mobile = $event"
      @submit="sendMobile"
    />
    <LoginMobileCode
      v-else
      :code="code"
      :loading="loading"
      :resend-timer-end="resendTimerEnd"
      :disable-resend="resendDisabled"
      :mobile="country + ' ' + mobile"
      @update:code="code = $event"
      @submit="verifyCode"
      @resend="resendCode"
      @edit="editMobile"
    />
    <FingerLogin v-if="step === 1" :mobile="country+mobile" />
    <b-alert v-if="loading || disableSend" >
      The login button is locked for logging in too much.
    </b-alert>
    <div v-if="error" class="mt-4 text-danger">❌ {{ error }}</div>
  </b-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { sendApi } from '@/plugins/api'
import { countriesCode } from '@/composables/countries'
import FingerLogin from '@/components/user/finger/FingerLogin.vue'
import LoginMobileInput from '@/components/user/login/LoginMobileInput.vue'
import LoginMobileCode from '@/components/user/login/LoginMobileCode.vue'
import router from '@/router'
const step = ref(1)
const country = ref('+98')
const mobile = ref('')
const code = ref('')
const error = ref('')
const loading = ref(false)
const disableSend = ref(false)
const resendDisabled = ref(false)
const sendTimerEnd = ref(null)
const resendTimerEnd = ref(null)
let sendTimerId = null
let resendTimerId = null

function startSendTimer(durationMs = 15000) {
  clearInterval(sendTimerId)
  let remaining = Math.ceil(durationMs / 1000)
  disableSend.value = true
  localStorage.setItem('sendTimerEnd', Date.now() + durationMs)
  sendTimerId = setInterval(() => {
    remaining--
    sendTimerEnd.value = remaining
    localStorage.setItem('sendTimerEnd', Date.now() + remaining * 1000)
    if (remaining <= 0) {
      clearInterval(sendTimerId)
      sendTimerId = null
      disableSend.value = false
      localStorage.removeItem('sendTimerEnd')
    }
  }, 1000)
}

function startResendTimer(durationMs = 60000) {
  // 5 min for max 3 sends
  clearInterval(resendTimerId)
  let remaining = Math.ceil(durationMs / 1000)
  resendDisabled.value = true
  localStorage.setItem('resendTimerEnd', Date.now() + durationMs)
  resendTimerId = setInterval(() => {
    remaining--
    resendTimerEnd.value = remaining
    localStorage.setItem('resendTimerEnd', Date.now() + remaining * 1000)
    if (remaining <= 0) {
      clearInterval(resendTimerId)
      resendTimerId = null
      resendDisabled.value = false
      localStorage.removeItem('resendTimerEnd')
    }
  }, 1000)
}

async function sendMobile() {
  if (disableSend.value) return
  error.value = ''
  if (!mobile.value) {
    error.value = 'incorrect mobile'
    return
  }
  loading.value = true
  const res = await sendApi({
    method: 'POST',
    url: '/login/send',
    data: { country: country.value, mobile: mobile.value }
  })
  loading.value = false
  if (res.error) {
    error.value = res.message
    if (res.message.includes('please wait')) startSendTimer(15000)
    if (res.message.includes('too many requests')) startResendTimer(300000)
  } else {
    step.value = 2
    localStorage.setItem('loginStep', '2')
    localStorage.setItem('loginMobile', mobile.value)
    startSendTimer()
    startResendTimer()
  }
}

async function resendCode() {
  if (resendDisabled.value) return
  if (!mobile.value) {
    step.value = 1
    localStorage.setItem('loginStep', '1')
    error.value = 'incorrect mobile'
    code.value = ''
    return
  }
  error.value = ''
  loading.value = true
  const res = await sendApi({
    method: 'POST',
    url: '/login/send',
    data: { country: country.value, mobile: mobile.value }
  })
  loading.value = false
  if (res.error) {
    error.value = res.message
    if (res.message.includes('too many requests')){
        startSendTimer(150000)
        startResendTimer(150000)
    } 
  } else{
    startResendTimer()
  } 
}

async function verifyCode() {
  if (code.value?.length != 6) return
  error.value = ''
  loading.value = true
  const res = await sendApi({
    method: 'POST',
    url: '/login/verify',
    data: { country: country.value, mobile: mobile.value, code: code.value }
  })
  loading.value = false
  if (res.error) {
    error.value = res.message
    if (res.message.includes('too many wrong attempts')) {
      step.value = 1
      code.value = ''
      startSendTimer(120000)
      localStorage.setItem('loginStep', '1')
    }
  } else {
    localStorage.setItem('loginStep', '1')
    localStorage.setItem('jwt', res.accessToken)
    localStorage.setItem('refresh', res.refreshToken)
    if (res.newUser) router.push({ path: '/register' })
    else router.push({ path: '/home' })
  }
}

function editMobile() {
  step.value = 1
  code.value = ''
  localStorage.setItem('loginStep', '1')
  startSendTimer()
}

onMounted(() => {
  const savedStep = localStorage.getItem('loginStep')
  const savedMobile = localStorage.getItem('loginMobile')
  const savedSendEnd = localStorage.getItem('sendTimerEnd')
  const savedResendEnd = localStorage.getItem('resendTimerEnd')
  if (savedMobile) {
    mobile.value = savedMobile
    step.value = savedStep ? parseInt(savedStep) : 1
  }
  if (savedSendEnd) {
    const remaining = Math.max(parseInt(savedSendEnd) - Date.now(), 0)
    if (remaining > 0) startSendTimer(remaining)
  }
  if (savedResendEnd) {
    const remainingResend = Math.max(parseInt(savedResendEnd) - Date.now(), 0)
    if (remainingResend > 0) startResendTimer(remainingResend)
  }
})
</script>
