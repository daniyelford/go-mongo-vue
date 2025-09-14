<template>
    {{ test }}
</template>
<script setup>
    import { sendApi } from '@/plugins/api';
    import { ref , onMounted } from 'vue';
    const test = ref('')
    onMounted(async () => {
        test.value = await sendApi({ method: 'GET', url: '/hello' })
    })
</script>