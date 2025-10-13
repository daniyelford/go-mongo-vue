import { createRouter, createWebHistory } from 'vue-router';
import { sendApi } from '@/plugins/api';
import Dashboard from '@/components/Dashboard.vue';
import Home from '@/components/Home.vue';
import Welcome from '@/components/Welcome.vue';
import NotFind from '@/components/NotFind.vue';
import BusinessList from '@/components/business/BusinessList.vue';
import { usersRoutes } from '@users/router'
const routes = [
  { path: '/', name: 'welcome', component: Welcome, meta: { isLogin: false } },
  { path: '/home', component: Home, name: 'home', meta: { isLogin: true, hasUserInfo: true } },
  { path: '/dashboard', component: Dashboard, name: 'dashboard', meta: { isLogin: true, hasUserInfo: true } },
  { path: '/business', component: BusinessList, name: 'business', meta: { isLogin: true, hasUserInfo: true } },
  { path: '/:pathMatch(.*)*', name: '404', component: NotFind },
  ...usersRoutes
];
const router = createRouter({
  history: createWebHistory(),
  routes,
});
async function checkTokenStatus(isLogin) {
  if(isLogin){
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
}
router.beforeEach(async (to, from, next) => {
  if (to.name) document.title = to.name;
  const mustBeLogin = to.meta.isLogin;
  const mustHasInfo = to.meta.hasUserInfo;
  const status = await checkTokenStatus(mustBeLogin);
  const userIsLogin = status?.loggedIn || false;
  const userInfo = status?.userHasInfo || false;
  if (userIsLogin && mustHasInfo && !userInfo) return next('/register');
  if (userIsLogin && !mustBeLogin) return next('/home');
  if (!userIsLogin && mustBeLogin) return next('/login');
  next();
});
export default router;