<script lang="ts">
    import { store, clearStorageAndReload } from "../../store.svelte";
    import { modalComp, orient } from "../../types";

    function logOut() {
        clearStorageAndReload();
    }
</script>

<div class="btn-box {$store.orientLoginBtns === orient.vert ? orient.vert : ''}">
    {#if $store.user?.token}
        <p class="username-txt">{$store.user.username}</p>
        <img
            class="user-img"
            src={$store.baseStorageUrl + $store.user.imageUrl}
            alt="user"
        />
        <button onclick={logOut}>Log Out</button>
    {:else}
        <button
            class={$store.orientLoginBtns === orient.vert ? "btn-width" : ""}
            onclick={() => ($store.showModal = modalComp.LoginForm)}
        >
            Log In
        </button>
        <button
            class={$store.orientLoginBtns === orient.vert ? "btn-width" : ""}
            onclick={() => ($store.showModal = modalComp.SignUpForm)}
        >
            Sign Up
        </button>
    {/if}
</div>

<style>
    .btn-box {
        display: flex;
        gap: 0.5rem;
        align-items: center;
    }
    .vertical {
        flex-direction: column;
    }
    .btn-width {
        width: 89px;
    }
    .user-img {
        margin: 0;
        width: 40px;
        height: 40px;
    }
    .username-txt {
        margin: 0;
    }
</style>
