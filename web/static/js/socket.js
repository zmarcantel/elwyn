//
// Utilities
//

function getHostname() {
    var host = window.location.host;
    var portIndex = host.indexOf(':');
    return (portIndex > 0) ? (host.substring(0, portIndex)) : host;
}

function addErrorAlert(message) {
    $('.main-container').prepend('<div class="alert alert-short alert-danger">' + message + '</div>');
}

function insertMessage(author, type, message, icon) {
    var tree = '<div class="message-row row">';

    if (icon !== '') {
        tree += '\
            <div class="identicon-wrapper col-xs-2 col-md-2 col-lg-1 col-xl-1">\
                <img class="identicon" src="' + icon + '"></img>\
            </div>';
    }

    tree += '\
            <div class="message-wrapper col-xs-10 col-md-10 col-lg-11 col-xl-11">\
                <p class="' + type + '"><i><b>' + author + ':</b> ' + message + '</i></p>\
            </div>\
        </div>';

    $('#channel-contents').append(tree);
}

//
// Socket manager
//


function ElwynSocket(address) {
    var self = this;
    this.host = address;
    this.socket = undefined;
    this.user = {
        username: undefined
    };

    if (this.host === undefined || this.host === '' || typeof this.host !== 'string') {
        console.log('ERROR: tried to connect websocket with invalid host');
        console.log('Defaulting to host of client code.');
        this.host = getHostname() + ':7654';
    }

    // open the web socket
    this.socket = new WebSocket('ws://' + this.host + '/chat');

    // socket open handler
    this.socket.onopen = function(event) {
        console.log("connected!");
    };

    // socket error handler
    this.socket.onerror = function(event) {
        console.log('Got error');
        addErrorAlert(event.data);
    };

    // message reception handler
    this.socket.onmessage = function(event) {
        var response = JSON.parse(event.data);

        if (response.action === 'ACK') {
            self.ack(response);
            return;
        } else if (response.action === "heartbeat") {
            self.heartbeat(response);
            return;
        }

        if (self.user === undefined || self.user.joined === false) {
            console.log('ERROR: client must not have acknowledged joining room');
            return;
        }

        insertMessage(
            response.mine ? 'Me' : response.sender,
            response.mine ? 'mine' :
                (response.sender === 'server' ? 'server-message' : 'general-message'),
            response.body,
            response.icon
        );
    };

    // socket closure handler
    this.socket.onclose = function(event) {
        console.log('Closed websocket');
        addErrorAlert('Lost connection... <a href="/">Reload</a>');
    };

    return this;
}

ElwynSocket.prototype.join = function() {
    var field = $('#username-field');
    if (field.val() === '') {
        $('.join-modal').prepend('<div class="alert alert-danger">You must supply a username</div>');
        return;
    }

    this.user.username = field.val();
    this.socket.send(JSON.stringify({
        action: 'join',
        sender: this.user.username
    }));
};

ElwynSocket.prototype.ack = function(response) {
    if (response.body === 'exists') {
        $('.join-modal').prepend('<div class="alert alert-danger">Username in use</div>');
    } else {
        $('.join-modal').modal('hide');
        this.user.joined = true;
        insertMessage('server', 'server-message', 'welcome, ' + this.user.username, '');
        return;
    }
};

ElwynSocket.prototype.heartbeat = function(response) {
    if (response.body !== "ping") {
        addErrorAlert('Received invalid heartbeat message. Potential server error. <a href="/">Reload?</a>');
        return;
    }

    this.socket.send(JSON.stringify({
        sender: this.user.username,
        action: "heartbeat",
        body: "pong"
    }));
    return;
};

ElwynSocket.prototype.send = function() {
    var message = $.trim($('#chat-input').val());
    if (message === '') { return; }

    this.socket.send(JSON.stringify({
        action: 'message',
        sender: this.user.username,
        body: message
    }));
};
