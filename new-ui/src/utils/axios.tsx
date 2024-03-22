import axios from "axios";
import Swal from 'sweetalert2'
import { OKTA, USERINFO } from "./endpoints";
import { useContext } from "react";
import {  useUserStore } from "../react-context/userNameContext";
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
            text: "Maybe your don't have permission to load the data",
          });
      }
      return Promise.reject(error);
    }
  );
  
  const openOkta = async () => {
    const url = await axios.get(OKTA);
    window.open(url.data);
  }

 
  export const  authVerify = async () => {
    return await axiosInstance
      .get(USERINFO)
      .then(async (response) => {
        if (response.status === 200) {
          useUserStore.getState().setUser(response.data);
          // user.setUser(response.data);
          return response.data;
        }
      })
      .catch(async (error) => {
        console.log(error);
       
        Swal.fire({
            icon: "error",
            text: "You are not authorized to access this page, you can click the button to login with Okta. or close the modal and switch the context",
            showConfirmButton: false, 
            html: '<div>You may need to login via Okta</div><button id="innerButton" class="swal2-confirm swal2-styled">Okta</button><button id="close"  class="swal2-danger swal2-styled" >Close</button>',
            didOpen: () => {
              const innerButton = document.getElementById('innerButton');
              if (innerButton) {
                innerButton.addEventListener('click', openOkta);
              }
              const closeButton = document.getElementById('close');
              if (closeButton) {
                closeButton.addEventListener('click', () => {
                  Swal.close();
                });
              }
          }
         });
         return "error";
      });
  }