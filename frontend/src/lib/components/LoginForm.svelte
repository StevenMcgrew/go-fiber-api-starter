<script lang="ts">
    import { S } from "../../store.svelte";

    function onsubmit(e: SubmitEvent) {
        e.preventDefault();
        e.stopPropagation();
    }
</script>

<div class="form-wrapper">
    <form {onsubmit} class="auth-form">
        <h3>Log In</h3>

        <label for="email"><b>Email Address</b></label>
        <input id="email" type="text" name="email" autocomplete="email" required />

        <label for="password"><b>Password</b></label>
        <input id="password" type="password" name="password" autocomplete="current-password" required />

        <div class="form-btn-box">
            <button type="button" onclick={() => (S.showModal = null)}>Cancel</button>
            <button type="submit">Log In</button>
        </div>
    </form>
</div>
