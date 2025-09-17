<template>
    <b-container class="p-4" style="max-width: 400px;">
        <h2 class="mb-4">Login</h2>
        <LoginMobileInput
        v-if="step === 1"
        :countries-code="countriesCode"
        :country="country"
        :mobile="mobile"
        :loading="loading"
        :timerEnd="timerEnd"
        :disableSend="disableSend"
        @update:country="country = $event"
        @update:mobile="mobile = $event"
        @submit="sendMobile"
        />
        <LoginMobileCode
        v-else
        :code="code"
        :loading="loading"
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
    let timerId = null
    let timerEnd = null
    async function resendCode() {
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
        } else {
            code.value = ''
        }
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
        } else {
            step.value = 2
            localStorage.setItem('loginStep', '2')
        }
    }
    async function verifyCode() {
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
        } else {
            token.value = res.token
            localStorage.setItem('jwt', res.token)
        }
    }
    function startDisableTimer() {
        disableSend.value = true
        const duration = 15000
        timerEnd = Date.now() + duration
        localStorage.setItem('sendCodeEnd', timerEnd)
        if (timerId) clearTimeout(timerId)
        timerId = setTimeout(() => {
            disableSend.value = false
            timerId = null
            localStorage.removeItem('sendCodeEnd')
        }, duration)
    }
    function editMobile() {
        step.value = 1
        code.value = ''
        startDisableTimer()
    }
    onMounted(() => {
        const savedStep = localStorage.getItem('loginStep')
        const savedEnd = localStorage.getItem('sendCodeEnd')
        if (savedStep) step.value = parseInt(savedStep)
        if (savedEnd) {
            const endTime = parseInt(savedEnd)
            const remaining = endTime - Date.now()
            if (remaining > 0) {
                disableSend.value = true
                timerId = setTimeout(() => {
                    disableSend.value = false
                    timerId = null
                    localStorage.removeItem('sendCodeEnd')
                }, remaining)
            }
        }
    })
</script>
