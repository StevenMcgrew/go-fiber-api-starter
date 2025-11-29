<script lang="ts">
    import { store, clearModal } from "../../store.svelte";
    import { submitForm } from "../../fetch";
    import { modalComp } from "../../types";

    let isLoading = false;
    let error: any = null;

    async function onsubmit(e: SubmitEvent) {
        e.preventDefault();

        isLoading = true;
        error = null;

        const form = e.currentTarget as HTMLFormElement;
        const formData = new FormData(form);
        const url = $store.baseFetchUrl + "/auth/reset-password/request";

        try {
            await submitForm(formData, url);
        } catch (err) {
            error = err;
        } finally {
            isLoading = false;
            if (error === null) {
                $store.showModal = modalComp.ResetPasswordForm;
            }
        }
    }
</script>

<div class="form-wrapper">
    <form {onsubmit} class="auth-form">
        <h3>Request Password Reset</h3>
        <p>
            To reset your password, a code will be emailed to you. Please enter your account's email address.
        </p>

        <label for="email"><b>Email Address</b></label>
        <input
            id="email"
            type="email"
            name="email"
            bind:value={$store.user.email}
            required
        />

        <div class="form-btn-box">
            <button type="button" onclick={clearModal}
                >Cancel</button
            >
            <button type="submit">Request Reset</button>
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
