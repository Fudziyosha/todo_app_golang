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
    const mobileSidebar = document.getElementById("appSidebar");
    const mobileSidebarToggle = document.getElementById("mobileSidebarToggle");
    const mobileSidebarBackdrop = document.getElementById("mobileSidebarBackdrop");
    const mobileAwareModals = document.querySelectorAll("#newListModal, #settingsModal");

    const setMobileSidebarState = (isOpen) => {
        if (!mobileSidebar || !mobileSidebarToggle || !mobileSidebarBackdrop) {
            return;
        }

        document.body.classList.toggle("mobile-sidebar-open", isOpen);
        mobileSidebarToggle.setAttribute("aria-expanded", String(isOpen));
        mobileSidebarBackdrop.hidden = !isOpen;
        mobileSidebarBackdrop.classList.toggle("is-visible", isOpen);
    };

    if (mobileSidebar && mobileSidebarToggle && mobileSidebarBackdrop) {
        mobileSidebarToggle.addEventListener("click", function () {
            setMobileSidebarState(true);
        });

        mobileSidebarBackdrop.addEventListener("click", function () {
            setMobileSidebarState(false);
        });

        mobileSidebar.querySelectorAll(".list-row__link, .nav-link").forEach((link) => {
            link.addEventListener("click", function () {
                if (window.innerWidth <= 767) {
                    setMobileSidebarState(false);
                }
            });
        });

        document.addEventListener("keydown", function (event) {
            if (event.key === "Escape") {
                setMobileSidebarState(false);
            }
        });

        window.addEventListener("resize", function () {
            if (window.innerWidth > 767) {
                setMobileSidebarState(false);
            }
        });

        mobileAwareModals.forEach((modal) => {
            modal.addEventListener("show.bs.modal", function () {
                setMobileSidebarState(false);
            });
        });
    }

    if (!panel || !toggleButton || !cancelButton) {
        return;
    }

    const listRows = Array.from(panel.querySelectorAll(".list-row"));
    const deleteForms = Array.from(panel.querySelectorAll(".list-row__delete-form"));

    const clearSelectedRows = () => {
        listRows.forEach((row) => row.classList.remove("list-row--selected"));
    };

    const getSelectedForms = () =>
        deleteForms.filter((form) => form.closest(".list-row")?.classList.contains("list-row--selected"));

    toggleButton.addEventListener("click", function () {
        if (panel.classList.contains("lists-panel--delete-mode")) {
            const selectedForms = getSelectedForms();

            if (selectedForms.length === 0) {
                return;
            }

            const submitForm = document.createElement("form");
            submitForm.method = "POST";
            submitForm.enctype = "multipart/form-data";
            submitForm.action = window.location.pathname + window.location.search;
            submitForm.style.display = "none";

            const actionInput = document.createElement("input");
            actionInput.type = "hidden";
            actionInput.name = "action";
            actionInput.value = "delete_list";
            submitForm.appendChild(actionInput);

            selectedForms.forEach((form) => {
                const listIdInput = form.querySelector('input[name="list_id"]');

                if (!listIdInput || !listIdInput.value) {
                    return;
                }

                const selectedInput = document.createElement("input");
                selectedInput.type = "hidden";
                selectedInput.name = "list_id";
                selectedInput.value = listIdInput.value;
                submitForm.appendChild(selectedInput);
            });

            document.body.appendChild(submitForm);
            submitForm.submit();

            return;
        }

        panel.classList.add("lists-panel--delete-mode");
    });

    cancelButton.addEventListener("click", function () {
        panel.classList.remove("lists-panel--delete-mode");
        clearSelectedRows();
    });

    listRows.forEach((row) => {
        row.addEventListener("click", function (event) {
            if (!panel.classList.contains("lists-panel--delete-mode")) {
                return;
            }

            event.preventDefault();
            row.classList.toggle("list-row--selected");
        });
    });
});
