<template>
    <b-form @submit.prevent="$emit('submit')">
        <b-form-group label="Enter Code" label-for="code-input" class="mb-3">
            <OtpInput
                v-model="innerCode"
                :length="6"
                :reset="resetNumber"
                :timer="15000"
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
    import OtpInput from "vue-otp-autofill";
    const resetNumber = ref(0)
    const timer = ref(60)
    const resendDisabled = ref(true)
    let interval = null
    const props = defineProps({
        code: String,
        loading: Boolean
    })
    const emits = defineEmits(['update:code', 'submit', 'resend', 'edit'])
    const innerCode = ref(props.code)
    watch(innerCode, val => emits('update:code', val))
    onMounted(() => startTimer())
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
        resetNumber.value++
        startTimer()
    }
    function editMobile() {
        emits('edit')
    }
</script>
