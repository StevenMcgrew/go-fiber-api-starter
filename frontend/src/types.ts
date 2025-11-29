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
    readonly OnlyModalText: string;
    readonly SignUpForm: string;
    readonly SignupVerificationForm: string;
    readonly UpdateEmailVerificationForm: string;
    readonly LoginForm: string;
    readonly LoggingOutMsg: string;
    readonly ForgotPasswordForm: string;
    readonly ResetPasswordForm: string;
}

export const modalComp: ModalComponents = {
    OnlyModalText: "OnlyModalText",
    SignUpForm: "SignUpForm",
    SignupVerificationForm: "SignupVerificationForm",
    UpdateEmailVerificationForm: "UpdateEmailVerificationForm",
    LoginForm: "LoginForm",
    LoggingOutMsg: "LoggingOutMsg",
    ForgotPasswordForm: "ForgotPasswordForm",
    ResetPasswordForm: "ResetPasswordForm",
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
    modalText: string;
    showToast: ToastData;
    newEmailAddress: string;
    user: User;
}