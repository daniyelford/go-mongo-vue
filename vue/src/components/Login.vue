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
        <div v-if="error" class="mt-4 text-danger">❌ {{ error }}</div>
        <div v-if="token" class="mt-4 text-success">
            ✅ Token received: <code class="d-block bg-light p-2 rounded">{{ token }}</code>
        </div>
    </b-container>
</template>
<script setup>
    import { ref, onMounted } from 'vue'
    import { sendApi } from '@/plugins/api'
    import { countriesCode } from '@/composables/countries'
    import LoginMobileInput from '@/components/login/LoginMobileInput.vue'
    import LoginMobileCode from '@/components/login/LoginMobileCode.vue'
    const step = ref(1)
    const country = ref('+98')
    const mobile = ref('')
    const code = ref('')
    const token = ref('')
    const error = ref('')
    const loading = ref(false)
    const disableSend = ref(false)
    const resendDisabled = ref(false)
    const sendTimerEnd = ref(null)
    const resendTimerEnd = ref(null)
    let sendTimerId = null
    let resendTimerId = null
    onMounted(() => {
        const savedStep = localStorage.getItem('loginStep')
        const savedMobile = localStorage.getItem('loginMobile')
        const savedSendEnd = localStorage.getItem('sendTimerEnd')
        const savedResendEnd = localStorage.getItem('resendTimerEnd')
        if (savedMobile) {
            mobile.value = savedMobile
            step.value = savedStep ? parseInt(savedStep) : 1
        } else {
            mobile.value = ''
            step.value = 1
            localStorage.removeItem('loginStep')
            localStorage.removeItem('sendTimerEnd')
            localStorage.removeItem('resendTimerEnd')
        }
        if (savedSendEnd){
            const remaining = Math.max(parseInt(savedSendEnd) - Date.now(), 0)
            if(remaining > 0) startSendTimer(remaining)
        }
        if (savedResendEnd){
            const remainingResend = Math.max(parseInt(savedResendEnd) - Date.now(), 0)
            if(remainingResend > 0) startResendTimer(remainingResend)
        }
    })
    function startSendTimer(durationMs = 15000) {
        clearInterval(sendTimerId)
        let remaining = Math.ceil(durationMs / 1000)
        disableSend.value = true
        localStorage.setItem('sendTimerEnd', Date.now() + durationMs)
        sendTimerId = setInterval(() => {
            remaining--
            localStorage.setItem('sendTimerEnd', Date.now() + remaining * 1000)
            sendTimerEnd.value=remaining
            if (remaining <= 0) {
                clearInterval(sendTimerId)
                sendTimerId = null
                disableSend.value = false
                localStorage.removeItem('sendTimerEnd')
            }
        }, 1000)
    }
    function startResendTimer(durationMs = 60000) {
        clearInterval(resendTimerId)
        let remaining = Math.ceil(durationMs / 1000)
        resendDisabled.value = true
        localStorage.setItem('resendTimerEnd', Date.now() + durationMs)
        resendTimerId = setInterval(() => {
            remaining--
            localStorage.setItem('resendTimerEnd', Date.now() + remaining * 1000)
            resendTimerEnd.value = remaining
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
        if (res.error) error.value = res.message
        else {
            step.value = 2
            localStorage.setItem('loginStep', '2')
            localStorage.setItem('loginMobile', mobile.value)
            startSendTimer()
            startResendTimer()
        }
    }
    async function verifyCode() {
        if(code.value?.length != 6) return
        error.value = ''
        loading.value = true
        const res = await sendApi({
            method: 'POST',
            url: '/login/verify',
            data: { country: country.value, mobile: mobile.value, code: code.value }
        })
        loading.value = false
        if (res.error) error.value = res.message
        else {
            token.value = res.token
            localStorage.setItem('jwt', res.token)
        }
    }
    async function resendCode() {
        if(resendDisabled.value) return
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
        if (res.error) error.value = res.message
        else startResendTimer()
    }
    function editMobile() {
        step.value = 1
        code.value = ''
        localStorage.setItem('loginStep', '1')
        startSendTimer()
    }
</script>
