<script lang="ts">
    import { S } from "../../store.svelte";
    import { modalComp, orient } from "../../types";

    function logOut() {
        // TODO: logout
    }
</script>

<div class="btn-box {S.orientLoginBtns === orient.vert ? orient.vert : ''}">
    {#if S.user?.token}
        <p>{S.user.username}</p>
        <img src={S.baseStorageUrl + S.user.imageUrl} alt="user">
        <button onclick={logOut}>Log Out</button>
    {:else}
        <button
            class={S.orientLoginBtns === orient.vert ? "btn-width" : ""}
            onclick={() => (S.showModal = modalComp.LoginForm)}
        >
            Log In
        </button>
        <button
            class={S.orientLoginBtns === orient.vert ? "btn-width" : ""}
            onclick={() => (S.showModal = modalComp.SignUpForm)}
        >
            Sign Up
        </button>
    {/if}
</div>

<style>
    .btn-box {
        display: flex;
        gap: 0.5rem;
    }
    .vertical {
        flex-direction: column;
    }
    .btn-width {
        width: 89px;
    }
</style>
