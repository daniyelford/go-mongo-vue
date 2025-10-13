import Login from '@users/components/user/Login.vue';
import Register from '@users/components/user/Register.vue';
import Setting from '@users/components/user/Setting.vue';
export const usersRoutes  = [
  { path: '/login', name: 'login', component: Login, meta: { isLogin: false } },
  { path: '/register', component: Register, name: 'register', meta: { isLogin: true } },
  { path: '/setting', component: Setting, name: 'user-setting', meta: { isLogin: true, hasUserInfo: true } },
];