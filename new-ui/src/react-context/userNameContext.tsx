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
