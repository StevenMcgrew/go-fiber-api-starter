export async function submitForm(formData: FormData, url: string) {
    try {
        const response = await fetch(url, {
            method: 'POST',
            body: formData
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
