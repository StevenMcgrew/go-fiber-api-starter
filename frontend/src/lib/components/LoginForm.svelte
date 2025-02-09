<script lang="ts">

function onsubmit(e: SubmitEvent) {
    e.preventDefault()
    e.stopPropagation()
}
</script>

<form {onsubmit} class="auth-form">
    <h3>Log In</h3>

    <label for="email"><b>Email Address</b></label>
    <input type="text" name="email" required>

    <label for="password"><b>Password</b></label>
    <input type="password" name="password" required>

    <div class="form-btn-box">
        <a href="/"><button type="button">Cancel</button></a>
        <button type="submit">Log In</button>
    </div>
</form>

<style>

</style>