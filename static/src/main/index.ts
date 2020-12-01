import { initSearchSocket } from "./search";

const registerTabs = () => {
    const tablinks = document.getElementsByClassName('tab-btn');
    for (const el of tablinks) el.addEventListener('click', function (this: HTMLElement) {
        const cl = 'active';
        document.querySelector(`.tab.${cl}`)?.classList.remove(cl);
        document.querySelector(`.tab-btn.${cl}`)?.classList.remove(cl);

        document.getElementById(this.dataset.tab?.toString() || '')?.classList.add(cl);
        this.classList.add(cl);
    });
};

const loadImage = async () => {
    const image = document.getElementById('itemImage');
    if (image === null) return;

    const staticURL = (() => {
        const host = window.location.host;
        const parts = host.split('.');
        return `//static.${parts.length > 2 ? [parts[parts.length - 2], parts[parts.length - 1]].join('.') : host}`;
    })();
    const imageID = image.dataset.id;
    const request = new Request(`${staticURL}/image/icon/1-1/${imageID}.png`);

    try {
        const response = await fetch(request);
        if (!response.ok) throw new Error(`Error while fetching image: ${response.status}`);

        const objectURL = URL.createObjectURL(await response.blob());
        const img = new Image();
        img.src = objectURL;
        img.onload = () => {
            const imgWidth = img.naturalHeight, imgHeight = img.naturalHeight;
            const boxWidth = image.offsetWidth, boxHeight = image.offsetHeight;

            if (imgHeight > boxHeight || imgWidth > boxWidth) image.style.backgroundSize = 'contain';

            image.style.backgroundImage = `url("${objectURL}")`;
        };
    } catch (err) {
        if (err.message === '404') {
            console.warn('Image can not be found');
        } else {
            console.error('Image load failed:', err.message);
        }
    }
};

const sortTables = () => {
    const tables = document.querySelectorAll<HTMLTableHeaderCellElement>('.sort-table.client-sort thead th');
    if (tables.length === 0) return;

    const getCellValue = (tr: HTMLTableRowElement, idx: number) =>
        (tr.children[idx] as HTMLElement).dataset.value ||
        (tr.children[idx] as HTMLElement).innerText ||
        tr.children[idx].textContent || '';

    const comparer = (idx: number, asc: boolean) =>
        (a: HTMLTableRowElement, b: HTMLTableRowElement) => ((v1: string, v2: string) =>
            !isNaN(Number(v1)) && !isNaN(Number(v2)) ?
                Number(v1) - Number(v2) :
                v1?.toString().localeCompare(v2 as string)
        )(getCellValue(asc ? a : b, idx), getCellValue(asc ? b : a, idx));


    for (const th of tables) {
        th.addEventListener('click', function (this: HTMLElement) {
            const table = this.closest('table');
            const body = table?.getElementsByTagName('tbody')[0];

            const clSorted = 'sorted-by';
            const clAsc = 'up';
            const clDesc = 'down';

            let isASC = this.classList.contains(clAsc);

            Array.from<HTMLTableRowElement>(body?.getElementsByTagName('tr') || [])
                .sort(comparer(Array.from(this.parentNode?.children || []).indexOf(th), isASC = !isASC))
                .forEach(tr => body?.appendChild(tr));

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

        th.style.cursor = 'pointer';
    }
};

const initInteractiveMap = async () => {
    const el = document.getElementById('map');
    if (el === null) return;

    const libPath = (document.getElementById('mapLib') as HTMLScriptElement | null)?.src;
    if (!libPath) {
        console.error("Library is missing");
        return;
    }

    try {
        const map = await import(libPath);
        await map.init(el);
        return map;
    } catch (err) {
        console.error(err);
    }

    return;
};

(async () => {
    registerTabs();
    loadImage();
    sortTables();
    const map = await initInteractiveMap();
    const form = document.getElementById("search") as HTMLFormElement | null;
    if (form) initSearchSocket(form, map);
})();

export { };
