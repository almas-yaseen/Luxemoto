var ctx = document.getElementById('chart-bar').getContext('2d');
var chart = new Chart(ctx, {
  type: 'bar',
  data: {
    labels: ['Premium cars', 'Mini cars', 'enquired Customers'],
    datasets: [{
      label: 'Count',
      data: [{{ .premiumcarstotal }}, {{ .minicarstotal }}, {{ .enquiredcustomers }}],
backgroundColor: 'rgba(75, 192, 192, 0.2)',
  borderColor: 'rgba(75, 192, 192, 1)',
    borderWidth: 1
  }]
},
options: {
  scales: {
    y: {
      beginAtZero: true
    }
  }
}
});
// Ensure the canvas ID matches


const ctxs = document.getElementById('brandDoughnutChart').getContext('2d');
console.log(ctxs)
const brandCounts = JSON.parse('{{ .brandCountsJ }}'); // Retrieve brandCounts from the template
const data = {
  labels: Object.keys(brandCounts),
  datasets: [{
    data: Object.values(brandCounts),
    backgroundColor: [
      'rgba(255, 99, 132, 0.2)',
      'rgba(54, 162, 235, 0.2)',
      'rgba(255, 206, 86, 0.2)',
      'rgba(75, 192, 192, 0.2)'
    ],
    borderColor: [
      'rgba(255, 99, 132, 1)',
      'rgba(54, 162, 235, 1)',
      'rgba(255, 206, 86, 1)',
      'rgba(75, 192, 192, 1)'
    ],
    borderWidth: 1
  }]
};

const config = {
  type: 'bar',
  data: data,
  options: {
    responsive: true,
    plugins: {
      legend: {
        position: 'top',
      },
      title: {
        display: true,
        text: 'Car Distribution by Brand'
      }
    }
  },
};

const brandDoughnutChart = new Chart(ctxs, config);