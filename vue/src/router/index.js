import { createRouter, createWebHistory } from 'vue-router';
import { sendApi } from '@/plugins/api';
import Login from '@/components/user/Login.vue';
import Register from '@/components/user/Register.vue';
import Dashboard from '@/components/Dashboard.vue';
import Home from '@/components/Home.vue';
import Welcome from '@/components/Welcome.vue';
import NotFind from '@/components/NotFind.vue';
import Setting from '@/components/user/Setting.vue';
const routes = [
  { path: '/', name: 'welcome', component: Welcome, meta: { isLogin: false } },
  { path: '/login', name: 'login', component: Login, meta: { isLogin: false } },
  { path: '/register', component: Register, name: 'register', meta: { isLogin: true } },
  { path: '/home', component: Home, name: 'home', meta: { isLogin: true, hasUserInfo: true } },
  { path: '/dashboard', component: Dashboard, name: 'dashboard', meta: { isLogin: true, hasUserInfo: true } },
  { path: '/setting', component: Setting, name: 'user-setting', meta: { isLogin: true, hasUserInfo: true } },
  { path: '/:pathMatch(.*)*', name: '404', component: NotFind }
];
const router = createRouter({
  history: createWebHistory(),
  routes,
});
async function checkTokenStatus() {
  try {
    const res = await sendApi({
      method: 'GET',
      url: '/auth/validate',
      autoCheckToken: true
    });
    if (res.error) return { loggedIn: false };
    return res;
  } catch {
    return { loggedIn: false };
  }
}
router.beforeEach(async (to, from, next) => {
  if (to.name) document.title = to.name;
  const mustBeLogin = to.meta.isLogin;
  const mustHasInfo = to.meta.hasUserInfo;
  const status = await checkTokenStatus();
  const userIsLogin = status.loggedIn || false;
  const userInfo = status.userHasInfo || false;
  if (userIsLogin && mustHasInfo && !userInfo) return next('/register');
  if (userIsLogin && !mustBeLogin) return next('/home');
  if (!userIsLogin && mustBeLogin) return next('/login');
  next();
});
export default router;