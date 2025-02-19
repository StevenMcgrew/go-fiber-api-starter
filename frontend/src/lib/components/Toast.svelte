<script lang="ts">
    import { store } from "../../store.svelte";

    $effect(() => {
        if ($store.showToast.text) {
            // After 2.9 sec (fadeout finished), reset toast
            setTimeout(function () {
                $store.showToast.text = "";
                $store.showToast.color = "grey";
            }, 2900);
        }
    });
</script>

<div
    class="toast {$store.showToast.text ? 'show' : ''} {$store.showToast.color}"
>
    {$store.showToast.text}
</div>

<style>
    .toast {
        visibility: hidden;
        min-width: 250px;
        margin-left: -125px;
        text-align: center;
        border-radius: 6px;
        padding: 16px;
        position: fixed;
        z-index: 1;
        left: 50%;
        bottom: 30px;
        border-width: 1px;
        border-style: solid;
    }

    .toast.show {
        visibility: visible;
        /* animate 0.5s, delay 2.5s */
        -webkit-animation:
            fadein 0.5s,
            fadeout 0.5s 2.5s;
        animation:
            fadein 0.5s,
            fadeout 0.5s 2.5s;
    }

    .red {
        color: #58151c;
        background-color: #f8d7da;
        border-color: #f1aeb5;
    }

    .green {
        color: #0a3622;
        background-color: #d1e7dd;
        border-color: #a3cfbb;
    }

    .yellow {
        color: #664d03;
        background-color: #fff3cd;
        border-color: #ffe69c;
    }

    .grey {
        color: #2b2f32;
        background-color: #e2e3e5;
        border-color: #c4c8cb;
    }

    @-webkit-keyframes fadein {
        from {
            bottom: 0;
            opacity: 0;
        }
        to {
            bottom: 30px;
            opacity: 1;
        }
    }

    @keyframes fadein {
        from {
            bottom: 0;
            opacity: 0;
        }
        to {
            bottom: 30px;
            opacity: 1;
        }
    }

    @-webkit-keyframes fadeout {
        from {
            bottom: 30px;
            opacity: 1;
        }
        to {
            bottom: 0;
            opacity: 0;
        }
    }

    @keyframes fadeout {
        from {
            bottom: 30px;
            opacity: 1;
        }
        to {
            bottom: 0;
            opacity: 0;
        }
    }
</style>
