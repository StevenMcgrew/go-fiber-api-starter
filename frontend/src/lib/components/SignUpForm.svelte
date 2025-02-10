<script lang="ts">
    function onsubmit(e: SubmitEvent) {
        e.preventDefault();
        e.stopPropagation();
    }
</script>

<form {onsubmit} class="auth-form">
    <h3>Sign Up</h3>

    <label for="email"><b>Email Address</b></label>
    <input type="text" name="email" required />

    <label for="username"><b>Username</b></label>
    <input type="text" name="username" required />

    <label for="password"><b>Password</b></label>
    <input type="password" name="password" required />

    <label for="passwordRepeat"><b>Repeat Password</b></label>
    <input type="password" name="passwordRepeat" required />

    <div class="form-btn-box">
        <a href="/"><button type="button">Cancel</button></a>
        <button type="submit">Sign Up</button>
    </div>
</form>

<style>
</style>
