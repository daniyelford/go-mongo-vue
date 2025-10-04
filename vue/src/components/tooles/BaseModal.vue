<template>
  <Transition name="fade">
    <div
      v-if="modelValue"
      class="fixed inset-0 flex items-center justify-center bg-black/50 z-50"
      @click.self="close"
    >
      <Transition name="slide-up">
        <div
          class="bg-white rounded-2xl shadow-lg w-full mx-4"
          :class="{
            'max-w-sm': size === 'sm',
            'max-w-md': size === 'md',
            'max-w-lg': size === 'lg',
          }"
        >
          <!-- Header -->
          <div class="flex justify-between items-center border-b px-4 py-3">
            <h3 class="text-lg font-semibold">{{ title }}</h3>
            <button @click="close" class="text-gray-400 hover:text-gray-600">
              ✕
            </button>
          </div>

          <!-- Body -->
          <div class="p-4">
            <slot />
          </div>

          <!-- Footer -->
          <div
            v-if="$slots.footer"
            class="flex justify-end border-t px-4 py-3 gap-2"
          >
            <slot name="footer" />
          </div>
        </div>
      </Transition>
    </div>
  </Transition>
</template>

<script setup>
const props = defineProps({
  modelValue: { type: Boolean, required: true },
  title: { type: String, default: '' },
  size: { type: String, default: 'md' }, // sm, md, lg
})

const emit = defineEmits(['update:modelValue'])

function close() {
  emit('update:modelValue', false)
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.25s ease, opacity 0.25s ease;
}
.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(40px);
  opacity: 0;
}
</style>
