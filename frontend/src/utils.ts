export const debounce = <F extends (...args: any[]) => any>(
    func: F,
    wait: number,
    immediate: boolean = false
): ((this: ThisParameterType<F>, ...args: Parameters<F>) => void) => {
    let timeout: ReturnType<typeof setTimeout> | null = null;

    return function (this: ThisParameterType<F>, ...args: Parameters<F>): void {
        const context = this;

        const later = () => {
            timeout = null;
            if (!immediate) func.apply(context, args);
        };

        const callNow = immediate && !timeout;

        if (timeout !== null) {
            clearTimeout(timeout);
        }

        timeout = setTimeout(later, wait);

        if (callNow) func.apply(context, args);
    };
};

export const getImgSrc = (baseImgUrl: string, userImageUrl: string): string => {
    if (userImageUrl) {
        return baseImgUrl + userImageUrl;
    }
    return baseImgUrl + "/default-profile-pic.png"
}

export function clickOutside(node: Node, callback: (node: Event) => void) {
    const handleClick = (event: Event) => {
        if (node &&
            !node.contains(event.target as Node) &&
            !event.defaultPrevented)
        {
            callback(event);
        }
    };

    document.addEventListener("click", handleClick);

    return {
        destroy() {
            document.removeEventListener("click", handleClick);
        },
    };
}