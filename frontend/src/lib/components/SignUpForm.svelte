<script lang="ts">
    import { S } from "../../store.svelte";

    function onsubmit(e: SubmitEvent) {
        e.preventDefault();
        e.stopPropagation();
    }
</script>

<div class="form-wrapper">
    <form {onsubmit} class="auth-form">
        <h3>Sign Up</h3>

        <label for="email"><b>Email Address</b></label>
        <input id="email" type="text" name="email" required />

        <label for="username"><b>Username</b></label>
        <input id="username" type="text" name="username" required />

        <label for="password"><b>Password</b></label>
        <input id="password" type="password" name="password" required />

        <label for="passwordRepeat"><b>Repeat Password</b></label>
        <input id="passwordRepeat" type="password" name="passwordRepeat" required />

        <div class="form-btn-box">
            <button type="button" onclick={() => (S.showModal = null)}>Cancel</button>
            <button type="submit">Sign Up</button>
        </div>
    </form>
</div>
