<template>
  <div v-if="has">
    <button @click="login">
      finger
    </button>
    <p v-if="message" class="mt-2">{{ message }}</p>
  </div>
</template>
<script setup>
import { onMounted, ref, watch } from 'vue';
import { sendApi } from '@/plugins/api';
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
        data: props.mobile
    });
    options.publicKey.challenge = Uint8Array.from(atob(options.publicKey.challenge), c => c.charCodeAt(0));
    options.publicKey.allowCredentials = options.publicKey.allowCredentials.map(cred => ({
      ...cred,
      id: Uint8Array.from(atob(cred.id), c => c.charCodeAt(0))
    }));
    const credential = await navigator.credentials.get(options);
    const res = await sendApi({
        method: "POST",
        url: "/login/fingerPrint/end",
        data: { token:credential.response,mobile:props.mobile}
    });
    const data = await res.json();
    if (data.accessToken) {
      message.value = "success";
    } else {
      message.value = "error";
    }
  } catch (e) {
    console.error(e);
    message.value = "error";
  }
}
watch(props.mobile,async(mobile)=>{
  try {
    const res = await sendApi({
      url:'/login/fingerPrint/has',
      method:'post',
      data:mobile
    })
    if(res.status) has.value=true
    console.log(res);
  } catch (error) {
    console.log(error);
    
  }
})
</script>