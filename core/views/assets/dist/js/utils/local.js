import { Routes } from "endpoints";

function LocalBrowserRefresh() {
    if (window.sessionStorage.getItem("remember-scroll")) {
        window.scroll({top: window.sessionStorage.getItem("remember-scroll")})
    }

    const eventSource = new EventSource(Routes.LocalBrowserRefresh)
    eventSource.onmessage = (msg) => {
        window.sessionStorage.setItem("remember-scroll", window.scrollY)
        window.location.reload()
    }

    eventSource.onerror = (err) => {
        console.error("eventSource failed:", err)

        eventSource.close();
    }
}

export { LocalBrowserRefresh }