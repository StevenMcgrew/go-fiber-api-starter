<script lang="ts">
    import { submitForm } from "../../fetch";
    import { store, switchToLocalStorage } from "../../store.svelte";
    import { type User, modalComp, toastColor } from "../../types";

    let isLoading = false;
    let error: any = null;
    let response: any = null;
    let emailInputValue = "";
    let passwordInputValue = "";
    let stayLoggedIn = false;

    function setUser(res: any, stay: boolean) {
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
        if (stay) {
            switchToLocalStorage();
        }
    }

    function customFormReset() {
        emailInputValue = "";
        passwordInputValue = "";
    }

    async function onsubmit(e: SubmitEvent) {
        e.preventDefault();

        isLoading = true;
        error = null;

        const form = e.currentTarget as HTMLFormElement;
        const formData = new FormData(form);
        const url = $store.baseFetchUrl + "/auth/login";

        try {
            response = await submitForm(formData, url);
        } catch (err) {
            error = err;
        } finally {
            isLoading = false;
            if (error === null) {
                customFormReset();
                setUser(response, stayLoggedIn);
                $store.showToast = {
                    color: toastColor.green,
                    text: "Logged In!",
                };
                $store.showModal = "";
            }
        }
    }
</script>

<div class="form-wrapper">
    <form {onsubmit} class="auth-form">
        <h3>Log In</h3>

        <label for="email"><b>Email Address</b></label>
        <input
            id="email"
            type="email"
            name="email"
            autocomplete="email"
            required
            bind:value={emailInputValue}
        />

        <label for="password"><b>Password</b></label>
        <input
            id="password"
            type="password"
            name="password"
            autocomplete="current-password"
            required
            bind:value={passwordInputValue}
        />

        <label class="ckbox-label">
            <input
                id="dummyToPreventWarning"
                type="checkbox"
                bind:checked={stayLoggedIn}
            /> Stay logged in
        </label>

        <div class="form-btn-box">
            <button type="button" onclick={() => ($store.showModal = "")}
                >Cancel</button
            >
            <button type="submit">Log In</button>
        </div>
        {#if isLoading}
            <p class="form-status-text">Submitting...</p>
        {:else if error}
            <p class="error-text form-status-text">
                Error: {error.message}
            </p>
        {/if}
        <button
            class="more-options-txt"
            onclick={() => ($store.showModal = modalComp.ForgotPasswordForm)}
            >Forgot password</button
        >
    </form>
</div>

<style>
    .ckbox-label {
        padding: 10px 0px 0px 5px;
        font-size: 14px;
    }
</style>
