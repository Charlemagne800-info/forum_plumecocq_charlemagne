const Hamburger = document.querySelector(".hamburger")
const navLinks = document.querySelector(".links")
    
Hamburger.addEventListener('click',()=>{
navLinks.classList.toggle('mobile_menu')
})