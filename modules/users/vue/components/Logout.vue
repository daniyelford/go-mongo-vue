<template>
    <a @click="logout">logout</a>
</template>
<script setup>
import { sendApi,removeTokensOut } from '@/plugins/api';
    const logout = async () => {
        localStorage.removeItem('loginStep')
        localStorage.removeItem('sendTimerEnd')
        localStorage.removeItem('resendTimerEnd')
        try {
            const res = await sendApi({
                method: "GET",
                url: "/auth/logout",
                autoCheckToken: true
            })
            if (!(res.error && !res.Logout)) {
                console.error('error1',res.error)
            }
        } catch (err) {
            console.error('error',err)
        }finally{
            removeTokensOut()
        }
    }
</script>