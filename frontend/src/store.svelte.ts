import { orient, toastColor } from "./types";
import type { User, ToastData, Store } from "./types";
import { writable, type Unsubscriber } from "svelte/store";

let unsubscribe: Unsubscriber;

const emptyUser: User = {
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

function createDefaultStore(): Store {
    const s: Store = {
        baseFetchUrl: "http://127.0.0.1:8080/api/v1",
        baseStorageUrl: "http://127.0.0.1:8080/temp-storage",
        orientLoginBtns: orient.horiz,
        showLoginBtns: true,
        showModal: "",
        showToast: emptyToast,
        newEmailAddress: "",
        user: emptyUser,
    }
    return s
}

function createBrowserStorage(key: string) {

    const storedLocal = localStorage.getItem(key);
    if (storedLocal) {
        // create store from localStorage
        const s: Store = JSON.parse(storedLocal);
        const store = writable(s);
        unsubscribe = store.subscribe((value) => {
            localStorage.setItem(key, JSON.stringify(value));
        });
        return store;
    }

    const storedSession = sessionStorage.getItem(key);
    if (storedSession) {
        // create store from sessionStorage
        const s: Store = JSON.parse(storedSession);
        const store = writable(s);
        unsubscribe = store.subscribe((value) => {
            sessionStorage.setItem(key, JSON.stringify(value));
        });
        return store;
    }

    // create default store in sessionStorage
    const s: Store = createDefaultStore();
    const store = writable(s);
    unsubscribe = store.subscribe((value) => {
        sessionStorage.setItem('store', JSON.stringify(value));
    });
    return store;
}

export let store = createBrowserStorage('store')

export function switchToLocalStorage() {
    const storedValue = sessionStorage.getItem('store');
    if (storedValue) {
        if (unsubscribe) { unsubscribe(); }
        const s: Store = JSON.parse(storedValue);
        store = writable(s);
        unsubscribe = store.subscribe((value) => {
            localStorage.setItem('store', JSON.stringify(value));
        });
    } else {
        console.error("Function 'switchToLocalStorage()' error: The key named 'store' was not found in sessionStorage.")
    }
}

export function clearStorageAndReload() {
    localStorage.clear();
    sessionStorage.clear();
    location.reload();
}