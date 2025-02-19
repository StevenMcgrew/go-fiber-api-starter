import { orient, toastColor } from "./types";
import type { Store, User, ToastData } from "./types";
import { writable } from "svelte/store";

export const emptyUser: User = {
    token: "",
    id: 0,
    email: "",
    username: "",
    role: "",
    status: "",
    imageUrl: "",
}

export const emptyToast: ToastData = {
    color: toastColor.grey,
    text: "",
}

export const initialStore: Store = {
    baseFetchUrl: "http://127.0.0.1:8080/api/v1",
    baseStorageUrl: "http://127.0.0.1:8080/temp-storage",
    orientLoginBtns: orient.horiz,
    showLoginBtns: true,
    showModal: "",
    showToast: emptyToast,
    stayLoggedIn: false,
    user: emptyUser,
}

function createPersistedStore(key: string, initialValue: Store) {
    const storedValue = localStorage.getItem(key);
    const initial: Store = storedValue ? JSON.parse(storedValue) : initialValue;
  
    const store = writable(initial);
  
    store.subscribe((value) => {
      localStorage.setItem(key, JSON.stringify(value));
    });
  
    return store;
  }
  
  export const store = createPersistedStore('store', initialStore);