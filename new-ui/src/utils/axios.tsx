import axios from "axios";

export const axiosInstance = axios.create();
// axiosInstance.interceptors.request.use(config => {
//     console.log('Request Headers:', config.headers);
//     return config;
//   });