import { create } from "zustand";
interface User {
  name: string;
  email?: string;
  groups?: string[];
}

interface UserStore {
  user: User | null;
  setUser: (userInfo: User) => void;
}

export const useUserStore = create<UserStore>((set) => ({
  user: { name: "User" },
  setUser: (userInfo: User) => set({ user: userInfo }),
}));

interface LeftSideBarStore {
    show: string;
    setShow: (show: string) => void;
    }


export const leftSideBarStore = create<LeftSideBarStore>((set) => ({
    show: "k8sViewer",
    setShow: (show: string) => set({ show }),
 }));

 interface K8SCurrentComponentStore {
    currentComponent: string;
    setCurrentComponent: (currentComponent: string) => void;
    }

export const k8sCurrentComponentStore = create<K8SCurrentComponentStore>((set) => ({
    currentComponent: "Pod",
    setCurrentComponent: (currentComponent: string) => set({ currentComponent }),
 }));    



