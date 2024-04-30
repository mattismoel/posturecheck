from sense_hat import SenseHat
import time
import requests

s = SenseHat()

LED_MATRIX_SIZE = 8
rotation = 0.0
max_angle = 30.0

time_before_signal = 5.0
time_since_overshoot = 0.0

blink_color = [255,0,0]
delta = 0.0

s.clear()

def blink(sensehat):
    pixels = []
    for i in range(LED_MATRIX_SIZE):
        for j in range(LED_MATRIX_SIZE):
            pixels.append(blink_color)
    sensehat.set_pixels(pixels)

def backpain():
    print("backpain")
    res = requests.post("http://localhost:8080/add")
    print(res.content)
    blink(s)
    quit()


while True:
    t = time.time()

    rotation = s.get_orientation_degrees()["pitch"]
    if rotation-270 >= max_angle:
        time_since_overshoot += delta

    if time_since_overshoot >= time_before_signal:
        time_since_overshoot = 0.0
        backpain()

    delta = time.time() - t


