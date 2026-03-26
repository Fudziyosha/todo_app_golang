const todoNewText = document.getElementById("todoNewText");
if (todoNewText) {
    const resizeTodoNewText = () => {
        const maxHeight = 160;

        todoNewText.style.setProperty("height", "0px", "important");
        todoNewText.style.setProperty("height", Math.min(todoNewText.scrollHeight, maxHeight) + "px", "important");
        todoNewText.style.overflowY = todoNewText.scrollHeight > maxHeight ? "auto" : "hidden";
    };

    ["input", "change", "focus"].forEach((eventName) => {
        todoNewText.addEventListener(eventName, resizeTodoNewText);
    });

    resizeTodoNewText();
    requestAnimationFrame(resizeTodoNewText);
    window.addEventListener("load", resizeTodoNewText);
}

document.querySelectorAll(".todo-inline-input.auto-grow").forEach((textarea) => {
    const resize = () => {
        textarea.style.setProperty("height", "auto", "important");
        textarea.style.setProperty("height", Math.min(textarea.scrollHeight, 200) + "px", "important");
        textarea.style.overflowY = textarea.scrollHeight > 200 ? "auto" : "hidden";
    };

    textarea.addEventListener("input", resize);
    textarea.addEventListener("change", resize);
    resize();
});

document.addEventListener("DOMContentLoaded", function () {
    const panel = document.getElementById("listsPanel");
    const toggleButton = document.getElementById("toggleDeleteListsBtn");
    const cancelButton = document.getElementById("cancelDeleteListsBtn");

    if (!panel || !toggleButton || !cancelButton) {
        return;
    }

    toggleButton.addEventListener("click", function () {
        panel.classList.add("lists-panel--delete-mode");
    });

    cancelButton.addEventListener("click", function () {
        panel.classList.remove("lists-panel--delete-mode");
    });
});