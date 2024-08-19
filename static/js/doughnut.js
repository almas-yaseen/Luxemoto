const brandLabels = {{ .brandLabels | toJson }};
const brandData = {{ .brandData | toJson }};
console.log("brands label come on my boss",brandLabels)

var ctx2 = document.getElementById('brandChart').getContext('2d');
var brandChart = new Chart(ctx2, {
  type: 'doughnut',
  data: {
    labels: brandLabels, // X-axis labels
    datasets: [{
      label: 'Number of Vehicles',
      data: brandData, // Data points
      backgroundColor: [
        'rgba(255, 99, 132, 0.2)',
        'rgba(54, 162, 235, 0.2)',
        'rgba(255, 206, 86, 0.2)',
        'rgba(75, 192, 192, 0.2)',
        'rgba(153, 102, 255, 0.2)',
        'rgba(255, 159, 64, 0.2)'
      ],
      borderColor: [
        'rgba(255, 99, 132, 1)',
        'rgba(54, 162, 235, 1)',
        'rgba(255, 206, 86, 1)',
        'rgba(75, 192, 192, 1)',
        'rgba(153, 102, 255, 1)',
        'rgba(255, 159, 64, 1)'
      ],
      borderWidth: 1
    }]
  },
  options: {
    responsive: true,
    plugins: {
      legend: {
        position: 'top'
      },
      title: {
        display: true,
        text: 'Car Distribution by Brand'
      }
    }
  }
});