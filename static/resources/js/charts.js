/* global Chart */
'use strict';

const bgMainColor = getComputedStyle(document.documentElement).getPropertyValue('--bg-main-color');
const fontMainColor = getComputedStyle(document.documentElement).getPropertyValue('--font-main-color');
const fontSecColor = getComputedStyle(document.documentElement).getPropertyValue('--font-sec-color');

Chart.defaults.global.defaultFontSize = 16;
Chart.defaults.global.defaultFontColor = fontMainColor;
Chart.defaults.global.defaultFontFamily = 'Bender';
Chart.defaults.global.elements.line.backgroundColor = fontSecColor;
Chart.defaults.global.elements.point.hitRadius = 15;
Chart.defaults.global.tooltips.backgroundColor = bgMainColor;
Chart.defaults.global.tooltips.titleFontColor = fontSecColor;

const ammoTypeChart = () => {
  const el = document.getElementById('ammoTypeChart');
  if (el === null) return;

  const ctx = el.getContext('2d');

  const options = {
    type: 'scatter',
    data: {
      labels: [],
      datasets: []
    },
    options: {
      legend: {
        display: false
      },
      tooltips: {
        callbacks: {
          label: (tooltipItem, data) => `${data.labels[tooltipItem.datasetIndex]} (PEN: ${tooltipItem.xLabel}, DMG: ${tooltipItem.yLabel})`
        }
      },
      scales: {
        xAxes: [{
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
        }],
        yAxes: [{
          scaleLabel: {
            display: true,
            labelString: 'Damage',
          },
          gridLines: {
            color: 'rgba(150, 136, 103, .1)',
            drawBorder: false
          }
        }]
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
      borderColor: fontMainColor,
      backgroundColor: fontMainColor,
      data: [{
        x: pen,
        y: count * dmg || dmg
      }]
    };
    options.data.labels.push(data.label);
    options.data.datasets.push(data);
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
