import {
    Chart,
    LineController,
    LinearScale,
    ScatterController,
    Title,
    Tooltip,
    Legend,
    ChartConfiguration,
    ScatterDataPoint,
    PointElement,
    LineElement,
    TooltipItem,
    defaults,
    LineOptions,
    LineHoverOptions,
    PointOptions,
    PointHoverOptions,
    FontSpec,
} from "chart.js";

Chart.register(
    LineController,
    LinearScale,
    PointElement,
    LineElement,
    ScatterController,
    Title,
    Tooltip,
    Legend
);

const getCSSVariable = (v: string) =>
    getComputedStyle(document.documentElement).getPropertyValue(v);

const bgMainColor = getCSSVariable("--bg-main-color").trim();
const fontMainColor = getCSSVariable("--font-main-color").trim();
const fontSecColor = getCSSVariable("--font-sec-color").trim();

defaults.color = fontMainColor;
defaults.font = <FontSpec>{
    ...defaults.font,
    size: 16,
    family: "Bender",
};
defaults.elements.line = <LineOptions & LineHoverOptions>{
    ...defaults.elements.line,
    backgroundColor: fontSecColor,
};
defaults.elements.point = <PointOptions & PointHoverOptions>{
    ...defaults.elements.point,
    hitRadius: 15,
    hoverRadius: 5,
};

interface CustomScatterPoint extends ScatterDataPoint {
    x: number;
    y: number;
    label: string;
}

const ammoRangeChart = () => {
    const ctx = document.getElementById("ammoRangeChart") as HTMLCanvasElement;
    if (ctx === null) return;

    const labels: string[] = [];
    const damage: number[] = [];
    const penetration: number[] = [];

    const config: ChartConfiguration<"line"> = {
        type: "line",
        data: {
            labels,
            datasets: [
                {
                    label: "Damage",
                    borderColor: fontMainColor,
                    backgroundColor: fontMainColor,
                    tension: 0.1,
                    data: damage,
                },
                {
                    label: "Penetration Power",
                    borderColor: fontMainColor,
                    backgroundColor: fontMainColor,
                    tension: 0.1,
                    borderDash: [5, 5],
                    data: penetration,
                },
            ],
        },
        options: {
            plugins: {
                legend: {
                    display: true,
                },
                tooltip: {
                    backgroundColor: bgMainColor,
                    titleColor: fontSecColor,
                    callbacks: {
                        label: ({
                            dataset,
                            formattedValue,
                        }: TooltipItem<"line">) => {
                            return `${formattedValue} ${dataset.label}`;
                        },
                    },
                },
            },
            interaction: {
                mode: "nearest",
            },
            scales: {
                x: {
                    title: {
                        display: true,
                        text: "Range (m)",
                    },
                    grid: {
                        color: "rgba(150, 136, 103, .1)",
                    },
                    border: {
                        display: false,
                    },
                    ticks: {
                        sampleSize: 10,
                    },
                    type: "linear",
                    position: "bottom",
                },
                y: {
                    title: {
                        display: true,
                        text: "Value",
                    },
                    grid: {
                        color: "rgba(150, 136, 103, .1)",
                    },
                    border: {
                        display: false,
                    },
                },
            },
        },
    };

    const range = document.querySelectorAll<HTMLTableRowElement>(
        ".item-table.ammo tr"
    );
    for (const distance of range) {
        const cells = distance.cells;
        if (!cells[0].dataset.value) continue;

        const label = cells[0].dataset.value;
        const dmg = parseFloat(cells[2].dataset.value ?? "");
        const pen = parseFloat(cells[3].dataset.value ?? "");

        labels.push(label);
        damage.push(dmg);
        penetration.push(pen);
    }

    const element = document.querySelector<HTMLElement>(".chart.ammo");
    if (element === null) return;

    const intersectionHandler = (entries: IntersectionObserverEntry[]) =>
        entries.forEach((entry) => {
            if (entry.isIntersecting) {
                new Chart(ctx, config);
                observer.unobserve(element);
            }
        });

    const intersectionOptions = { threshold: 0.3 };

    const observer = new IntersectionObserver(
        intersectionHandler,
        intersectionOptions
    );

    observer.observe(element);
};

ammoRangeChart();

const ammoTypeChart = () => {
    const ctx = document.getElementById("ammoTypeChart") as HTMLCanvasElement;
    if (ctx === null) return;

    const points: CustomScatterPoint[] = [];

    const config: ChartConfiguration<"scatter"> = {
        type: "scatter",
        data: {
            labels: [],
            datasets: [
                {
                    borderColor: fontMainColor,
                    backgroundColor: fontMainColor,
                    data: points,
                },
            ],
        },
        options: {
            plugins: {
                legend: {
                    display: false,
                },
                tooltip: {
                    backgroundColor: bgMainColor,
                    titleColor: fontSecColor,
                    callbacks: {
                        label: ({
                            dataset,
                            dataIndex,
                        }: TooltipItem<"scatter">) => {
                            const data = dataset.data[
                                dataIndex
                            ] as CustomScatterPoint;
                            return `${data.label} (PEN: ${data.x}, DMG: ${data.y})`;
                        },
                    },
                },
            },
            scales: {
                x: {
                    title: {
                        display: true,
                        text: "Penetration",
                    },
                    grid: {
                        color: "rgba(150, 136, 103, .1)",
                    },
                    border: {
                        display: false,
                    },
                    type: "linear",
                    position: "bottom",
                },
                y: {
                    title: {
                        display: true,
                        text: "Damage",
                    },
                    grid: {
                        color: "rgba(150, 136, 103, .1)",
                    },
                    border: {
                        display: false,
                    },
                },
            },
        },
    };

    const types = document.querySelectorAll<HTMLTableRowElement>(
        ".item-table.ammo tr"
    );
    for (const type of types) {
        if (!type.dataset.name) continue;

        const pen = parseInt(type.dataset.penetration ?? "");
        const dmg = parseInt(type.dataset.damage ?? "");
        const count = parseInt(type.dataset.projectilecount ?? "");
        const data: CustomScatterPoint = {
            x: pen,
            y: count * dmg || dmg,
            label: type.dataset.name,
        };

        points.push(data);
    }

    const element = document.querySelector<HTMLElement>(".chart.ammo");
    if (element === null) return;

    const intersectionHandler = (entries: IntersectionObserverEntry[]) =>
        entries.forEach((entry) => {
            if (entry.isIntersecting) {
                new Chart(ctx, config);
                observer.unobserve(element);
            }
        });

    const intersectionOptions = { threshold: 0.3 };

    const observer = new IntersectionObserver(
        intersectionHandler,
        intersectionOptions
    );

    observer.observe(element);
};

ammoTypeChart();
