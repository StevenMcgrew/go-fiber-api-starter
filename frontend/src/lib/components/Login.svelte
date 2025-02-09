<script lang="ts">
    import ModalDialog from "./ModalDialog.svelte";
    import { orientLoginBtns, Orient } from "../../store";

    let showLoginForm: any = $state()

    function loginClicked(e: MouseEvent) {
        showLoginForm = true
    }
</script>

<div class="btn-box {$orientLoginBtns === Orient.Vert ? Orient.Vert : ''}">
    <a href="/login">
        <button class={$orientLoginBtns === Orient.Vert ? "btn-width" : ""} onclick={loginClicked}>Log In</button>
    </a>
    <a href="/signup">
        <button class={$orientLoginBtns === Orient.Vert ? "btn-width" : ""}>Sign Up</button>
    </a>
</div>

<ModalDialog/>

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
