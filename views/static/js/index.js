// get element by class
const june = document.getElementById('june');

june.addEventListener('mouseover', (event) => {
    // set src to gif
    june.src = '/imgs/june.gif';
    
});

june.addEventListener('mouseleave', () => {
    // set src to png
    june.src = '/imgs/june.png';
});
