<template>
  <div class="p-4">
    <button @click="register" class="bg-blue-600 text-white px-4 py-2 rounded">
      finger
    </button>
    <p v-if="message" class="mt-2">{{ message }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { sendApi } from '@/plugins/api';
const message = ref("");

async function register() {
  try {
    const options = await sendApi({
        method: "POST",
        autoCheckToken: true,
        url: "/register/fingerPrint/start",
    });
    options.publicKey.challenge = Uint8Array.from(atob(options.publicKey.challenge), c => c.charCodeAt(0));
    options.publicKey.user.id = Uint8Array.from(atob(options.publicKey.user.id), c => c.charCodeAt(0));
    const credential = await navigator.credentials.create(options);
    const res = await sendApi({
        method: "POST",
        autoCheckToken: true,
        url: "/register/fingerPrint/end",
        data: credential.response,
    });
    message.value = res.message || "success";
  } catch (e) {
    console.error(e);
    message.value = "error";
  }
}
</script>
