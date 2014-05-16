//
// Text input
//

var shiftDown = false;
$("#chat-input").keyup(function(event){
    if(event.keyCode === 16){
        shiftDown = false;
    }
});
$("#chat-input").keydown(function(event){
    if(event.keyCode === 16){
        shiftDown = true;
    } else if (event.keyCode === 13) {
        if (shiftDown) {
            var input = $('#chat-input');
            var currentRows = parseInt(input.attr('rows'));
            input.attr('rows', currentRows + 1);
        } else {
            $("#send-button").click();
        }
    }
});
$("#chat-input").keyup(function(event) {
    if (event.keyCode === 13 && shiftDown === false) {
        var inputBox = $('#chat-input');
        var inputBoxClone = $('.expanding-clone');
        inputBox.val('');
        inputBox.change();
        inputBoxClone.val('');
        inputBox.focus();
    }
});


//
// Join modal handler
//
$("#username-field").keyup(function(event){
    if(event.keyCode === 13){
        $("#join-button").click();
        $('#chat-input').focus();
    }
});


//
// Window resizing
//

function resize() {
    $('.contents-row').css('height', Math.ceil(window.innerHeight * 0.8) + 'px');
}


//
// Autorun/init functions
//
(function() {
    $('.join-modal').modal();
    $('#username-field').focus();
    $('#chat-input').expanding();

    resize();
    $( window ).resize(resize);
})();
