import type { Component } from "svelte";

type Orientation = {
    readonly horiz: string;
    readonly vert: string;
}

export const orient: Orientation = {
    horiz: "horizontal",
    vert: "vertical",
}

export type User = {
    token: string;
    id: number;
    email: string;
    username: string;
    role: string;
    status: string;
    imageUrl: string;
}

export type Store = {
    baseFetchUrl: string;
    baseStorageUrl: string;
    orientLoginBtns: string;
    showLoginBtns: boolean;
    showLoginForm: boolean;
    showSignUpForm: boolean;
    showModal: Component | null;
    user: User | null;
}