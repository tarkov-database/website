interface Item {
    id: string;
    name: string;
    shortName?: string;
    parent?: string;
    type: ItemType;
}

enum ItemType {
    Item,
    Location,
}

interface SocketRequest {
    id: number;
    term: string;
    filter?: {
        item?: string;
        location?: string;
    };
    items: boolean;
    locations: boolean;
}

interface SocketResponse {
    id: number;
    items?: Array<Item>;
    error?: Error;
}

const errErrorClosure = new Error("Search socket was closed with an error");

const wait = (time: number) =>
    new Promise<void>((resolve) => setTimeout(() => resolve(), time));

const regexMeta = new RegExp(/[.*+?^${}()|[\]\\]/, "g");
const quoteMeta = (str: string) => str.replace(regexMeta, "\\$&");

const regexFilterPrefix = new RegExp(
    /(?:^\s*(?<key>\w+):)(?:\s*(?<value>\w+)\s)?(?:\s*(?<term>.*))/,
    "i"
);
const getFilterPrefix = (str: string) => regexFilterPrefix.exec(str);

const showElement = (el: HTMLElement | null) =>
    el?.classList.replace("hide", "show");
const hideElement = (el: HTMLElement | null) =>
    el?.classList.replace("show", "hide");

export const initSearchSocket = async (
    element: HTMLFormElement,
): Promise<void> => {
    const input = element.querySelector<HTMLInputElement>(
        'input[type="search"]'
    );
    const suggBox = element.querySelector<HTMLElement>(".suggestion");
    const suggInline = element.querySelector<HTMLElement>(".inline-suggestion");

    let socket: WebSocket;
    let idleTimeout: number | undefined;
    let noReconnect = false;
    let qCount = 0;

    let lastTerm = "";

    const updateBoxSuggestions = (items: Array<Item>, term: string) => {
        const newUl = document.createElement("ul");

        for (const item of items) {
            const a = document.createElement("a");
            const li = document.createElement("li");
            const text = document.createElement("span");
            const div = document.createElement("div");

            const highlightMatches = (str: string) => {
                const regex = new RegExp(`(${term})`, "gi");
                const matches = str.match(regex);
                if (matches)
                    for (const m of matches) {
                        str = str.replace(m, `<b>${m}</b>`);
                    }

                return str;
            };

            text.innerHTML = highlightMatches(item.name);
            a.title = item.name;

            switch (item.type) {
                case 0:
                    a.href = `/item/${item.parent}/${item.id}`;
                    div.className = `icon ${item.parent}`;
                    break;
            }

            div.innerHTML = "&nbsp;";

            a.addEventListener(
                "keydown",
                function (this: HTMLAnchorElement, e: KeyboardEvent) {
                    const current = this.parentNode as HTMLLIElement;

                    let next: Element | null = null;
                    switch (e.key) {
                        case "ArrowUp":
                            e.preventDefault();
                            next = current.previousElementSibling;
                            if (!next) input?.focus();
                            break;
                        case "ArrowDown":
                            e.preventDefault();
                            next = current.nextElementSibling;
                            if (!next) input?.focus();
                            break;
                        case "Escape":
                            current.blur();
                            hideElement(suggBox);
                            break;
                        default:
                            return;
                    }

                    next?.querySelector("a")?.focus();
                }
            );

            a.appendChild(div);
            a.append(text);
            li.appendChild(a);
            newUl.appendChild(li);
        }

        const ul = suggBox?.querySelector("ul");
        if (ul != null) suggBox?.replaceChild(newUl, ul);
    };

    const updateInlineSuggestion = (val: string) => {
        if (suggInline === null) return;
        suggInline.innerHTML = val?.replace(" ", "&nbsp;");
        suggInline.dataset.value = val;
    };

    const filters = new Filters();

    const openSocket = (
        msgListener: (this: WebSocket, ev: MessageEvent<any>) => any
    ) => {
        const host = window.location.host;
        const path = "search/ws";

        const proto = window.location.protocol === "https:" ? "wss" : "ws";
        if (proto === "ws") console.warn("Insecure WebSocket protocol is used");

        let socket;
        try {
            socket = new WebSocket(`${proto}://${host}/${path}`);
        } catch (err) {
            return Promise.reject(err);
        }

        socket.addEventListener("open", () => {
            console.info("Search socket opened");
        });

        socket.addEventListener("close", (e) => {
            if (e.wasClean) {
                const msg = "Search socket closed";
                console.info(e.reason ? `${msg}: ${e.reason}` : msg);
                if (e.code !== 1000) noReconnect = true;
                if (e.code === 1012)
                    wait(1000).then(() => (noReconnect = false));
            } else {
                const msg = "Search socket closed unexpectedly";
                console.error(e.reason ? `${msg}: ${e.reason}` : msg);
                noReconnect = true;
            }
        });

        socket.addEventListener("error", (e) => {
            console.error("Search socket error: %s", e);
        });

        socket.addEventListener("message", msgListener);

        return socket;
    };

    const connect = async () => {
        if (!socket || socket.readyState > 1) {
            if (noReconnect) return Promise.reject(errErrorClosure);
            try {
                socket = await openSocket(onMessage);
            } catch (err) {
                return Promise.reject(err);
            }
        }

        while (socket.readyState === 0) await wait(50);

        clearTimeout(idleTimeout);
        idleTimeout = setTimeout(
            () => socket.close(1000, "Idle timeout"),
            30 * 1000
        );

        return;
    };

    const onMessage = function (this: WebSocket, event: MessageEvent<any>) {
        const data: SocketResponse = JSON.parse(event.data);

        if (data.error !== null) {
            console.error(data.error);
            return;
        }

        if (data.id < qCount) return;

        if (!data.items || data.items.length === 0) {
            hideElement(suggBox);
            return;
        }

        const term = quoteMeta(lastTerm).replace(" ", "|");

        updateBoxSuggestions(data.items, term);
        showElement(suggBox);
    };

    const onInput = async function (this: HTMLInputElement) {
        const count = qCount;

        try {
            await connect();
        } catch (err) {
            if (err !== errErrorClosure) console.error(err);
            return;
        }

        const val = this.value;

        this.setCustomValidity("");

        if (!this.validity.valid) {
            hideElement(suggBox);
            hideElement(suggInline);
            return;
        }

        const keyCompletion = (val: string) => {
            const match = filters.firstMatch(val);
            if (match) {
                const v = this.value.replace(val, `${match}:`);
                updateInlineSuggestion(v);
                showElement(suggInline);
            } else {
                hideElement(suggInline);
            }
        };

        const valueCompletion = (key: string, val: string) => {
            const match = filters.firstMatchIn(key, val);
            if (match) {
                const v = this.value.replace(val, `${match} `);
                updateInlineSuggestion(v);
                showElement(suggInline);
            } else {
                hideElement(suggInline);
            }
        };

        const prefix = getFilterPrefix(val);

        if (prefix?.groups && filters.contains(prefix.groups.key)) {
            const key = prefix.groups["key"],
                value = prefix.groups["value"],
                term = prefix.groups["term"];

            if (!value && term.length >= 2) {
                valueCompletion(key, term);
            } else {
                hideElement(suggInline);
            }
        } else {
            keyCompletion(val.trim());
        }

        let currentTerm = "";

        const filter: { [key: string]: string } = {};

        if (prefix?.groups) {
            const key = prefix.groups["key"],
                value = prefix.groups["value"],
                term = prefix.groups["term"];

            if (!filters.contains(key) || !filters.includes(key, value)) {
                this.setCustomValidity("Filter is invalid");
                return;
            }

            currentTerm = term;
            filter[key] = value;
        } else {
            currentTerm = val.trim();
        }

        if (currentTerm.length < 3) {
            this.setCustomValidity("Term is too short (min. 3 characters)");
            return;
        }

        if (currentTerm.length > 32) {
            this.setCustomValidity("Term is too long (max. 32 characters)");
            return;
        }

        if (currentTerm === lastTerm) {
            showElement(suggBox);
            return;
        }

        if (count < qCount) return;

        qCount++;

        const data: SocketRequest = {
            id: qCount,
            term: currentTerm,
            items: !filter.location,
            locations: !filter.item,
        };

        if (Object.keys(filter).length) data.filter = filter;

        socket.send(JSON.stringify(data));

        lastTerm = currentTerm;
    };

    const onInputFocusIn = async function (this: HTMLInputElement) {
        try {
            await connect();
        } catch (err) {
            if (err !== errErrorClosure) console.error(err);
            return;
        }

        if (this.dataset.valid === "true" && suggBox?.querySelector("ul > li"))
            showElement(suggBox);
    };

    const onInputKeydown = function (this: HTMLElement, event: KeyboardEvent) {
        const selFirstBoxSugg = () => {
            const next = suggBox?.querySelector<HTMLLinkElement>(
                "ul > li:first-child > a"
            );
            if (next) next.focus();
        };
        const selLastBoxSugg = () => {
            const next = suggBox?.querySelector<HTMLLinkElement>(
                "ul > li:last-child > a"
            );
            if (next) next.focus();
        };
        const applyInlineSugg = () => {
            if (input && suggInline) {
                const { dataset, classList } = suggInline;
                if (classList.contains("show") && dataset.value)
                    input.value = dataset.value;
            }
        };

        switch (event.key) {
            case "ArrowDown":
                event.preventDefault();
                selFirstBoxSugg();
                break;
            case "ArrowUp":
                event.preventDefault();
                selLastBoxSugg();
                break;
            case "Tab":
                event.preventDefault();
                applyInlineSugg();
                hideElement(suggInline);
                break;
            case "Escape":
                hideElement(suggBox);
                break;
        }
    };

    input?.addEventListener("keydown", onInputKeydown);
    input?.addEventListener("focusin", onInputFocusIn);
    input?.addEventListener("input", onInput);

    const onDocKeydown = function (this: Document, event: KeyboardEvent) {
        if (document.activeElement?.nodeName === "INPUT") return;
        switch (event.key) {
            case "F3":
            case "/":
            case "s":
                event.preventDefault();
                input?.focus();
        }
    };

    document.addEventListener("keydown", onDocKeydown);
};

