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
                        drawBorder: false,
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
                        drawBorder: false,
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
