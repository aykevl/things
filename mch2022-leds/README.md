# TinyGo example for the MCH2022 badge

This is a very simple example how you can write a program for the MCH2022 badge using TinyGo.

Compile it using the following command:

    tinygo build -o main.bin -target=mch2022

This creates a `main.bin` file, which is ready to be sent to the badge. You can do it like this:

    python3 mch2022-template-app/tools/webusb_push.py "TinyGo LED example" main.bin --run

You can find the file `webusb_push.py` [here](https://github.com/badgeteam/mch2022-tools/blob/bef7edfe709f89d9d601de7dde61b31fe5317854/webusb_push.py).
