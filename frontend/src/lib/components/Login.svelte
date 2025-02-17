<script lang="ts">
    import { S, emptyUser } from "../../store.svelte";
    import { modalComp, orient } from "../../types";

    function logOut() {
        S.user = emptyUser;
    }
</script>

<div class="btn-box {S.orientLoginBtns === orient.vert ? orient.vert : ''}">
    {#if S.user?.token}
        <p class="username-txt">{S.user.username}</p>
        <img
            class="user-img"
            src={S.baseStorageUrl + S.user.imageUrl}
            alt="user"
        />
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
