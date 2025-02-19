<script lang="ts">
    import { store } from "../../store.svelte";

    let showStore = $state(false);
    let flattenedStore = $state([]);
    let positions = ["top-left", "top-right", "bottom-right", "bottom-left"];
    let position = $state("bottom-right");

    function handleKeydown(e: KeyboardEvent) {
        if (e.ctrlKey && e.shiftKey && e.key === "?") {
            showStore = !showStore;
            return;
        }
        if (e.ctrlKey && e.shiftKey && e.code === "ArrowRight") {
            let curIdx = positions.indexOf(position);
            if (curIdx > positions.length - 2) {
                position = positions[0];
                return;
            }
            position = positions[curIdx + 1];
        }
    }

    function flattenObject(obj: any, parentKey = "", result: any = {}) {
        for (let key in obj) {
            if (obj.hasOwnProperty(key)) {
                const newKey = parentKey ? `${parentKey}.${key}` : key;
                if (
                    typeof obj[key] === "object" &&
                    obj[key] !== null &&
                    !Array.isArray(obj[key])
                ) {
                    flattenObject(obj[key], newKey, result); // Recursively flatten nested objects
                } else {
                    result[newKey] = obj[key]; // Add the key-value pair to the result
                }
            }
        }
        return result;
    }

    $effect(() => {
        flattenedStore = flattenObject($store);
        localStorage.setItem("store", JSON.stringify($store));
    });
</script>

<svelte:window on:keydown={handleKeydown} />
<div class="store-view-box {showStore ? 'show' : 'hide'} {position}">
    <p><kbd>Ctrl</kbd><kbd>Shift</kbd><kbd>🡒</kbd> to move panel</p>
    <button onclick={() => showStore = false}>×</button>
    <ul>
        {#each Object.entries(flattenedStore) as item}
            <li>{item[0]}: <span class="success-text">{item[1]}</span></li>
        {/each}
    </ul>
</div>

<style>
    .store-view-box {
        position: fixed;
        background-color: rgba(0, 0, 0, 0.8);
        z-index: 1000;
        font-size: 10px;
        padding: 8px;
        border: solid 1px dimgrey;
        border-radius: 4px;
        max-width: 506px;
    }
    button {
        position: absolute;
        right: 6px;
        top: 6px;
        padding: 0px 5px 3px 5px;
        font-size: 1.5rem;
        line-height: 23px;
    }
    .top-left {
        top: 0;
        left: 0;
    }
    .top-right {
        top: 0;
        right: 0;
    }
    .bottom-right {
        bottom: 0;
        right: 0;
    }
    .bottom-left {
        bottom: 0;
        left: 0;
    }
    kbd {
        margin-right: 1px;
        background-color: #4f4f4f;
    }
    p {
        display: inline-block;
        margin-bottom: 4px;
    }
    ul {
        list-style: none;
        padding-left: 0;
        margin: 0;
    }
    .show {
        display: block;
    }
    .hide {
        display: none;
    }
</style>
