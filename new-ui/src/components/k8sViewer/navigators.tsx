import { useState } from "react";
import { Service } from "./service";
import { Pod } from "./pod";
import { Configmap } from "./configmap";
import { Secret } from "./secret"
import { Deployment } from "./deployment";
import { Ingress } from "./ingress";
import { Resource } from "./resource";
import { LOGOUT } from "../../utils/endpoints";
import { Dropdown } from "flowbite-react";
import { axiosInstance } from "../../utils/axios";
import { k8sCurrentComponentStore, useUserStore } from "../../react-context/shareContext";
import { UserInfoModal } from "./userInfoModal";
export const Navigator = () => {
  const user = useUserStore((state) => state.user);
  const menuMap: { [key: string]: any } = {
    Pod: <Pod />,
    Deployment: <Deployment />,
    Service: <Service />,
    Configmap: <Configmap />,
    Secret: <Secret />,
    Ingress: <Ingress />,
    Resource: <Resource />,

  };

  // const [currentComponent, setCurrentComponent] = useState("Resource");
  const [userInfoModal, setUserInfoModal] = useState(false);
  const currentComponent = k8sCurrentComponentStore((state) => state.currentComponent);
  const setCurrentComponent = k8sCurrentComponentStore((state) => state.setCurrentComponent);

  const renderComponent = () => {
    return menuMap[currentComponent];
  };

  const signOut = () => {
    axiosInstance.get(LOGOUT).then(() => {
      window.location.reload();
    });
  };

  const userInfo = () => {
    setUserInfoModal(true);
  };

  return (
    <div className="flex h-screen">
      <UserInfoModal
        show={userInfoModal}
        setShow={() => {
          setUserInfoModal(false);
        }}
      />
      <div className="flex-grow">
        <nav className="bg-white border-gray-200 ">
          <div className="max-w-screen-xxl flex flex-wrap items-center justify-between mx-auto p-4">
            <a className="flex items-center space-x-3 rtl:space-x-reverse">
              <img src="/assets/kubernetes.svg" className="h-8" alt="Logo" />
              <span className="self-center text-2xl font-semibold whitespace-nowrap dark:text-white">
                Kubernetes Management
              </span>
            </a>
            <div
              className="w-full md:block md:w-auto"
              id="navbar-dropdown"
            >
              <ul className="flex flex-col font-medium p-4 md:p-0 mt-4 border border-gray-100 rounded-lg bg-gray-50 md:space-x-8 rtl:space-x-reverse md:flex-row md:mt-0 md:border-0 md:bg-white dark:bg-gray-800 md:dark:bg-gray-900 dark:border-gray-700">
                {Object.entries(menuMap).map(([key]) => (
                  <li key={key}>
                    <a
                      className={`block py-2 px-3 text-gray-900 rounded hover:bg-gray-100 md:hover:bg-transparent md:border-0 md:hover:text-blue-700 md:p-0 dark:text-white md:dark:hover:text-blue-500 dark:hover:bg-gray-700 dark:hover:text-white md:dark:hover:bg-transparent ${
                        currentComponent === key ? "text-yellow-700" : ""
                      }`}
                      onClick={() => {
                        
                        setCurrentComponent(key);
                      }}
                    >
                      {key}
                    </a>
                  </li>
                ))}
                <li>
                  {user?.name ? (
                    <Dropdown label={user?.name} inline size="sm">
                      <Dropdown.Item onClick={userInfo}>Settings</Dropdown.Item>
                      <Dropdown.Item onClick={signOut}>Sign out</Dropdown.Item>
                    </Dropdown>
                  ) : (
                    "User"
                  )}
                </li>
              </ul>
            </div>
          </div>
        </nav>
        <div className="flex-grow w-[calc(100vw-4rem)">
        {renderComponent()}
        </div>
      </div>
    </div>
  );
};
