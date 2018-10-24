Twitch chat text to speech
===

Reads twitch chat out loud for you.

* Fetches chat text using tmi.js library: https://docs.tmijs.org/
* Executes `espeak.exe` for each chat line received
* Using `espeak.exe` means it works only on Windows but could easily be ported
* Bad audio quality
* Supports reading multiple simultaneous chats for the true chaotic experience

### Usage

If I recall correctly, you should fill this part in `text2speech.js` file with your login, password and desired twitch streams:

```Javascript
var options = {
	// ...
    identity: {
        username: "YOUR_TWITCH_USERNAME",
        password: "oauth:YOUR TWITCH OAUTH PASSWORD" // Generaete your's in: http://twitchapps.com/tmi/. See https://docs.tmijs.org/ for more info.
    },
    channels: ["#stream1", "#stream2"] // Can't remember if # prefix is mandatory, but it works with it so.
};
```

Then run it in console using node:

`node text2speech.js`

### License

IDGAF