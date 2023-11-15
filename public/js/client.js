let main = document.querySelector('main')
let message = document.querySelector('#message').innerHTML
let sse = new EventSource("http://127.0.0.1:5973/chat");

function getRandomInt(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

sse.onmessage = function(e) {
    let msg = JSON.parse(e.data);
    let cls

    switch (getRandomInt(1, 3)) {
        case 1:
            cls = 'terran'
            break;
        case 2:
            cls = 'protoss'
            break;
        case 3:
            cls = 'zerg'
            break;
        default:
            cls = 'tarran'
            break;
    }

    messageSend({
        'class': cls,
        'nickname': msg.user.nickname,
        'text': msg.text,
    })
};

function messageSend(v) {
    main.insertAdjacentHTML("afterbegin", messageFormat(v))
}

function messageFormat(v) {
    let msg = message

    for (let key in v) {
        if (v.hasOwnProperty(key)) {
            msg = msg.replace('{{'+key+'}}', v[key])
        }
    }

    return msg
}