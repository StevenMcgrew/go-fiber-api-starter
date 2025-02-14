<script lang="ts">
    import { submitSignUp } from "../../fetch";
    import { S } from "../../store.svelte";
    import { modalComp, type User } from "../../types";

    let formData = {
        email: "",
        username: "",
        password: "",
        passwordRepeat: "",
    };

    let isLoading = false;
    let error: any = null;
    let response: any = null;

    function resetFormData() {
        formData.email = ""
        formData.username = ""
        formData.password = ""
        formData.passwordRepeat = ""
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
        S.user = user;
    }

    async function onsubmit(e: SubmitEvent) {
        e.preventDefault();
        isLoading = true;
        error = null;

        try {
            response = await submitSignUp(formData);
        } catch (err) {
            error = err;
        } finally {
            isLoading = false;
            if (error === null) {
                resetFormData()
                setUser(response)
                S.showModal = modalComp.VerificationForm
            }
        }
    }
</script>

<div class="form-wrapper">
    <form {onsubmit} class="auth-form">
        <h3>Sign Up</h3>

        <label for="email"><b>Email Address</b></label>
        <input
            bind:value={formData.email}
            id="email"
            type="text"
            name="email"
            autocomplete="email"
            required
        />

        <label for="username"><b>Username</b></label>
        <input
            bind:value={formData.username}
            id="username"
            type="text"
            name="username"
            autocomplete="username"
            required
        />

        <label for="password"><b>Password</b></label>
        <input
            bind:value={formData.password}
            id="password"
            type="password"
            name="password"
            required
        />

        <label for="passwordRepeat"><b>Repeat Password</b></label>
        <input
            bind:value={formData.passwordRepeat}
            id="passwordRepeat"
            type="password"
            name="passwordRepeat"
            required
        />

        <div class="form-btn-box">
            <button type="button" onclick={() => (S.showModal = "")}
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
    </form>
</div>
