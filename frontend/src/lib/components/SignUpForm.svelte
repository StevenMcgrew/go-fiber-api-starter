<script lang="ts">
    import { submitSignUp } from "../../fetch";
    import { S } from "../../store.svelte";

    let formData = {
        email: "",
        username: "",
        password: "",
        passwordRepeat: "",
    };

    let data: Object | null = null;
    let isLoading: boolean = false;
    let error: any = null;
    let didSucceed: boolean = false;

    async function onsubmit(e: SubmitEvent) {
        e.preventDefault();
        isLoading = true;
        data = null
        error = null

        try {
            data = await submitSignUp(formData);
            // TODO: do something with data
        } catch (err) {
            error = err;
        } finally {
            isLoading = false;
            if (error === null) {
                didSucceed = true;
                setTimeout(() => {
                    S.showModal = null;
                }, 1000);
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
            <button type="button" onclick={() => (S.showModal = null)}
                >Cancel</button
            >
            <button type="submit">Sign Up</button>
        </div>
        {#if isLoading}
            <p class="form-status-text">Submitting...</p>
        {:else if error}
            <p class="error-text form-status-text">Error: {error.message}</p>
        {:else if didSucceed}
            <p class="success-text form-status-text">Success!</p>
        {/if}
    </form>
</div>
