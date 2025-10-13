<template>
    <b-card class="shadow-sm rounded-4">
        <div class="d-flex justify-content-between align-items-center mb-2">
            <b-card-title class="m-0">📄 Posts</b-card-title>
            <b-button variant="success" size="sm" @click="toggleAddModal">➕ Add Post</b-button>
        </div>
        <b-list-group flush v-if="posts.length>0">
        <b-list-group-item
          v-for="post in posts"
          :key="post._id"
          class="d-flex justify-content-between align-items-center">
          <div v-if="post.media?.length" class="mt-2">
            <swiper :slides-per-view="1" navigation pagination>
                <swiper-slide v-for="(item, index) in post.media" :key="index">
                    <img v-if="item.type.startsWith('image/')" :src="item.url" class="img-fluid"/>
                    <video v-else-if="item.type.startsWith('video/')" :src="item.url" controls class="w-100"/>
                </swiper-slide>
            </swiper>
        </div>
          <div>
            <h5 class="mb-1">{{ post.title }}</h5>
            <small class="text-muted">{{ post.content }}</small>
          </div>
          <div v-if="post.self">
            <b-button size="sm" variant="outline-warning" @click="editPost(post)">edit</b-button>
            <b-button size="sm" variant="outline-danger" class="ms-2" @click="deletePost(post._id)">delete</b-button>
          </div>
        </b-list-group-item>
        </b-list-group>
        <b-alert v-else show variant="danger">
            You have no posts yet.
        </b-alert>
        <div class="text-center mt-3">
          <b-button
            variant="primary"
            @click="loadMore"
            :disabled="loading || finished"
          >
            {{ loading ? 'loading...' : finished ? 'end' : 'more' }}
          </b-button>
        </div>
    </b-card>
    <BaseModal v-model="show" size="lg" title="add posts">
        <AddPost @add="toggleAddModal"/>
    </BaseModal>
    <EditPost
    v-model:show="showEditModal"
    :post-data="editingPost"
    @updated="refreshPosts"/>
</template>
<script setup>
    import { ref, onMounted, onUnmounted  } from 'vue'
    import { sendApi } from '@/plugins/api';
    import BaseModal from '@/components/tooles/BaseModal.vue';
    import AddPost from '@/components/post/AddPost.vue';
    import EditPost from '@/components/post/EditPost.vue';
    import { Swiper, SwiperSlide } from 'swiper/vue';
    import 'swiper/swiper-bundle.css';
    const show = ref(false)
    const posts = ref([])
    const page = ref(1)
    const limit = 10
    const loading = ref(false)
    const finished = ref(false)
    const editingPost = ref(null)
    const showEditModal = ref(false)
    let refreshInterval = null
    const toggleAddModal = () => {
        show.value = !show.value
        refreshPosts()
    }
    const mapPosts = (rawPosts) => {
        return rawPosts.map(p => ({
            _id: p.ID,
            userID: p.UserID,
            title: p.Title,
            content: p.Content,
            media: p.media || [],
            self: p.self
        }))
    }
    const loadMore = () => loadPosts()
    const deletePost = async (id) => {
        if (!confirm('آیا مطمئنی می‌خواهی این پست را حذف کنی؟')) return
        try {
            await sendApi({
                url:'/posts/delete',
                data: { id: id },
                method:'DELETE',
                autoCheckToken:true
            })
            posts.value = posts.value.filter(p => p._id !== id)
        } catch (e) {
            console.error('خطا در حذف پست:', e)
        }
    }
    const editPost = (post) => {
        editingPost.value = { ...post }
        showEditModal.value = true
    }
    const loadPosts = async () => {
        if (loading.value || finished.value) return
        loading.value = true
        try {
            const res = await sendApi({
                url:'/posts/all',
                method:'POST',
                data:{
                    page: page.value,
                    limit
                },
                autoCheckToken:true
            })
            const newPosts = mapPosts(res.data || [])
            if (newPosts.length < limit) finished.value = true
            const existingIds = new Set(posts.value.map(p => p._id))
            newPosts.forEach(p => {
                if (!existingIds.has(p._id)) posts.value.push(p)
            })
            page.value++
        } catch (e) {
            console.error('خطا در دریافت پست‌ها:', e)
        } finally {
            loading.value = false
        }
    }
    const refreshPosts = async () => {
        try {
            const res = await sendApi({
                url:'/posts/all',
                method:'POST',
                data:{
                    page: 1,
                    limit: page.value * limit
                },
                autoCheckToken:true
            })            
            const freshPosts = mapPosts(res.data || [])
            const existingIds = new Set(posts.value.map(p => p._id))
            freshPosts.forEach(p => {
                if (!existingIds.has(p._id)) posts.value.unshift(p)
            })
        } catch (e) {
            console.error('خطا در رفرش پست‌ها:', e)
        }
    }
    onMounted(() => {
        loadPosts()
        refreshInterval = setInterval(refreshPosts, 15000)
    })
    onUnmounted(() => {
        if (refreshInterval){
            clearInterval(refreshInterval)
            refreshInterval = null
        } 
    })
</script>