type FilterType = "item" | "location";

interface Filter {
    type: FilterType;
    available: boolean;
    values: string[];
}

class Filters {
    private data: Map<string, Filter>;

    constructor() {
        this.data = new Map();

        this.data.set("item", {
            type: "item",
            available: true,
            values: [
                "ammunition",
                "armor",
                "backpack",
                "barrel",
                "barter",
                "bipod",
                "charge",
                "clothing",
                "common",
                "container",
                "device",
                "firearm",
                "food",
                "foregrip",
                "gasblock",
                "goggles",
                "grenade",
                "handguard",
                "headphone",
                "key",
                "launcher",
                "magazine",
                "map",
                "medical",
                "melee",
                "auxiliary",
                "mod-other",
                "money",
                "mount",
                "muzzle",
                "pistolgrip",
                "receiver",
                "sight-special",
                "sight",
                "stock",
                "tacticalrig",
            ],
        });
    }

    contains(key: string): boolean {
        return this.data.get(key)?.available ?? false;
    }

    includes(key: string, value: string): boolean {
        const entry = this.data.get(key);
        return (
            (entry?.available ?? false) &&
            (entry?.values.includes(value) ?? false)
        );
    }

    firstMatch(key: string): string {
        const match = Array.from(this.data).find(
            (v) => v[1].available && v[0].startsWith(key)
        );
        return match ? match[0] : "";
    }

    firstMatchIn(key: string, value: string): string | null {
        const entry = this.data.get(key);
        return entry?.available ?? false
            ? entry?.values.find((v) => v.startsWith(value)) ?? null
            : null;
    }
}
