<template>
  <b-modal
    v-model="innerValue"
    :title="title"
    :size="size"
    centered
    hide-footer
    no-close-on-backdrop
    no-close-on-esc
    :body-class="'p-4'"
    :dialog-class="'rounded-2xl shadow-lg'"
    :header-class="'border-b'"
    :header-bg-variant="'white'"
    :body-bg-variant="'white'"
    :footer-bg-variant="'white'"
    :content-class="'rounded-2xl overflow-hidden'"
  >
    <slot />
    <template #footer>
      <div v-if="$slots.footer" class="d-flex justify-content-end gap-2 border-top pt-3">
        <slot name="footer" />
      </div>
    </template>
  </b-modal>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: { type: Boolean, required: true },
  title: { type: String, default: '' },
  size: { type: String, default: 'md' },
})

const emit = defineEmits(['update:modelValue'])

const innerValue = ref(props.modelValue)

watch(() => props.modelValue, val => (innerValue.value = val))
watch(innerValue, val => emit('update:modelValue', val))
</script>

<style scoped>
.modal.fade .modal-dialog {
  transition: transform 0.25s ease-out, opacity 0.25s ease-out;
  transform: translateY(30px);
  opacity: 0;
}
.modal.show .modal-dialog {
  transform: translateY(0);
  opacity: 1;
}
</style>
