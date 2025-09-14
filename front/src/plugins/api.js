import axios from 'axios';
const api = axios.create({
  baseURL: '/api',
  timeout: 5000,
});
export async function sendApi({ method = 'get', url = '', data = {}, headers = {} }) {
  try {
    const response = await api({
      method,
      url,
      data,
      headers
    });
    return response.data;
  } catch (error) {
    console.error('API Error:', error);
    return { error: true, message: error.response?.data || error.message };
  }
}
export default api;