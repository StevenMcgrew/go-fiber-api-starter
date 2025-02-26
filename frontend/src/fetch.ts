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
            throw new Error(data.error)
        }
        return data
    } catch (err: any) {
        throw err
    }
}
