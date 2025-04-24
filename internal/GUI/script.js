panel = 1

function tabButtonPress(ref) {
    switch (ref.title) {
        case "Driving":
            panel = 1
        case "Settings":
            panel = 2
        case "Connections":
            panel = 3
    }
    document.querySelectorAll(".sidemenu_button").forEach((elem) => {
        elem.setAttribute("selected", "false")
    })
    document.querySelectorAll(".panel_container").forEach((elem) => {
        elem.setAttribute("selected", "false")
    })
    document.querySelector(`.panel_container[title="${ref.title}"`).setAttribute("selected", "true")
    ref.setAttribute("selected", "true")
}