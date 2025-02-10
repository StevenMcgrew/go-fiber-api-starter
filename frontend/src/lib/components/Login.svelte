<script lang="ts">
    import { S, orient } from "../../store.svelte";
    import LoginForm from "./LoginForm.svelte";
    import SignUpForm from "./SignUpForm.svelte";
</script>

<div class="btn-box {S.orientLoginBtns === orient.vert ? orient.vert : ''}">
    <button class={S.orientLoginBtns === orient.vert ? "btn-width" : ""}
            onclick={() => (S.showModal = LoginForm)}>
        Log In
    </button>
    <button class={S.orientLoginBtns === orient.vert ? "btn-width" : ""}
            onclick={() => (S.showModal = SignUpForm)}>
        Sign Up
    </button>
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
