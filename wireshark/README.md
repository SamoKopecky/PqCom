# About

This a custom dissector for wireshark. However this dissector only works when the PqCom protocol is running on the default port `4040`.

## How to use

Run wireshark with

```sh
wireshark -X lua_script:dissector.lua
```

in this directory and search for the `PQCOM` protocol in the wireshark filter field.
