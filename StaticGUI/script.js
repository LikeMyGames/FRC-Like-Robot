const socket = new WebSocket("ws://localhost:8080/ws");

socket.onmessage = (event) => {
    let data = JSON.parse(event.data);
    // console.log("Backend says:", data);
    if (data.system_logger != undefined) {
        systemLog(data.system_logger)
    }
    if (data.robot_status != undefined) {
        refreshRobotStatus(data.robot_status)
    }
};

socket.onopen = () => {
    socket.send(JSON.stringify({ message: "GUI Connected" }));
};

var panel = 1
var filter = ["all"];
var logger_settings = {
    pause: false
}

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

function changeSystemLoggerFilter(filterVal) {
    switch (filterVal) {
        case "all":
            filter = ["all"]
            document.querySelector(`.system_logger_filter_option[title="log"]`).setAttribute("selected", "false")
            document.querySelector(`.system_logger_filter_option[title="success"]`).setAttribute("selected", "false")
            document.querySelector(`.system_logger_filter_option[title="warn"]`).setAttribute("selected", "false")
            document.querySelector(`.system_logger_filter_option[title="error"]`).setAttribute("selected", "false")
            document.querySelector(`.system_logger_filter_option[title="all"]`).setAttribute("selected", "true")
            break
        case "log":
            if (!filter.includes(filterVal)) {
                filter.push(filterVal)
                document.querySelector(`.system_logger_filter_option[title="log"]`).setAttribute("selected", "true")
                document.querySelector(`.system_logger_filter_option[title="all"]`).setAttribute("selected", "false")
            } else {
                let index = filter.indexOf(filterVal)
                filter = filter.slice(0, index).concat(filter.slice(index + 1))
                document.querySelector(`.system_logger_filter_option[title="log"]`).setAttribute("selected", "false")
            }
            break
        case "success":
            if (!filter.includes(filterVal)) {
                filter.push(filterVal)
                document.querySelector(`.system_logger_filter_option[title="success"]`).setAttribute("selected", "true")
                document.querySelector(`.system_logger_filter_option[title="all"]`).setAttribute("selected", "false")
            } else {
                let index = filter.indexOf(filterVal)
                filter = filter.slice(0, index).concat(filter.slice(index + 1))
                document.querySelector(`.system_logger_filter_option[title="success"]`).setAttribute("selected", "false")
            }
            break
        case "warn":
            if (!filter.includes(filterVal)) {
                filter.push(filterVal)
                document.querySelector(`.system_logger_filter_option[title="warn"]`).setAttribute("selected", "true")
                document.querySelector(`.system_logger_filter_option[title="all"]`).setAttribute("selected", "false")
            } else {
                let index = filter.indexOf(filterVal)
                filter = filter.slice(0, index).concat(filter.slice(index + 1))
                document.querySelector(`.system_logger_filter_option[title="warn"]`).setAttribute("selected", "false")
            }
            break
        case "error":
            if (!filter.includes(filterVal)) {
                filter.push(filterVal)
                document.querySelector(`.system_logger_filter_option[title="error"]`).setAttribute("selected", "true")
                document.querySelector(`.system_logger_filter_option[title="all"]`).setAttribute("selected", "false")
            } else {
                let index = filter.indexOf(filterVal)
                filter = filter.slice(0, index).concat(filter.slice(index + 1))
                document.querySelector(`.system_logger_filter_option[title="error"]`).setAttribute("selected", "false")
            }
            break
    }
    if (filter.includes("log") && filter.includes("success") && filter.includes("warn") && filter.includes("error")) {
        changeSystemLoggerFilter("all")
    }
}

function systemLog(system_logger) {
    if (logger_settings.pause) {
        return
    }
    let appendStr = ""
    switch (system_logger.type) {
        case "log":
            appendStr = "LOG - "
            break;
        case "success":
            appendStr = "SUCCESS - "
            break;
        case "warn":
            appendStr = "WARNING - "
            break;
        case "error":
            appendStr = "ERROR - "
            break;
    }
    document.querySelector(".system_logger_display_text").appendChild(elementFromHTML(`
        <h4 class="${system_logger.type}">${appendStr}${system_logger.message}</h4>
    `));
}

function systemLogPause() {
    logger_settings.pause = !logger_settings.pause
}

function systemLogRefresh() {
    location.reload()
}

function systemLogClear() {
    document.querySelector(".system_logger_display_text").innerHTML = "";
}

function elementFromHTML(html) {
    let template = document.createElement("template");
    template.innerHTML = html.trim();
    return template.content.firstElementChild;
}

function refreshRobotStatus(robot_status) {
    switch (robot_status.type) {
        case "comms":
            document.querySelector(`.robot_status_indicator_display[for="comms"]`).setAttribute("value", robot_status.value);
            break;
        case "code":
            document.querySelector(`.robot_status_indicator_display[for="code"]`).setAttribute("value", robot_status.value);
            break;
        case "sticks":
            document.querySelector(`.robot_status_indicator_display[for="sticks"]`).setAttribute("value", robot_status.value);
            break;
    }
}