<template>
<Logout/>
<RouterLink :to="{path:'/setting'}">setting</RouterLink>
<RouterLink :to="{path:'/business'}">business</RouterLink>
</template>
<script setup>
import Logout from '@users/components/Logout.vue';
</script>