/* global Chart */

const getCSSVariable = v => getComputedStyle(document.documentElement).getPropertyValue(v);

const bgMainColor = getCSSVariable('--bg-main-color').trim();
const fontMainColor = getCSSVariable('--font-main-color').trim();
const fontSecColor = getCSSVariable('--font-sec-color').trim();

Chart.defaults.font.size = 16;
Chart.defaults.font.color = fontMainColor;
Chart.defaults.font.family = 'Bender';
Chart.defaults.elements.line.backgroundColor = fontSecColor;
Chart.defaults.elements.point.hitRadius = 15;
Chart.defaults.elements.point.hoverRadius = 5;
Chart.defaults.plugins.tooltip.backgroundColor = bgMainColor;
Chart.defaults.plugins.tooltip.titleFontColor = fontSecColor;

const ammoTypeChart = () => {
  const el = document.getElementById('ammoTypeChart');
  if (el === null) return;

  const ctx = el.getContext('2d');

  const options = {
    type: 'scatter',
    data: {
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
          label: ({ dataIndex, dataPoint, dataset }) => `${dataset.data[dataIndex].label} (PEN: ${dataPoint.x}, DMG: ${dataPoint.y})`
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

  const types = document.querySelectorAll('.item-table.ammo tr');
  for (const type of types) {
    if (!type.dataset.name) continue;
    const pen = parseInt(type.dataset.penetration);
    const dmg = parseInt(type.dataset.damage);
    const count = parseInt(type.dataset.projectilecount);
    const data = {
      label: type.dataset.name,
      x: pen,
      y: count * dmg || dmg
    };
    options.data.datasets[0].data.push(data);
  }

  const element = document.querySelector('.chart.ammo');

  const intersectionHandler = entries => entries.forEach(entry => {
    if (entry.isIntersecting) {
      new Chart(ctx, options);
      observer.unobserve(element);
    }
  });

  const intersectionOptions = { threshold: 0.3 };

  const observer = new IntersectionObserver(intersectionHandler, intersectionOptions);

  observer.observe(element);
};

ammoTypeChart();
