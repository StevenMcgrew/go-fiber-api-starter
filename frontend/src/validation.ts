export const getEmailValidationWarnings = (email: string) => {
    let warnings = ""
    const len = email.length
    if (len < 3 || len > 320) {
        warnings += "Must be 3 to 320 characters. "
    }
    if (!email.includes("@")) {
        warnings += "Must include @ symbol."
    }
    return warnings
}

export const getUsernameValidationWarnings = (username: string) => {
    let warnings = ""
    const len = username.length
    if (len < 3 || len > 20) {
        warnings += "Must be 3 to 20 characters. "
    }
    if (!hasAllowedUsernameChars(username)) {
        warnings += "Must be alphanumeric or underscore."
    }
    return warnings
}

export const getPasswordValidationWarnings = (password: string) => {
    let warnings = ""
    const len = password.length
    if (len < 8 || len > 72) {
        warnings += "Must be 8 to 72 characters. "
    }
    return warnings
}

function hasAllowedUsernameChars(username: string): boolean {
    // Regular expression to match only alphanumeric characters and underscores
    const usernameRegex = /^[a-zA-Z0-9_]+$/;
    return usernameRegex.test(username);
  }