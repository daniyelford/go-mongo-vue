<template>
    <b-form @submit.prevent="$emit('submit')">
        <div class="d-flex align-items-center mb-3">
            <b-form-group label="Country" label-for="country-select" class="flex-grow-4">
                <v-select
                id="country-select"
                :searchable="true"
                :options="props.countriesCode"
                :reduce="opt => opt.value"
                label="text"
                v-model="innerCountry"
                placeholder="Select code"
                />
            </b-form-group>
            <b-form-group label="Mobile" label-for="mobile-input" class="flex-grow-6">
                <b-form-input
                id="mobile-input"
                v-model="innerMobile"
                type="text"
                placeholder="Enter mobile number"
                />
            </b-form-group>
        </div>
        <b-button type="submit" variant="primary" :disabled="props.loading || props.disableSend">
            {{ props.loading
                ? 'Sending...'
                : (props.disableSend
                    ? `Send Code (${props.timerEnd || 0}s)`
                    : 'Send Code'
                )
            }}
        </b-button>
    </b-form>
</template>
<script setup>
    import { ref, watch } from 'vue'
    import vSelect from 'vue-select'
    import 'vue-select/dist/vue-select.css'
    const props = defineProps({
        countriesCode: Array,
        country: String,
        mobile: String,
        loading: Boolean,
        disableSend: Boolean,
        timerEnd: Number
    })
    const emits = defineEmits(['update:country', 'update:mobile', 'submit'])
    const innerCountry = ref(props.country)
    const innerMobile = ref(props.mobile)
    watch(innerCountry, val => emits('update:country', val))
    watch(innerMobile, val => emits('update:mobile', val))
</script>
<style>
    .flex-grow-4 { flex: 4; }
    .flex-grow-6 { flex: 6; }
    .vs__selected-options { padding: 0 !important; flex-wrap: nowrap !important; }
    input.vs__search { padding: 0 !important; width: 0 !important; max-width: 0 !important; }
</style>