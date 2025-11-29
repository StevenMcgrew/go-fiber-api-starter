import { store } from "./store.svelte";
import { modalComp } from "./types";

export const submitForm = async (
    formData: FormData,
    url: string,
    _method: string = "POST",
    token: string = "none"
) => {
    try {
        const response = await fetch(url, {
            method: _method,
            body: formData,
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });
        const data = await response.json();
        if (data.status !== "success") {
            if (data.error.includes("expired")) { // JWT expired
                data.error = "Try again after logging in";
                store.update(currentStore => {
                    currentStore.modalText = "⚠️ Your session has expired. Please log in before continuing.";
                    currentStore.showModal = modalComp.LoginForm;
                    return currentStore;
                })
            }
            throw new Error(data.error)
        }
        return data
    } catch (err: any) {
        throw err
    }
}
