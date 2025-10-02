<template>
    <b-form @submit.prevent="$emit('submit')">
        {{ props.mobile }}
        <b-button size="sm" variant="secondary" @click="editMobile">Edit Mobile</b-button>
        <b-form-group label="Enter Code" label-for="code-input" class="mb-3">
            <Vue3OtpInput
                ref="otpInput"
                v-model:value="code"
                :num-inputs="6"
                @on-change="onCodeChange"
                :should-auto-focus="true"
                input-classes="otp-input"
                separator="-"
            />
        </b-form-group>
        <div class="d-flex justify-content-between mb-3">
            <b-button size="sm" variant="outline-primary" @click="resendCode" :disabled="props.disableResend||props.loading">
                {{ props.loading
                    ? 'Resending...'
                    : (props.disableResend
                        ? `Resend Code (${props.resendTimerEnd || 0}s)`
                        : 'Resend Code'
                    )
                }}
            </b-button>
        </div>
        <b-button type="submit" variant="success" :disabled="props.loading || !complate" :class="{ 'btn-disabled-custom': props.loading || !complate }">
            {{ props.loading ? 'Verifying...' : 'Login' }}
        </b-button>
    </b-form>
</template>
<script setup>
    import { ref, watch, onMounted } from 'vue'
    import Vue3OtpInput from "vue3-otp-input"
    const props = defineProps({
        code: String,
        loading: Boolean,
        mobile: String,
        resendTimerEnd: Number,
        disableResend: Boolean
    })
    const complate=ref(false)
    const otpInput = ref(null)
    const code = ref('')
    const emits = defineEmits(['update:code', 'submit', 'resend', 'edit'])
    function resendCode() { emits('resend') }
    function editMobile() { emits('edit') }
    function onCodeChange(val) {
        emits('update:code', Array.isArray(val) ? val.join('') : val)
        if(val.length==6) complate.value = true
        else complate.value = false
    }
    watch(() => props.code, val => code.value = val || '')
    onMounted(() => {
        if ("OTPCredential" in window) {
            navigator.credentials.get({ otp: { transport: ["sms"] }, signal: new AbortController().signal })
            .then(otp => { if (otp && otp.code) innerCode.value = otp.code })
            .catch(err => console.warn("WebOTP failed:", err))
        }
    })
</script>
<style>
    .otp-input { width: 40px; height: 40px; text-align: center; font-size: 20px; border: 1px solid #ccc; border-radius: 6px; }
    .btn-disabled-custom{ cursor: not-allowed; pointer-events: none; opacity: 0.6; }
</style>