#!/usr/bin/python3

import time

import paho.mqtt.client as mqtt
import serial.tools.list_ports

CLOUD_SERIALS = [
    '2E8A:000A', # pico
]

EFFECTS = [
    'white',
    'party',
    'forest',
    'ocean',
    'lightning',
]

def get_port():
    for port in serial.tools.list_ports.comports():
        if port.vid is None or port.pid is None:
            continue
        if '%04X:%04X' % (port.vid, port.pid) in CLOUD_SERIALS:
            return port.device
    return None # no port found

class Cloud:
    def __init__(self, serial, client):
        self.serial = serial
        self.client = client
        self.on = False
        self.effect = 'party'
        self.brightness = 8

    def set_state(self, on):
        if on == 'ON':
            on = True
        elif on == 'OFF':
            on = False
        elif type(on) != bool:
            print('Unknown state:', on)
            return
        self.on = on
        self.update_cloud()
        self.client.publish('cloud/light/state', {True: 'ON', False: 'OFF'}[on], retain=True)

    def set_effect(self, effect):
        if effect not in EFFECTS:
            print('Unknown effect:', effect)
            return
        self.effect = effect
        self.update_cloud()
        self.client.publish('cloud/light/effect', effect, retain=True)

    def set_brightness(self, brightness):
        if brightness < 1 or brightness > 10:
            print('Brightness out of range:', brightness)
            return
        self.brightness = brightness
        self.serial.write(b'b%d' % (self.brightness - 1))
        self.client.publish('cloud/light/brightness', str(brightness), retain=True)

    def update_cloud(self):
        if not self.on:
            self.serial.write(b'D')
        elif self.effect == 'white':
            self.serial.write(b'W')
        elif self.effect == 'party':
            self.serial.write(b'P')
        elif self.effect == 'forest':
            self.serial.write(b'F')
        elif self.effect == 'ocean':
            self.serial.write(b'O')
        elif self.effect == 'lightning':
            self.serial.write(b'L')
        else:
            print('Unknown state/effect:', self.on, self.effect)

    def on_connect(self):
        print('Connected to MQTT broker.')
        self.client.subscribe('cloud/light/switch')
        self.client.subscribe('cloud/light/set_effect')
        self.client.subscribe('cloud/light/set_brightness')
        self.set_state(self.on)
        self.set_effect(self.effect)
        self.set_brightness(self.brightness)
        self.client.publish('cloud/light/online', 'online', retain=True)

    def on_message(self, userdata, msg):
        if msg.topic == 'cloud/light/switch':
            self.set_state(msg.payload.decode())
            return
        if msg.topic == 'cloud/light/set_effect':
            self.set_effect(msg.payload.decode())
            return
        if msg.topic == 'cloud/light/set_brightness':
            self.set_brightness(int(msg.payload.decode()))
            return
        print('Unknown message:', msg.topic, msg.payload)


def main():
    while True:
        try:
            # Open serial port.
            port = get_port()
            if not port:
                print('Looking for port...')
                while not port:
                    time.sleep(5)
                    port = get_port()
            print('Opening port:', port)
            ser = serial.Serial(port)

            # Connect to MQTT broker.
            client = mqtt.Client()
            cloud = Cloud(ser, client)
            client.on_connect = lambda client, userdata, flags, rc: cloud.on_connect()
            client.on_message = lambda client, userdata, msg: cloud.on_message(userdata, msg)
            client.will_set('cloud/light/online', 'offline', retain=True)
            client.connect('localhost')
            client.loop_forever()
        except serial.serialutil.SerialException:
            print('Lost serial connection.')
            client.publish('cloud/light/online', 'offline', retain=True)
            client.disconnect()
            client.loop_stop()

if __name__ == '__main__':
   main()
