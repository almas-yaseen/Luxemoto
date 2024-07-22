!function(e){"use strict";if(e(".menu-item.has-submenu .menu-link").on("click",function(s){s.preventDefault(),e(this).next(".submenu").is(":hidden")&&e(this).parent(".has-submenu").siblings().find(".submenu").slideUp(200),e(this).next(".submenu").slideToggle(200)}),e("[data-trigger]").on("click",function(s){s.preventDefault(),s.stopPropagation();var n=e(this).attr("data-trigger");e(n).toggleClass("show"),e("body").toggleClass("offcanvas-active"),e(".screen-overlay").toggleClass("show")}),e(".screen-overlay, .btn-close").click(function(s){e(".screen-overlay").removeClass("show"),e(".mobile-offcanvas, .show").removeClass("show"),e("body").removeClass("offcanvas-active")}),e(".btn-aside-minimize").on("click",function(){window.innerWidth<768?(e("body").removeClass("aside-mini"),e(".screen-overlay").removeClass("show"),e(".navbar-aside").removeClass("show"),e("body").removeClass("offcanvas-active")):e("body").toggleClass("aside-mini")}),e(".select-nice").length&&e(".select-nice").select2(),e("#offcanvas_aside").length){const e=document.querySelector("#offcanvas_aside");new PerfectScrollbar(e)}e(".darkmode").on("click",function(){e("body").toggleClass("dark")})}(jQuery);


                // Function to handle file input change event
                            

                document.getElementById("productImg").addEventListener("change", previewImages);

                const form = document.getElementById("productForm");


                document.addEventListener("DOMContentLoaded", function () {
                    let formData = [];
                    const productForm = document.getElementById("productForm");
                    productForm.addEventListener("submit", async (event) => {
                        if (!validateForm()) {
                            event.preventDefault();
                        } else {
                            event.preventDefault();
                            const checkboxes = document.querySelectorAll('input[type="checkbox"]');
                            const imageData = new FormData();

                            for (let i = 0; i < checkboxes.length; i++) {
                                if (!checkboxes[i].checked) {
                                    imageData.append("productImg", document.getElementById("productImg").files[i]);
                                }
                            }

                            const productName = document.getElementById("productName").value;
                            const productDescription = document.getElementById("productDiscription").value;
                            const regularPrice = document.getElementById("regularPrice").value;
                            const sellingPrice = document.getElementById("sellingPrice").value;
                            const productStatus = document.getElementById("productStatus").value;
                            const category = document.getElementById("category").value;
                            const productBrand = document.getElementById("productBrand").value;
                            const productSize = document.getElementById("productSize").value;
                            const productStock = document.getElementById("productStock").value;
                            imageData.append("productName", productName);
                            imageData.append("productDiscription", productDescription);
                            imageData.append("regularPrice", regularPrice);
                            imageData.append("sellingPrice", sellingPrice);
                            imageData.append("productStatus", productStatus);
                            imageData.append("productCategory", category);
                            imageData.append("productBrand", productBrand);
                            imageData.append("productSize", productSize);
                            imageData.append("productStock", productStock);

                            const response = await fetch(`/admin/postadd-product`, {
                                method: 'POST',
                                body: imageData,
                            });

                            const data = await response.json();
                            if (data.message) {
                                Swal.fire({
                                    position: "top-center",
                                    icon: "success",
                                    title: data.message,
                                    showConfirmButton: false,
                                    timer: 1500
                                });
                                setTimeout(() => {
                                    window.location.href = "/admin/product-list"
                                }, 1500);

                            } else if (data.error) {
                                Swal.fire({
                                    icon: "error",
                                    title: data.error,
                                    text: "OOps" || "Unknown Error",
                                    footer: '<a href="#">Why do I have this issue?</a>',
                                    timer: 1500
                                });
                            }

                        }
                    });

                    function validateForm() {
                        const productName = document.getElementById("productName").value;
                        const productDescription = document.getElementById("productDiscription").value;
                        const regularPrice = document.getElementById("regularPrice").value;
                        const sellingPrice = document.getElementById("sellingPrice").value;
                        const productStatus = document.getElementById("productStatus").value;
                        const category = document.getElementById("category").value;
                        const productBrand = document.getElementById("productBrand").value;
                        const productSize = document.getElementById("productSize").value;
                        const productStock = document.getElementById("productStock").value;
                        const checkboxes = document.querySelectorAll('input[type="checkbox"]');
                        let anyCheckboxChecked = 0;

                        if (isNaN(regularPrice) || isNaN(sellingPrice) || isNaN(productSize) || isNaN(productStock)) {
                            Swal.fire({
                                icon: 'error',
                                title: 'Oops...',
                                text: 'Price, Size, and Stock must be numeric values!',
                            });
                            return false;
                        }

                        // Validate business rules
                        if (regularPrice <= 0 || sellingPrice <= 0) {
                            Swal.fire({
                                icon: 'error',
                                title: 'Oops...',
                                text: 'Price must be more than zero!',
                            });
                            return false;
                        }

                        if (sellingPrice >= regularPrice) {
                            Swal.fire({
                                icon: 'error',
                                title: 'Oops...',
                                text: 'Selling Price must be lesser than the Actual Price!',
                            });
                            return false;
                        }

                        if (productSize <= 0 || productStock < 0) {
                            Swal.fire({
                                icon: 'error',
                                title: 'Oops...',
                                text: 'Stock must be above zero or zero',
                            });
                            return false;
                        }

                        if (!productName.trim() || !productDescription.trim() || !regularPrice.trim() || !sellingPrice.trim() || !productStatus.trim() || !category.trim() || !productBrand.trim() || !productSize.trim() || !productStock.trim()) {
                            Swal.fire({
                                icon: 'error',
                                title: 'Oops...',
                                text: 'All fields are required!',
                            });
                            return false;
                        }
                        for (const checkbox of checkboxes) {
                            if (checkbox.checked) {
                                anyCheckboxChecked = anyCheckboxChecked + 1;
                            }
                        }
                        if (checkboxes.length === anyCheckboxChecked) {
                            Swal.fire({
                                icon: 'error',
                                title: 'All Image Deselected',
                                text: 'Please leave atleast one image',
                            });
                            return false;
                        }

                        return true; // Submit the form
                    }
                });