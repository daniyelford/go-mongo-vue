import router from '@/router';
import axios from 'axios';
const api = axios.create({
  baseURL: '/api',
  timeout: 5000,
});
function getToken() {
  return localStorage.getItem('jwt');
}
function getRefreshToken() {
  return localStorage.getItem('refresh');
}
function setToken(token) {
  localStorage.setItem('jwt', token);
}
function removeTokens() {
  localStorage.removeItem('jwt');
  localStorage.removeItem('refresh');
  router.push({path:'/login'})
}
export async function sendApi({ method = 'get', url = '', data = {}, headers = {}, autoCheckToken = false }) {
  try {
    let token = getToken();
    if (autoCheckToken && token) {
      const res = await api({
        method: 'GET',
        url: '/auth/validate',
        headers: { Authorization: `Bearer ${token}` }
      });
      if (res.data?.expired) {
        const refreshToken = getRefreshToken();
        if (!refreshToken) {
          removeTokens();
          throw new Error('Session expired');
        }
        const refreshRes = await api({
          method: 'POST',
          url: '/token/refresh',
          data: { refreshToken }
        });
        if (!refreshRes.data?.accessToken) {
          removeTokens();
          throw new Error('Session expired');
        }
        token = refreshRes.data.accessToken;
        setToken(token);
      }
    }
    const response = await api({
      method,
      url,
      data,
      headers: autoCheckToken && token ? { ...headers, Authorization: `Bearer ${token}` } : headers
    });
    return response.data;
  } catch (error) {
    console.error('API Error:', error);
    return { error: true, message: error.response?.data || error.message };
  }
}
export function removeTokensOut() {
  removeTokens()
}
export default api;