<script lang="ts">
    import Router from "svelte-spa-router";
    import { routes } from "./routes";
    import { store, emptyUser } from "./store.svelte";
    import Header from "./lib/components/Header.svelte";
    import ModalDialog from "./lib/components/ModalDialog.svelte";
    import StoreView from "./lib/components/StoreView.svelte";
    import Toast from "./lib/components/Toast.svelte";

    function onbeforeunload() {
        if (!$store.stayLoggedIn) {
            $store.user = emptyUser
        }
    }

</script>

<svelte:window {onbeforeunload} />
<Header />
<main>
    <Router {routes} />
</main>
<ModalDialog />
<Toast />
<StoreView />
