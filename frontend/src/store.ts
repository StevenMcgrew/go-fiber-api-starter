import { writable, type Writable } from "svelte/store";

type Orientation = {
    readonly Horiz: string;
    readonly Vert: string;
}

export const Orient: Orientation = {
    Horiz: "horizontal",
    Vert: "vertical"
}

export const orientLoginBtns: Writable<string> = writable(Orient.Horiz)
export const showLoginBtns: Writable<boolean> = writable(true)