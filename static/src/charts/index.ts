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
    ChartDataset,
    defaults
} from "chart.js";

Chart.register(LineController,
    LinearScale,
    PointElement,
    LineElement,
    ScatterController,
    Title,
    Tooltip,
    Legend,
);

const getCSSVariable = (v: string) => getComputedStyle(document.documentElement).getPropertyValue(v);

const bgMainColor = getCSSVariable('--bg-main-color').trim();
const fontMainColor = getCSSVariable('--font-main-color').trim();
const fontSecColor = getCSSVariable('--font-sec-color').trim();

defaults.font.size = 16;
defaults.font.color = fontMainColor;
defaults.font.family = 'Bender';
defaults.elements.line.backgroundColor = fontSecColor;
defaults.elements.point.hitRadius = 15;
defaults.elements.point.hoverRadius = 5;
defaults.plugins.tooltip.backgroundColor = bgMainColor;
defaults.plugins.tooltip.titleFont.color = fontSecColor;

interface CustomScatterPoint extends ScatterDataPoint {
    x: number;
    y: number;
    label: string;
}

const ammoTypeChart = () => {
    const ctx = document.getElementById('ammoTypeChart') as HTMLCanvasElement;
    if (ctx === null) return;

    const config: ChartConfiguration<"scatter"> = {
        type: 'scatter',
        data: {
            labels: [],
            datasets: [
                {
                    borderColor: fontMainColor,
                    backgroundColor: fontMainColor,
                    data: []
                }
            ]
        },
        options: {
            legend: {
                display: false
            },
            tooltips: {
                callbacks: {
                    label: ({ dataset, dataPoint, dataIndex }: TooltipItem) =>
                        `${((dataset as ChartDataset<"scatter">).data[dataIndex] as CustomScatterPoint).label} (PEN: ${dataPoint.x}, DMG: ${dataPoint.y})`
                }
            },
            scales: {
                x: {
                    scaleLabel: {
                        display: true,
                        labelString: 'Penetration',
                    },
                    gridLines: {
                        color: 'rgba(150, 136, 103, .1)',
                        drawBorder: false
                    },
                    type: 'linear',
                    position: 'bottom'
                },
                y: {
                    scaleLabel: {
                        display: true,
                        labelString: 'Damage',
                    },
                    gridLines: {
                        color: 'rgba(150, 136, 103, .1)',
                        drawBorder: false
                    }
                }
            }
        }
    };

    const types = document.querySelectorAll<HTMLTableRowElement>('.item-table.ammo tr');
    for (const type of types) {
        if (!type.dataset.name) continue;

        const pen = parseInt(type.dataset.penetration ?? '');
        const dmg = parseInt(type.dataset.damage ?? '');
        const count = parseInt(type.dataset.projectilecount ?? '');
        const data: CustomScatterPoint = {
            x: pen,
            y: count * dmg || dmg,
            label: type.dataset.name,
        };

        config.data.datasets[0].data.push(data);
    }

    const element = document.querySelector<HTMLElement>('.chart.ammo');
    if (element === null) return;

    const intersectionHandler = (entries: IntersectionObserverEntry[]) => entries.forEach(entry => {
        if (entry.isIntersecting) {
            new Chart(ctx, config);
            observer.unobserve(element);
        }
    });

    const intersectionOptions = { threshold: 0.3 };

    const observer = new IntersectionObserver(intersectionHandler, intersectionOptions);

    observer.observe(element);
};

ammoTypeChart();
