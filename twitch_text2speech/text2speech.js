var tmi = require("tmi.js");

// eSpeak\command_line\espeak.exe "hello world"
// speakText("Kappa");
function speakText(text) {
    var exec = require('child_process').exec;
    exec("eSpeak\\command_line\\espeak.exe -g 10 -v en+f4 --path=\"eSpeak\" \""+text+"\"", function callback(error, stdout, stderr){
        console.log(text);
    });
}

var options = {
    options: {
        debug: false
    },
    connection: {
        reconnect: true
    },
    identity: {
        username: "YOUR_TWITCH_USERNAME",
        password: "oauth:YOUR TWITCH OAUTH PASSWORD" // Generaete your's in: http://twitchapps.com/tmi/. See https://docs.tmijs.org/ for more info.
    },
    channels: ["#stream1", "#stream2"]
};

var client = new tmi.client(options);

client.on("chat", function (channel, userstate, message, self) {
    // Don't listen to my own messages..
    if (self) return;
    speakText(userstate.username + " says: " + message + "");
});

client.connect();
