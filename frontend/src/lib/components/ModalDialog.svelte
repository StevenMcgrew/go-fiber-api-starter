<script lang="ts">
    import { store } from "../../store.svelte";
    import { modalComp } from "../../types";
    import ForgotPasswordForm from "./ForgotPasswordForm.svelte";
    import LoginForm from "./LoginForm.svelte";
    import ResetPasswordForm from "./ResetPasswordForm.svelte";
    import SignUpForm from "./SignUpForm.svelte";
    import VerificationForm from "./VerificationForm.svelte";

    let dialog: HTMLDialogElement;

    $effect(() => {
        if ($store.showModal === "") {
            dialog.close();
        } else {
            dialog.showModal();
        }
    });
</script>

<dialog bind:this={dialog}>
    <button class="close-btn" onclick={() => ($store.showModal = "")}>×</button>
    {#if $store.showModal == modalComp.LoginForm}
        <LoginForm />
    {:else if $store.showModal == modalComp.SignUpForm}
        <SignUpForm />
    {:else if $store.showModal == modalComp.VerificationForm}
        <VerificationForm />
    {:else if $store.showModal == modalComp.ForgotPasswordForm}
        <ForgotPasswordForm />
    {:else if $store.showModal == modalComp.ResetPasswordForm}
        <ResetPasswordForm />
    {:else}
        <p>Missing content</p>
    {/if}
</dialog>

<style>
    .close-btn {
        position: absolute;
        right: 8px;
        top: 8px;
        padding: 0px 5px 3px 5px;
        font-size: 2rem;
        line-height: 26px;
    }
</style>
