<script lang="ts">

    let dialog: HTMLDialogElement;

    $effect(() => {
        // if (showModal === true) {
        //     dialog.showModal()
        // }
    })
</script>

<dialog bind:this={dialog}>
    <button>Close</button>
    <p>This modal dialog has a groovy backdrop!</p>
</dialog>
