<template>
    <b-card class="shadow-sm rounded-4">
        <div class="d-flex justify-content-between align-items-center mb-2">
            <b-card-title class="m-0">📄 Posts</b-card-title>
            <b-button variant="success" size="sm" @click="addPost">➕ Add Post</b-button>
        </div>
        <b-row v-if="posts.length > 0" class="g-3">
            <b-col
                v-for="post in posts"
                :key="post._id"
                cols="12" md="6" lg="4"
            >
                <b-card class="h-100 shadow-sm">
                <div v-if="post.media?.length" class="mb-2">
                    <swiper :slides-per-view="1" navigation pagination>
                    <swiper-slide
                        v-for="(item, index) in post.media"
                        :key="index"
                    >
                        <img
                        v-if="item.type.startsWith('image/')"
                        :src="item.url"
                        class="img-fluid rounded"
                        style="object-fit: cover; width: 100%; height: 250px;"
                        />
                        <video
                        v-else-if="item.type.startsWith('video/')"
                        :src="item.url"
                        controls
                        class="w-100 rounded"
                        style="height: 250px; object-fit: cover;"
                        />
                    </swiper-slide>
                    </swiper>
                </div>
                <h5 class="mb-1">{{ post.title }}</h5>
                <p class="text-muted small mb-2" style="min-height: 50px;">
                    {{ post.content }}
                </p>
                <div class="d-flex justify-content-between align-items-center">
                    <small class="text-muted">🕓 {{ formatDate(post.createdAt) }}</small>
                    <div v-if="post.self">
                    <b-button size="sm" variant="outline-warning" @click="editPost(post)">
                        ✏️
                    </b-button>
                    <b-button size="sm" variant="outline-danger" class="ms-2" @click="deletePost(post._id)">
                        🗑️
                    </b-button>
                    </div>
                </div>
                </b-card>
            </b-col>
        </b-row>
        <b-alert v-else show variant="danger">
            You have no posts yet.
        </b-alert>
        <div ref="infiniteScrollTrigger" class="text-center my-3 w-100">
            <b-spinner v-if="loading" small></b-spinner>
            <small v-else-if="finished" class="text-muted">📍 No more posts</small>
        </div>
    </b-card>
    <BaseModal v-model="showAdd" size="lg" title="add posts">
        <AddPost @add="add"/>
    </BaseModal>
    <BaseModal v-model="showEdit" size="lg" title="Edit Post">
        <EditPost
        :post-data="editingPost"
        @edit="edit"/>
    </BaseModal>
</template>
<script setup>
    import { ref, onMounted, onUnmounted  } from 'vue'
    import { sendApi } from '@/plugins/api';
    import BaseModal from '@/components/tooles/BaseModal.vue';
    import AddPost from '@/components/post/AddPost.vue';
    import EditPost from '@/components/post/EditPost.vue';
    import { Swiper, SwiperSlide } from 'swiper/vue';
    import 'swiper/swiper-bundle.css';
    const showAdd = ref(false)
    const showEdit = ref(false)
    const posts = ref([])
    const page = ref(1)
    const limit = 10
    const loading = ref(false)
    const finished = ref(false)
    const editingPost = ref(null)
    const infiniteScrollTrigger = ref(null)
    let observer = null
    let refreshInterval = null
    const add = () => {
        showAdd.value = false
        refreshPosts()
    }
    const mapPosts = (rawPosts) => {
        return rawPosts.map(p => ({
            _id: p.ID,
            userID: p.UserID,
            title: p.Title,
            content: p.Content,
            media: p.media || [],
            self: p.self,
            createdAt: p.CreatedAt ? new Date(p.CreatedAt) : null
        }))
    }
    const formatDate = (d) => {
        if (!d) return ''
        let dateObj = d instanceof Date ? d : new Date(d)
        if (isNaN(dateObj.getTime())) {
            return ''
        }
    //   return dateObj.toLocaleString('fa-IR')
        return dateObj.toLocaleString()
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
    const edit = () => {
        showEdit.value = false
        editingPost.value = null
        refreshPosts()
    }
    const editPost = (post) => {
        showEdit.value = true
        editingPost.value = { ...post }
    }
    const addPost = () => {
        showAdd.value=true
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
            const freshPosts = mapPosts(res.data || []);
            const existingIds = new Set(posts.value.map(p => p._id));
            freshPosts.forEach(p => {
                const existingIndex = posts.value.findIndex(ep => ep._id === p._id);
                if (existingIndex === -1) {
                    posts.value.unshift(p);
                } else {
                    const old = posts.value[existingIndex];
                    const hasChanged = JSON.stringify(old) !== JSON.stringify(p);
                    if (hasChanged) {
                        posts.value[existingIndex] = p;
                    }
                }
            });
        } catch (e) {
            console.error('خطا در رفرش پست‌ها:', e)
        }
    }
    onMounted(() => {
        loadPosts()
        observer = new IntersectionObserver((entries) => {
            const entry = entries[0]
            if (entry.isIntersecting && !loading.value && !finished.value) {
                loadPosts()
            }
        }, { threshold: 0.5 })
        if (infiniteScrollTrigger.value) {
            observer.observe(infiniteScrollTrigger.value)
        }
        refreshInterval = setInterval(refreshPosts, 15000)
    })
    onUnmounted(() => {
        if (observer && infiniteScrollTrigger.value) {
            observer.unobserve(infiniteScrollTrigger.value)
        }
        if (refreshInterval){
            clearInterval(refreshInterval)
            refreshInterval = null
        } 
    })
</script>