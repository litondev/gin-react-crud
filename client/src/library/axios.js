import axios from "axios";

axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';
axios.defaults.baseURL = process.env.REACT_APP_API_URL;

axios.interceptors.request.use(config => { 	
  	config.headers['Authorization'] = localStorage.getItem('user-token') ? "Bearer "+localStorage.getItem("user-token") : null
  	return config;
});

window.$axios = axios;
