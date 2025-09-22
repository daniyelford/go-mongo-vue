<template>
    <a @click="logout">logout</a>
</template>
<script setup>
import { sendApi } from '@/plugins/api';
import router from '@/router';
    const logout = async () => {
        localStorage.removeItem('loginStep')
        localStorage.removeItem('sendTimerEnd')
        localStorage.removeItem('resendTimerEnd')
        const token = localStorage.getItem('jwt')
        if (token) {
            try {
                const res = await sendApi({
                    method: "GET",
                    url: "/auth/logout",
                    headers: { Authorization: `Bearer ${token}` }
                })
                if (!res.error && res.Logout) {
                    router.push({path:'/login'})
                }else{
                    console.error('error1',res.error)
                }
            } catch (err) {
                console.error('error',err)
            }finally{
                localStorage.removeItem('jwt')
                localStorage.removeItem('refresh')
            }
        }
        router.push({path:'/login'})
    }
</script>