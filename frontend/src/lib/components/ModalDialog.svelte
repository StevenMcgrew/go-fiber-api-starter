<script lang="ts">
    import { S } from "../../store.svelte";
    import LoginForm from "./LoginForm.svelte";
    import SignUpForm from "./SignUpForm.svelte";

    let dialog: HTMLDialogElement;

    $effect(() => {
        if (S.showModal !== null) {
            dialog.showModal();
        }
    });

    function close() {
        dialog.close();
    }
</script>

<dialog bind:this={dialog}>
    <button onclick={close}>Close</button>
    {#if S.showModal == LoginForm}
        <LoginForm />
    {:else if S.showModal == SignUpForm}
        <SignUpForm />
    {:else}
        <p>No content</p>
    {/if}
</dialog>
