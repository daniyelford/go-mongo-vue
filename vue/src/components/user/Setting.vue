<template>
  <b-container class="my-4">
    <b-card title="تنظیمات کاربر" class="shadow-sm">
      <div class="text-center mb-3">
        <b-avatar
          :src="preview || user.avatar"
          size="6rem"
          class="mb-2"
        ></b-avatar>
        <div>
          <b-form-file
            accept="image/*"
            @change="onFileChange"
            placeholder="انتخاب عکس پروفایل"
          ></b-form-file>
        </div>
      </div>

      <b-form-group label="نام" label-for="nameInput">
        <b-form-input
          id="nameInput"
          v-model="form.name"
          placeholder="نام خود را وارد کنید"
        ></b-form-input>
      </b-form-group>

      <b-form-group label="نام خانوادگی" label-for="familyInput">
        <b-form-input
          id="familyInput"
          v-model="form.family"
          placeholder="نام خانوادگی خود را وارد کنید"
        ></b-form-input>
      </b-form-group>
    </b-card>
  </b-container>
</template>

<script setup>
import { sendApi } from '@/plugins/api'
import { onMounted, ref } from 'vue'

const user = ref({
  avatar: '/images/default-avatar.png',
  name: 'دانیال',
  family: 'فرد'
})

const form = ref({
  name: user.value.name,
  family: user.value.family
})

const preview = ref(null)

function onFileChange(e) {
  const file = e.target.files[0]
  if (file) {
    preview.value = URL.createObjectURL(file)

  }
}

onMounted(async () => {
    try {
        const res = await sendApi({
            method: "POST",
            autoCheckToken: true,
            url: "/user/info",
            data: data,
        });
        if(res.success){
            console.log(res.user);
            
        }
    } catch (error) {
        console.log(error);
    }
})
</script>


