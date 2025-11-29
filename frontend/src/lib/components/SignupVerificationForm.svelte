<script lang="ts">
    import { submitForm } from "../../fetch";
    import { store, clearModal } from "../../store.svelte";
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
        const url = $store.baseFetchUrl + "/auth/verify-email";

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
                    text: "Verified! You are now logged in.",
                };
                clearModal();
            }
        }
    }

    async function sendCodeAgain() {
        isLoading = true;
        error = null;

        const formData = new FormData();
        formData.append("email", $store.user.email);
        const url = $store.baseFetchUrl + "/auth/resend-email-verification";

        try {
            await submitForm(formData, url);
        } catch (err) {
            error = err;
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="form-wrapper">
    <form {onsubmit} class="auth-form">
        <h3>Verification</h3>
        <p>
            During sign-up a verification code is sent to your email address.
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
        <p class="send-again-txt">
            If you don't receive the email after a few minutes, check your junk
            mail folder or...
        </p>
        <button
            type="button"
            onclick={sendCodeAgain}
            class="more-options-txt send-again-btn">Send code again</button
        >
    </form>
</div>

<style>
    .send-again-txt {
        margin: 2rem 0rem 0rem 0rem;
        font-size: 10px;
    }
    .send-again-btn {
        margin-top: 0.3rem;
    }
</style>
