<script lang="ts">
    import { link } from "svelte-spa-router";
    import { store, clearStorageAndReload } from "../../store.svelte";
    import { modalComp, orient } from "../../types";
    import { getImgSrc } from "../../utils";

    function logOut() {
        $store.showModal = modalComp.LoggingOutMsg;
        clearStorageAndReload();
    }

</script>

<div
    class="btn-box {$store.orientLoginBtns === orient.vert ? orient.vert : ''}"
>
    {#if $store.user?.token}
        <p class="username-txt">{$store.user.username}</p>
        <a use:link href="/account">
            <img
                class="user-img"
                src={getImgSrc($store.baseStorageUrl, $store.user.imageUrl)}
                alt="user"
            />
        </a>
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
    .username-txt {
        margin: 0;
    }
</style>
