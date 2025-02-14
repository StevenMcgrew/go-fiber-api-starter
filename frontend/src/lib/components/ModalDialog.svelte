<script lang="ts">
    import { S } from "../../store.svelte";
    import { modalComp } from "../../types";
    import LoginForm from "./LoginForm.svelte";
    import SignUpForm from "./SignUpForm.svelte";
    import VerificationForm from "./VerificationForm.svelte";

    let dialog: HTMLDialogElement;

    $effect(() => {
        if (S.showModal === "") {
            dialog.close();
        } else {
            dialog.showModal();
        }
    });
</script>

<dialog bind:this={dialog}>
    <button class="close-btn" onclick={() => (S.showModal = "")}>×</button>
    {#if S.showModal == modalComp.LoginForm}
        <LoginForm />
    {:else if S.showModal == modalComp.SignUpForm}
        <SignUpForm />
    {:else if S.showModal == modalComp.VerificationForm}
        <VerificationForm />
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
