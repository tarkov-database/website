const registerTabs = () => {
    const tablinks = document.getElementsByClassName("tab-btn");
    for (const el of tablinks)
        el.addEventListener("click", function (this: HTMLElement) {
            const cl = "active";
            document.querySelector(`.tab.${cl}`)?.classList.remove(cl);
            document.querySelector(`.tab-btn.${cl}`)?.classList.remove(cl);

            document
                .getElementById(this.dataset.tab?.toString() || "")
                ?.classList.add(cl);
            this.classList.add(cl);
        });
};

const initImageView = () => {
    const opener = document.getElementsByClassName(
        "open-image-view"
    ) as HTMLCollectionOf<HTMLElement>;
    if (opener.length === 0) return;

    function remove(element: HTMLElement) {
        element.remove();
        document.body.classList.remove("modal-opened");
    }

    for (const el of opener) {
        el.addEventListener("click", function (this: HTMLElement) {
            const url = this.dataset.largeUrl;
            if (url === undefined) return;

            const view = document.createElement("div");
            view.classList.add("image-view");
            view.addEventListener("click", function (this: HTMLDivElement) {
                remove(this);
            });
            document.body.addEventListener(
                "keydown",
                function (this: HTMLElement, e: KeyboardEvent) {
                    if (e.key == "Escape") {
                        remove(view);
                    }
                }
            );

            const img = document.createElement("img");
            img.src = url;
            view.appendChild(img);

            document.body.classList.add("modal-opened");
            document.body.appendChild(view);
        });

        el.style.cursor = "pointer";
    }
};

const sortTables = () => {
    const tables = document.querySelectorAll<HTMLTableHeaderCellElement>(
        ".sort-table.client-sort thead th"
    );
    if (tables.length === 0) return;

    const getCellValue = (tr: HTMLTableRowElement, idx: number) =>
        (tr.children[idx] as HTMLElement).dataset.value ||
        (tr.children[idx] as HTMLElement).innerText ||
        tr.children[idx].textContent ||
        "";

    const comparer =
        (idx: number, asc: boolean) =>
        (a: HTMLTableRowElement, b: HTMLTableRowElement) =>
            ((v1: string, v2: string) =>
                !isNaN(Number(v1)) && !isNaN(Number(v2))
                    ? Number(v1) - Number(v2)
                    : v1?.toString().localeCompare(v2 as string))(
                getCellValue(asc ? a : b, idx),
                getCellValue(asc ? b : a, idx)
            );

    for (const th of tables) {
        th.addEventListener("click", function (this: HTMLElement) {
            const table = this.closest("table");
            const body = table?.getElementsByTagName("tbody")[0];

            const clSorted = "sorted-by";
            const clAsc = "up";
            const clDesc = "down";

            let isASC = this.classList.contains(clAsc);

            Array.from<HTMLTableRowElement>(
                body?.getElementsByTagName("tr") || []
            )
                .sort(
                    comparer(
                        Array.from(this.parentNode?.children || []).indexOf(th),
                        (isASC = !isASC)
                    )
                )
                .forEach((tr) => body?.appendChild(tr));

            const sorted = table?.getElementsByClassName(clSorted)[0];
            if (th !== sorted) {
                sorted?.classList.remove(clSorted);
                this.classList.add(clSorted);
            }

            if (isASC) {
                if (this.classList.contains(clDesc)) {
                    this.classList.replace(clDesc, clAsc);
                } else {
                    this.classList.add(clAsc);
                }
            } else {
                if (this.classList.contains(clAsc)) {
                    this.classList.replace(clAsc, clDesc);
                } else {
                    this.classList.add(clDesc);
                }
            }
        });

        th.style.cursor = "pointer";
    }
};

const initListFilter = () => {
    const el = document.getElementById("listFilter");
    if (el === null) return;

    const changeFilter = function (this: HTMLSelectElement) {
        const params = new URLSearchParams(window.location.search);
        const name = this.name,
            value = this.value;

        if (!value || value === "all") {
            params.delete(name);
        } else {
            params.set(name, value);
        }

        const pageKey = "p";
        if (params.has(pageKey)) params.delete(pageKey);

        window.location.search = params.toString();
    };

    const ul = el.querySelectorAll("ul > li");
    for (const li of ul) {
        const sel = li.getElementsByTagName("select")[0];
        sel.addEventListener("change", changeFilter);
    }
};

(async () => {
    registerTabs();
    sortTables();
    initListFilter();
    initImageView();
})();

export {};
