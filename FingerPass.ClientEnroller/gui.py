import requests
import subprocess
import io

import tkinter as tk
from tkinter import ttk
from PIL import Image as PILImage
from PIL import ImageTk

from wand.image import Image as WandImage

image_buffers = [None, None, None, None, None]
labels = []
images_tk = [None, None, None, None, None]

def submit():
    print(f"Name: {name_var.get()}")
    print(f"Email: {email_var.get()}")
    print(f"Program: {program_var.get()}")


def image_converter(filename: str):
    source = WandImage(filename=filename)
    # destination = source.convert("png")
    return io.BytesIO(source.make_blob("png"))


def read_prints():
    for i in range(5):
        ret_val = subprocess.call(["./urfp.1", "--create-image"])
        if ret_val == 0:
            buff = image_converter("img.pgm")
            image_buffers.append(buff)
            image_tk = ImageTk.PhotoImage(PILImage.open(buff).resize((100,100)))
            images_tk[i] = image_tk
            labels[i].configure(image=image_tk)

# App setup
root = tk.Tk()
root.title("Simple Registration Form")
root.geometry("630x500")
root.configure(bg="#f0f4f8")

# Style
style = ttk.Style()
style.configure("TLabel", background="#f0f4f8", font=("Segoe UI", 12))
style.configure("TEntry", font=("Segoe UI", 12))
style.configure("TButton", font=("Segoe UI", 12, "bold"))

# Variables
name_var = tk.StringVar()
email_var = tk.StringVar()
program_var = tk.StringVar()
images = []

# Form frame
form_frame = ttk.Frame(root, padding="20 20 20 10")
form_frame.pack(fill="x", padx=20, pady=10)

ttk.Label(form_frame, text="Name:").grid(column=0, row=0, sticky="w", pady=5)
ttk.Entry(form_frame, textvariable=name_var, width=40).grid(column=1, row=0, pady=5)

ttk.Label(form_frame, text="Email:").grid(column=0, row=1, sticky="w", pady=5)
ttk.Entry(form_frame, textvariable=email_var, width=40).grid(column=1, row=1, pady=5)

ttk.Label(form_frame, text="Program:").grid(column=0, row=2, sticky="w", pady=5)
ttk.Entry(form_frame, textvariable=program_var, width=40).grid(column=1, row=2, pady=5)

ttk.Button(form_frame, text="Submit", command=submit).grid(
    column=1, row=3, pady=15, sticky="e"
)
ttk.Button(form_frame, text="Read", command=read_prints).grid(
    column=0, row=3, pady=15, sticky="e"
)

# Placeholder frame
placeholder_frame = ttk.Frame(root, padding="20 10 20 20")
placeholder_frame.pack(fill="both", expand=True, padx=20)

canvas = tk.Canvas(placeholder_frame, bg="#f0f4f8", highlightthickness=0)
canvas.pack(fill="both", expand=True)

# Draw placeholder rectangles (just visuals for now)
padding = 10
width = 100
height = 80

blank_pilimage = PILImage.new("RGB", (20, 20), color="#ddd")
blank_tk = ImageTk.PhotoImage(blank_pilimage)

for i in range(5):
    lbl = tk.Label(canvas, image=blank_tk)
    lbl.grid(row=0, column=i, padx=5, pady=5)
    labels.append(lbl)


root.mainloop()

