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
    orientLoginBtns: string;
    showLoginBtns: boolean;
    showLoginForm: boolean;
    showSignUpForm: boolean;
    showModal: Component | null;
}

export const S: Store = $state({
    orientLoginBtns: orient.horiz,
    showLoginBtns: true,
    showLoginForm: false,
    showSignUpForm: false,
    showModal: null
})