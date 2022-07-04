function activePopup(e, i){
    e.preventDefault();
    popups[i].classList.add('active');
}

function deActivePopup(i){
    popups[i].classList.remove('active');
}

    let openBtns = document.getElementsByClassName("button_pop");
    let closeBtns = document.getElementsByClassName("pop_up_close");
    let popups = document.getElementsByClassName("pop_up");

    for(let i = 0; i<closeBtns.length; i++){
    closeBtns[i].addEventListener("click", function(){deActivePopup(i)});

    openBtns[i].addEventListener("click", function(e){activePopup(e, i)});
}

window.onkeydown = function( event ) {
    for(let i = 0; i<closeBtns.length; i++){
        if ( event.keyCode === 27 ) {
            popups[i].classList.remove('active');
        }
    }
};
