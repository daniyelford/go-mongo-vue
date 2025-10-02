<template>
  <div>
    <button v-if="has" @click="login">
      finger
    </button>
    <p v-if="message" class="mt-2">{{ message }}</p>
  </div>
</template>
<script setup>
import { ref, watchEffect, defineProps } from 'vue';
import router from '@/router'
import { sendApi } from '@/plugins/api';
import { bufferToBase64Url,base64UrlToUint8ArrayReg } from '@/plugins/base64';
const message = ref("");
const has = ref(false);
const props= defineProps({
    mobile: String,
})
async function login() {
  try {
    const options = await sendApi({
        method: "POST",
        url: "/login/fingerPrint/start",
        data: { mobile: props.mobile}
    });
    options.publicKey.challenge = base64UrlToUint8ArrayReg(options.publicKey.challenge);
    options.publicKey.allowCredentials = options.publicKey.allowCredentials.map(cred => ({
      ...cred,
      id: base64UrlToUint8ArrayReg(cred.id)
    }));
    const credential = await navigator.credentials.get(options);
    const token = {
      id: credential.id,
      rawId: bufferToBase64Url(credential.rawId),
      type: credential.type,
      response: {
        authenticatorData: bufferToBase64Url(credential.response.authenticatorData),
        clientDataJSON: bufferToBase64Url(credential.response.clientDataJSON),
        signature: bufferToBase64Url(credential.response.signature),
        userHandle: credential.response.userHandle ? bufferToBase64Url(credential.response.userHandle) : null
      }
    };
    const res = await sendApi({
        method: "POST",
        url: "/login/fingerPrint/end",
        headers:{
          "Content-Type": "application/json",
          "X-Mobile":props.mobile
        },
        data: token
    });
    if (res.accessToken) {
      localStorage.setItem('jwt', res.accessToken)
      localStorage.setItem('refresh', res.refreshToken)
      router.push({ path: '/home' })
    } else {
      message.value = "error";
    }
  } catch (e) {
    console.error(e);
    message.value = "error";
  }
}
async function checkFinger() {
  if (!props.mobile) return;
  try {
    const res = await sendApi({
      url: '/login/fingerPrint/has',
      method: 'POST',
      data: { mobile: props.mobile }
    });
    has.value = !!res.status;
  } catch (error) {
    console.log(error);
  }
}
watchEffect(() => {
  checkFinger();
});
</script>