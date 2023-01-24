(function () {
    var forms = document.querySelectorAll('.needs-validation');
    forms.forEach(f => addValidationListener(f));
})();

function addValidationListener(form) {
    form.addEventListener(
        'submit', 
        function (event) {
            if (!form.checkValidity()) {
                event.preventDefault()
                event.stopPropagation()
                makeButtonRun(true);
            } else {
                makeButtonRun(false);
            }
            form.classList.add('was-validated')
        }, 
        false
    );
}

const randInt = (max) => {
  return Math.floor(Math.random() * (max + 1));
};

function makeButtonRun(runAway) {
    console.log(runAway)
    const button = document.getElementById("runaway");
    if (run) {
        button.addEventListener('mouseover', run);
        button.addEventListener('click', run)
    } else {
        button.removeEventListener('mouseover', run);
        button.removeEventListener('click', run);
    }

}

function run() {
    console.log(this)
    const left = randInt(window.innerWidth - this.offsetWidth);
    console.log(left)
    anime({
      targets: this,
      ['left']: `${left}px`,
      easing: "easeOutCirc"
    }).play();
}