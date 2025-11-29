<script lang="ts">
    import { store, clearModal } from "../../store.svelte";
    import { modalComp } from "../../types";
    import ForgotPasswordForm from "./ForgotPasswordForm.svelte";
    import LoginForm from "./LoginForm.svelte";
    import LoggingOutMsg from "./LoggingOutMsg.svelte";
    import ResetPasswordForm from "./ResetPasswordForm.svelte";
    import SignUpForm from "./SignUpForm.svelte";
    import SignupVerificationForm from "./SignupVerificationForm.svelte";
    import UpdateEmailVerificationForm from "./UpdateEmailVerificationForm.svelte";

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
    <button class="close-btn" onclick={clearModal}>×</button>

    {#if $store.modalText}
        <p class="modal-text">{$store.modalText}</p>
    {/if}

    {#if $store.showModal == modalComp.LoginForm}
        <LoginForm />
    {:else if $store.showModal == modalComp.LoggingOutMsg}
        <LoggingOutMsg />
    {:else if $store.showModal == modalComp.SignUpForm}
        <SignUpForm />
    {:else if $store.showModal == modalComp.SignupVerificationForm}
        <SignupVerificationForm />
    {:else if $store.showModal == modalComp.UpdateEmailVerificationForm}
        <UpdateEmailVerificationForm />
    {:else if $store.showModal == modalComp.ForgotPasswordForm}
        <ForgotPasswordForm />
    {:else if $store.showModal == modalComp.ResetPasswordForm}
        <ResetPasswordForm />
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

    .modal-text {
        padding: 40px 20px;
        margin: 0;
        max-width: 292px;
    }
</style>
