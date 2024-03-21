import axios from "axios";
import Swal from 'sweetalert2'
import { OKTA, USERINFO } from "./endpoints";
export const axiosInstance = axios.create(
    {
        timeout: 3500
    }
);
axiosInstance.interceptors.response.use(
    (response) => {
      return response;
    },
    (error) => {
      if (error.code === 'ECONNABORTED' && error.message.includes('timeout')) {
        Swal.fire({
            icon: "error",
            html: '<div style="text-align:left">Following steps may fixed the timeout issue: <br>  1: access okta listen url like <a class="text-blue-700" href="http://localhost:8000", target=_"blank">http://localhost:8000</a> for okta login <br> 2:check with administator for permissions <br> 3: check cluter health maybe it is down</div>',
          });
      }
      return Promise.reject(error);
    }
  );

  export const  authVerify = async () => {
    return await axiosInstance
      .get(USERINFO)
      .then(async (response) => {
        if (response.status === 200) {
          return response.data;
        }
        if (response.status === 401) {
            console.log("Unauthorized");

            const url =  await axios.get(OKTA);
            console.log(url);
        }
      })
      .catch(async (error) => {
        console.log(error);
   
          console.log("Unauthorized");
          const url =  await axios.get(OKTA);
          window.location.href=url.data;
        
      });
  }