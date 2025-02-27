<script lang="ts">
    import { submitForm } from "../../fetch";
    import { store } from "../../store.svelte";
    import { toastColor } from "../../types";

    let isLoading = false;
    let error: any = null;
    let response: any = null;

    async function onsubmit(e: SubmitEvent) {
        e.preventDefault();

        isLoading = true;
        error = null;

        const form = e.currentTarget as HTMLFormElement;
        const formData = new FormData(form);
        const url = `${$store.baseFetchUrl}/users/${$store.user.id}/change-email/update`;
        const method = "PATCH"
        const token = $store.user.token

        try {
            response = await submitForm(formData, url, method, token);
        } catch (err) {
            error = err;
        } finally {
            isLoading = false;
            if (error === null) {
                form.reset();
                $store.newEmailAddress = "";
                $store.user.email = response.data.email
                $store.showToast = {
                    color: toastColor.green,
                    text: "Email address updated!",
                };
                $store.showModal = "";
            }
        }
    }
</script>

<div class="form-wrapper">
    <form {onsubmit} class="auth-form">
        <h3>Verification</h3>
        <p>
            A one-time passcode was sent to the new email address.
            Please check the new email and enter the code here:
        </p>

        <label for="email"><b>New Email Address</b></label>
        <input
            id="email"
            type="email"
            name="email"
            bind:value={$store.newEmailAddress}
            required
        />

        <label for="otp"><b>One-Time Passcode</b></label>
        <input
            id="otp"
            name="otp"
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
