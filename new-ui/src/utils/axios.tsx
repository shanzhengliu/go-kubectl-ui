import axios from "axios";

export const axiosInstance = axios.create(
    {
        timeout: 4000
    }
);
axiosInstance.interceptors.response.use(
    (response) => {
      return response;
    },
    (error) => {
      if (error.code === 'ECONNABORTED' && error.message.includes('timeout')) {
        alert('Following things may fixed the timeout issue: \n 1: access okta listen url like http://localhost:8000 for okta login \n 2:check with administator for permissions \n 3: check cluter health maybe it is down' );
      }
      return Promise.reject(error);
    }
  );