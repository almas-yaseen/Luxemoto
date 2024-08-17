function addImageInput() {
    const container = document.getElementById('imageInputsContainer');
    const newInput = document.createElement('input');
    newInput.type = 'file';
    newInput.name = 'images[]';
    newInput.classList.add('form-control-file', 'image-input');
    newInput.accept = 'image/*';
    container.appendChild(newInput);
}
// Event listener to handle click on the "Add Another Image" button
document.getElementById('addImageInput').addEventListener('click', function () {
    addImageInput();
});// Event listener to handle change in file inputs
document.getElementById('imageInputsContainer').addEventListener('change', function (event) {
    if (event.target && event.target.classList.contains('image-input')) {
        const files = event.target.files;
        const previewContainer = document.getElementById('imagePreview');
        for (let i = 0; i < files.length; i++) {
            const file = files[i];
            const reader = new FileReader();
            reader.onload = function () {
                const image = new Image();
                image.src = reader.result;
                image.style = 'width: 50px; height: 50px;'
                image.classList.add('square-image');
                previewContainer.appendChild(image);
            }
            reader.readAsDataURL(file);
        }
    }
});