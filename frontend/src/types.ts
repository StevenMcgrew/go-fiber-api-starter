type Orientation = {
    readonly horiz: string;
    readonly vert: string;
}

export const orient: Orientation = {
    horiz: "horizontal",
    vert: "vertical",
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

export type Store = {
    baseFetchUrl: string;
    baseStorageUrl: string;
    orientLoginBtns: string;
    showLoginBtns: boolean;
    showModal: string;
    user: User;
}