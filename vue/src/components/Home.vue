<template>
<Logout/>
<RouterLink :to="{path:'/setting'}">setting</RouterLink>
<RouterLink :to="{path:'/business'}">business</RouterLink>
</template>
<script setup>
import Logout from './Logout.vue';
</script>