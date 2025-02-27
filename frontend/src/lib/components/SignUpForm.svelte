<script lang="ts">
    import { submitForm } from "../../fetch";
    import { store } from "../../store.svelte";
    import { modalComp, type User } from "../../types";
    import { debounce } from "../../utils.js";
    import { getEmailValidationWarnings, getPasswordValidationWarnings, getUsernameValidationWarnings } from "../../validation";

    let isLoading = false;
    let error: any = null;
    let response: any = null;

    let emailWarnings = "";
    let usernameWarnings = "";
    let passwordWarnings = "";
    let passwordRepeatWarnings = "";

    async function handleEmailInput(email: string) {
        emailWarnings = getEmailValidationWarnings(email);
        if (!emailWarnings) {
            const formData = new FormData();
            formData.append("email", email);
            const url = $store.baseFetchUrl + "/users/email/availability";
            try {
                const response = await submitForm(formData, url);
                if (response.data === false) {
                    emailWarnings = "Already in use by another user.";
                }
            } catch (err: any) {
                emailWarnings = err.message;
            }
        }
    }

    async function handleUsernameInput(username: string) {
        usernameWarnings = getUsernameValidationWarnings(username)
        if (!usernameWarnings) {
            const formData = new FormData();
            formData.append("username", username);
            const url = $store.baseFetchUrl + "/users/username/availability";
            try {
                const response = await submitForm(formData, url);
                if (response.data === false) {
                    usernameWarnings = "Already in use by another user.";
                }
            } catch (err: any) {
                usernameWarnings = err.message;
            }
        }
    }

    function handlePasswordInput(password: string) {
        passwordWarnings = getPasswordValidationWarnings(password)
    }

    const delayedHandleEmailInput = debounce(handleEmailInput, 1000);
    const delayedHandleUsernameInput = debounce(handleUsernameInput, 1000);
    const delayedHandlePasswordInput = debounce(handlePasswordInput, 1000);


    async function oninput_email(e: Event) {
        emailWarnings = ""
        const el = e.currentTarget as HTMLFormElement;
        delayedHandleEmailInput(el.value);
    }

    async function oninput_username(e: Event) {
        usernameWarnings = ""
        const el = e.currentTarget as HTMLFormElement;
        delayedHandleUsernameInput(el.value) 
    }

    function oninput_password(e: Event) {
        passwordWarnings = ""
        const el = e.currentTarget as HTMLFormElement;
        delayedHandlePasswordInput(el.value)
    }

    function setUser(res: any) {
        const user: User = {
            token: "",
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
        const url = $store.baseFetchUrl + "/users";

        try {
            response = await submitForm(formData, url);
        } catch (err) {
            error = err;
        } finally {
            isLoading = false;
            if (error === null) {
                form.reset();
                setUser(response);
                $store.showModal = modalComp.SignupVerificationForm;
            }
        }
    }
</script>

<div class="form-wrapper">
    <form {onsubmit} class="auth-form">
        <h3>Sign Up</h3>

        <label for="email"><b>Email Address</b></label>
        <input
            class="form-input"
            id="email"
            type="email"
            name="email"
            autocomplete="email"
            required
            oninput={oninput_email}
        />
        <p class="error-text form-input-warning">{emailWarnings}</p>

        <label for="username"><b>Username</b></label>
        <input
            class="form-input"
            id="username"
            type="text"
            name="username"
            autocomplete="username"
            required
            oninput={oninput_username}
        />
        <p class="error-text form-input-warning">{usernameWarnings}</p>

        <label for="password"><b>Password</b></label>
        <input
            class="form-input"
            id="password"
            type="password"
            name="password"
            required
            oninput={oninput_password}
        />
        <p class="error-text form-input-warning">{passwordWarnings}</p>

        <label for="passwordRepeat"><b>Repeat Password</b></label>
        <input
            class="form-input"
            id="passwordRepeat"
            type="password"
            name="passwordRepeat"
            required
        />
        <p class="error-text form-input-warning">{passwordRepeatWarnings}</p>

        <div class="form-btn-box">
            <button type="button" onclick={() => ($store.showModal = "")}
                >Cancel</button
            >
            <button type="submit">Sign Up</button>
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
            onclick={() => ($store.showModal = modalComp.SignupVerificationForm)}
            >Enter a verification code</button
        >
    </form>
</div>
