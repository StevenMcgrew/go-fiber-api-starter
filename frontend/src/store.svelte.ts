import type { Component } from "svelte";

type Orientation = {
    readonly horiz: string;
    readonly vert: string;
}

export const orient: Orientation = {
    horiz: "horizontal",
    vert: "vertical"
}

type Store = {
    baseFetchUrl: string;
    orientLoginBtns: string;
    showLoginBtns: boolean;
    showLoginForm: boolean;
    showSignUpForm: boolean;
    showModal: Component | null;
}

export const S: Store = $state({
    baseFetchUrl: "http://127.0.0.1:8080/api/v1",
    orientLoginBtns: orient.horiz,
    showLoginBtns: true,
    showLoginForm: false,
    showSignUpForm: false,
    showModal: null
})