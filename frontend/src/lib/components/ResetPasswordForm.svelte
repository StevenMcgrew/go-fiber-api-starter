<script lang="ts">
    import { store } from "../../store.svelte";
    import { submitForm } from "../../fetch";
    import { type User, toastColor } from "../../types";

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
        $store.user = user;
    }

    async function onsubmit(e: SubmitEvent) {
        e.preventDefault();

        isLoading = true;
        error = null;

        const form = e.currentTarget as HTMLFormElement;
        const formData = new FormData(form);
        const url = $store.baseFetchUrl + "/auth/reset-password/update";

        try {
            response = await submitForm(formData, url);
        } catch (err) {
            error = err;
        } finally {
            isLoading = false;
            if (error === null) {
                form.reset();
                setUser(response);
                $store.showToast = {
                    color: toastColor.green,
                    text: "Password Reset! You are now logged in.",
                };
                $store.showModal = "";
            }
        }
    }
</script>

<div class="form-wrapper">
    <form {onsubmit} class="auth-form">
        <h3>Reset Password</h3>
        <p>
            A reset code was sent to your email address.
            Please check your email and enter the code here:
        </p>

        <label for="email"><b>Email Address</b></label>
        <input
            id="email"
            type="email"
            name="email"
            bind:value={$store.user.email}
            required
        />

        <label for="resetCode"><b>Reset Code</b></label>
        <input
            id="resetCode"
            name="resetCode"
            type="text"
            required
            maxlength="6"
        />

        <label for="newPassword"><b>New Password</b></label>
        <input id="newPassword" type="password" name="newPassword" required />

        <label for="repeatNewPassword"><b>Repeat New Password</b></label>
        <input
            id="repeatNewPassword"
            type="password"
            name="repeatNewPassword"
            required
        />

        <div class="form-btn-box">
            <button type="button" onclick={() => ($store.showModal = "")}
                >Cancel</button
            >
            <button type="submit">Reset Password</button>
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
