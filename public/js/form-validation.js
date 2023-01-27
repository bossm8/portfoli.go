(function () {
    const form = document.getElementById('contact-form');
    const button = document.getElementById("runaway");

    form.addEventListener('submit', validate);

    function validate(event) {
        if (!form.checkValidity()) {
            event.preventDefault()
            event.stopPropagation()
            button.addEventListener('mouseover', run);
            button.addEventListener('click', run);
        }
        this.classList.add('was-validated')
    }

    const randInt = (max) => {
      return Math.floor(Math.random() * (max + 1));
    };

    function run() {
        if (form.checkValidity()) {
            button.removeEventListener('click', run);
            button.removeEventListener('mouseover', run);
            return
        }
        console.log(button.parentElement.clientWidth)
        const left = randInt(button.parentElement.clientWidth - this.offsetWidth);
        anime({
          targets: this,
          ['left']: `${left}px`,
          easing: "easeOutCirc"
        }).play();
    }

})();


