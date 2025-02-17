type Orientation = {
    readonly horiz: string;
    readonly vert: string;
}

export const orient: Orientation = {
    horiz: "horizontal",
    vert: "vertical",
}

type ToastColor = {
    readonly red: string;
    readonly green: string;
    readonly yellow: string;
    readonly grey: string;
}

export const toastColor: ToastColor = {
    red: "red",
    green: "green",
    yellow: "yellow",
    grey: "grey",
}

type ModalComponents = {
    readonly SignUpForm: string;
    readonly VerificationForm: string;
    readonly LoginForm: string;
}

export const modalComp: ModalComponents = {
    SignUpForm: "SignUpForm",
    VerificationForm: "VerificationForm",
    LoginForm: "LoginForm",
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

export type ToastData = {
    color: string,
    text: string
}

export type Store = {
    baseFetchUrl: string;
    baseStorageUrl: string;
    orientLoginBtns: string;
    showLoginBtns: boolean;
    showModal: string;
    showToast: ToastData;
    user: User;
}