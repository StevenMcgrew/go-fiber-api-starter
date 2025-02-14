import { orient } from "./types";
import type { Store, User } from "./types";

const initialUser: User = {
    token: "",
    id: 0,
    email: "",
    username: "",
    role: "",
    status: "",
    imageUrl: "",
}

export const S: Store = $state({
    baseFetchUrl: "http://127.0.0.1:8080/api/v1",
    baseStorageUrl: "http://127.0.0.1:8080/temp-storage",
    orientLoginBtns: orient.horiz,
    showLoginBtns: true,
    showModal: "",
    user: initialUser,
})

function updateObjectValues<T extends object>(obj: T, updates: Partial<T>) {
    for (const [key, value] of Object.entries(updates) as Array<[keyof T, T[keyof T]]>) {
        if (key in obj) {
            obj[key] = value;
        }
    }
}

if (!localStorage.getItem("store")) {
    localStorage.setItem("store", JSON.stringify(S))
} else {
    let ls_store = localStorage.getItem("store")
    if (ls_store) {
        let _store = JSON.parse(ls_store)
        updateObjectValues(S, _store)
    }
}