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
import { base64UrlToUint8Array } from '@/plugins/base64';
const message = ref("");
async function register() {
  try {
    const options = await sendApi({
        method: "POST",
        autoCheckToken: true,
        url: "/register/fingerPrint/start",
    });
    options.publicKey.challenge = base64UrlToUint8Array(options.publicKey.challenge)
    options.publicKey.user.id = base64UrlToUint8Array(options.publicKey.user.id)
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
