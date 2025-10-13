<template>
  <b-container class="mt-5" style="max-width: 500px;">
    <b-card title="register" class="shadow-sm">
      <b-form @submit.prevent="handleRegister">
        <b-form-group label="profile photo" label-for="photo">
            <input
            type="file"
            ref="fileInput"
            style="display:none"
            accept="image/*"
            @change="onFileChange"
            />
            <div class="text-center mb-3" @click="fileInput.click()">
                <img v-if="preview" :src="preview" alt="preview" class="img-thumbnail" style="max-height: 200px;" />
                <svg v-else class="img-thumbnail" style="max-height: 200px;" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0,0,256,256"><g fill="#000000" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-weight="none" font-size="none" text-anchor="none" style="mix-blend-mode: normal"><g transform="scale(5.33333,5.33333)"><path d="M24,2c-5.02628,0 -9,4.54905 -9,10c0,5.45095 3.97372,10 9,10c5.02628,0 9,-4.54905 9,-10c0,-5.45095 -3.97372,-10 -9,-10zM24,4c3.81029,0 7,3.50983 7,8c0,4.49017 -3.18971,8 -7,8c-3.81029,0 -7,-3.50983 -7,-8c0,-4.49017 3.18971,-8 7,-8zM24,24c-8.78758,0 -16.5961,5.62919 -19.375,13.96484l-0.25781,0.76953c-0.85461,2.56251 1.09393,5.26563 3.79492,5.26563h31.67578c2.70031,0 4.64842,-2.70186 3.79492,-5.26367l-0.25781,-0.77148c-2.77888,-8.33664 -10.58742,-13.96484 -19.375,-13.96484zM24,26c7.93442,0 14.96939,5.07029 17.47852,12.59766l0.25586,0.76953c0.4405,1.32219 -0.5028,2.63281 -1.89648,2.63281h-31.67578c-1.39501,0 -2.33983,-1.30932 -1.89844,-2.63281l0.25781,-0.76953c2.5091,-7.52634 9.5441,-12.59766 17.47852,-12.59766z"></path></g></g></svg>
            </div>
        </b-form-group>
        <b-form-group label="name" label-for="name">
          <b-form-input
            id="name"
            v-model="form.name"
            required
          ></b-form-input>
        </b-form-group>

        <b-form-group label="family" label-for="family">
          <b-form-input
            id="family"
            v-model="form.family"
            required
          ></b-form-input>
        </b-form-group>

        <b-button type="submit" variant="primary" block>
          save
        </b-button>
      </b-form>
    <div v-if="error" class="mt-4 text-danger">❌ {{ error }}</div>
    </b-card>
  </b-container>
</template>

<script setup>
import { ref } from "vue";
import { sendApi } from "@/plugins/api";
import router from "@/router";
const error = ref('')
const form = ref({
  name: "",
  family: "",
  photo: null,
});
const preview = ref(null);
const fileInput = ref(null);
function compressImage(file, maxWidth = 800, quality = 0.7) {
  return new Promise((resolve, reject) => {
    try {
      const img = new Image();
      const reader = new FileReader();
      reader.onload = (e) => {
        img.src = e.target.result;
      };
      img.onload = () => {
        const canvas = document.createElement("canvas");
        const ctx = canvas.getContext("2d");
        let width = img.width;
        let height = img.height;
        if (width > maxWidth) {
          height = (height * maxWidth) / width;
          width = maxWidth;
        }
        canvas.width = width;
        canvas.height = height;
        ctx.drawImage(img, 0, 0, width, height);
        canvas.toBlob(
          (blob) => {
            if (!blob) return reject("compression failed");
            resolve(new File([blob], file.name, { type: "image/jpeg" }));
          },
          "image/jpeg",
          quality
        );
      };
      reader.readAsDataURL(file);
    } catch (err) {
      reject(err);
    }
  });
}
const onFileChange = async (event) => {
  const file = event.target.files[0];
  if (file) {
    try {
      const compressed = await compressImage(file, 800, 0.7);
      form.value.photo = compressed;
      preview.value = URL.createObjectURL(compressed);
    } catch (err) {
      console.error("Image compression error:", err);
      alert("Image compression error:")
    }
  }
};
const handleRegister = async () => {
  const data = new FormData();
  data.append("name", form.value.name);
  data.append("family", form.value.family);
  if (form.value.photo) {
    data.append("photo", form.value.photo);
  }
  try {
    const res = await sendApi({
      method: "POST",
      url: "/register/save",
      data: data,
      autoCheckToken: true
    });
    if (res.success) {
      form.value = { name: "", family: "", photo: null };
      preview.value = null;

      router.push({ path: "/home" });
    } else {
        error.value=res?.message?.error
    }
  } catch (err) {
        error.value=err
    }
};
</script>