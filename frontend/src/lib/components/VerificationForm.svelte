<script lang="ts">
    import { submitForm } from "../../fetch";
    import { S } from "../../store.svelte";
    import type { User } from "../../types";

    let isLoading = false;
    let error: any = null;
    let response: any = null;

    function setUser(res: any) {
        const user: User = {
            token: res.data.token,
            id: res.data.id,
            email: res.data.email,
            username: res.data.username,
            role: res.data.role,
            status: res.data.status,
            imageUrl: res.data.imageUrl,
        };
        S.user = user;
    }

    async function onsubmit(e: SubmitEvent) {
        e.preventDefault();

        isLoading = true;
        error = null;

        const form = e.currentTarget as HTMLFormElement
        const formData = new FormData(form)
        const url = S.baseFetchUrl + "/auth/verify-email"

        try {
            response = await submitForm(formData, url);
        } catch (err) {
            error = err;
        } finally {
            isLoading = false;
            if (error === null) {
                form.reset()
                setUser(response);
                S.showModal = "";
            }
        }
    }
</script>

<div class="form-wrapper">
    <form {onsubmit} class="auth-form">
        <h3>Verification</h3>
        <p>
            A verification code was sent to your email address. Please check
            your email and enter the code here:
        </p>

        <label for="email"><b>Email Address</b></label>
        <input
            id="email"
            type="text"
            name="email"
            value={S.user?.email}
            required
        />

        <label for="verificationCode"><b>Verification Code</b></label>
        <input
            id="verificationCode"
            name="verificationCode"
            type="text"
            required
            maxlength="6"
        />

        <div class="form-btn-box">
            <button type="submit">Verify</button>
        </div>
        {#if isLoading}
            <p class="form-status-text">Submitting...</p>
        {:else if error}
            <p class="error-text form-status-text">
                Error: {error.message}
            </p>
        {/if}
    </form>
</div>
