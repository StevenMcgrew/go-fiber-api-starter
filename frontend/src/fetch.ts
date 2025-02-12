import { S } from "./store.svelte";

export async function submitSignUp(formData: Object) {
    const url = S.baseFetchUrl + "/users"
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formData)
        });
        const data = await response.json();
        if (data.status !== "success") {
            throw new Error(data.error)
        }
        return data
    } catch (err: any) {
        throw err
    }
}

export async function submitVerification(formData:Object) {
    const url = S.baseFetchUrl + "/auth/verify-email"
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formData)
        });
        const data = await response.json();
        if (data.status !== "success") {
            throw new Error(data.error)
        }
        return data
    } catch (err: any) {
        throw err
    }
}
