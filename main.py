from sense_hat import SenseHat
import time
import requests

# Størrelsen af LED-matrixen på den fysiske sense hat.
LED_MATRIX_SIZE = 8

# Farven som sense hat blinker, når backpain registreres.
BLINK_COLOR = [255,0,0]

# Den maksimale vinkel, før at et backpain kan registereres.
MAX_ANGLE = 30.0

# Mængde af tid, at brugeren skal være over 'max_angle' før at et backpain
# registreres.
REGISTRATION_DURATION = 5.0

# sense hat's rotation. benytter sig udelukkende af dens 'pitch' del.
rotation = 0.0


# Mængde af tid siden brugeren overskred 'max_angle'. Denne nulstilles når
# brugeren igen er i ordenlig holdning.
time_since_overshoot = 0.0

# Tidsforkellen mellem to frames. Benyttes til korrekt udregning af
# tidsændringer.
delta = 0.0

# Er brugeren over 'max_angle'?
is_user_down = False

# Opsætning af sensehat-variablen 's'.
s = SenseHat()
s.clear()

# Definerer hvordan at en input sensehat skal blinke.
# Hvert "cell" i vores LED-matrix sættes til vores 'blink_color'-variabel.
def blink(sensehat):
    cells = []
    for _ in range(LED_MATRIX_SIZE):
        for _ in range(LED_MATRIX_SIZE):
            cells.append(BLINK_COLOR)

    sensehat.set_pixels(cells)


# Definerer funktionalitet for håndtering af backpain-registrering.
# Et POST request sendes til web-serveren.
def backpain():
    requests.post("http://localhost:8080/add")
    blink(s)


while True:
    t = time.time()

    rotation = s.get_orientation_degrees()["pitch"] - 270
    # Hvis brugeren har dårlig holdning, men retter sig op, sidder er brugeren
    # ikke længere nede.
    if is_user_down:
        if rotation <= MAX_ANGLE:
            is_user_down = False

    # Hvis brugeren ikke er nede, og brugeren overskrider den maksimale vinkel
    # starter vi med at tælle med 'time_since_overshoot'. Hvis
    # 'time_since_overshoot' overskrider vores 'registration_duration',
    # registrerer vi backpain.
    else:
        if rotation >= MAX_ANGLE:
            time_since_overshoot += delta

        if time_since_overshoot >= REGISTRATION_DURATION:
            time_since_overshoot = 0.0
            is_user_down = True
            backpain()

    # Udregning af delta - benyttes til at tælle 'time_since_overshoot'.
    delta = time.time() - t
