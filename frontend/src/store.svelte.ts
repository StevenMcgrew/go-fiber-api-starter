import type { Store } from "./types";
import { orient } from "./types";
import VerificationForm from "./lib/components/VerificationForm.svelte";

export const S: Store = $state({
    baseFetchUrl: "http://127.0.0.1:8080/api/v1",
    baseStorageUrl: "http://127.0.0.1:8080/temp-storage",
    orientLoginBtns: orient.horiz,
    showLoginBtns: true,
    showLoginForm: false,
    showSignUpForm: false,
    showModal: VerificationForm,
    user: null,
})