import sys
import subprocess
import base64
import pygame
import requests
import json

from pygame.locals import QUIT

from wand.image import Image

pygame.init()

WIDTH = 300
HEIGHT = 300

RED = pygame.Color(255, 0, 0)
GREEN = pygame.Color(0, 191, 0)
BLUE = pygame.Color(0, 0, 255)

UR_HOST = "http://192.168.5.105:8080/"

session = requests.Session()

def to_base64(filename):
    with open(filename, "rb") as f:
        contents = f.read()
        return base64.b64encode(contents).decode("ascii")


def match_user(filename):
    payload = {"template": to_base64(filename)}
    return session.post(f"{UR_HOST}/user/fetch", json=payload)


def game_loop(fullscreen: bool):

    if fullscreen:
        screen = pygame.display.set_mode((0, 0), pygame.FULLSCREEN)
    else:
        screen = pygame.display.set_mode((WIDTH, HEIGHT))


    while True:
        screen.fill(BLUE)
        pygame.display.update()
        for event in pygame.event.get():
            if event.type == QUIT:
                pygame.quit()
                sys.exit()

        retval = create_fp_image("dummy")
        if retval == 0:
            image_converter("img.pgm")
            response = match_user("img.pgm.png").json()
            if len(response["usuario"]) > 0:
                screen.fill(GREEN)
            else:
                screen.fill(RED)

            pygame.display.update()
            pygame.time.wait(2000)

        pygame.display.update()


def image_converter(filename: str):
    source = Image(filename=filename)
    destination = source.convert("png")
    destination.save(filename=f"{filename}.png")


def create_fp_image(fp_filename):
    return subprocess.call(["./urfp.1", "--create-image"])


if __name__ == "__main__":
    fullscreen = False
    if len(sys.argv) == 2:
        fullscreen = sys.argv[1] == "f"

    game_loop(fullscreen)

