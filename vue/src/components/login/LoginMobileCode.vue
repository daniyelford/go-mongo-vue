<template>
    <b-form @submit.prevent="$emit('submit')">
        <b-form-group label="Enter Code" label-for="code-input" class="mb-3">
            <Vue3OtpInput
                v-model="innerCode"
                :num-inputs="6"
                :should-auto-focus="true"
                input-classes="otp-input"
                separator="-"
            />
        </b-form-group>
        <div class="d-flex justify-content-between mb-3">
            <b-button size="sm" variant="secondary" @click="editMobile">Edit Mobile</b-button>
            <b-button
                size="sm"
                variant="outline-primary"
                @click="resendCode"
                :disabled="resendDisabled">
                {{ resendDisabled ? `Resend in ${timer}s` : 'Resend Code' }}
            </b-button>
        </div>
        <b-button type="submit" variant="success" :disabled="props.loading">
            {{ props.loading ? 'Verifying...' : 'Login' }}
        </b-button>
    </b-form>
</template>
<script setup>
    import { ref, watch, onMounted } from 'vue'
    import Vue3OtpInput from "vue3-otp-input"
    const timer = ref(60)
    const resendDisabled = ref(true)
    let interval = null
    const props = defineProps({
        code: String,
        loading: Boolean
    })
    const emits = defineEmits(['update:code', 'submit', 'resend', 'edit'])
    function startTimer() {
        resendDisabled.value = true
        timer.value = 60
        clearInterval(interval)
        interval = setInterval(() => {
            if (timer.value > 0) {
                timer.value--
            } else {
                resendDisabled.value = false
                clearInterval(interval)
            }
        }, 1000)
    }
    function resendCode() {
        emits('resend')
        startTimer()
    }
    function editMobile() {
        emits('edit')
    }
    onMounted(() => {
        startTimer()
        if ("OTPCredential" in window) {
            navigator.credentials.get({
                otp: { transport: ["sms"] },
                signal: new AbortController().signal
            }).then(otp => {
            if (otp && otp.code) {
                innerCode.value = otp.code 
            }
            }).catch(err => {
                console.warn("WebOTP failed:", err)
            })
        }
    })
    const innerCode = ref(props.code)
    watch(innerCode, val => emits('update:code', val))
</script>
<style>
    .otp-input {
        width: 40px;
        height: 40px;
        text-align: center;
        font-size: 20px;
        border: 1px solid #ccc;
        border-radius: 6px;
    }
</style>