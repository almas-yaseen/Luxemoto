document.getElementById('bannerImageInput').addEventListener('change', function(event) {
    const imagePreview = document.getElementById('PreviousImage');
    imagePreview.innerHTML = ''; // Clear previous preview

    const file = event.target.files[0];
    if (file) {
        const reader = new FileReader();
        reader.onload = function(e) {
            const img = document.createElement('img');
            img.src = e.target.result;
            img.style = 'width:50px;height:50px'
            imagePreview.appendChild(img);
        };
        reader.readAsDataURL(file);
    }
});



document.addEventListener('DOMContentLoaded', function() {
const vehicleTypeSelect = document.getElementById('vehicleType');
const bannerImageField = document.getElementById('bannerImageField');

vehicleTypeSelect.addEventListener('change', function() {
console.log("Vehicle Type Selected:", vehicleTypeSelect.value); // This should log the selected value
if (vehicleTypeSelect.value === 'Premium') {
    bannerImageField.style.display = 'block';  // Show the banner image field
} else {
    bannerImageField.style.display = 'none';   // Hide the banner image field
}
});

// Trigger change event on page load to ensure the correct field is shown/hidden
vehicleTypeSelect.dispatchEvent(new Event('change'));
});