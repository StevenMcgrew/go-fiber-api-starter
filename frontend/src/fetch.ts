export const submitForm = async (formData: FormData, url: string, _method: string = "POST") => {
    try {
        const response = await fetch(url, {
            method: _method,
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
