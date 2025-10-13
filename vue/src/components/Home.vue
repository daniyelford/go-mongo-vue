<template>
    <div class="container py-4">
    <b-navbar type="dark" variant="primary" class="rounded-3 mb-4">
      <b-navbar-brand href="#">brand</b-navbar-brand>
      <b-navbar-nav class="ms-auto">
        <RouterLink class="nav-link text-white" :to="{ path: '/setting' }">setting</RouterLink>
        <Logout class="nav-link text-white"/>
      </b-navbar-nav>
    </b-navbar>
    <AllPosts />
    </div>
</template>
<script setup>
  import Logout from '@/components/Logout.vue';
  import AllPosts from '@/components/post/AllPosts.vue';
</script>